package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

var (
	// DefaultRelativePacketTimeoutHeight is the default packet timeout height (in blocks) relative
	// to the current block height of the counterparty chain provided by the client state. The
	// timeout is disabled when set to 0.
	DefaultRelativePacketTimeoutHeight = "0-1000"

	// DefaultRelativePacketTimeoutTimestamp is the default packet timeout timestamp (in nanoseconds)
	// relative to the current block timestamp of the counterparty chain provided by the client
	// state. The timeout is disabled when set to 0. The default is currently set to a 10 minute
	// timeout.
	DefaultRelativePacketTimeoutTimestamp = uint64((time.Duration(10) * time.Minute).Nanoseconds())
)

// NewIBCSwapPacketData constructs a new IBCSwapPacketData instance
func NewIBCSwapPacketData(
	mType SwapMessageType,
	data []byte,
) IBCSwapPacketData {
	return IBCSwapPacketData{
		Type: mType,
		Data: data,
	}
}

// ValidateBasic is used for validating the token swap.
func (m IBCSwapPacketData) ValidateBasic() error {
	if m.Data == nil || len(m.Data) == 0 {
		return ErrInvalidLengthPacket
	}
	return nil
}

// GetBytes is a helper for serialising
func (m IBCSwapPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&m))
}
