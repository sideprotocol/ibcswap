package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgWithdrawRequest{}

func NewMsgWithdraw(creator string, sender string, denomOut string) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		Sender:   sender,
		DenomOut: denomOut,
	}
}

func (msg *MsgWithdrawRequest) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawRequest) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgWithdrawRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
