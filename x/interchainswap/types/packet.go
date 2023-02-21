package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAtomicSwapPacketData(
	mType MessageType,
	data []byte,
) IBCSwapDataPacket {
	return IBCSwapDataPacket{
		Type: mType,
		Data: data,
	}
}

// ValidateBasic is used for validating the token swap.
func (pd IBCSwapDataPacket) ValidateBasic() error {
	if pd.Data == nil || len(pd.Data) == 0 {
		return ErrInvalidLengthPacket
	}
	return nil
}

// GetBytes is a helper for serialising
func (pd IBCSwapDataPacket) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&pd))
}
