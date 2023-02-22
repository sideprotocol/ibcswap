package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate message
	err := host.PortIdentifierValidator(msg.SourcePort)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to create pool because of %s")
	}

	err = host.ChannelIdentifierValidator(msg.SourceChannel)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to create pool because of %s")
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to create pool because of %s")
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

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut()

	err = k.SendIBCSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)

	if err != nil {
		return nil, err
	}
	return &types.MsgCreatePoolResponse{}, nil
}
