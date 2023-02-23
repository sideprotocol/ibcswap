package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/31-atomic-swap/types"
)

func (k Keeper) SendSwapPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	swapPacket types.AtomicSwapPacketData,
) error {

	if err := swapPacket.ValidateBasic(); err != nil {
		return err
	}

	if !k.GetSwapEnabled(ctx) {
		return types.ErrSendDisabled
	}

	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	fmt.Println("Timestamp", timeoutTimestamp)
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

	defer func() {
		//if sendingCoin.Amount.IsInt64() {
		//	telemetry.SetGaugeWithLabels(
		//		[]string{"tx", "msg", "ibc", "swap"},
		//		float32(sendingCoin.Amount.Int64()),
		//		[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, "fullDenomPath")},
		//	)
		//}
	}()

	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.AtomicSwapPacketData) error {

	switch data.Type {
	case types.MAKE_SWAP:
		var msg types.MsgMakeSwapRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		if err := k.OnReceivedMake(ctx, packet, &msg); err != nil {
			return err
		}

		return nil

	case types.TAKE_SWAP:
		var msg types.MsgTakeSwapRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		if err2 := k.OnReceivedTake(ctx, packet, &msg); err2 != nil {
			return err2
		} else {
			return nil
		}

	case types.CANCEL_SWAP:
		var msg types.MsgCancelSwapRequest

		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}
		if err2 := k.OnReceivedCancel(ctx, packet, &msg); err2 != nil {
			return err2
		} else {
			return nil
		}

	default:
		return types.ErrUnknownDataPacket
	}

	ctx.EventManager().EmitTypedEvents(&data)

	return nil
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.AtomicSwapPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		switch data.Type {
		case types.TAKE_SWAP:
			var msg types.MsgTakeSwapRequest

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			// check order status
			if order, ok := k.GetLimitOrder(ctx, msg.OrderId); ok {
				k.executeLimitOrderMatch(ctx, order, &msg, StepAcknowledgement)
			} else if orderOtc, ok2 := k.GetOTCOrder(ctx, msg.OrderId); ok2 {
				k.executeOTCOrderMatch(ctx, orderOtc, &msg, StepAcknowledgement)
			} else {
				return types.ErrOrderDoesNotExists
			}
			break

		case types.CANCEL_SWAP:
			var msg types.MsgCancelSwapRequest

			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}
			if err2 := k.executeCancel(ctx, &msg, StepAcknowledgement); err2 != nil {
				return err2
			} else {
				return nil
			}
			break
		}
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.AtomicSwapPacketData) error {
	return k.refundPacketToken(ctx, packet, data)
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.AtomicSwapPacketData) error {

	ctx.Logger().Debug("refundPacketToken: %s", data)

	return nil
}
