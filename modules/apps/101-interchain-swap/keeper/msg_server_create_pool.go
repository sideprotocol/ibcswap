package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate message
	err := host.PortIdentifierValidator(msg.SourcePort)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	err = host.ChannelIdentifierValidator(msg.SourceChannel)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	pool := types.NewInterchainLiquidityPool(
		ctx,
		k.bankKeeper,
		msg.Denoms,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)

	poolData, err := json.Marshal(pool)
	if err != nil {
		return nil, err
	}

	localAssetCount := 0
	for _, denom := range msg.Denoms {
		if k.bankKeeper.HasSupply(ctx, denom) {
			localAssetCount += 1
		}
	}

	// should have 1 native asset on the chain
	if localAssetCount < 1 {
		return nil, types.ErrNumberOfLocalAsset
	}

	// construct IBC data packet
	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_CREATE,
		Data: poolData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)

	if err != nil {
		return nil, err
	}
	return &types.MsgCreatePoolResponse{
		PoolId: pool.PoolId,
	}, nil
}
