package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate msg
	err := msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "failed to swap due to %s", err)
	}

	//poolID := types.GetPoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "pool not found: %s", types.ErrNotFoundPool)
	}

	if pool.Status != types.PoolStatus_ACTIVE {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "pool not ready for swap: %s", types.ErrNotReadyForSwap)
	}

	// Lock swap-in token to the swap module
	err = k.LockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*msg.TokenIn))
	if err != nil {
		return nil, err
	}

	amm := *types.NewInterchainMarketMaker(&pool)

	var tokenOut *sdk.Coin
	var msgType types.SwapMessageType

	switch msg.SwapType {
	case types.SwapMsgType_LEFT:
		msgType = types.LEFT_SWAP
		tokenOut, err = amm.LeftSwap(*msg.TokenIn, msg.TokenOut.Denom)
		if err != nil {
			return nil, err
		}
	case types.SwapMsgType_RIGHT:
		msgType = types.RIGHT_SWAP
		tokenOut, err = amm.RightSwap(*msg.TokenIn, *msg.TokenOut)
		if err != nil {
			return nil, err
		}
	default:
		return nil, types.ErrInvalidSwapType
	}

	if tokenOut.Amount.LTE(sdk.NewInt(0)) {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "token amount is non-positive: %s", tokenOut.Amount)
	}

	// Slippage checking
	factor := types.MaximumSlippage - msg.Slippage
	expected := msg.TokenOut.Amount.Mul(sdk.NewIntFromUint64(uint64(factor))).Quo(sdk.NewIntFromUint64(types.MaximumSlippage))
	if tokenOut.Amount.LT(expected) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnSwapReceived, "slippage check failed! expected: %v, output: %v, factor: %d", expected, tokenOut, factor)
	}

	msg.TokenOut = tokenOut
	// Construct the IBC data packet
	swapData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        msgType,
		Data:        swapData,
		StateChange: &types.StateChange{Out: []*sdk.Coin{tokenOut}},
	}

	timeoutHeight, timeoutTimestamp := types.GetDefaultTimeOut(&ctx)

	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutTimestamp = msg.TimeoutTimeStamp
	}

	_, err = k.SendIBCSwapPacket(
		ctx,
		msg.Port,
		msg.Channel,
		timeoutHeight,
		timeoutTimestamp,
		packet,
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapResponse{
		SwapType: msg.SwapType,
		Tokens:   []*sdk.Coin{msg.TokenIn, tokenOut},
	}, nil
}
