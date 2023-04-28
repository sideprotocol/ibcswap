package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Validate msg
	err := msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "failed to swap because of %s", err)
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, types.GetPoolId([]string{
		msg.TokenIn.Denom, msg.TokenOut.Denom,
	}))

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", types.ErrNotFoundPool)
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_READY {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", types.ErrNotReadyForSwap)
	}

	// Lock swap-in token to the swap module
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*msg.TokenIn))
	if err != nil {
		return nil, err
	}

	// Construct the IBC data packet
	swapData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

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

	if tokenOut.Amount.LTE(math.NewInt(0)) {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because", tokenOut.Amount)
	}

	// Slippage checking
	factor := types.MaximumSlippage - msg.Slippage
	expected := msg.TokenOut.Amount.Mul(sdk.NewIntFromUint64(uint64(factor))).Quo(sdk.NewIntFromUint64(types.MaximumSlippage))
	if tokenOut.Amount.LT(expected) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnSwapReceived, "doesn't meet slippage for swap!, expect: %v, output: %v, factor:%d", expected, tokenOut, factor)
	}

	packet := types.IBCSwapPacketData{
		Type:        msgType,
		Data:        swapData,
		StateChange: &types.StateChange{Out: []*sdk.Coin{tokenOut}},
	}

	timeOutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)

	err = k.SendIBCSwapPacket(
		ctx,
		pool.EncounterPartyPort,
		pool.EncounterPartyChannel,
		timeOutHeight,
		timeoutStamp,
		packet,
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapResponse{
		SwapType: msg.SwapType,
		Tokens:   []*sdk.Coin{msg.TokenIn, msg.TokenOut},
	}, nil
}
