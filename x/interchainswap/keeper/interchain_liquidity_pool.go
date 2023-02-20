package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

// SetInterchainLiquidityPool set a specific interchainLiquidityPool in the store from its index
func (k Keeper) SetInterchainLiquidityPool(ctx sdk.Context, interchainLiquidityPool types.InterchainLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	b := k.cdc.MustMarshal(&interchainLiquidityPool)
	store.Set(types.InterchainLiquidityPoolKey(
		interchainLiquidityPool.PoolId,
	), b)
}

// GetInterchainLiquidityPool returns a interchainLiquidityPool from its index
func (k Keeper) GetInterchainLiquidityPool(
	ctx sdk.Context,
	poolId string,

) (val types.InterchainLiquidityPool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

	b := store.Get(types.InterchainLiquidityPoolKey(
		poolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInterchainLiquidityPool removes a interchainLiquidityPool from the store
func (k Keeper) RemoveInterchainLiquidityPool(
	ctx sdk.Context,
	poolId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	store.Delete(types.InterchainLiquidityPoolKey(
		poolId,
	))
}

// GetAllInterchainLiquidityPool returns all interchainLiquidityPool
func (k Keeper) GetAllInterchainLiquidityPool(ctx sdk.Context) (list []types.InterchainLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InterchainLiquidityPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
