package types

const (
	// ModuleName defines the IBC swap name
	ModuleName = "interchainswap"

	// Version defines the current version the IBC swap
	// module supports
	Version = "ics101-1"

	// PortID is the default port id that swap module binds to
	PortID = ModuleName

	// StoreKey is the store key string for IBC swap
	StoreKey = ModuleName

	// RouterKey is the message route for IBC swap
	RouterKey = ModuleName

	// QuerierRoute is the querier route for IBC swap
	QuerierRoute = ModuleName

	Multiplier      = 1e18
	MaximumSlippage = 10000
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = []byte{0x01}
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
