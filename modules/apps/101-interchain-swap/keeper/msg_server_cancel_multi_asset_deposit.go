package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) CancelMultiAssetDeposit(ctx context.Context, msg *types.MsgCancelMultiAssetDepositRequest) (*types.MsgCancelMultiAssetDepositResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_ = sdkCtx
	// // Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	_, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	order, found := k.GetMultiDepositOrder(sdkCtx, msg.PoolId, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFoundMultiDepositOrder, ":%s", types.ErrCancelOrder)
	}

	if order.Status == types.OrderStatus_COMPLETE {
		return nil, errorsmod.Wrapf(types.ErrAlreadyCompletedOrder, ":%s", types.ErrCancelOrder)
	}
	if msg.Creator != order.SourceMaker {
		return nil, errorsmod.Wrapf(types.ErrNotEnoughPermission, ":%s", types.ErrCancelOrder)
	}

	
	cancelOrderData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}
	// save order in source chain
	packet := types.IBCSwapPacketData{
		Type: types.CANCEL_MULTI_DEPOSIT,
		Data: cancelOrderData,
		StateChange: &types.StateChange{
			MultiDepositOrderId: order.Id,
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)
	// Use input timeoutHeight, timeoutStamp
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

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeCancelPool,
			sdk.Attribute{
				Key:   types.AttributeKeyPoolId,
				Value: msg.PoolId,
			},
			sdk.Attribute{
				Key:   types.AttributeKeyMultiDepositOrderId,
				Value: msg.OrderId,
			},
		))

	return &types.MsgCancelMultiAssetDepositResponse{
		PoolId:  msg.PoolId,
		OrderId: msg.OrderId,
	}, nil
}
