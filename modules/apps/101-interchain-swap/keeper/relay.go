package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) SendIBCSwapPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	swapPacket types.IBCSwapDataPacket,
) error {

	if err := swapPacket.ValidateBasic(); err != nil {
		return err
	}

	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return errorsmod.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return errorsmod.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	packet := channeltypes.NewPacket(
		swapPacket.GetBytes(),
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp)

	if err := k.ics4Wrapper.SendPacket(ctx, channelCap, packet); err != nil {
		return err
	}

	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IBCSwapDataPacket) ([]byte, error) {
	switch data.Type {
	case types.MessageType_CREATE:
		var msg types.MsgCreatePoolRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		pooId, err := k.OnCreatePoolReceived(ctx, &msg, packet.DestinationPort, packet.DestinationChannel)
		if err != nil {
			return nil, err
		}
		data, err := types.ModuleCdc.Marshal(&types.MsgCreatePoolResponse{PoolId: *pooId}) //types.ModuleCdc.Marshal(&types.MsgCreatePoolResponse{PoolId: *pooId})
		return data, err
	case types.MessageType_DEPOSIT:
		var msg types.MsgDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnDepositReceived(ctx, &msg)
		if err != nil {
			return nil, err
		}
		data, err := types.ModuleCdc.Marshal(res) //types.ModuleCdc.Marshal(res)
		return data, err

	case types.MessageType_WITHDRAW:
		var msg types.MsgWithdrawRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err2 := k.OndWithdrawReceive(ctx, &msg)
		if err2 != nil {
			return nil, err2
		}

		data, err := types.ModuleCdc.Marshal(res)
		return data, err

	case types.MessageType_LEFTSWAP, types.MessageType_RIGHTSWAP:
		var msg types.MsgSwapRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnSwapReceived(ctx, &msg)
		if err != nil {
			return nil, err
		}
		data, err := types.ModuleCdc.Marshal(res) //types.ModuleCdc.Marshal(res)
		return data, err
	default:
		return nil, types.ErrUnknownDataPacket
	}
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket, ack channeltypes.Acknowledgement) error {
	logger := k.Logger(ctx)
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		switch data.Type {
		case types.MessageType_CREATE:
			var msg types.MsgCreatePoolRequest
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug(err.Error())
				return err
			}
			k.OnCreatePoolAcknowledged(ctx, &msg)
		case types.MessageType_DEPOSIT:
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				"Deposit Ack 1:",
				sdk.NewAttribute(
					"Deposit",
					data.String(),
				),
			))
			var msg types.MsgDepositRequest
			var res types.MsgDepositResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug("Deposit:packet:", err.Error())
				return err
			}

			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				logger.Debug("Deposit:ack:", err.Error())
				return err
			}
			ctx.EventManager().EmitEvent(sdk.NewEvent(
				"Deposit Ack2:",
				sdk.NewAttribute(
					"Deposit",
					data.String(),
				),
			))
			if err := k.OnSingleDepositAcknowledged(ctx, &msg, &res); err != nil {
				logger.Debug("Deposit:Single", err.Error())
				return err
			}
		case types.MessageType_WITHDRAW:
			var msg types.MsgWithdrawRequest
			var res types.MsgWithdrawResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				return err
			}
			if err := k.OnWithdrawAcknowledged(ctx, &msg, &res); err != nil {
				return err
			}
		case types.MessageType_LEFTSWAP, types.MessageType_RIGHTSWAP:
			var msg types.MsgSwapRequest
			var res types.MsgSwapResponse

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				return err
			}
			if err := k.OnSwapAcknowledged(ctx, &msg, &res); err != nil {
				return err
			}
		}
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket) error {
	return k.refundPacketToken(ctx, packet, data)
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket) error {

	var token sdk.Coin
	var sender string
	switch data.Type {
	case types.MessageType_DEPOSIT:
		var msg types.MsgDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.Tokens[0]
		sender = msg.Sender
	case types.MessageType_WITHDRAW:
		var msg types.MsgWithdrawRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.PoolCoin
		sender = msg.Sender
	case types.MessageType_RIGHTSWAP:
		var msg types.MsgSwapRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.TokenIn
		sender = msg.Sender
	default:
		return types.ErrUnknownDataPacket
	}
	escrowAccount := types.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, escrowAccount.String(), sdk.AccAddress(sender), sdk.NewCoins(token))
	return nil
}
