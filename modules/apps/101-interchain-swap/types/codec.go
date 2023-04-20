package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePoolRequest{}, "interchainswap/CreatePool", nil)
	cdc.RegisterConcrete(&MsgDepositRequest{}, "interchainswap/Deposit", nil)
	cdc.RegisterConcrete(&MsgDoubleDepositRequest{}, "interchainswap/DoubleDeposit", nil)
	cdc.RegisterConcrete(&RemoteDeposit{}, "interchainswap/RemoteDeposit", nil)
	cdc.RegisterConcrete(&LocalDeposit{}, "interchainswap/LocalDeposit", nil)
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
		&MsgDoubleDepositRequest{},
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

func DeserializeDepositTx(cdc codec.BinaryCodec, data []byte) (*banktypes.MsgSend, error) {
	// only ProtoCodec is supported
	if _, ok := cdc.(*codec.ProtoCodec); !ok {
		return nil, ErrInvalidMsg
	}

	var bankMsg banktypes.MsgSend
	if err := cdc.Unmarshal(data, &bankMsg); err != nil {
		return nil, err
	}
	return &bankMsg, nil
}
