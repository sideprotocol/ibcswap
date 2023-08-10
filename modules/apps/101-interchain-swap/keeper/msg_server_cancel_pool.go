package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) CancelPool(ctx context.Context, msg *types.MsgCancelPoolRequest) (*types.MsgCancelPoolResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_, connected := k.GetCounterPartyChainID(sdkCtx, msg.SourcePort, msg.SourceChannel)
	if !connected {
		return nil, errormod.Wrapf(types.ErrFailedMakePool, "%s", types.ErrConnection)
	}

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errormod.Wrapf(types.ErrFailedMakePool, ":%s", types.ErrAlreadyExistPool)
	}

	if pool.Status == types.PoolStatus_ACTIVE {
		return nil, errormod.Wrapf(types.ErrFailedDeposit, ":%s", "pool is in active")
	}

	// Validate message
	err := host.PortIdentifierValidator(msg.SourcePort)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedSwap, "due to %s", err)
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, errormod.Wrapf(types.ErrCancelPool, "due to %s", err)
	}

	// Move initial funds to liquidity pool

	cancelPoolData := types.ModuleCdc.MustMarshalJSON(msg)
	rawStateChange := types.ModuleCdc.MustMarshalJSON(&types.StateChange{
		PoolId:        msg.PoolId,
		SourceChainId: sdkCtx.ChainID(),
	})

	// Construct IBC data packet
	packet := types.IBCSwapPacketData{
		Type:        types.CANCEL_POOL,
		Data:        cancelPoolData,
		StateChange: rawStateChange,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)

	// use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	_, err = k.SendIBCSwapPacket(sdkCtx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	// emit events
	k.EmitEvent(
		sdkCtx, types.EventValueActionCancelPool, msg.PoolId, msg.Creator,
	)
	return &types.MsgCancelPoolResponse{
		PoolId: msg.PoolId,
	}, nil
}
