package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDoubleDeposit = "deposit"

var _ sdk.Msg = &MsgDepositRequest{}

func NewMsgDoubleDeposit(poolId string, sender string, tokens []*sdk.Coin, subTx *CPDepositTx) *MsgDoubleDepositRequest {
	return &MsgDoubleDepositRequest{
		PoolId:      poolId,
		Sender:      sender,
		Tokens:      tokens,
		CpDepositTx: subTx,
	}
}

func (msg *MsgDoubleDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgDoubleDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDoubleDepositRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDoubleDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDoubleDepositRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	if len(msg.Tokens) == 0 {
		return errorsmod.Wrapf(ErrInvalidTokenLength, "invalid token length (%d)", len(msg.Tokens))
	}
	denoms := map[string]int{}
	for _, token := range msg.Tokens {
		if _, ok := denoms[token.Denom]; ok {
			return errorsmod.Wrapf(ErrFailedDeposit, "because of %s", ErrInvalidDecimalPair)
		}
		denoms[token.Denom] = 1
		if token.Amount.Equal(sdk.NewInt(0)) {
			return errorsmod.Wrapf(ErrFailedDeposit, "because of %s", ErrInvalidAmount)
		}
	}

	return nil
}
