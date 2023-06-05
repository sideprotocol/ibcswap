package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreatePoolRequest{}, "interchainswap/CreatePool", nil)
	cdc.RegisterConcrete(&MsgSingleAssetDepositRequest{}, "interchainswap/Deposit", nil)
	cdc.RegisterConcrete(&MsgMultiAssetDepositRequest{}, "interchainswap/DoubleDeposit", nil)
	cdc.RegisterConcrete(&DepositSignature{}, "interchainswap/DepositSignature", nil)
	cdc.RegisterConcrete(&MsgSingleAssetWithdrawRequest{}, "interchainswap/SingleWithdraw", nil)
	cdc.RegisterConcrete(&MsgMultiAssetWithdrawRequest{}, "interchainswap/MultiWithdraw", nil)
	cdc.RegisterConcrete(&MsgSwapRequest{}, "interchainswap/Swap", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreatePoolRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSingleAssetDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgMultiAssetDepositRequest{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgSingleAssetWithdrawRequest{},
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
