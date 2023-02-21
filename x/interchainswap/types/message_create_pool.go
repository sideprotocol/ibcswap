package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePoolRequest{}

func NewMsgCreatePool(creator string, sourcePort string, sourceChannel string, sender string, weight string) *MsgCreatePoolRequest {
	return &MsgCreatePoolRequest{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Sender:        sender,
		Weight:        weight,
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
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	//validation message
	if len(msg.Denoms) != 2 {
		return ErrInvalidDenomPair
	}

	if len(msg.Decimals) != 2 {
		return ErrInvalidDecimalPair
	}
	if len(strings.Split(msg.Weight, ":")) != 2 {
		return ErrInvalidWeightPair
	}
	return nil
}
