package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePoolRequest{}, "interchainswap/CreatePool", nil)
	cdc.RegisterConcrete(&MsgDepositRequest{}, "interchainswap/Deposit", nil)
	cdc.RegisterConcrete(&MsgWithdrawRequest{}, "interchainswap/Withdraw", nil)
	cdc.RegisterConcrete(&MsgSwapRequest{}, "interchainswap/Swap", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePoolRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgWithdrawRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSwapRequest{},
	)
	// this line is used by starport scaffolding # 3

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
