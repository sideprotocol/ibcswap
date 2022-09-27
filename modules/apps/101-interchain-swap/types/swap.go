package types

import sdk "github.com/cosmos/cosmos-sdk/types"

func NewPoolAsset(denom string, weight uint32, decimal uint32) *PoolAsset {
	balance := sdk.NewCoin(denom, sdk.ZeroInt())
	return &PoolAsset{
		Side:     PoolSide_POOL_SIDE_PENDING, // will need to be updated later
		Balance:  &balance,
		Weight:   weight,
		Decimals: decimal,
	}
}

func NewBalancerLiquidityPool(denoms []string, decimals []uint32, weight string) BalancerLiquidityPool {

	numOfPairs := len(denoms)
	id := GeneratePoolId(denoms)
	assets := make([]*PoolAsset, numOfPairs)

	weights, _ := ParseWeight(weight, numOfPairs)

	for i := 0; i < numOfPairs; i++ {
		assets[i] = NewPoolAsset(denoms[i], weights[i], decimals[i])
	}

	return BalancerLiquidityPool{
		Id:     id,
		Assets: assets,
	}
}

// Validate check if liquidity pool is valid
func (m *BalancerLiquidityPool) Validate() error {
	return nil
}

func (m *BalancerLiquidityPool) UpdateAssetPoolSide(denom string, side PoolSide) {
	for _, asset := range m.Assets {
		if asset.Balance.Denom == denom {
			asset.Side = side
		}
	}
}
