package keeper

import (
	"context"
	"encoding/binary"

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

	// Use the store with the mapping of poolId to its count
	store := ctx.KVStore(k.storeKey)
	poolIdToCountStore := prefix.NewStore(store, types.PoolIdToCountKeyPrefix)

	// A counter for pagination
	var counter int

	pageRes, err := query.Paginate(poolIdToCountStore, req.Pagination, func(key []byte, value []byte) error {
		counter++

		// Decode the count for the given poolId
		count := binary.BigEndian.Uint64(value)

		// Get the actual pool using the count
		poolStore := prefix.NewStore(store, types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
		poolBytes := poolStore.Get(GetInterchainLiquidityPoolKey(count))
		if poolBytes == nil {
			return nil
		}

		var interchainLiquidityPool types.InterchainLiquidityPool
		if err := k.cdc.Unmarshal(poolBytes, &interchainLiquidityPool); err != nil {
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

func (k Keeper) InterchainLiquidityMyPoolAll(goCtx context.Context, req *types.QueryAllInterchainLiquidityMyPoolRequest) (*types.QueryAllInterchainLiquidityPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var interchainLiquidityPools []types.InterchainLiquidityPool
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Use the store with the mapping of poolId to its count
	store := ctx.KVStore(k.storeKey)
	poolIdToCountStore := prefix.NewStore(store, types.PoolIdToCountKeyPrefix)

	pageRes, err := query.Paginate(poolIdToCountStore, req.Pagination, func(key []byte, value []byte) error {
		// Decode the count for the given poolId
		count := binary.BigEndian.Uint64(value)

		// Get the actual pool using the count
		poolStore := prefix.NewStore(store, types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
		poolBytes := poolStore.Get(GetInterchainLiquidityPoolKey(count))
		if poolBytes == nil {
			return nil
		}

		var interchainLiquidityPool types.InterchainLiquidityPool
		if err := k.cdc.Unmarshal(poolBytes, &interchainLiquidityPool); err != nil {
			return err
		}

		// Check if the creator is either SourceCreator or DestinationCreator
		if interchainLiquidityPool.SourceCreator == req.Creator || interchainLiquidityPool.DestinationCreator == req.Creator {
			interchainLiquidityPools = append(interchainLiquidityPools, interchainLiquidityPool)
		}
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
