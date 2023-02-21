package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDeposit = "deposit"

var _ sdk.Msg = &MsgDepositRequest{}

func NewMsgDeposit(poolId string, sender string) *MsgDepositRequest {
	return &MsgDepositRequest{
		PoolId: poolId,
		Sender: sender,
	}
}

func (msg *MsgDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDepositRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDepositRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if len(msg.Tokens) == 0 {
		return ErrInvalidLength
	}
	return nil
}
