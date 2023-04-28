package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeposit = "deposit"

var _ sdk.Msg = &MsgSingleDepositRequest{}

func NewMsgSingleDeposit(poolId string, sender string, token *sdk.Coin) *MsgSingleDepositRequest {
	return &MsgSingleDepositRequest{
		PoolId: poolId,
		Sender: sender,
		Token:  token,
	}
}

func (msg *MsgSingleDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgSingleDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgSingleDepositRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSingleDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSingleDepositRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if msg.Token == nil {
		return errorsmod.Wrapf(ErrInvalidMessage, "invalid token length (%d)", msg.Token)
	}

	if msg.Token.Amount.Equal(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrFailedDeposit, "because of %s", ErrInvalidAmount)
	}
	return nil
}
