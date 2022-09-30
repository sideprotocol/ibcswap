package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC Interchain sentinel errors
var (
	ErrInvalidPacketTimeout    = sdkerrors.Register(ModuleName, 2, "invalid packet timeout")
	ErrInvalidDenomForTransfer = sdkerrors.Register(ModuleName, 3, "invalid denomination for cross-chain transfer")
	ErrInvalidVersion          = sdkerrors.Register(ModuleName, 4, "invalid ICS101 version")
	ErrInvalidAmount           = sdkerrors.Register(ModuleName, 5, "invalid token amount")
	ErrSendDisabled            = sdkerrors.Register(ModuleName, 7, "swap from this chain are disabled")
	ErrReceiveDisabled         = sdkerrors.Register(ModuleName, 8, "swap to this chain are disabled")
	ErrMaxTransferChannels     = sdkerrors.Register(ModuleName, 9, "max transfer channels")
	ErrInvalidCodec            = sdkerrors.Register(ModuleName, 10, "codec is not supported")
	ErrUnknownDataPacket       = sdkerrors.Register(ModuleName, 11, "data packet is not supported")
	ErrInvalidPairLength       = sdkerrors.Register(ModuleName, 18, "invalid pair length")
	ErrInvalidWeightOfPool     = sdkerrors.Register(ModuleName, 19, "invalid weights of pool")
	ErrInvalidPoolId           = sdkerrors.Register(ModuleName, 20, "invalid pool id")
	ErrInvalidToken            = sdkerrors.Register(ModuleName, 21, "invalid token")
	ErrNoNativeTokenInPool     = sdkerrors.Register(ModuleName, 22, "at least 1 native token is required")
	ErrTokenNotInPool          = sdkerrors.Register(ModuleName, 23, "not found token in pool")
	ErrAmountInsufficient      = sdkerrors.Register(ModuleName, 24, "amount insufficient")
)
