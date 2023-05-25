package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InterchainSwapEstimate(goCtx context.Context, req *types.QuerySwapEstimateRequest) (*types.QuerySwapEstimateResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, req.PoolId)
	if !found {
		return nil, status.Error(codes.InvalidArgument, "doesn't exist pool")
	}

	fee := k.GetSwapFeeRate(sdkCtx)
	amm := types.NewInterchainMarketMaker(&pool, fee)

	var out *sdk.Coin
	var err error
	if req.SwapType == types.SwapMsgType_LEFT {
		out, err = amm.LeftSwap(*req.TokenIn, "out")
		if err != nil {
			return nil, err
		}
	} else {
		out, err := amm.RightSwap(*req.TokenIn, sdk.NewCoin("out", sdk.NewInt(0)))
		if err != nil {
			return nil, err
		}
		return &types.QuerySwapEstimateResponse{
			TokenOut: out,
		}, nil
	}
	return &types.QuerySwapEstimateResponse{
		TokenOut: out,
	}, nil

}
