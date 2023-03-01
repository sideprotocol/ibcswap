package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgWithdrawRequest{}

func NewMsgWithdraw(sender string, poolCoin *sdk.Coin, denomOut string) *MsgWithdrawRequest {
	return &MsgWithdrawRequest{
		Sender:   sender,
		PoolCoin: poolCoin,
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
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if strings.TrimSpace(msg.DenomOut) == "" {
		return errorsmod.Wrapf(ErrEmptyDenom, "none exist denom (%s)", err)
	}
	if msg.PoolCoin == nil || msg.PoolCoin.Amount.LTE(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", "")
	}
	return nil
}
