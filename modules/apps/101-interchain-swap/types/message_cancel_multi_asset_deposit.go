package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCancelMultiAssetDeposit = "cancel_multi_asset_deposit"

var _ sdk.Msg = &MsgCancelMultiAssetDepositRequest{}

func NewMsgCancelMultiAssetDeposit(
	sourcePort,
	sourceChannel,
	creator,
	poolId,
	orderId string,
) *MsgCancelMultiAssetDepositRequest {

	return &MsgCancelMultiAssetDepositRequest{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Creator:       creator,
		PoolId:        poolId,
		OrderId:       orderId,
	}
}

func (msg *MsgCancelMultiAssetDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgCancelMultiAssetDepositRequest) Type() string {
	return TypeMsgCancelPool
}

func (msg *MsgCancelMultiAssetDepositRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelMultiAssetDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelMultiAssetDepositRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	if msg.SourceChannel == "" {
		return ErrMissedIBCParams
	}
	if msg.PoolId == "" {
		return ErrInvalidPoolId
	}
	if msg.SourcePort == "" {
		msg.SourcePort = PortID
	}

	return nil
}
