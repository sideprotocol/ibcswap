package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// SetMultiDepositOrder set a specific multiDepositOrder in the store from its index
func (k Keeper) SetLatestOrderId(ctx sdk.Context, poolId, sourceMaker, orderId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderCountKeyPrefix))
	store.Set(types.MultiDepositOrderPrefixKey(
		poolId+sourceMaker,
	), []byte(orderId))
}

// GetMultiDepositOrder returns a multiDepositOrder from its index
func (k Keeper) GetLatestMultiDepositOrderId(
	ctx sdk.Context,
	poolId,
	sourceMaker string,

) (val string, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.MultiDepositOrderCountKeyPrefix))
	b := store.Get(types.MultiDepositOrderPrefixKey(
		poolId + sourceMaker,
	))
	if b == nil {
		return val, false
	}
	return string(b), true
}

// SetMultiDepositOrder set a specific multiDepositOrder in the store from its index
func (k Keeper) SetMultiDepositOrder(ctx sdk.Context, multiDepositOrder types.MultiAssetDepositOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(multiDepositOrder.PoolId+types.MultiDepositOrderKeyPrefix))
	b := k.cdc.MustMarshal(&multiDepositOrder)
	store.Set(types.MultiDepositOrderPrefixKey(
		multiDepositOrder.Id,
	), b)
	k.SetLatestOrderId(ctx, multiDepositOrder.PoolId, multiDepositOrder.SourceMaker, multiDepositOrder.Id)
}

// GetMultiDepositOrder returns a multiDepositOrder from its index
func (k Keeper) GetMultiDepositOrder(
	ctx sdk.Context,
	poolId,
	orderId string,

) (val types.MultiAssetDepositOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))

	b := store.Get(types.MultiDepositOrderPrefixKey(
		orderId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveMultiDepositOrder removes a multiDepositOrder from the store
func (k Keeper) RemoveMultiDepositOrder(
	ctx sdk.Context,
	poolId string,
	orderId string,

) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	store.Delete(types.MultiDepositOrderPrefixKey(
		orderId,
	))

}

// GetAllMultiDepositOrder returns all multiDepositOrder
func (k Keeper) GetAllMultiDepositOrder(ctx sdk.Context, poolId string) (list []types.MultiAssetDepositOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.MultiAssetDepositOrder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}
