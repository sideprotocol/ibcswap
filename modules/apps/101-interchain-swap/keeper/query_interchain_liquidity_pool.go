package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) InterchainLiquidityPoolAll(goCtx context.Context, req *types.QueryAllInterchainLiquidityPoolRequest) (*types.QueryAllInterchainLiquidityPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var interchainLiquidityPools []types.InterchainLiquidityPool
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	interchainLiquidityPoolStore := prefix.NewStore(store, types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

	pageRes, err := query.Paginate(interchainLiquidityPoolStore, req.Pagination, func(key []byte, value []byte) error {
		var interchainLiquidityPool types.InterchainLiquidityPool
		if err := k.cdc.Unmarshal(value, &interchainLiquidityPool); err != nil {
			return err
		}

		interchainLiquidityPools = append(interchainLiquidityPools, interchainLiquidityPool)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInterchainLiquidityPoolResponse{InterchainLiquidityPool: interchainLiquidityPools, Pagination: pageRes}, nil
}

func (k Keeper) InterchainLiquidityPool(goCtx context.Context, req *types.QueryGetInterchainLiquidityPoolRequest) (*types.QueryGetInterchainLiquidityPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetInterchainLiquidityPool(
		ctx,
		req.PoolId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInterchainLiquidityPoolResponse{InterchainLiquidityPool: val}, nil
}
