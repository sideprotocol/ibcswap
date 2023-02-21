package types

import (
	"strings"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/types"
)

func NewInterchainLiquidityPool(
	ctx types.Context,
	store BankKeeper,

	denoms []string,
	decimals []uint32,
	weight string,
	portId string,
	channelId string,
) *InterchainLiquidityPool {

	//generate poolId
	poolId := GetPoolId(denoms)

	weightSize := len(strings.Split(weight, ":"))
	denomSize := len(denoms)
	decimalSize := len(decimals)
	assets := []*PoolAsset{}

	if denomSize == weightSize && decimalSize == weightSize {
		for _, denom := range denoms {
			side := PoolSide_NATIVE
			if !store.HasSupply(ctx, denom) {
				side = PoolSide_REMOTE
			}
			asset := PoolAsset{
				Side: side,
				Balance: &types.Coin{
					Amount: math.NewInt(0),
					Denom:  denom,
				},
			}
			assets = append(assets, &asset)
		}
	}

	return &InterchainLiquidityPool{
		PoolId: poolId,
		Assets: assets,
		Supply: &types.Coin{
			Amount: math.NewInt(0),
			Denom:  poolId,
		},
		Status:                PoolStatus_POOL_STATUS_INITIAL,
		EncounterPartyPort:    portId,
		EncounterPartyChannel: channelId,
	}
}

func (ilp *InterchainLiquidityPool) FindAssetByDenom(denom string) (*PoolAsset, error) {
	for _, asset := range ilp.Assets {
		if asset.Balance.Denom == denom {
			return asset, nil
		}
	}
	return nil, ErrNotFoundDenomInPool
}

func (ilp *InterchainLiquidityPool) UpdateAssetPoolSide(denom string, side PoolSide) (*PoolAsset, error) {
	for _, asset := range ilp.Assets {
		if asset.Balance.Denom == denom {
			asset.Side = side
		}
	}
	return nil, ErrNotFoundDenomInPool
}
