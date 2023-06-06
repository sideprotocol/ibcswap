package types_test

import (
	"testing"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/stretchr/testify/require"
)

func TestGenesisState_Validate(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		genState *types.GenesisState
		valid    bool
	}{
		{
			desc:     "default is valid",
			genState: types.DefaultGenesis(),
			valid:    true,
		},
		{
			desc: "valid genesis state",
			genState: &types.GenesisState{
				PortId: types.PortID,
				InterchainLiquidityPoolList: []types.InterchainLiquidityPool{
					{
						Id: "0",
					},
					{
						Id: "1",
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
				// this line is used by starport scaffolding # types/genesis/validField
			},
			valid: true,
		},
		{
			desc: "duplicated interchainLiquidityPool",
			genState: &types.GenesisState{
				InterchainLiquidityPoolList: []types.InterchainLiquidityPool{
					{
						Id: "0",
					},
					{
						Id: "0",
					},
				},
			},
			valid: false,
		},
		{
			desc: "duplicated interchainMarketMaker",
			genState: &types.GenesisState{
				InterchainMarketMakerList: []types.InterchainMarketMaker{
					{
						PoolId: "0",
					},
					{
						PoolId: "0",
					},
				},
			},
			valid: false,
		},
		// this line is used by starport scaffolding # types/genesis/testcase
	} {
		t.Run(tc.desc, func(t *testing.T) {
			err := tc.genState.Validate()
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
