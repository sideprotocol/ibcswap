package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgMultiAssetWithdrawRequest{}

func NewMsgWithdraw(localSender, remoteSender string, localPoolCoin *sdk.Coin, remotePoolCoin *sdk.Coin) *MsgMultiAssetWithdrawRequest {
	return &MsgMultiAssetWithdrawRequest{
		LocalWithdraw: &WithdrawRequest{
			Sender:   localSender,
			PoolCoin: localPoolCoin,
		},

		RemoteWithdraw: &WithdrawRequest{
			Sender:   localSender,
			PoolCoin: remotePoolCoin,
		},
	}
}

func (msg *MsgMultiAssetWithdrawRequest) Route() string {
	return RouterKey
}

func (msg *MsgMultiAssetWithdrawRequest) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgMultiAssetWithdrawRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.LocalWithdraw.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMultiAssetWithdrawRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMultiAssetWithdrawRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.LocalWithdraw.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.RemoteWithdraw.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.LocalWithdraw.PoolCoin == nil || msg.LocalWithdraw.PoolCoin.Amount.LTE(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", msg.LocalWithdraw.PoolCoin.Amount)
	}

	if msg.RemoteWithdraw.PoolCoin == nil || msg.RemoteWithdraw.PoolCoin.Amount.LTE(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", msg.RemoteWithdraw.PoolCoin.Amount)
	}
	return nil
}
