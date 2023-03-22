package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePoolRequest{}

func NewMsgCreatePool(sourcePort string, sourceChannel string, sender string, weight string, denoms []string, decimals []uint32, initialLiquidity uint64) *MsgCreatePoolRequest {
	return &MsgCreatePoolRequest{
		SourcePort:      sourcePort,
		SourceChannel:   sourceChannel,
		Sender:          sender,
		Weight:          weight,
		Denoms:          denoms,
		Decimals:        decimals,
		InitalLiquidity: initialLiquidity,
	}
}

func (msg *MsgCreatePoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgCreatePoolRequest) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePoolRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidAddress
	}

	denomSize := len(msg.Denoms)
	//validation message
	if denomSize != 2 {
		return ErrInvalidDenomPair
	}

	if len(msg.Decimals) != 2 {
		return ErrInvalidDecimalPair
	}
	if len(strings.Split(msg.Weight, ":")) != 2 {
		return ErrInvalidWeightPair
	}

	if msg.InitalLiquidity <= 0 {
		return ErrEmptyInitialLiquidity
	}
	return nil
}
