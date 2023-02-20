package keeper

import (
	"context"
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	//abortTransactionUnless(host.PortIdentifierValidator(msg.SourcePort))

	pool := types.NewInterchainLiquidityPool(
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
	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_CREATE,
		Data: poolData,
	}

	rawPacket, err := json.Marshal(packet)
	if err != nil {
		return nil, err
	}

	//
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(msg.SourcePort, msg.SourceChannel))
	if !ok {
		return nil, nil
	}

	timeoutHeight := clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: 10,
	}
	timeoutStamp := time.Now().UTC().Unix()

	_, err = k.channelKeeper.SendPacket(ctx, channelCap, msg.SourcePort, msg.SourceChannel, timeoutHeight, uint64(timeoutStamp), rawPacket)
	if err != nil {
		return nil, err
	}
	return &types.MsgCreatePoolResponse{}, nil
}
