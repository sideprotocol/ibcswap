package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDoubleDeposit = "deposit"

var _ sdk.Msg = &MsgSingleAssetDepositRequest{}

func NewMsgMultiAssetDeposit(poolId string, senders []string, tokens []*sdk.Coin, sig []byte) *MsgMultiAssetDepositRequest {
	return &MsgMultiAssetDepositRequest{
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

func (msg *MsgMultiAssetDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgMultiAssetDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgMultiAssetDepositRequest) GetSigners() []sdk.AccAddress {
	signers := []sdk.AccAddress{}
	creator, err := sdk.AccAddressFromBech32(msg.LocalDeposit.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, creator)
	return signers
}

func (msg *MsgMultiAssetDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMultiAssetDepositRequest) ValidateBasic() error {

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
