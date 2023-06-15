package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgMultiAssetWithdrawRequest{}

func NewMsgMultiAssetWithdraw(poolId string, sourceReceiver, destinationReceiver string, poolToken *sdk.Coin) *MsgMultiAssetWithdrawRequest {
	return &MsgMultiAssetWithdrawRequest{
		PoolId:               poolId,
		Receiver:             sourceReceiver,
		CounterPartyReceiver: destinationReceiver,
		PoolToken:            poolToken,
	}
}

func (msg *MsgMultiAssetWithdrawRequest) Route() string {
	return RouterKey
}

func (msg *MsgMultiAssetWithdrawRequest) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgMultiAssetWithdrawRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Receiver)
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
	_, err := sdk.AccAddressFromBech32(msg.Receiver)
	if err != nil {
		return ErrInvalidAddress
	}
	if msg.PoolToken.Amount.LTE(sdk.NewInt(0)) {
		return ErrInvalidAmount
	}
	return nil
}
