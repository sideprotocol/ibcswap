package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCancelPool = "cancel_pool"

var _ sdk.Msg = &MsgMakePoolRequest{}

func NewMsgCancelPool(
	sourcePort string,
	sourceChannel string,
	creator string,
	poolID string,
) *MsgCancelPoolRequest {

	return &MsgCancelPoolRequest{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Creator:       creator,
		PoolId:        poolID,
	}
}

func (msg *MsgCancelPoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgCancelPoolRequest) Type() string {
	return TypeMsgCancelPool
}

func (msg *MsgCancelPoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelPoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelPoolRequest) ValidateBasic() error {
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
