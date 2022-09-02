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

// NewAtomicSwapPacketData constructs a new AtomicSwapPacketData instance
func NewAtomicSwapPacketData(
	mType SwapMessageType,
	data []byte,
	memo string,
) AtomicSwapPacketData {
	return AtomicSwapPacketData{
		Type: mType,
		Data: data,
		Memo: memo,
	}
}

// ValidateBasic is used for validating the token swap.
func (pd AtomicSwapPacketData) ValidateBasic() error {
	if pd.Data == nil || len(pd.Data) == 0 {
		return ErrInvalidLengthPacket
	}
	return nil
}

// GetBytes is a helper for serialising
func (pd AtomicSwapPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&pd))
}
