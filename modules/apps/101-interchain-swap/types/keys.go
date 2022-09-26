package types

import (
	"crypto/sha256"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName defines the IBC swap name
	ModuleName = "swap"

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
)

var (
	// PortKey defines the key to store the port ID in store
	PortKey = []byte{0x01}
	// LimitOrderBookKey defines the key to store the denomination trace info in store
	LimitOrderBookKey = []byte{0x02}
	// OTCOrderBookKey defines the key to store the denomination trace info in store
	OTCOrderBookKey = []byte{0x03}
)

// GetEscrowAddress returns the escrow address for the specified channel.
// The escrow address follows the format as outlined in ADR 028:
// https://github.com/cosmos/cosmos-sdk/blob/master/docs/architecture/adr-028-public-key-addresses.md
func GetEscrowAddress(portID, channelID string) sdk.AccAddress {
	// a slash is used to create domain separation between port and channel identifiers to
	// prevent address collisions between escrow addresses created for different channels
	contents := fmt.Sprintf("%s/%s", portID, channelID)

	// ADR 028 AddressHash construction
	preImage := []byte(Version)
	preImage = append(preImage, 0)
	preImage = append(preImage, contents...)
	hash := sha256.Sum256(preImage)
	return hash[:20]
}
