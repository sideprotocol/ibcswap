package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// You may need to adjust the function signature, return types, and parameter types based on your module's implementation
func (k Keeper) InterchainMultiDepositOrder(ctx context.Context, req *types.QueryGetInterchainMultiDepositOrderRequest) (*types.QueryGetInterchainMultiDepositOrderResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order, found := k.GetMultiDepositOrder(sdkCtx, req.PoolId, req.OrderId)

	if !found {
		return nil, types.ErrNotFoundMultiDepositOrder
	}
	return &types.QueryGetInterchainMultiDepositOrderResponse{
		Order: &order,
	}, nil
}

// You may need to adjust the function signature, return types, and parameter types based on your module's implementation
func (k Keeper) InterchainMultiDepositOrdersAll(ctx context.Context, req *types.QueryAllInterchainMultiDepositOrdersRequest) (*types.QueryAllInterchainMultiDepositOrdersResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	orders := k.GetAllMultiDepositOrder(sdkCtx, req.PoolId)
	ordersPtr := make([]*types.MultiAssetDepositOrder, len(orders))
	for i := range orders {
		ordersPtr[i] = &orders[i]
	}
	return &types.QueryAllInterchainMultiDepositOrdersResponse{
		Orders: ordersPtr,
	}, nil
}

// You may need to adjust the function signature, return types, and parameter types based on your module's implementation
func (k Keeper) InterchainLatestMultiDepositOrder(ctx context.Context, req *types.QueryLatestInterchainMultiDepositOrderRequest) (*types.QueryGetInterchainMultiDepositOrderResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	order, found := k.GetLatestMultiDepositOrder(sdkCtx, req.PoolId)
	if !found {
		return nil, types.ErrNotFoundMultiDepositOrder
	}
	return &types.QueryGetInterchainMultiDepositOrderResponse{
		Order: &order,
	}, nil
}
