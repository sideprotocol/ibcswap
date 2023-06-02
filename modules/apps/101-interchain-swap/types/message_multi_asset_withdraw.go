package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgWithdraw = "withdraw"

var _ sdk.Msg = &MsgMultiAssetWithdrawRequest{}

func NewMsgMultiAssetWithdraw(sourceSender, targetSender string, sourcePoolToken *sdk.Coin, targetPoolToken *sdk.Coin) *MsgMultiAssetWithdrawRequest {
	return &MsgMultiAssetWithdrawRequest{

		Withdraws: []*WithdrawAsset{
			{
				Receiver: sourceSender,
				Balance:  sourcePoolToken,
			},
			{
				Receiver: targetSender,
				Balance:  targetPoolToken,
			},
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
	creator, err := sdk.AccAddressFromBech32(msg.Withdraws[0].Receiver)
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
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidAddress
	}

	for _, asset := range msg.Withdraws {
		_, err := sdk.AccAddressFromBech32(asset.Receiver)
		if err != nil {
			return ErrInvalidAddress
		}
		if asset.Balance.Amount.Equal(sdk.NewInt(0)) {
			return ErrInvalidAmount
		}
	}
	return nil
}
