package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTakeMultiAssetDeposit = "take_multi_asset_deposit"

var _ sdk.Msg = &MsgTakeMultiAssetDepositRequest{}

func NewMsgTakeMultiAssetDeposit(sender, poolId string, orderId uint64) *MsgTakeMultiAssetDepositRequest {
	return &MsgTakeMultiAssetDepositRequest{
		Sender:  sender,
		PoolId:  poolId,
		OrderId: orderId,
	}
}

func (msg *MsgTakeMultiAssetDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgTakeMultiAssetDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgTakeMultiAssetDepositRequest) GetSigners() []sdk.AccAddress {
	signers := []sdk.AccAddress{}

	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, creator)
	return signers
}

func (msg *MsgTakeMultiAssetDepositRequest) GetSignBytes() []byte {
	marshaledMsg := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(marshaledMsg)
}

func (msg *MsgTakeMultiAssetDepositRequest) ValidateBasic() error {
	return nil
}
