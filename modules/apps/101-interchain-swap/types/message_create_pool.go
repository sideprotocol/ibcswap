package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgCreatePool = "create_pool"

var _ sdk.Msg = &MsgCreatePoolRequest{}

func NewMsgCreatePool(
	sourcePort string,
	sourceChannel string,
	creator string,
	counterPartyCreator string,
	counterPartySig []byte,
	sourceLiquidity,
	targetLiquidity PoolAsset,
) *MsgCreatePoolRequest {

	return &MsgCreatePoolRequest{
		SourcePort:          sourcePort,
		SourceChannel:       sourceChannel,
		Creator:             creator,
		CounterPartyCreator: counterPartyCreator,
		Liquidity: []*PoolAsset{
			&sourceLiquidity,
			&targetLiquidity,
		},
		CounterPartySig: counterPartySig,
	}
}

func (msg *MsgCreatePoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgCreatePoolRequest) Type() string {
	return TypeMsgCreatePool
}

func (msg *MsgCreatePoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePoolRequest) GetLiquidityDenoms() []string {
	denoms := []string{}
	for _, asset := range msg.Liquidity {
		denoms = append(denoms, asset.Balance.Denom)
	}
	return denoms
}

func (msg *MsgCreatePoolRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	_, err = sdk.AccAddressFromBech32(msg.CounterPartyCreator)
	if err != nil {
		return ErrInvalidAddress
	}

	tokenCount := len(msg.Liquidity)
	// Validation message
	if tokenCount != 2 {
		return ErrInvalidDenomPair
	}

	if err := ValidateLiquidityBasic(msg.Liquidity); err != nil {
		return err
	}
	return nil
}
