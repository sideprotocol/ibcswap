package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// IBC transfer sentinel errors
var (
	ErrInvalidPacketTimeout        = sdkerrors.Register(ModuleName, 2, "invalid packet timeout")
	ErrInvalidDenomForTransfer     = sdkerrors.Register(ModuleName, 3, "invalid denomination for cross-chain transfer")
	ErrInvalidVersion              = sdkerrors.Register(ModuleName, 4, "invalid ICS31 version")
	ErrInvalidAmount               = sdkerrors.Register(ModuleName, 5, "invalid token amount")
	ErrTraceNotFound               = sdkerrors.Register(ModuleName, 6, "denomination trace not found")
	ErrSendDisabled                = sdkerrors.Register(ModuleName, 7, "swap from this chain are disabled")
	ErrReceiveDisabled             = sdkerrors.Register(ModuleName, 8, "swap to this chain are disabled")
	ErrMaxTransferChannels         = sdkerrors.Register(ModuleName, 9, "max transfer channels")
	ErrInvalidCodec                = sdkerrors.Register(ModuleName, 10, "codec is not supported")
	ErrUnknownDataPacket           = sdkerrors.Register(ModuleName, 11, "data packet is not supported")
	ErrOrderDoesNotExists          = sdkerrors.Register(ModuleName, 12, "Make Order does not exist")
	ErrOrderCanceled               = sdkerrors.Register(ModuleName, 13, "Order has been canceled")
	ErrOrderCompleted              = sdkerrors.Register(ModuleName, 14, "Order has completed already")
	ErrOrderDenominationMismatched = sdkerrors.Register(ModuleName, 15, "denomination are not matched")
	ErrOrderInsufficientAmount     = sdkerrors.Register(ModuleName, 16, "amount of taker token is insufficient")
	ErrOrderPermissionIsNotAllowed = sdkerrors.Register(ModuleName, 17, "sender is not the owner of the order")
	ErrNotFoundChannel             = sdkerrors.Register(ModuleName, 18, "did not find channel")
	ErrFailedMakeSwap              = sdkerrors.Register(ModuleName, 19, "Failed to make swap")
	ErrInvalidTakerReceiverAddress = sdkerrors.Register(ModuleName, 20, "Taker Address has counter party chain address")
	ErrInvalidOrderStatus          = sdkerrors.Register(ModuleName, 21, "invalid order status")
	ErrInvalidSellToken            = sdkerrors.Register(ModuleName, 22, "invalid sell token")
	ErrInvalidTakerAddress         = sdkerrors.Register(ModuleName, 23, "invalid taker address")
	ErrAlreadyOrderTook            = sdkerrors.Register(ModuleName, 24, "already order took")
	ErrNotFoundOrder               = sdkerrors.Register(ModuleName, 25, "did not find order")
)
