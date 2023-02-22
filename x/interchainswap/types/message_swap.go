package types

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid recipient address (%s)", err)
	}

	if msg.TokenIn == nil || msg.TokenOut == nil || strings.TrimSpace(msg.TokenIn.Denom) == "" || strings.TrimSpace(msg.TokenOut.Denom) == "" {
		return errorsmod.Wrapf(ErrEmptyDenom, "missed token denoms (%s)", err)
	}
	if msg.TokenIn.Amount.LTE(math.NewInt(0)) || msg.TokenOut.Amount.LTE(math.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid token amounts (%s)", err)
	}
	if msg.Slippage == 0 {
		return ErrInvalidSlippage
	}

	return nil
}
