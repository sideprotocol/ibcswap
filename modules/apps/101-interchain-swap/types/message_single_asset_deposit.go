package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDeposit = "deposit"

var _ sdk.Msg = &MsgSingleAssetDepositRequest{}

func NewMsgSingleAssetDeposit(poolId string, sender string, token *sdk.Coin) *MsgSingleAssetDepositRequest {
	return &MsgSingleAssetDepositRequest{
		PoolId: poolId,
		Sender: sender,
		Token:  token,
	}
}

func (msg *MsgSingleAssetDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgSingleAssetDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgSingleAssetDepositRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSingleAssetDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSingleAssetDepositRequest) ValidateBasic() error {
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
