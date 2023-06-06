package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) SingleAssetWithdraw(ctx context.Context, msg *types.MsgSingleAssetWithdrawRequest) (*types.MsgSingleAssetWithdrawResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolCoin.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "pool not found: %s", types.ErrNotFoundPool)
	}

	amm := *types.NewInterchainMarketMaker(&pool)
	denomOut, _ := pool.FindDenomBySide(types.PoolAssetSide_SOURCE)
	out, err := amm.SingleWithdraw(*msg.PoolCoin, *denomOut)
	if err != nil {
		return nil, err
	}
	// Construct the IBC data packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.SINGLE_WITHDRAW,
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

	err = k.SendIBCSwapPacket(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}
	return &types.MsgSingleAssetWithdrawResponse{}, nil
}
