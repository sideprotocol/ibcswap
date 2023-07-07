package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const TypeMsgMakePool = "make_pool"

var _ sdk.Msg = &MsgMakePoolRequest{}

func NewMsgMakePool(
	sourcePort string,
	sourceChannel string,
	creator string,
	counterPartyCreator string,
	sourceLiquidity,
	targetLiquidity PoolAsset,
	swapFee uint32,
) *MsgMakePoolRequest {

	return &MsgMakePoolRequest{
		SourcePort:          sourcePort,
		SourceChannel:       sourceChannel,
		Creator:             creator,
		CounterPartyCreator: counterPartyCreator,
		Liquidity: []*PoolAsset{
			&sourceLiquidity,
			&targetLiquidity,
		},
		SwapFee: swapFee,
	}
}

func (msg *MsgMakePoolRequest) Route() string {
	return RouterKey
}

func (msg *MsgMakePoolRequest) Type() string {
	return TypeMsgMakePool
}

func (msg *MsgMakePoolRequest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgMakePoolRequest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgMakePoolRequest) GetLiquidityDenoms() []string {
	denoms := []string{}
	for _, asset := range msg.Liquidity {
		denoms = append(denoms, asset.Balance.Denom)
	}
	return denoms
}

func (msg *MsgMakePoolRequest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return ErrInvalidAddress
	}
	if err := ValidateLiquidityBasic(msg.Liquidity); err != nil {
		return err
	}
	if msg.SwapFee < 0 || msg.SwapFee > 10000 {
		return ErrInvalidSwapFee
	}
	if msg.SourceChannel == "" {
		return ErrMissedIBCParams
	}
	if msg.SourcePort == "" {
		msg.SourcePort = PortID
	}
	return nil
}
