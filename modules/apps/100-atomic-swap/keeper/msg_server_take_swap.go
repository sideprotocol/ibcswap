package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// TakeSwap is the step 5 (Lock Order & Lock Token) of the atomic swap: https://github.com/liangping/ibc/blob/atomic-swap/spec/app/ics-100-atomic-swap/ibcswap.png
// This method lock the order (set a value to the field "Taker") and lock Token
func (k Keeper) TakeSwap(goCtx context.Context, msg *types.TakeSwapMsg) (*types.MsgTakeSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// pre-check
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return nil, types.ErrOrderDoesNotExists
	}

	sourceChannel := extractSourceChannelForTakerMsg(order.Path)
	sourcePort := extractSourcePortForTakerMsg(order.Path)

	escrowAddr := types.GetEscrowAddress(sourcePort, sourceChannel)

	if order.Status != types.Status_SYNC && order.Status != types.Status_INITIAL {
		return nil, errors.New("order is not in valid state")
	}

	// Make sure the maker's buy token matches the taker's sell token
	if order.Maker.BuyToken.Denom != msg.SellToken.Denom && !order.Maker.BuyToken.Amount.Equal(msg.SellToken.Amount) {
		return &types.MsgTakeSwapResponse{}, errors.New("invalid sell token")
	}

	// Checks if the order has already been taken
	if order.Takers != nil {
		return &types.MsgTakeSwapResponse{}, errors.New("order has already been taken")
	}

	// If `desiredTaker` is set, only the desiredTaker can accept the order.
	if order.Maker.DesiredTaker != "" && order.Maker.DesiredTaker != msg.TakerAddress {
		return &types.MsgTakeSwapResponse{}, errors.New("invalid taker address")
	}

	takerAddr, err := sdk.AccAddressFromBech32(msg.TakerAddress)
	if err != nil {
		return &types.MsgTakeSwapResponse{}, err
	}

	_, err = sdk.AccAddressFromBech32(msg.TakerReceivingAddress)
	if err == nil {
		return &types.MsgTakeSwapResponse{}, types.ErrFailedMakeSwap
	}

	balance := k.bankKeeper.GetBalance(ctx, takerAddr, msg.SellToken.Denom)
	if balance.Amount.LT(msg.SellToken.Amount) {
		return &types.MsgTakeSwapResponse{}, errors.New("insufficient balance")
	}

	// Locks the sellToken to the escrow account
	if err = k.bankKeeper.SendCoins(ctx, takerAddr, escrowAddr, sdk.NewCoins(msg.SellToken)); err != nil {
		return &types.MsgTakeSwapResponse{}, err
	}

	packet := types.AtomicSwapPacketData{
		Type: types.TAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}

	if _, err := k.SendSwapPacket(ctx, sourcePort, sourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	// Update order state
	// Mark that the order has been occupied
	order.Takers = msg
	k.SetAtomicOrder(ctx, order)
	ctx.EventManager().EmitTypedEvents(msg)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTakeSwap,
			sdk.Attribute{
				Key:   types.AttributeOrderId,
				Value: order.Id,
			},
		),
	)

	return &types.MsgTakeSwapResponse{}, nil
}
