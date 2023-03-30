package keeper

import (
	"context"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gogo/protobuf/proto"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// MakeSwap is called when the maker wants to make atomic swap. The method create new order and lock tokens.
// This is the step 1 (Create order & Lock Token) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap.
func (k Keeper) MakeSwap(goCtx context.Context, msg *types.MakeSwapMsg) (*types.MsgMakeSwapResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	msgByte, err0 := proto.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err1 != nil {
		return nil, err1
	}

	balance := k.bankKeeper.GetBalance(ctx, sender, msg.SellToken.Denom)
	if balance.Amount.BigInt().Cmp(msg.SellToken.Amount.BigInt()) < 0 {
		return &types.MsgMakeSwapResponse{}, errors.New("insufficient balance")
	}

	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)

	// lock sell token into module
	if err := k.bankKeeper.SendCoins(
		ctx, sender, escrowAddr, sdk.NewCoins(msg.SellToken),
	); err != nil {
		return nil, err
	}

	packet := types.AtomicSwapPacketData{
		Type: types.MAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}

	order := createOrder(ctx, msg, k.channelKeeper)

	if err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgMakeSwapResponse{OrderId: order.Id}, nil
}