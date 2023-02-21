package keeper

import (
	"context"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)

	if !found {
		return nil, nil
	}

	_, found = k.GetInterchainLiquidityPool(ctx, msg.DenomOut)

	if !found {
		return nil, nil
	}

	escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	k.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, string(escrowAccount), sdk.AccAddress(msg.Sender), sdk.NewCoins(*msg.PoolCoin))

	//Send packet
	rawMsgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_ = types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(pool.EncounterPartyPort, pool.EncounterPartyChannel))
	if !ok {
		return nil, nil
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut()

	k.channelKeeper.SendPacket(ctx, channelCap, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), rawMsgData)
	return &types.MsgWithdrawResponse{}, nil
}
