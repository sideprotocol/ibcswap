package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// MakeSwap is called when the maker wants to make atomic swap. The method create new order and lock tokens.
// This is the step 1 (Create order & Lock Token) of the atomic swap: https://github.com/cosmos/ibc/tree/main/spec/app/ics-100-atomic-swap
func (k Keeper) MakeSwap(goCtx context.Context, msg *types.MakeSwapMsg) (*types.MsgMakeSwapResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err1 != nil {
		return nil, err1
	}

	balance := k.bankKeeper.GetBalance(ctx, sender, msg.SellToken.Denom)
	if balance.Amount.LT(msg.SellToken.Amount) {
		return &types.MsgMakeSwapResponse{}, errors.New("insufficient balance")
	}

	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)

	// lock sell token into module
	if err := k.bankKeeper.SendCoins(
		ctx, sender, escrowAddr, sdk.NewCoins(msg.SellToken),
	); err != nil {
		return nil, err
	}

	order, err := createOrder(ctx, msg, k.channelKeeper)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedMakeSwap, "due to %s", err)
	}

	packet := types.AtomicSwapPacketData{
		Type:    types.MAKE_SWAP,
		Data:    msgByte,
		OrderId: order.Id,
		Path:    order.Path,
		Memo:    "",
	}

	if _, err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	k.SetAtomicOrder(ctx, *order)
	//ctx.EventManager().EmitTypedEvents(msg)

	sdk.NewEvent(
		types.EventTypeMakeSwap,
		sdk.Attribute{
			Key:   types.AttributeOrderId,
			Value: order.Id,
		},
	)
	return &types.MsgMakeSwapResponse{OrderId: order.Id}, nil
}
