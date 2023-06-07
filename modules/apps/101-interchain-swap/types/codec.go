package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgMakePoolRequest{}, "interchainswap/MakePool", nil)
	cdc.RegisterConcrete(&MsgTakePoolRequest{}, "interchainswap/TakePool", nil)
	cdc.RegisterConcrete(&MsgSingleAssetDepositRequest{}, "interchainswap/Deposit", nil)
	cdc.RegisterConcrete(&MsgMakeMultiAssetDepositRequest{}, "interchainswap/MakeMultiAssetDeposit", nil)
	cdc.RegisterConcrete(&MsgTakeMultiAssetDepositRequest{}, "interchainswap/TakeMultiAssetDeposit", nil)
	cdc.RegisterConcrete(&MsgMultiAssetWithdrawRequest{}, "interchainswap/MultiWithdraw", nil)
	cdc.RegisterConcrete(&MsgSwapRequest{}, "interchainswap/Swap", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMakePoolRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTakePoolRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSingleAssetDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMakeMultiAssetDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgTakeMultiAssetDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMultiAssetWithdrawRequest{},
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
