package types

import "encoding/binary"

var _ binary.ByteOrder

const (
	// MultiDepositOrderKeyPrefix is the prefix to retrieve all MultiDepositOrder
	MultiDepositOrderKeyPrefix             = "MultiDepositOrder/value/"
	MultiDepositOrderCountKeyPrefix        = "MultiDepositOrderCount/value/"
	MultiDepositOrderIDByCreatorsKeyPrefix = "MultiDepositOrderIDByCreator/value/"
)

// MultiDepositOrderPrefixKey returns the store key to retrieve a MultiDepositOrder from the index fields
func MultiDepositOrderPrefixKey(
	orderId string,
) []byte {
	var key []byte

	orderIdBytes := []byte(orderId)
	key = append(key, orderIdBytes...)
	key = append(key, []byte("/")...)

	return key
}
