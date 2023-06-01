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
		return nil, errormod.Wrapf(types.ErrFailedSwap, "due to %s", err)
	}

	localAssetCount := 0
	for _, token := range msg.Tokens {
		if k.bankKeeper.HasSupply(sdkCtx, token.Denom) {
			localAssetCount += 1
		}
	}

	// Should have 1 native asset on the chain
	if localAssetCount < 1 {
		return nil, types.ErrNumberOfLocalAsset
	}

	// Check if user owns initial liquidity or not
	senderAddress := sdk.MustAccAddressFromBech32(msg.Sender)
	holdingNativeCoin := k.bankKeeper.GetBalance(sdkCtx, senderAddress, msg.Tokens[0].Denom)
	if holdingNativeCoin.Amount.LT(msg.Tokens[0].Amount) {
		return nil, types.ErrEmptyInitialLiquidity
	}

	// Move initial funds to liquidity pool
	err = k.LockTokens(sdkCtx, msg.SourcePort, msg.SourceChannel, senderAddress, sdk.NewCoins(*msg.Tokens[0]))

	if err != nil {
		return nil, err
	}
	msg.ChainId = sdkCtx.ChainID()
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
	poolId := types.GetPoolIdWithTokens(msg.Tokens)

	return &types.MsgCreatePoolResponse{
		PoolId: poolId,
	}, nil
}
