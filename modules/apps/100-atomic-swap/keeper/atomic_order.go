package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

// // / Atomic orders
// // GetAtomicOrder returns the OTCOrder for the swap module.
// func (k Keeper) GetAtomicOrder(ctx sdk.Context, orderId string) (types.Order, bool) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)

// 	// Retrieve the full key from the secondary index
// 	indexStore := prefix.NewStore(store, []byte(types.OrderBookIndexKey))
// 	fullKey := indexStore.Get([]byte(orderId))
// 	if fullKey == nil {
// 		return types.Order{}, false
// 	}

// 	bz := store.Get(fullKey)
// 	if bz == nil {
// 		return types.Order{}, false
// 	}

// 	order := k.MustUnmarshalOrder(bz)
// 	return order, true
// }

// // HasAtomicOrder checks if a the key with the given id exists on the store.
// func (k Keeper) HasAtomicOrder(ctx sdk.Context, orderId string) bool {
// 	key, err := hex.DecodeString(orderId)
// 	if err != nil {
// 		return false
// 	}
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
// 	return store.Has(key)
// }

// // SetAtomicOrder sets a new OTCOrder to the store.
// func (k Keeper) SetAtomicOrder(ctx sdk.Context, order types.Order) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
// 	// Use the creation timestamp as a prefix to the order ID.
// 	key := fmt.Sprintf(order.Id)
// 	bz := k.MustMarshalOrder(order)
// 	store.Set([]byte(key), bz)
// }

// GetAuctionCount get the total number of auction
func (k Keeper) GetAtomicOrderCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKeyCountKey)
	byteKey := types.OTCOrderBookKeyIndexKey
	bz := store.Get(byteKey)
	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}

	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// GetAuctionCount get the total number of auction
func (k Keeper) GetAtomicOrderCountByOrderId(ctx sdk.Context, orderId string) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKeyIndexKey)
	byteKey := []byte(orderId)
	bz := store.Get(byteKey)
	// Count doesn't exist: no element
	if bz == nil {
		return 0
	}
	// Parse bytes
	return binary.BigEndian.Uint64(bz)
}

// SetAuctionCount set the total number of auction
func (k Keeper) SetAtomicOrderCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKeyCountKey)
	byteKey := types.OTCOrderBookKeyIndexKey
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// SetAuctionCount set the total number of auction
func (k Keeper) SetAtomicOrderCountToOrderID(ctx sdk.Context, orderId string, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKeyIndexKey)
	byteKey := []byte(orderId)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, count)
	store.Set(byteKey, bz)
}

// AppendAuction appends a auction in the store with a new id and update the count
func (k Keeper) AppendAtomicOrder(
	ctx sdk.Context,
	order types.Order,
) uint64 {
	// Create the auction
	count := k.GetAtomicOrderCount(ctx)
	// Set the ID of the appended value
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	appendedValue := k.cdc.MustMarshal(&order)
	store.Set(GetOrderIDBytes(count), appendedValue)
	// Update auction count
	k.SetAtomicOrderCountToOrderID(ctx, order.Id, count)
	k.SetAtomicOrderCount(ctx, count+1)
	return count
}

// SetAuction set a specific auction in the store
func (k Keeper) SetAtomicOrder(ctx sdk.Context, order types.Order) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	b := k.cdc.MustMarshal(&order)
	id := k.GetAtomicOrderCountByOrderId(ctx, order.Id)
	store.Set(GetOrderIDBytes(id), b)
}

// GetAuction returns a auction from its id
func (k Keeper) GetAtomicOrder(ctx sdk.Context, orderId string) (val types.Order, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	id := k.GetAtomicOrderCountByOrderId(ctx, orderId)
	b := store.Get(GetOrderIDBytes(id))
	if b == nil {
		return val, false
	}
	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAuction removes a auction from the store
func (k Keeper) RemoveOrder(ctx sdk.Context, orderId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	id := k.GetAtomicOrderCountByOrderId(ctx, orderId)
	store.Delete(GetOrderIDBytes(id))
}

// GetAllAuction returns all auction
func (k Keeper) GetAllOrder(ctx sdk.Context) (list []types.Order) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var val types.Order
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}
	return
}

// GetAuctionIDBytes returns the byte representation of the ID
func GetOrderIDBytes(id uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, id)
	return bz
}

// GetAuctionIDFromBytes returns ID in uint64 format from a byte array
func GetBidIDFromBytes(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}
