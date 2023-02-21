package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
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

	// // get the next sequence
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
		timeoutTimestamp,
	)

	if err := k.ics4Wrapper.SendPacket(ctx, channelCap, packet); err != nil {
		return err
	}
	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IBCSwapDataPacket) error {
	switch data.Type {
	case types.MessageType_CREATE:
		var msg types.MsgCreatePoolRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}

		_, err := k.OnCreatePoolReceived(ctx, &msg, packet.DestinationPort, packet.DestinationChannel)
		if err != nil {
			return err
		}

		return nil

	case types.MessageType_DEPOSIT:
		var msg types.MsgDepositRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		_, err := k.OnDepositReceived(ctx, &msg)
		if err != nil {
			return err
		} else {
			return nil
		}

	case types.MessageType_WITHDRAW:
		var msg types.MsgWithdrawRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		if _, err2 := k.OndWithdrawReceive(ctx, &msg); err2 != nil {
			return err2
		} else {
			return nil
		}

	case types.MessageType_LEFTSWAP, types.MessageType_RIGHTSWAP:
		var msg types.MsgSwapRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		if _, err2 := k.OnSwapReceived(ctx, &msg); err2 != nil {
			return err2
		} else {
			return nil
		}

	default:
		return types.ErrUnknownDataPacket
	}
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		switch data.Type {
		case types.MessageType_CREATE:
			var msg types.MsgCreatePoolRequest
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			k.onCreatePoolAcknowledged(ctx, &msg)
		case types.MessageType_DEPOSIT:
			var msg types.MsgDepositRequest

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			//k.OnDepositReceived()
			//k.onSingleDepositAcknowledged(ctx, &msg,)
			break
		case types.MessageType_WITHDRAW:
			var msg types.MsgWithdrawRequest

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			//k.executeWithdrawRequest(ctx, &msg, StepAcknowledgement)
			break
		case types.MessageType_LEFTSWAP, types.MessageType_RIGHTSWAP:
			var msg types.MsgSwapRequest

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			//k.executeSwapRequest(ctx, &msg, StepAcknowledgement)
			break
		}
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket) error {
	return k.refundPacketToken(ctx, packet, data)
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapDataPacket) error {

	ctx.Logger().Debug("refundPacketToken: %s", data)

	return nil
}
