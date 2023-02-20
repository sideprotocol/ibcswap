package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwap = "swap"

var _ sdk.Msg = &MsgSwapRequest{}

func NewMsgSwap(sender string, slippage uint64, recipient string) *MsgSwapRequest {
	return &MsgSwapRequest{
		Sender:    sender,
		Slippage:  slippage,
		Recipient: recipient,
	}
}

func (msg *MsgSwapRequest) Route() string {
	return RouterKey
}

func (msg *MsgSwapRequest) Type() string {
	return TypeMsgSwap
}

func (msg *MsgSwapRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSwapRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSwapRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
