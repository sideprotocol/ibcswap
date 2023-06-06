package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// SetInterchainMarketMaker set a specific interchainMarketMaker in the store from its index
func (k Keeper) SetInterchainMarketMaker(ctx sdk.Context, interchainMarketMaker types.InterchainMarketMaker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainMarketMakerKeyPrefix))
	b := k.cdc.MustMarshal(&interchainMarketMaker)
	store.Set(types.InterchainMarketMakerKey(
		interchainMarketMaker.PoolId,
	), b)
}

// GetInterchainMarketMaker returns a interchainMarketMaker from its index
func (k Keeper) GetInterchainMarketMaker(
	ctx sdk.Context,
	poolId string,

) (val types.InterchainMarketMaker, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainMarketMakerKeyPrefix))

	b := store.Get(types.InterchainMarketMakerKey(
		poolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInterchainMarketMaker removes a interchainMarketMaker from the store
func (k Keeper) RemoveInterchainMarketMaker(
	ctx sdk.Context,
	poolId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainMarketMakerKeyPrefix))
	store.Delete(types.InterchainMarketMakerKey(
		poolId,
	))
}

// GetAllInterchainMarketMaker returns all interchainMarketMaker
func (k Keeper) GetAllInterchainMarketMaker(ctx sdk.Context) (list []types.InterchainMarketMaker) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainMarketMakerKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.InterchainMarketMaker
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
