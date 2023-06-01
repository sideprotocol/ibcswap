package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// CancelSwap is the step 10 (Cancel Request) of the atomic swap: https://github.com/cosmos/ibc/tree/main/spec/app/ics-100-atomic-swap.
// It is executed on the Maker chain. Only the maker of the order can cancel the order.
func (k Keeper) CancelSwap(goCtx context.Context, msg *types.CancelSwapMsg) (*types.MsgCancelSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgbyte, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return &types.MsgCancelSwapResponse{}, types.ErrOrderDoesNotExists
	}

	// Make sure the sender is the maker of the order.
	// TODO: New Message logic for cancel swap is incomplete. (using command line)
	// We should verify/make sure that msg.MakerAddress and sender of TX is same. Otherwise this should fail
	if order.Maker.MakerAddress != msg.MakerAddress {
		return &types.MsgCancelSwapResponse{}, fmt.Errorf("sender is not the maker of the order")
	}

	// Make sure the order is in a valid state for cancellation
	if order.Status != types.Status_SYNC && order.Status != types.Status_INITIAL {
		return &types.MsgCancelSwapResponse{}, fmt.Errorf("order is not in a valid state for cancellation")
	}

	packet := types.AtomicSwapPacketData{
		Type: types.CANCEL_SWAP,
		Data: msgbyte,
		Memo: "",
	}

	if err := k.SendSwapPacket(ctx, order.Maker.SourcePort, order.Maker.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgCancelSwapResponse{}, nil
}
