package keeper

import (
	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

/// Balancer Liquidity Pool

// GetBalancerPool returns the BalancerLiquidityPool for the swap module.
func (k Keeper) GetBalancerPool(ctx sdk.Context, poolId string) (types.BalancerLiquidityPool, bool) {
	key, err := hex.DecodeString(poolId)
	if err != nil {
		return types.BalancerLiquidityPool{}, false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancerPoolKey)
	bz := store.Get(key)
	if bz == nil {
		return types.BalancerLiquidityPool{}, false
	}

	order := k.MustUnmarshalBalancerLiquidityPool(bz)
	return order, true
}

// HasBalancerPool checks if a the key with the given id exists on the store.
func (k Keeper) HasBalancerPool(ctx sdk.Context, poolId string) bool {
	key, err := hex.DecodeString(poolId)
	if err != nil {
		return false
	}
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancerPoolKey)
	return store.Has(key)
}

// SetBalancerPool sets a new BalancerLiquidityPool to the store.
func (k Keeper) SetBalancerPool(ctx sdk.Context, pool types.BalancerLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.BalancerPoolKey)
	bz := k.MustMarshalBalancerLiquidityPool(pool)
	store.Set([]byte(pool.Id), bz)
}

// IterateBalancerPools iterates over the BalancerLiquidityPools in the store
// and performs a callback function.
func (k Keeper) IterateBalancerPools(ctx sdk.Context, cb func(order types.BalancerLiquidityPool) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.BalancerPoolKey)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {

		pool := k.MustUnmarshalBalancerLiquidityPool(iterator.Value())
		if cb(pool) {
			break
		}
	}
}
