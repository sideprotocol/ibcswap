package keeper

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

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

// MoveOrderToBottom moves the atomic order with the given ID to the bottom of the list.
func (k Keeper) MoveOrderToBottom(ctx sdk.Context, orderId string) error {
	// Step 1: Retrieve the item based on the given ID.
	order, found := k.GetAtomicOrder(ctx, orderId)
	if !found {
		return types.ErrNotFoundOrder
	}
	// Step 2: Remove the item from its current position.
	k.RemoveOrder(ctx, orderId)
	// Step 3: Append the item to the end.
	count := k.GetAtomicOrderCount(ctx)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)
	bz := k.cdc.MustMarshal(&order)
	store.Set(GetOrderIDBytes(count), bz)
	k.SetAtomicOrderCount(ctx, count+1)
	return nil
}

func (k Keeper) TrimExcessOrders(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.OTCOrderBookKey)

	totalCount := k.GetAtomicOrderCount(ctx)
	if totalCount <= types.MaxOrderCount {
		return
	}
	// Calculate number of items to be removed
	excess := totalCount - types.MaxOrderCount
	for i := uint64(0); i < excess; i++ {
		// As items are appended, to remove from the bottom, we need to remove the items
		// starting from totalCount - i (i.e., the last item in the list, then the second last, etc.)
		idToRemove := totalCount - i - 1
		store.Delete(GetOrderIDBytes(idToRemove))
		k.SetAtomicOrderCount(ctx, totalCount-i-1)
	}
}
