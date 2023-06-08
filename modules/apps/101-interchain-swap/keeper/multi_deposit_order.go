package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// SetInterchainLiquidityPool set a specific interchainLiquidityPool in the store from its index
func (k Keeper) SetMultiDepositOrder(ctx sdk.Context, order types.MultiAssetDepositOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderKeyPrefix))
	b := k.cdc.MustMarshal(&order)
	store.Set(types.MultiDepositOrderPrefixKey(
		order.Id,
	), b)
}

// GetInterchainLiquidityPool returns a interchainLiquidityPool from its index
func (k Keeper) GetMultiDepositOrder(
	ctx sdk.Context,
	orderId string,

) (val types.MultiAssetDepositOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderKeyPrefix))

	b := store.Get(types.MultiDepositOrderPrefixKey(
		orderId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveInterchainLiquidityPool removes a interchainLiquidityPool from the store
func (k Keeper) RemoveMultiDepositOrder(
	ctx sdk.Context,
	orderId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderKeyPrefix))
	store.Delete(types.MultiDepositOrderPrefixKey(
		orderId,
	))
}

// GetAllInterchainLiquidityPool returns all interchainLiquidityPool
func (k Keeper) GetAllMultiDepositOrder(ctx sdk.Context) (list []types.MultiAssetDepositOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MultiAssetDepositOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
