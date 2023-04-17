package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDoubleDeposit = "deposit"

var _ sdk.Msg = &MsgDepositRequest{}

func NewMsgDoubleDeposit(poolId string, senders []string, tokens []*sdk.Coin, sig []byte) *MsgDoubleDepositRequest {
	return &MsgDoubleDepositRequest{
		PoolId:                  poolId,
		Senders:                 senders,
		Tokens:                  tokens,
		EncounterPartySignature: sig,
	}
}

func (msg *MsgDoubleDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgDoubleDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgDoubleDepositRequest) GetSigners() []sdk.AccAddress {
	signers := []sdk.AccAddress{}
	creator, err := sdk.AccAddressFromBech32(msg.Senders[0])
	if err != nil {
		panic(err)
	}
	signers = append(signers, creator)
	return signers
}

func (msg *MsgDoubleDepositRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDoubleDepositRequest) ValidateBasic() error {

	for _, sender := range msg.Senders {
		_, err := sdk.AccAddressFromBech32(sender)
		// senderPrefix, _, err := bech32.Decode(sender)
		// if sdk.GetConfig().GetBech32AccountAddrPrefix() != senderPrefix && index == 0 {
		// 	return errorsmod.Wrapf(ErrInvalidAddress, "first address has to be this chain address (%s)", err)
		// }
		if err != nil {
			return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address (%s)", err)
		}
	}

	if len(msg.Senders) != 2 {
		return errorsmod.Wrapf(ErrInvalidTokenLength, "invalid token length (%d)", len(msg.Tokens))
	}

	if len(msg.Tokens) < 2 {
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
