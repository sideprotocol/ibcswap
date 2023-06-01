package keeper

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// SetInterchainLiquidityPool set a specific interchainLiquidityPool in the store from its index
func (k Keeper) SetInterchainLiquidityPool(ctx sdk.Context, interchainLiquidityPool types.InterchainLiquidityPool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.InterchainLiquidityPoolKeyPrefix))
	b := k.cdc.MustMarshal(&interchainLiquidityPool)
	store.Set(types.InterchainLiquidityPoolKey(
		interchainLiquidityPool.PoolId,
	), b)
}

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
