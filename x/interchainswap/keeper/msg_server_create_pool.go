package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message

	//abortTransactionUnless(host.PortIdentifierValidator(msg.SourcePort))

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
