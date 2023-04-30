package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSingleWithdraw = "single withdraw"

var _ sdk.Msg = &MsgSingleAssetDepositRequest{}

func NewMsgSingleAssetWithdraw(sender, denomOut string, poolCoin *sdk.Coin) *MsgSingleAssetWithdrawRequest {
	return &MsgSingleAssetWithdrawRequest{
		Sender:   sender,
		PoolCoin: poolCoin,
		DenomOut: denomOut,
	}
}

func (msg *MsgSingleAssetWithdrawRequest) Route() string {
	return RouterKey
}

func (msg *MsgSingleAssetWithdrawRequest) Type() string {
	return TypeMsgWithdraw
}

func (msg *MsgSingleAssetWithdrawRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSingleAssetWithdrawRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSingleAssetWithdrawRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.PoolCoin == nil || msg.PoolCoin.Amount.LTE(sdk.NewInt(0)) {
		return errorsmod.Wrapf(ErrInvalidAmount, "invalid pool coin amount (%s)", msg.PoolCoin.Amount)
	}
	return nil
}
