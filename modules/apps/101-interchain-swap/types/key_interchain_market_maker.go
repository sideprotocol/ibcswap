package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// InterchainMarketMakerKeyPrefix is the prefix to retrieve all InterchainMarketMaker
	InterchainMarketMakerKeyPrefix = "InterchainMarketMaker/value/"
)

// InterchainMarketMakerKey returns the store key to retrieve a InterchainMarketMaker from the index fields
func InterchainMarketMakerKey(
	poolId string,
) []byte {
	var key []byte

	poolIdBytes := []byte(poolId)
	key = append(key, poolIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
