package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/interchainswap module sentinel errors
var (
	ErrSample               = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidRequest       = errorsmod.Register(ModuleName, 1101, "invalid request!")
	ErrInvalidType          = errorsmod.Register(ModuleName, 1102, "invalid packet type!")
	ErrUnknownRequest       = errorsmod.Register(ModuleName, 1103, "unknown request!")
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrNotFoundPool         = errorsmod.Register(ModuleName, 1502, "did not find pool!")
	ErrInvalidAmount        = errorsmod.Register(ModuleName, 1503, "Invalid token amount!")
	ErrInvalidLength        = errorsmod.Register(ModuleName, 1504, "Invalid token length!")
	ErrUnknownDataPacket    = errorsmod.Register(ModuleName, 1504, "unknown packet!")
	ErrInvalidAddress       = errorsmod.Register(ModuleName, 1505, "invalid address!")
	//relevant for pool
	ErrNotFoundDenomInPool = errorsmod.Register(ModuleName, 1507, "not find denom in the pool!")
	ErrInvalidDenomPair    = errorsmod.Register(ModuleName, 1508, "invalid denom pair!")
	ErrInvalidDecimalPair  = errorsmod.Register(ModuleName, 1509, "invalid decimal pair!")
	ErrInvalidWeightPair   = errorsmod.Register(ModuleName, 1510, "invalid weight pair!")
	ErrEmptyDenom          = errorsmod.Register(ModuleName, 1511, "dropped denom!")
	ErrInvalidSlippage     = errorsmod.Register(ModuleName, 1511, "invalid slippage!")
)
