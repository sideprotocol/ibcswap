package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) CreatePool(ctx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate message
	err := host.PortIdentifierValidator(msg.SourcePort)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedSwap, "due to %s", err)
	}

	err = host.ChannelIdentifierValidator(msg.SourceChannel)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedSwap, "due to %s", err)
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedCreatePool, "due to %s", err)
	}

	if !k.bankKeeper.HasSupply(sdkCtx, msg.Liquidity[0].Balance.Denom) {
		return nil, errormod.Wrapf(types.ErrFailedCreatePool, "due to %s", types.ErrInvalidLiquidity)
	}

	// Check if user owns initial liquidity or not
	senderAddress := sdk.MustAccAddressFromBech32(msg.Creator)

	sourceLiquidity := k.bankKeeper.GetBalance(sdkCtx, senderAddress, msg.Liquidity[0].Balance.Denom)
	if sourceLiquidity.Amount.LT(msg.Liquidity[0].Balance.Amount) {
		return nil, types.ErrEmptyInitialLiquidity
	}

	// Move initial funds to liquidity pool
	err = k.LockTokens(sdkCtx, msg.SourcePort, msg.SourceChannel, senderAddress, sdk.NewCoins(*msg.Liquidity[0].Balance))

	if err != nil {
		return nil, err
	}

	poolData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// Construct IBC data packet
	packet := types.IBCSwapPacketData{
		Type: types.CREATE_POOL,
		Data: poolData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)

	// use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	err = k.SendIBCSwapPacket(sdkCtx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}
	poolId := types.GetPoolId(msg.GetLiquidityDenoms())
	return &types.MsgCreatePoolResponse{
		PoolId: poolId,
	}, nil
}
