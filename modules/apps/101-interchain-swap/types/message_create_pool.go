package types

import (
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePoolRequest{}

func NewMsgCreatePool(sourcePort string, sourceChannel string, sender string, weight string, tokens []*sdk.Coin, decimals []uint32) *MsgCreatePoolRequest {
	return &MsgCreatePoolRequest{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		Sender:        sender,
		Weight:        weight,
		Tokens:        tokens,
		Decimals:      decimals,
	}
}

func (msg *MsgCreatePoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgCreatePoolRequest) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePoolRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return ErrInvalidAddress
	}

	denomSize := len(msg.Tokens)
	//validation message
	if denomSize != 2 {
		return ErrInvalidDenomPair
	}

	if len(msg.Decimals) != 2 {
		return ErrInvalidDecimalPair
	}

	weights := strings.Split(msg.Weight, ":")
	if len(weights) != 2 {
		return ErrInvalidWeightPair
	}

	totalWeight := 0
	for _, weight := range weights {
		w,_ := strconv.Atoi(weight)
		totalWeight += w
	}

	if totalWeight != 100 {
		return ErrInvalidWeightPair
	}

	for _, token := range msg.Tokens {
		if token.Amount.LTE(sdk.NewInt(0)) {
			return ErrInvalidAmount
		}
	}
	return nil
}
