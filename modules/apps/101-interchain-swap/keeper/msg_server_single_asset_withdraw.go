package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) SingleAssetWithdraw(ctx context.Context, msg *types.MsgSingleAssetWithdrawRequest) (*types.MsgSingleAssetWithdrawResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// Check denom
	if !k.bankKeeper.HasSupply(sdkCtx, msg.DenomOut) {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "invalid denom in local withdraw message: %s", msg.DenomOut)
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolCoin.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "pool not found: %s", types.ErrNotFoundPool)
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_READY {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "pool not ready for swap: %s", types.ErrNotReadyForSwap)
	}

	fee := k.GetSwapFeeRate(sdkCtx)
	amm := *types.NewInterchainMarketMaker(&pool, fee)

	out, err := amm.SingleWithdraw(*msg.PoolCoin, msg.DenomOut)

	if err != nil {
		return nil, err
	}

	// Construct the IBC data packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.MULTI_WITHDRAW,
		Data: rawMsgData,
		StateChange: &types.StateChange{
			Out:        []*sdk.Coin{out},
			PoolTokens: []*sdk.Coin{msg.PoolCoin},
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)
	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	err = k.SendIBCSwapPacket(sdkCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}
	return &types.MsgSingleAssetWithdrawResponse{}, nil
}
