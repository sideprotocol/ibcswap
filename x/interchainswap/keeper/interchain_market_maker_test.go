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

func createNInterchainMarketMaker(keeper *keeper.Keeper, ctx sdk.Context, n int) []types.InterchainMarketMaker {
	items := make([]types.InterchainMarketMaker, n)
	for i := range items {
		items[i].PoolId = strconv.Itoa(i)

		keeper.SetInterchainMarketMaker(ctx, items[i])
	}
	return items
}

func TestInterchainMarketMakerGet(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainMarketMaker(keeper, ctx, 10)
	for _, item := range items {
		rst, found := keeper.GetInterchainMarketMaker(ctx,
			item.PoolId,
		)
		require.True(t, found)
		require.Equal(t,
			nullify.Fill(&item),
			nullify.Fill(&rst),
		)
	}
}
func TestInterchainMarketMakerRemove(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainMarketMaker(keeper, ctx, 10)
	for _, item := range items {
		keeper.RemoveInterchainMarketMaker(ctx,
			item.PoolId,
		)
		_, found := keeper.GetInterchainMarketMaker(ctx,
			item.PoolId,
		)
		require.False(t, found)
	}
}

func TestInterchainMarketMakerGetAll(t *testing.T) {
	keeper, ctx := keepertest.InterchainswapKeeper(t)
	items := createNInterchainMarketMaker(keeper, ctx, 10)
	require.ElementsMatch(t,
		nullify.Fill(items),
		nullify.Fill(keeper.GetAllInterchainMarketMaker(ctx)),
	)
}
