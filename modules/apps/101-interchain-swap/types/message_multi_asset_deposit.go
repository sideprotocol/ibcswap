package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgDoubleDeposit = "deposit"

var _ sdk.Msg = &MsgSingleAssetDepositRequest{}

func NewMsgMultiAssetDeposit(poolId string, senders []string, tokens []*sdk.Coin, sig []byte) *MsgMultiAssetDepositRequest {
	return &MsgMultiAssetDepositRequest{
		PoolId: poolId,
		Deposits: []*DepositAsset{
			{
				Sender:  senders[0],
				Balance: tokens[0],
			},
			{
				Sender:    senders[1],
				Balance:   tokens[1],
				Signature: sig,
			},
		},
	}
}

func (msg *MsgMultiAssetDepositRequest) Route() string {
	return RouterKey
}

func (msg *MsgMultiAssetDepositRequest) Type() string {
	return TypeMsgDeposit
}

func (msg *MsgMultiAssetDepositRequest) GetSigners() []sdk.AccAddress {
	signers := []sdk.AccAddress{}

	creator, err := sdk.AccAddressFromBech32(msg.Deposits[0].Sender)
	if err != nil {
		panic(err)
	}
	signers = append(signers, creator)
	return signers
}

func (msg *MsgMultiAssetDepositRequest) GetSignBytes() []byte {
	marshaledMsg := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(marshaledMsg)
}

func (msg *MsgMultiAssetDepositRequest) ValidateBasic() error {
	// Check address
	for _, deposit := range msg.Deposits {
		_, err := sdk.AccAddressFromBech32(deposit.Sender)
		if err != nil {
			return ErrInvalidAddress
		}
		if deposit.Balance.Amount.Equal(sdk.NewInt(0)) {
			return ErrInvalidAmount
		}
	}
	return nil
}
