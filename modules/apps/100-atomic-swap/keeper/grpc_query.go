package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) GetAllOrders(goCtx context.Context, request *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	orderStore := prefix.NewStore(ctx.KVStore(q.storeKey), types.OTCOrderBookKey)
	var orders []*types.Order
	pageRes, err := query.Paginate(orderStore, request.Pagination, func(key []byte, value []byte) error {
		var order types.Order
		if err := q.cdc.Unmarshal(value, &order); err != nil {
			return err
		}
		orders = append(orders, &order)
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}

	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q Keeper) GetAllOrdersByType(ctx context.Context, request *types.QueryOrdersByRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)
	store := clientCtx.KVStore(q.storeKey)
	orderStore := prefix.NewStore(store, types.OTCOrderBookKey)
	var orders []*types.Order

	pageRes, err := query.Paginate(orderStore, request.Pagination, func(key, value []byte) error {
		var order types.Order
		if err := q.cdc.Unmarshal(value, &order); err != nil {
			return err
		}
		acc := q.authKeeper.GetAccount(clientCtx, sdk.MustAccAddressFromBech32(order.Maker.MakerAddress))
		if (acc != nil && request.OrderType == types.OrderType_SellToBuy) ||
			(acc == nil && request.OrderType == types.OrderType_BuyToSell) {
			orders = append(orders, &order)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}
	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q Keeper) GetSubmittedOrders(ctx context.Context, request *types.QuerySubmittedOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)
	store := clientCtx.KVStore(q.storeKey)
	orderStore := prefix.NewStore(store, types.OTCOrderBookKey)
	var orders []*types.Order

	pageRes, err := query.Paginate(orderStore, request.Pagination, func(key, value []byte) error {

		var order types.Order
		if err := q.cdc.Unmarshal(value, &order); err != nil {
			return err
		}
		if order.Maker != nil && order.Maker.MakerAddress == request.MakerAddress {
			orders = append(orders, &order)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}
	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q Keeper) GetTookOrders(ctx context.Context, request *types.QueryTookOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)
	store := clientCtx.KVStore(q.storeKey)
	orderStore := prefix.NewStore(store, types.OTCOrderBookKey)

	var orders []*types.Order

	pageRes, err := query.Paginate(orderStore, request.Pagination, func(key, value []byte) error {
		var order types.Order
		if err := q.cdc.Unmarshal(value, &order); err != nil {
			return err
		}
		if order.Takers != nil && order.Takers.TakerAddress == request.TakerAddress {
			orders = append(orders, &order)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}
	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

func (q Keeper) GetPrivateOrders(ctx context.Context, request *types.QueryPrivateOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)
	store := clientCtx.KVStore(q.storeKey)
	orderStore := prefix.NewStore(store, types.OTCOrderBookKey)
	var orders []*types.Order

	pageRes, err := query.Paginate(orderStore, request.Pagination, func(key, value []byte) error {
		var order types.Order
		if err := q.cdc.Unmarshal(value, &order); err != nil {
			return err
		}
		if order.Maker.DesiredTaker == request.DesireAddress {
			orders = append(orders, &order)
		}
		return nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "paginate: %v", err)
	}
	return &types.QueryOrdersResponse{Orders: orders, Pagination: pageRes}, nil
}

// Params implements the Query/Params gRPC method
func (q Keeper) Params(c context.Context, _ *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	params := q.GetParams(ctx)

	return &types.QueryParamsResponse{
		Params: &params,
	}, nil
}

// EscrowAddress implements the EscrowAddress gRPC method
func (q Keeper) EscrowAddress(c context.Context, req *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	addr := types.GetEscrowAddress(req.PortId, req.ChannelId)

	return &types.QueryEscrowAddressResponse{
		EscrowAddress: addr.String(),
	}, nil
}
