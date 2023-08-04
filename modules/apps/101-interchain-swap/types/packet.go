package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewInterchainSwapPacketData(
	mType SwapMessageType,
	data []byte,
) IBCSwapPacketData {
	return IBCSwapPacketData{
		Type: mType,
		Data: data,
	}
}

// ValidateBasic is used for validating the token swap.
func (pd IBCSwapPacketData) ValidateBasic() error {
	if pd.Data == nil || len(pd.Data) == 0 {
		return ErrInvalidLengthPacket
	}
	return nil
}

// GetBytes is a helper for serialising
func (pd IBCSwapPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&pd))
}

type AckData[T any, U any] struct {
	Req T
	Res U
}

func (s *StateChange) FindOutByDenom(denom string) (*sdk.Coin, error) {
	for _, out := range s.Out {
		if out.Denom == denom {
			return out, nil
		}
	}
	return nil, ErrInvalidDenom
}
