package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// GetMultiDepositOrderCount get the total number of multiDepositOrder
func (k Keeper) GetMultiDepositOrderCount(ctx sdk.Context, poolId string) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(poolId + types.MultiDepositOrderCountKeyPrefix)
	bz := store.Get(byteKey)

	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetMultiDepositOrderCount set the total number of multiDepositOrder
func (k Keeper) SetMultiDepositOrderCount(ctx sdk.Context, poolId string, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(poolId + types.MultiDepositOrderCountKeyPrefix)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// SetMultiDepositOrderCount set the total number of multiDepositOrder
func (k Keeper) SetMultiDepositOrderLatestOrderByCreators(ctx sdk.Context, poolId, sourceMaker, destinationMaker string, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(poolId + sourceMaker + destinationMaker + types.MultiDepositOrderIDByCreatorsKeyPrefix)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// SetMultiDepositOrderCount set the total number of multiDepositOrder
func (k Keeper) GetMultiDepositOrderLatestOrderByCreators(ctx sdk.Context, poolId, sourceMaker, destinationMaker string) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), []byte{})
	byteKey := types.KeyPrefix(poolId + sourceMaker + destinationMaker + types.MultiDepositOrderIDByCreatorsKeyPrefix)
	bz := store.Get(byteKey)
	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}
	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// AppendMultiDepositOrder appends a multiDepositOrder in the store with a new id and update the count
func (k Keeper) AppendMultiDepositOrder(
	ctx sdk.Context,
	poolId string,
	multiDepositOrder types.MultiAssetDepositOrder,
) uint64 {
	// Create the multiDepositOrder
	count := k.GetMultiDepositOrderCount(ctx, poolId)
	// Set the ID of the appended value
	multiDepositOrder.Id = count

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	appendedValue := k.cdc.MustMarshal(&multiDepositOrder)
	store.Set(GetMultiDepositOrderIDBytes(multiDepositOrder.Id), appendedValue)

	// Update multiDepositOrder count
	k.SetMultiDepositOrderCount(ctx, poolId, count+1)
	k.SetMultiDepositOrderLatestOrderByCreators(ctx, poolId, multiDepositOrder.SourceMaker, multiDepositOrder.DestinationTaker, count)
	return count
}

// SetMultiDepositOrder set a specific multiDepositOrder in the store
func (k Keeper) SetMultiDepositOrder(ctx sdk.Context, poolId string, multiDepositOrder types.MultiAssetDepositOrder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	b := k.cdc.MustMarshal(&multiDepositOrder)
	store.Set(GetMultiDepositOrderIDBytes(multiDepositOrder.Id), b)
}

// GetMultiDepositOrder returns a multiDepositOrder from its id
func (k Keeper) GetMultiDepositOrder(ctx sdk.Context, poolId string, id uint64) (val types.MultiAssetDepositOrder, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	b := store.Get(GetMultiDepositOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetMultiDepositOrder returns a multiDepositOrder from its id
func (k Keeper) GetLatestMultiDepositOrder(ctx sdk.Context, poolId string) (val types.MultiAssetDepositOrder, found bool) {
	id := k.GetMultiDepositOrderCount(ctx, poolId)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	b := store.Get(GetMultiDepositOrderIDBytes(id - 1))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)

	return val, true
}

// GetMultiDepositOrder returns a multiDepositOrder from its id
func (k Keeper) RemoveLatestMultiDepositOrder(ctx sdk.Context, poolId string) {
	id := k.GetMultiDepositOrderCount(ctx, poolId)
	order, found := k.GetMultiDepositOrder(ctx, poolId, id)
	if found {
		k.SetMultiDepositOrderLatestOrderByCreators(ctx, poolId, order.SourceMaker, order.DestinationTaker, id-1)
		k.RemoveMultiDepositOrder(ctx, poolId, id)
		k.SetMultiDepositOrderCount(ctx, poolId, id-1)
	}
}

// RemoveMultiDepositOrder removes a multiDepositOrder from the store
func (k Keeper) RemoveMultiDepositOrder(ctx sdk.Context, poolId string, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId+types.MultiDepositOrderKeyPrefix))
	store.Delete(GetMultiDepositOrderIDBytes(id))
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

// GetMultiDepositOrderIDBytes returns the byte representation of the ID
func GetMultiDepositOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetMultiDepositOrderIDFromBytes returns ID in uint64 format from a byte array
func GetMultiDepositOrderIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
