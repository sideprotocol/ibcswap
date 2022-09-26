package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"
)

var (
	StepSend            = 1
	StepReceive         = 2
	StepAcknowledgement = 3
)
var _ types.MsgServer = Keeper{}

func (k Keeper) MakeSwap(goCtx context.Context, msg *types.MsgMakeSwapRequest) (*types.MsgMakeSwapResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	msgByte, err0 := types.ModuleCdc.Marshal(types.NewMakerFromMsg(msg))
	if err0 != nil {
		return nil, err0
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err1 != nil {
		return nil, err1
	}

	// lock sell token into module
	if err := k.bankKeeper.SendCoinsFromAccountToModule(
		ctx, sender, types.ModuleName, sdk.NewCoins(msg.SellToken),
	); err != nil {
		return nil, err
	}

	if len(msg.DesiredTaker) > 0 {
		order := types.NewOTCOrder(msg, msg.SourceChannel)
		k.SetOTCOrder(ctx, order)
	} else {
		order := types.NewLimitOrder(msg, msg.SourceChannel)
		k.SetLimitOrder(ctx, order)
	}

	packet := types.AtomicSwapPacketData{
		Type: types.MAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}
	if err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgMakeSwapResponse{}, nil
}

func (k Keeper) TakeSwap(goCtx context.Context, msg *types.MsgTakeSwapRequest) (*types.MsgTakeSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// pre-check
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgByte, err0 := types.ModuleCdc.Marshal(types.NewTakerFromMsg(msg))
	if err0 != nil {
		return nil, err0
	}

	// check order status
	if order, ok := k.GetLimitOrder(ctx, msg.OrderId); ok {
		k.executeLimitOrderMatch(ctx, order, msg, StepSend)
	} else if orderOtc, ok2 := k.GetOTCOrder(ctx, msg.OrderId); ok2 {
		k.executeOTCOrderMatch(ctx, orderOtc, msg, StepSend)
	} else {
		return nil, types.ErrOrderDoesNotExists
	}

	packet := types.AtomicSwapPacketData{
		Type: types.TAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}

	if err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgTakeSwapResponse{}, nil
}

func (k Keeper) CancelSwap(goCtx context.Context, msg *types.MsgCancelSwapRequest) (*types.MsgCancelSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgbyte, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.AtomicSwapPacketData{
		Type: types.CANCEL_SWAP,
		Data: msgbyte,
		Memo: "",
	}

	if err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgCancelSwapResponse{}, nil
}

// See createOutgoingPacket in spec:https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay

func (k Keeper) executeLimitOrderMatch(ctx sdk.Context, order types.LimitOrder, msg *types.MsgTakeSwapRequest, step int) error {

	switch order.Status {
	case types.Status_CANCEL:
		return types.ErrOrderCanceled
	case types.Status_COMPLETE:
		return types.ErrOrderCompleted
	default:
		// continue
	}
	if order.Maker.BuyToken.Denom != msg.SellToken.Denom {
		return types.ErrOrderDenominationMismatched
	}
	// calculate how much is filled
	sum := sdk.NewCoin(order.Maker.BuyToken.Denom, sdk.NewInt(0))
	order.Takers = append(order.Takers, types.NewTakerFromMsg(msg))
	for _, taker := range order.Takers {
		sum.Add(taker.SellToken)
	}

	if sum.IsGTE(order.Maker.BuyToken) {
		order.FillStatus = types.FillStatus_COMPLETE_FILL
		order.Status = types.Status_COMPLETE
		order.CompleteTimestamp = msg.CreateTimestamp
	} else {
		order.FillStatus = types.FillStatus_PARTIAL_FILL
	}

	switch step {
	case StepSend:
		// executed on the destination chain of make swap,
		// lock the taker's sell coin into the module
		sender, err1 := sdk.AccAddressFromBech32(msg.TakerAddress)
		if err1 != nil {
			return err1
		}

		// lock sell token into module
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx, sender, types.ModuleName, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	case StepReceive:
		// StepReceive
		// executed on the source chain of make swap,
		// send maker's sell token to taker
		receiver, err1 := sdk.AccAddressFromBech32(msg.TakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// 可能存在精度问题
		amount := int64(float64(msg.SellToken.Amount.Int64()) / float64(order.Maker.BuyToken.Amount.Int64()) * float64(order.Maker.SellToken.Amount.Int64()))

		// send maker's sell token from module to taker
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(sdk.NewCoin(order.Maker.SellToken.Denom, sdk.NewInt(amount))),
		); err != nil {
			return err
		}
		break
	case StepAcknowledgement:
		// executed on the destination chain of the swap,
		// send taker's sell token to maker
		receiver, err1 := sdk.AccAddressFromBech32(order.Maker.MakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// send taker's sell token from module to maker
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	}

	// save updated status
	k.SetLimitOrder(ctx, order)

	return nil
}

func (k Keeper) executeOTCOrderMatch(ctx sdk.Context, order types.OTCOrder, msg *types.MsgTakeSwapRequest, step int) error {

	switch order.Status {
	case types.Status_CANCEL:
		return types.ErrOrderCanceled
	case types.Status_COMPLETE:
		return types.ErrOrderCompleted
	default:
		// continue
	}
	if order.Maker.BuyToken.Denom != msg.SellToken.Denom {
		return types.ErrOrderDenominationMismatched
	}

	if msg.SellToken.IsGTE(order.Maker.BuyToken) {
		order.Takers = types.NewTakerFromMsg(msg)
		order.Status = types.Status_COMPLETE
		order.CompleteTimestamp = msg.CreateTimestamp
	} else {
		return types.ErrOrderInsufficientAmount
	}

	switch step {
	case StepSend:
		// executed on the destination chain of make swap,
		// lock the taker's sell coin into the module
		sender, err1 := sdk.AccAddressFromBech32(msg.TakerAddress)
		if err1 != nil {
			return err1
		}

		// lock sell token into module
		if err := k.bankKeeper.SendCoinsFromAccountToModule(
			ctx, sender, types.ModuleName, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	case StepReceive:
		// StepReceive
		// executed on the source chain of make swap,
		// send maker's sell token to taker
		receiver, err1 := sdk.AccAddressFromBech32(msg.TakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// send maker's sell token from module to taker
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(order.Maker.SellToken),
		); err != nil {
			return err
		}
		break
	case StepAcknowledgement:
		// executed on the destination chain of the swap,
		// send taker's sell token to maker
		receiver, err1 := sdk.AccAddressFromBech32(order.Maker.MakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// send taker's sell token from module to maker
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(
			ctx, types.ModuleName, receiver, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	}

	// save updated status
	k.SetOTCOrder(ctx, order)

	return nil
}

func (k Keeper) executeCancel(ctx sdk.Context, msg *types.MsgCancelSwapRequest, step int) error {
	// check order status
	if order, ok := k.GetLimitOrder(ctx, msg.OrderId); ok {
		if order.Maker.MakerAddress != msg.MakerAddress {
			return types.ErrOrderPermissionIsNotAllowed
		}
		switch order.Status {
		case types.Status_CANCEL:
			return types.ErrOrderCanceled
		case types.Status_COMPLETE:
			return types.ErrOrderCompleted
		default:
			// continue
			if step != StepSend {
				order.CancelTimestamp = msg.CreateTimestamp
				order.Status = types.Status_CANCEL
			}
		}

	} else if orderOtc, ok2 := k.GetOTCOrder(ctx, msg.OrderId); ok2 {
		if orderOtc.Maker.MakerAddress != msg.MakerAddress {
			return types.ErrOrderPermissionIsNotAllowed
		}
		switch orderOtc.Status {
		case types.Status_CANCEL:
			return types.ErrOrderCanceled
		case types.Status_COMPLETE:
			return types.ErrOrderCompleted
		default:
			// continue
			if step != StepSend {
				order.CancelTimestamp = msg.CreateTimestamp
				order.Status = types.Status_CANCEL
			}
		}
	} else {
		return types.ErrOrderDoesNotExists
	}

	return nil
}

// the following methods are executed On Destination chain.

func (k Keeper) OnReceivedMake(ctx sdk.Context, packet channeltypes.Packet, msg *types.MsgMakeSwapRequest) error {

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	if len(msg.DesiredTaker) > 0 {
		order := types.NewOTCOrder(msg, packet.DestinationChannel)
		k.SetOTCOrder(ctx, order)
	} else {
		order := types.NewLimitOrder(msg, packet.DestinationChannel)
		k.SetLimitOrder(ctx, order)
	}

	ctx.EventManager().EmitTypedEvents(msg)
	return nil
}

func (k Keeper) OnReceivedTake(ctx sdk.Context, packet channeltypes.Packet, msg *types.MsgTakeSwapRequest) error {

	if err := msg.ValidateBasic(); err != nil {
		return err
	}

	// check order status
	if order, ok := k.GetLimitOrder(ctx, msg.OrderId); ok {
		k.executeLimitOrderMatch(ctx, order, msg, StepReceive)
	} else if orderOtc, ok2 := k.GetOTCOrder(ctx, msg.OrderId); ok2 {
		k.executeOTCOrderMatch(ctx, orderOtc, msg, StepReceive)
	} else {
		return types.ErrOrderDoesNotExists
	}

	ctx.EventManager().EmitTypedEvents(msg)
	return nil
}

func (k Keeper) OnReceivedCancel(ctx sdk.Context, packet channeltypes.Packet, msg *types.MsgCancelSwapRequest) error {
	if err := msg.ValidateBasic(); err != nil {
		return err
	}
	// check order status

	if err := k.executeCancel(ctx, msg, StepReceive); err != nil {
		return err
	}

	ctx.EventManager().EmitTypedEvents(msg)
	return nil
}
