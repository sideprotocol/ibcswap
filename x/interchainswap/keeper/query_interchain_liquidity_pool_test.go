package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/sideprotocol/ibcswap/v4/testutil/keeper"
	"github.com/sideprotocol/ibcswap/v4/testutil/nullify"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestInterchainLiquidityPoolQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInterchainLiquidityPool(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetInterchainLiquidityPoolRequest
		response *types.QueryGetInterchainLiquidityPoolResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetInterchainLiquidityPoolRequest{
				PoolId: msgs[0].PoolId,
			},
			response: &types.QueryGetInterchainLiquidityPoolResponse{InterchainLiquidityPool: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetInterchainLiquidityPoolRequest{
				PoolId: msgs[1].PoolId,
			},
			response: &types.QueryGetInterchainLiquidityPoolResponse{InterchainLiquidityPool: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetInterchainLiquidityPoolRequest{
				PoolId: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.InterchainLiquidityPool(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestInterchainLiquidityPoolQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNInterchainLiquidityPool(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllInterchainLiquidityPoolRequest {
		return &types.QueryAllInterchainLiquidityPoolRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InterchainLiquidityPoolAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InterchainLiquidityPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InterchainLiquidityPool),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.InterchainLiquidityPoolAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.InterchainLiquidityPool), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.InterchainLiquidityPool),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.InterchainLiquidityPoolAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.InterchainLiquidityPool),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.InterchainLiquidityPoolAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
