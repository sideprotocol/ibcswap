package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers the necessary x/ibc interchain swap interfaces and concrete types
// on the provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePoolRequest{}, "cosmos-sdk/MsgCreatePool", nil)
	cdc.RegisterConcrete(&MsgSingleDepositRequest{}, "cosmos-sdk/MsgSingleDeposit", nil)
	cdc.RegisterConcrete(&MsgWithdrawRequest{}, "cosmos-sdk/MsgWithdraw", nil)
	cdc.RegisterConcrete(&MsgLeftSwapRequest{}, "cosmos-sdk/MsgLeftSwap", nil)
	cdc.RegisterConcrete(&MsgRightSwapRequest{}, "cosmos-sdk/MsgRight", nil)
}

// RegisterInterfaces register the ibc interchain swap module interfaces to protobuf
// Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePoolRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSingleDepositRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgLeftSwapRequest{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgRightSwapRequest{})

	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgCreatePoolResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSingleDepositResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgWithdrawResponse{})
	registry.RegisterImplementations((*sdk.Msg)(nil), &MsgSwapResponse{})

	//msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/ibc interchain swap module codec. Note, the codec
	// should ONLY be used in certain instances of tests and for JSON encoding.
	//
	// The actual codec used for serialization should be provided to x/ibc transfer and
	// defined at the application level.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc is a amino codec created to support amino json compatible msgs.
	AminoCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}
