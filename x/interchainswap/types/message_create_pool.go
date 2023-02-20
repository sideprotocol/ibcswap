package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
