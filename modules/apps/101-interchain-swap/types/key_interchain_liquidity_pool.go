package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// InterchainLiquidityPoolKeyPrefix is the prefix to retrieve all InterchainLiquidityPool
	InterchainLiquidityPoolKeyPrefix = "InterchainLiquidityPool/value/"
)

// InterchainLiquidityPoolKey returns the store key to retrieve a InterchainLiquidityPool from the index fields
func InterchainLiquidityPoolKey(
	poolId string,
) []byte {
	var key []byte

	poolIdBytes := []byte(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}



