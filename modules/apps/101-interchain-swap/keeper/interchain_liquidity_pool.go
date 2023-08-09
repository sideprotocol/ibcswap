package keeper

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// SetInterchainLiquidityPool set a specific interchainLiquidityPool in the store from its index
func (k Keeper) SetInitialPoolAssets(ctx sdk.Context, poolId string, tokens sdk.Coins) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId))

	// Convert each sdk.Coin to a string in the format "Denom:Amount"
	var coinStrings []string
	for _, coin := range tokens {
		coinStrings = append(coinStrings, coin.Denom+":"+coin.Amount.String())
	}

	// Join the coin strings with a comma
	coinsString := strings.Join(coinStrings, ",")

	// Convert the final string to a byte slice
	b := []byte(coinsString)
	store.Set(types.InterchainLiquidityPoolKey(poolId), b)
}

func (k Keeper) GetInitialPoolAssets(ctx sdk.Context, poolId string) sdk.Coins {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(poolId))

	b := store.Get(types.InterchainLiquidityPoolKey(poolId))
	if b == nil {
		return sdk.NewCoins()
	}

	// Convert bytes to string
	coinsString := string(b)

	// Split the coins string into coin strings
	coinStrings := strings.Split(coinsString, ",")

	var tokens sdk.Coins
	// Convert each coin string back to sdk.Coin
	for _, coinString := range coinStrings {
		parts := strings.Split(coinString, ":")
		if len(parts) != 2 {
			return sdk.NewCoins()
		}

		denom := parts[0]
		amount, ok := sdk.NewIntFromString(parts[1])
		if !ok {
			return sdk.NewCoins()
		}

		coin := sdk.NewCoin(denom, amount)
		tokens = append(tokens, coin)
	}
	return tokens
}

// // GetInterchainLiquidityPool returns a interchainLiquidityPool from its index
// func (k Keeper) GetInterchainLiquidityPool(
// 	ctx sdk.Context,
// 	poolId string,

// ) (val types.InterchainLiquidityPool, found bool) {

// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
// 	b := store.Get(types.InterchainLiquidityPoolKey(
// 		poolId,
// 	))
// 	if b == nil {
// 		return val, false
// 	}

// 	k.cdc.MustUnmarshal(b, &val)
// 	return val, true
// }

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

// func (k Keeper) SetInterchainLiquidityPool(ctx sdk.Context, interchainLiquidityPool types.InterchainLiquidityPool) {
// 	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

// 	// Get current pool count
// 	poolCount := k.GetPoolCount(ctx)

// 	// Increment the count
// 	poolCount++

// 	// Set the new count
// 	k.SetPoolCount(ctx, poolCount)

// 	// Marshal the pool and set in store
// 	b := k.cdc.MustMarshal(&interchainLiquidityPool)
// 	store.Set(GetInterchainLiquidityPoolKey(poolCount), b)

// 	// Check if we exceed max pools
// 	if poolCount > types.MaxPoolCount {
// 		// Remove the oldest pool
// 		store.Delete(GetInterchainLiquidityPoolKey(poolCount - types.MaxPoolCount))
// 	}
// }

func (k Keeper) GetPoolCount(ctx sdk.Context) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	b := store.Get(types.CurrentPoolCountKey)
	if b == nil {
		return 0
	}
	return binary.BigEndian.Uint64(b)
}

func (k Keeper) SetPoolCount(ctx sdk.Context, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, count)
	store.Set(types.CurrentPoolCountKey, b)
}

func GetInterchainLiquidityPoolKey(count uint64) []byte {
	return []byte(fmt.Sprintf("%020d", count))
}

// GetAllInterchainLiquidityPool returns all interchainLiquidityPool
func (k Keeper) GetAllInterchainLiquidityPool(ctx sdk.Context) (list []types.InterchainLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

	// Start from the latest pool and move to the oldest
	poolCount := k.GetPoolCount(ctx)
	for i := poolCount; i >= 1 && (poolCount-i) < types.MaxPoolCount; i-- {
		b := store.Get(GetInterchainLiquidityPoolKey(i))
		if b == nil {
			continue
		}
		var val types.InterchainLiquidityPool
		k.cdc.MustUnmarshal(b, &val)
		list = append(list, val)
	}
	return
}

// Sets the mapping between poolId and its count index
func (k Keeper) SetPoolIdToCountMapping(ctx sdk.Context, poolId string, count uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolIdToCountKeyPrefix)
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, count)
	store.Set([]byte(poolId), b)
}

// Gets the count index of the poolId
func (k Keeper) GetCountByPoolId(ctx sdk.Context, poolId string) (count uint64, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.PoolIdToCountKeyPrefix)
	b := store.Get([]byte(poolId))
	if b == nil {
		return 0, false
	}
	return binary.BigEndian.Uint64(b), true
}

// Modified SetInterchainLiquidityPool
func (k Keeper) SetInterchainLiquidityPool(ctx sdk.Context, interchainLiquidityPool types.InterchainLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

	// Get current pool count
	poolCount := k.GetPoolCount(ctx)

	// Increment the count
	poolCount++

	// Set the new count
	k.SetPoolCount(ctx, poolCount)

	// Set the poolId to count mapping
	k.SetPoolIdToCountMapping(ctx, interchainLiquidityPool.Id, poolCount)

	// Marshal the pool and set in store
	b := k.cdc.MustMarshal(&interchainLiquidityPool)
	store.Set(GetInterchainLiquidityPoolKey(poolCount), b)

	// Check if we exceed max pools
	if poolCount > types.MaxPoolCount {
		// Remove the oldest pool
		store.Delete(GetInterchainLiquidityPoolKey(poolCount - types.MaxPoolCount))
	}
}

// Modified GetInterchainLiquidityPool
func (k Keeper) GetInterchainLiquidityPool(ctx sdk.Context, poolId string) (val types.InterchainLiquidityPool, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))

	count, found := k.GetCountByPoolId(ctx, poolId)
	if !found {
		return val, false
	}

	b := store.Get(GetInterchainLiquidityPoolKey(count))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}
