package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) TakePool(ctx context.Context, msg *types.MsgTakePoolRequest) (*types.MsgTakePoolResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_ = sdkCtx

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedTakePool, "due to %", types.ErrNotFoundPool)
	}

	if pool.OriginatingChainId == sdkCtx.ChainID() {
		return nil, errorsmod.Wrapf(types.ErrFailedTakePool, "due to %", "same chain")
	}

	creatorAddr := sdk.MustAccAddressFromBech32(msg.Creator)

	asset, err := pool.FindAssetBySide(types.PoolAssetSide_SOURCE)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedTakePool, "due to %", err)
	}

	liquidity := k.bankKeeper.GetBalance(sdkCtx, creatorAddr, asset.Denom)
	if liquidity.Amount.LTE(sdk.NewInt(0)) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "due to %s", types.ErrInEnoughAmount)
	}

	// Move initial funds to liquidity pool
	err = k.LockTokens(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, creatorAddr, sdk.NewCoins(*asset))

	rawMsg, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// Construct IBC data packet
	packet := types.IBCSwapPacketData{
		Type: types.MAKE_POOL,
		Data: rawMsg,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)

	// use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	err = k.SendIBCSwapPacket(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgTakePoolResponse{
		PoolId: msg.PoolId,
	}, nil
}
