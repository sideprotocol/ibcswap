package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) MakePool(ctx context.Context, msg *types.MsgMakePoolRequest) (*types.MsgMakePoolResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	counterPartyChainId, connected := k.GetCounterPartyChainID(sdkCtx, msg.SourcePort, msg.SourceChannel)
	if !connected {
		return nil, errormod.Wrapf(types.ErrFailedMakePool, "%s", types.ErrConnection)
	}

	denoms := []string{}
	for _, liquidity := range msg.Liquidity {
		denoms = append(denoms, liquidity.Balance.Denom)
	}
	poolId := types.GetPoolId(sdkCtx.ChainID(), counterPartyChainId, denoms)
	_, found := k.GetInterchainLiquidityPool(sdkCtx, poolId)

	if found {
		return nil, errormod.Wrapf(types.ErrFailedMakePool, "%s", types.ErrAlreadyExistPool)
	}

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
		return nil, errormod.Wrapf(types.ErrFailedMakePool, "due to %s", err)
	}

	if !k.bankKeeper.HasSupply(sdkCtx, msg.Liquidity[0].Balance.Denom) {
		return nil, errormod.Wrapf(types.ErrFailedMakePool, "due to %s", types.ErrInvalidLiquidity)
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
		Type: types.MAKE_POOL,
		Data: poolData,
		StateChange: &types.StateChange{
			PoolId:        poolId,
			SourceChainId: sdkCtx.ChainID(),
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)

	// use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	_,err = k.SendIBCSwapPacket(sdkCtx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}
	return &types.MsgMakePoolResponse{
		PoolId: poolId,
	}, nil
}
