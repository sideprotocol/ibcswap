package keeper

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) SendIBCSwapPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	swapPacket types.IBCSwapPacketData,
) error {

	if err := swapPacket.ValidateBasic(); err != nil {
		return err
	}

	if err := swapPacket.ValidateBasic(); err != nil {
		return err
	}

	if !k.GetSwapEnabled(ctx) {
		return types.ErrSwapEnabled
	}

	_, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return errorsmod.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	// get the next sequence
	_, found = k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
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

	if _, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, swapPacket.GetBytes()); err != nil {
		return err
	}

	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.IBCSwapPacketData) ([]byte, error) {
	switch data.Type {
	case types.MAKE_POOL:
		var msg types.MsgMakePoolRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}

		if strings.TrimSpace(data.StateChange.PoolId) == "" {
			return nil, types.ErrEmptyPoolId
		}

		poolId, err := k.OnMakePoolReceived(ctx, &msg, data.StateChange.PoolId)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(&types.MsgMakePoolResponse{PoolId: *poolId})
		return resData, err

	case types.TAKE_POOL:
		var msg types.MsgTakePoolRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		poolId, err := k.OnTakePoolReceived(ctx, &msg)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(&types.MsgMakePoolResponse{PoolId: *poolId})
		return resData, err

	case types.SINGLE_DEPOSIT:
		var msg types.MsgSingleAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnSingleAssetDepositReceived(ctx, &msg, data.StateChange)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(res)
		return resData, err

	case types.MAKE_MULTI_DEPOSIT:
		var msg types.MsgMakeMultiAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnMakeMultiAssetDepositReceived(ctx, &msg, data.StateChange)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(res)
		return resData, err

	case types.TAKE_MULTI_DEPOSIT:
		var msg types.MsgTakeMultiAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnTakeMultiAssetDepositReceived(ctx, &msg, data.StateChange)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(res)
		return resData, err

	case types.MULTI_WITHDRAW:
		var msg types.MsgMultiAssetWithdrawRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnMultiAssetWithdrawReceived(ctx, &msg, data.StateChange)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(res)
		return resData, err

	case types.LEFT_SWAP, types.RIGHT_SWAP:
		var msg types.MsgSwapRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return nil, err
		}
		res, err := k.OnSwapReceived(ctx, &msg, data.StateChange)
		if err != nil {
			return nil, err
		}
		resData, err := types.ModuleCdc.Marshal(res)
		return resData, err

	default:
		return nil, types.ErrUnknownDataPacket
	}
}

// OnAcknowledgementPacket processes the packet acknowledgement and performs actions based on the acknowledgement type
func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapPacketData, ack channeltypes.Acknowledgement) error {
	logger := k.Logger(ctx)

	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		switch data.Type {
		case types.MAKE_POOL:
			var msg types.MsgMakePoolRequest
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug(err.Error())
				return err
			}
			err := k.OnMakePoolAcknowledged(ctx, &msg, data.StateChange.PoolId)
			if err != nil {
				return err
			}

		case types.TAKE_POOL:
			var msg types.MsgTakePoolRequest
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug(err.Error())
				return err
			}
			err := k.OnTakePoolAcknowledged(ctx, &msg)
			if err != nil {
				return err
			}
			return nil
		case types.SINGLE_DEPOSIT:
			var msg types.MsgSingleAssetDepositRequest
			var res types.MsgSingleAssetDepositResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug("Deposit:packet:", err.Error())
				return err
			}

			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				logger.Debug("Deposit:ack:", err.Error())
				return err
			}
			if err := k.OnSingleAssetDepositAcknowledged(ctx, &msg, &res); err != nil {
				logger.Debug("Deposit:Single", err.Error())
				return err
			}

		case types.MAKE_MULTI_DEPOSIT:

			var msg types.MsgMakeMultiAssetDepositRequest
			var res types.MsgMultiAssetDepositResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				logger.Debug("MakeMultiDeposit:packet:", err.Error())
				return err
			}

			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				logger.Debug("MakeMultiDeposit:ack:", err.Error())
				return err
			}
			if err := k.OnMakeMultiAssetDepositAcknowledged(ctx, &msg, &res); err != nil {
				logger.Debug("MakeMultiDeposit:Single", err.Error())
				return err
			}
		case types.TAKE_MULTI_DEPOSIT:
			var msg types.MsgTakeMultiAssetDepositRequest
			var res types.MsgTakePoolResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				return err
			}
			if err := k.OnTakeMultiAssetDepositAcknowledged(ctx, &msg); err != nil {
				logger.Debug("TakeMultiDeposit:Single", err.Error())
				return err
			}
			return nil

		case types.MULTI_WITHDRAW:
			var msg types.MsgMultiAssetWithdrawRequest
			var res types.MsgMultiAssetWithdrawResponse
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			if err := types.ModuleCdc.Unmarshal(ack.GetResult(), &res); err != nil {
				return err
			}
			if err := k.OnMultiWithdrawAcknowledged(ctx, &msg, &res); err != nil {
				return err
			}
		case types.LEFT_SWAP, types.RIGHT_SWAP:
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

// OnTimeoutPacket processes a timeout packet and refunds the tokens
func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapPacketData) error {
	return k.refundPacketToken(ctx, packet, data)
}

// refundPacketToken refunds tokens in case of a timeout
func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapPacketData) error {
	var token sdk.Coin
	var sender string

	switch data.Type {
	case types.MAKE_POOL:
		var msg types.MsgMakePoolRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		// Refund initial liquidity
		sender = msg.Creator
		token = *msg.Liquidity[0].Balance

	case types.SINGLE_DEPOSIT:
		var msg types.MsgSingleAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.Token
		sender = msg.Sender
	case types.MAKE_MULTI_DEPOSIT:
		var msg types.MsgMakeMultiAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.Deposits[0].Balance
		sender = msg.Deposits[0].Sender
	case types.TAKE_MULTI_DEPOSIT:
		var msg types.MsgTakeMultiAssetDepositRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		//token = *msg.
		//sender = msg.Sender
	case types.MULTI_WITHDRAW:
		var msg types.MsgMultiAssetWithdrawRequest
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		token = *msg.PoolToken
		sender = msg.Receiver
	case types.RIGHT_SWAP:
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
	err := k.bankKeeper.SendCoins(ctx, escrowAccount, sdk.AccAddress(sender), sdk.NewCoins(token))
	return err
}
