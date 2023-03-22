package types

// DONTCOVER

import (
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/interchainswap module sentinel errors
var (
	ErrSample               = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrInvalidRequest       = errorsmod.Register(ModuleName, 1101, "invalid request")
	ErrInvalidType          = errorsmod.Register(ModuleName, 1102, "invalid packet type")
	ErrUnknownRequest       = errorsmod.Register(ModuleName, 1103, "unknown request")
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 1500, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 1501, "invalid version")
	ErrNotFoundPool         = errorsmod.Register(ModuleName, 1502, "did not find pool")
	ErrInvalidAmount        = errorsmod.Register(ModuleName, 1503, "Invalid token amount")
	ErrInvalidTokenLength   = errorsmod.Register(ModuleName, 1504, "Invalid token length")
	ErrUnknownDataPacket    = errorsmod.Register(ModuleName, 1505, "unknown packet")
	ErrInvalidAddress       = errorsmod.Register(ModuleName, 1506, "invalid address")
	ErrInvalidPort          = errorsmod.Register(ModuleName, 1507, "invalid port")
	ErrInvalidChannel       = errorsmod.Register(ModuleName, 1508, "invalid channel")
	//relevant for pool
	ErrNotFoundDenomInPool = errorsmod.Register(ModuleName, 1509, "not find denom in the pool")
	ErrInvalidDenomPair    = errorsmod.Register(ModuleName, 1510, "invalid denom pair")
	ErrInvalidDecimalPair  = errorsmod.Register(ModuleName, 1511, "invalid decimal pair")
	ErrInvalidWeightPair   = errorsmod.Register(ModuleName, 1512, "invalid weight pair")
	ErrEmptyDenom          = errorsmod.Register(ModuleName, 1513, "dropped denom")
	ErrInvalidSlippage     = errorsmod.Register(ModuleName, 1514, "invalid slippage")
	ErrNotReadyForSwap     = errorsmod.Register(ModuleName, 1515, "pool is not ready for swap")
	ErrNumberOfLocalAsset  = errorsmod.Register(ModuleName, 1516, "should have 1 native asset on the chain")
	ErrNotNativeDenom      = errorsmod.Register(ModuleName, 1517, "invalid native denom")

	//msg srv errors
	ErrFailedCreatePool = errorsmod.Register(ModuleName, 1518, "failed to create pool")
	ErrFailedDeposit    = errorsmod.Register(ModuleName, 1519, "failed to deposit")
	ErrFailedWithdraw   = errorsmod.Register(ModuleName, 1520, "failed to withdraw")
	ErrFailedSwap       = errorsmod.Register(ModuleName, 1521, "failed to interchain swap")

	ErrFailedOnCreatePoolReceived = errorsmod.Register(ModuleName, 1522, "failed to treat create pool msg!")
	ErrFailedOnDepositReceived    = errorsmod.Register(ModuleName, 1523, "failed to treat deposit msg!")
	ErrFailedOnWithdrawReceived   = errorsmod.Register(ModuleName, 1524, "failed to treat withdraw msg!")
	ErrFailedOnSwapReceived       = errorsmod.Register(ModuleName, 1525, "failed to treat swap msg!")

	ErrFailedOnCreatePoolAck = errorsmod.Register(ModuleName, 1526, "failed to treat create pool ack!")
	ErrFailedOnDepositAck    = errorsmod.Register(ModuleName, 1527, "failed to treat deposit ack!")
	ErrFailedOnWithdrawAck   = errorsmod.Register(ModuleName, 1528, "failed to treat withdraw ack!")
	ErrFailedOnSwapAck       = errorsmod.Register(ModuleName, 1529, "failed to treat swap ack!")

	ErrMaxTransferChannels            = errorsmod.Register(ModuleName, 1530, "max transfer channels")
	ErrAlreadyExistPool               = errorsmod.Register(ModuleName, 1531, "already exist pool!")
	ErrSwapEnabled                    = errorsmod.Register(ModuleName, 1532, "swap is disabled!")
	ErrEmptyInitialLiquidity          = errorsmod.Register(ModuleName, 1533, "creator don't have liquidity!")
	ErrIncorrectInitialLiquidityDenom = errorsmod.Register(ModuleName, 1534, "initial liquidity denom has to match with one of pool denom pair")
	ErrInvalidSwapType                = errorsmod.Register(ModuleName, 1535, "invalid swap type!")

	ErrInvalidDenom = errorsmod.Register(ModuleName, 1536, "invalid denom")
)
