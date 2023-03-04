package keeper

import (
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// UnmarshalBalancerLiquidityPool attempts to decode and return an BalancerLiquidityPool object from
// raw encoded bytes.
func (k Keeper) UnmarshalBalancerLiquidityPool(bz []byte) (types.BalancerLiquidityPool, error) {
	var pool types.BalancerLiquidityPool
	if err := k.cdc.Unmarshal(bz, &pool); err != nil {
		return types.BalancerLiquidityPool{}, err
	}

	return pool, nil
}

// MustUnmarshalBalancerLiquidityPool attempts to decode and return an BalancerLiquidityPool object from
// raw encoded bytes. It panics on error.
func (k Keeper) MustUnmarshalBalancerLiquidityPool(bz []byte) types.BalancerLiquidityPool {
	var pool types.BalancerLiquidityPool
	k.cdc.MustUnmarshal(bz, &pool)
	return pool
}

// MarshalBalancerLiquidityPool attempts to encode an BalancerLiquidityPool object and returns the
// raw encoded bytes.
func (k Keeper) MarshalBalancerLiquidityPool(pool types.BalancerLiquidityPool) ([]byte, error) {
	return k.cdc.Marshal(&pool)
}

// MustMarshalBalancerLiquidityPool attempts to encode an BalancerLiquidityPool object and returns the
// raw encoded bytes. It panics on error.
func (k Keeper) MustMarshalBalancerLiquidityPool(pool types.BalancerLiquidityPool) []byte {
	return k.cdc.MustMarshal(&pool)
}
