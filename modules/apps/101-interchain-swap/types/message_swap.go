package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSwap = "swap"

var _ sdk.Msg = &MsgSwapRequest{}

func NewMsgSwap(swapType SwapMsgType, sender string, slippage uint64, recipient string, tokenIn, tokenOut *sdk.Coin) *MsgSwapRequest {
	return &MsgSwapRequest{
		Sender:    sender,
		Slippage:  slippage,
		Recipient: recipient,
		TokenIn:   tokenIn,
		TokenOut:  tokenOut,
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
	if msg.SwapType != SwapMsgType_LEFT && msg.SwapType != SwapMsgType_RIGHT {
		return ErrInvalidSwapType
	}
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
	if msg.TokenIn.Amount.LTE(sdk.NewInt(0)) || msg.TokenOut.Amount.LTE(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid token amounts (%s)", err)
	}
	if msg.Slippage == 0 {
		return ErrInvalidSlippage
	}

	return nil
}
