package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

var _ types.QueryServer = Keeper{}

func (q Keeper) Orders(ctx context.Context, request *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)

	var orders []*types.Order
	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		orders = append(orders, &order)
		return false
	})
	return &types.QueryOrdersResponse{Orders: orders}, nil
}

func (q Keeper) GetAllOrders(ctx context.Context, request *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)

	var orders []*types.Order
	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		orders = append(orders, &order)
		return false
	})
	return &types.QueryOrdersResponse{Orders: orders}, nil
}

func (q Keeper) GetAllOrdersByType(ctx context.Context, request *types.QueryOrdersByRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)
	var orders []*types.Order

	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		_, err := q.authKeeper.GetAccount(clientCtx, sdk.MustAccAddressFromBech32(order.Maker.MakerAddress))
		if (err == nil && request.OrderType == types.OrderType_SellToBuy) ||
			(err != nil && request.OrderType == types.OrderType_BuyToSell) {
			orders = append(orders, &order)
		}
		return false
	})

	return &types.QueryOrdersResponse{Orders: orders}, nil
}

func (q Keeper) GetSubmittedOrders(ctx context.Context, request *types.QuerySubmittedOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)

	var orders []*types.Order
	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		if order.Maker.MakerAddress == request.MakerAddress {
			orders = append(orders, &order)
		}
		return false
	})
	return &types.QueryOrdersResponse{Orders: orders}, nil
}

func (q Keeper) GetTookOrders(ctx context.Context, request *types.QueryTookOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)

	var orders []*types.Order
	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		if order.Takers.TakerAddress == request.TakerAddress {
			orders = append(orders, &order)
		}
		return false
	})
	return &types.QueryOrdersResponse{Orders: orders}, nil
}

func (q Keeper) GetPrivateOrders(ctx context.Context, request *types.QueryPrivateOrdersRequest) (*types.QueryOrdersResponse, error) {
	clientCtx := sdk.UnwrapSDKContext(ctx)

	var orders []*types.Order
	q.IterateAtomicOrders(clientCtx, func(order types.Order) bool {
		if order.Maker.DesiredTaker == request.DesireAddress {
			orders = append(orders, &order)
		}
		return false
	})
	return &types.QueryOrdersResponse{Orders: orders}, nil
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
