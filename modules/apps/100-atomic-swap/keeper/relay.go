package keeper

import (
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
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

	_, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	//destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	//destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	_, found2 := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found2 {
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

	//packet := channeltypes.NewPacket(
	//	swapPacket.GetBytes(),
	//	sequence,
	//	sourcePort,
	//	sourceChannel,
	//	destinationPort,
	//	destinationChannel,
	//	timeoutHeight,
	//	timeoutTimestamp,
	//)

	_, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, swapPacket.GetBytes())
	if err != nil {
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
		var msg types.SwapMaker
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}

		if err := k.OnReceivedMake(ctx, packet, &msg); err != nil {
			return err
		}

	case types.TAKE_SWAP:
		var msg types.SwapTaker
		if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
			return err
		}

		if err2 := k.OnReceivedTake(ctx, packet, &msg); err2 != nil {
			return err2
		} else {
			return nil
		}

	case types.CANCEL_SWAP:
		var msg types.CancelSwapMsg
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
		case types.MAKE_SWAP:
			// This is the step 4 (Acknowledge Make Packet) of the atomic swap: https://github.com/liangping/ibc/blob/atomic-swap/spec/app/ics-100-atomic-swap/ibcswap.png
			// This logic is executed when Taker chain acknowledge the make swap packet.
			var msg types.SwapMaker
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}

			// check order status
			o := types.NewAtomicOrder(&msg, msg.SourceChannel)
			order, ok := k.GetAtomicOrder(ctx, o.Id)
			if !ok {
				return types.ErrOrderDoesNotExists
				//return nil
			}
			order.Status = types.Status_SYNC
			k.SetAtomicOrder(ctx, order)
			return nil

		case types.TAKE_SWAP:
			// This is the step 9 (Transfer Take Token & Close order): https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap
			// The step is executed on the Taker chain.
			takeMsg := &types.SwapTaker{}
			if err := types.ModuleCdc.Unmarshal(data.Data, takeMsg); err != nil {
				return err
			}

			order, _ := k.GetAtomicOrder(ctx, takeMsg.OrderId)
			escrowAddr := types.GetEscrowAddress(types.PortID, packet.SourceChannel)
			makerReceivingAddr, err := sdk.AccAddressFromBech32(order.Maker.MakerReceivingAddress)
			if err != nil {
				return err
			}

			if err = k.bankKeeper.SendCoins(ctx, escrowAddr, makerReceivingAddr, sdk.NewCoins(takeMsg.SellToken)); err != nil {
				return err
			}

			order.Status = types.Status_COMPLETE
			order.Takers = &types.SwapTaker{
				OrderId:               takeMsg.OrderId,
				SellToken:             takeMsg.SellToken,
				TakerAddress:          takeMsg.TakerAddress,
				TakerReceivingAddress: takeMsg.TakerReceivingAddress,
				CreateTimestamp:       takeMsg.CreateTimestamp,
			}
			order.CompleteTimestamp = takeMsg.CreateTimestamp
			k.SetAtomicOrder(ctx, order)
			return nil
		case types.CANCEL_SWAP:
			// This is the step 14 (Cancel & refund) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap
			// It is executed on the Maker chain.

			var msg types.CancelSwapMsg
			if err := types.ModuleCdc.Unmarshal(data.Data, &msg); err != nil {
				return err
			}

			order, _ := k.GetAtomicOrder(ctx, msg.OrderId)
			escrowAddr := types.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)
			makerAddr, err := sdk.AccAddressFromBech32(order.Maker.MakerAddress)
			if err != nil {
				return err
			}

			if err = k.bankKeeper.SendCoins(ctx, escrowAddr, makerAddr, sdk.NewCoins(order.Maker.SellToken)); err != nil {
				return err
			}
			order.Status = types.Status_CANCEL
			order.CancelTimestamp = msg.CreateTimestamp
			k.SetAtomicOrder(ctx, order)
			return nil
		default:
			return errors.New("unknown data packet")
		}
	}
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data *types.AtomicSwapPacketData) error {
	return k.refundPacketToken(ctx, packet, data)
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.AtomicSwapPacketData) error {
	swapPacket := &types.AtomicSwapPacketData{}
	if err := swapPacket.Unmarshal(packet.GetData()); err != nil {
		return err
	}

	escrowAddr := types.GetEscrowAddress(packet.SourcePort, packet.SourceChannel)

	switch swapPacket.Type {
	case types.MAKE_SWAP:
		// This is the step 3.2 (Refund) of the atomic swap: https://github.com/liangping/ibc/blob/atomic-swap/spec/app/ics-100-atomic-swap/ibcswap.png
		// This logic will be executed when Relayer sends make swap packet to the taker chain, but the request timeout
		// and locked tokens form the first step (see the picture on the link above) MUST be returned to the account of
		// the maker on the maker chain.
		makeMsg := &types.MakeSwapMsg{}
		if err := makeMsg.Unmarshal(swapPacket.Data); err != nil {
			return err
		}

		makerAddr, err := sdk.AccAddressFromBech32(makeMsg.MakerAddress)
		if err != nil {
			return err
		}

		// send tokens back to maker
		err = k.bankKeeper.SendCoins(ctx, escrowAddr, makerAddr, sdk.NewCoins(makeMsg.SellToken))
		if err != nil {
			return err
		}

		orderID := types.GenerateOrderId(packet)
		order, found := k.GetAtomicOrder(ctx, orderID)
		if !found {
			return fmt.Errorf("order not found for ID %s", orderID)
		}
		order.Status = types.Status_CANCEL
		k.SetAtomicOrder(ctx, order)

	case types.TAKE_SWAP:
		// This is the step 7.2 (Unlock order and refund) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap
		// This step is executed on the Taker chain when Take Swap request timeout.
		takeMsg := &types.TakeSwapMsg{}
		if err := takeMsg.Unmarshal(swapPacket.Data); err != nil {
			return err
		}

		takerAddr, err := sdk.AccAddressFromBech32(takeMsg.TakerAddress)
		if err != nil {
			return err
		}

		// send tokens back to taker
		err = k.bankKeeper.SendCoins(ctx, escrowAddr, takerAddr, sdk.NewCoins(takeMsg.SellToken))
		if err != nil {
			return err
		}

		orderID := types.GenerateOrderId(packet)
		order, found := k.GetAtomicOrder(ctx, orderID)
		if !found {
			return fmt.Errorf("order not found for ID %s", orderID)
		}
		order.Takers = nil // release the occupation
		k.SetAtomicOrder(ctx, order)
	case types.CANCEL_SWAP:
		// do nothing, only send tokens back when cancel msg is acknowledged.
	default:
		return errors.New("unknown data packet")
	}

	return nil
}
