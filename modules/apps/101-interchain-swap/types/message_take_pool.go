package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgTakePool = "take_pool"

var _ sdk.Msg = &MsgTakePoolRequest{}

func NewMsgTakePool(creator, poolId, port, channel string) *MsgTakePoolRequest {
	return &MsgTakePoolRequest{
		Creator: creator,
		PoolId:  poolId,
		Port:    port,
		Channel: channel,
	}
}

func (msg *MsgTakePoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgTakePoolRequest) Type() string {
	return TypeMsgMakePool
}

func (msg *MsgTakePoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTakePoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTakePoolRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	if msg.Channel == "" {
		return ErrMissedIBCParams
	}
	if msg.Port == "" {
		msg.Port = PortID
	}
	return nil
}
