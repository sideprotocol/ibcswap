package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/sideprotocol/ibcswap/v4/testutil/keeper"
	"github.com/sideprotocol/ibcswap/v4/testutil/nullify"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/keeper"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
	"github.com/stretchr/testify/require"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func createNInterchainLiquidityPool(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.InterchainLiquidityPool {
	items := make([]types.InterchainLiquidityPool, n)
	for i := range items {
		items[i].PoolId = strconv.Itoa(i)

		keeper.SetInterchainLiquidityPool(ctx, items[i])
	}
	return items
}

func TestInterchainLiquidityPoolGet(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainLiquidityPool(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInterchainLiquidityPool(ctx,
			item.PoolId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestInterchainLiquidityPoolRemove(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainLiquidityPool(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInterchainLiquidityPool(ctx,
			item.PoolId,
		)
		_, found := keeper.GetInterchainLiquidityPool(ctx,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestInterchainLiquidityPoolGetAll(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainLiquidityPool(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllInterchainLiquidityPool(ctx)),
	)
}
