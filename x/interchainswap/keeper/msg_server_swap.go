package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	pool, found := k.GetInterchainLiquidityPool(ctx, types.GetPoolId([]string{
		msg.TokenIn.Denom, msg.TokenOut.Denom,
	}))

	if !found {
		return nil, nil
	}

	//lock swap-in token to the swap module
	escrowAddr := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(msg.Sender), escrowAddr.String(), sdk.NewCoins(*msg.TokenIn))

	//constructs the IBC data packet
	rawMsgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	var msgType types.MessageType
	switch msg.SwapType {
	case types.SwapMsgType_LEFT:
		msgType = types.MessageType_LEFTSWAP
	case types.SwapMsgType_RIGHT:
		msgType = types.MessageType_RIGHTSWAP
	default:

	}

	packet := types.IBCSwapDataPacket{
		Type: msgType,
		Data: rawMsgData,
	}

	timeOutHeight, timeoutStamp := types.GetDefaultTimeOut()

	err = k.SendIBCSwapPacket(
		ctx,
		pool.EncounterPartyPort,
		pool.EncounterPartyChannel,
		timeOutHeight,
		timeoutStamp,
		packet,
	)
	if err != nil {
		return nil, err
	}
	return &types.MsgSwapResponse{
		SwapType: msg.SwapType,
		Tokens:   []*sdk.Coin{msg.TokenIn, msg.TokenOut},
	}, nil
}
