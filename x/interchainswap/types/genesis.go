package types

import (
	"fmt"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
)

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		PortId:                      PortID,
		InterchainLiquidityPoolList: []InterchainLiquidityPool{},
		InterchainMarketMakerList:   []InterchainMarketMaker{},
		// this line is used by starport scaffolding # genesis/types/default
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := host.PortIdentifierValidator(gs.PortId); err != nil {
		return err
	}
	// Check for duplicated index in interchainLiquidityPool
	interchainLiquidityPoolIndexMap := make(map[string]struct{})

	for _, elem := range gs.InterchainLiquidityPoolList {
		index := string(InterchainLiquidityPoolKey(elem.PoolId))
		if _, ok := interchainLiquidityPoolIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for interchainLiquidityPool")
		}
		interchainLiquidityPoolIndexMap[index] = struct{}{}
	}
	// Check for duplicated index in interchainMarketMaker
	interchainMarketMakerIndexMap := make(map[string]struct{})

	for _, elem := range gs.InterchainMarketMakerList {
		index := string(InterchainMarketMakerKey(elem.PoolId))
		if _, ok := interchainMarketMakerIndexMap[index]; ok {
			return fmt.Errorf("duplicated index for interchainMarketMaker")
		}
		interchainMarketMakerIndexMap[index] = struct{}{}
	}
	// this line is used by starport scaffolding # genesis/types/validate

	return gs.Params.Validate()
}
