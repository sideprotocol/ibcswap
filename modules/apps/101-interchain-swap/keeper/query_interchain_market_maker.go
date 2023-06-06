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

func (k Keeper) InterchainMarketMakerAll(goCtx context.Context, req *types.QueryAllInterchainMarketMakerRequest) (*types.QueryAllInterchainMarketMakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var interchainMarketMakers []types.InterchainMarketMaker
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	interchainMarketMakerStore := prefix.NewStore(store, types.KeyPrefix(types.InterchainMarketMakerKeyPrefix))

	pageRes, err := query.Paginate(interchainMarketMakerStore, req.Pagination, func(key []byte, value []byte) error {
		var interchainMarketMaker types.InterchainMarketMaker
		if err := k.cdc.Unmarshal(value, &interchainMarketMaker); err != nil {
			return err
		}

		interchainMarketMakers = append(interchainMarketMakers, interchainMarketMaker)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllInterchainMarketMakerResponse{InterchainMarketMaker: interchainMarketMakers, Pagination: pageRes}, nil
}

func (k Keeper) InterchainMarketMaker(goCtx context.Context, req *types.QueryGetInterchainMarketMakerRequest) (*types.QueryGetInterchainMarketMakerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	val, found := k.GetInterchainMarketMaker(
		ctx,
		req.PoolId,
	)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetInterchainMarketMakerResponse{InterchainMarketMaker: val}, nil
}
