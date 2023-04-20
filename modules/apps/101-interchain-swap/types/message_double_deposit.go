package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDoubleDeposit = "deposit"

var _ sdk.Msg = &MsgDepositRequest{}

func NewMsgDoubleDeposit(poolId string, senders []string, tokens []*sdk.Coin, sig []byte) *MsgDoubleDepositRequest {
	return &MsgDoubleDepositRequest{
		PoolId: poolId,
		LocalDeposit: &LocalDeposit{
			Sender: senders[0],
			Token:  tokens[0],
		},
		RemoteDeposit: &RemoteDeposit{
			Sender:    senders[1],
			Token:     tokens[1],
			Signature: sig,
		},
	}
}

func (msg *MsgDoubleDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgDoubleDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDoubleDepositRequest) GetSigners() []sdk.AccAddress {
	signers := []sdk.AccAddress{}
	creator, err := sdk.AccAddressFromBech32(msg.LocalDeposit.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, creator)
	return signers
}

func (msg *MsgDoubleDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDoubleDepositRequest) ValidateBasic() error {

	_, err := sdk.AccAddressFromBech32(msg.LocalDeposit.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.RemoteDeposit.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
