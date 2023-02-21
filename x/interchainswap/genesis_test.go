package interchainswap_test

import (
	"testing"

	keepertest "github.com/sideprotocol/ibcswap/v4/testutil/keeper"
	"github.com/sideprotocol/ibcswap/v4/testutil/nullify"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),
		PortId: types.PortID,
		InterchainLiquidityPoolList: []types.InterchainLiquidityPool{
			{
				PoolId: "0",
			},
			{
				PoolId: "1",
			},
		},
		InterchainMarketMakerList: []types.InterchainMarketMaker{
			{
				PoolId: "0",
			},
			{
				PoolId: "1",
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.InterchainswapKeeper(t)
	interchainswap.InitGenesis(ctx, *k, genesisState)
	got := interchainswap.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.PortId, got.PortId)

	require.ElementsMatch(t, genesisState.InterchainLiquidityPoolList, got.InterchainLiquidityPoolList)
	require.ElementsMatch(t, genesisState.InterchainMarketMakerList, got.InterchainMarketMakerList)
	// this line is used by starport scaffolding # genesis/test/assert
}
