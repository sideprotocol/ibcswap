package keeper

import (
	"context"
	"encoding/json"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, nil
	}

	// Deposit token to Escrow account
	var coins sdk.Coins
	for _, denom := range msg.Tokens {
		balance := k.bankViewKeeper.GetBalance(ctx, sdk.AccAddress(msg.Sender), denom.Denom)
		if balance.Amount.Equal(math.NewInt(0)) {
			return nil, nil
		}
		coin := sdk.Coin{
			Denom:  denom.Denom,
			Amount: denom.Amount,
		}
		coins.Add(coin)
	}

	k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.Sender), types.ModuleName, coins)

	//Send packet
	rawMsgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_ = types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}
	//channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(msg., msg.SourceChannel))
	// if !ok {
	// 	return nil, nil
	// }
	// k.channelKeeper.SendPacket(ctx)
	return &types.MsgDepositResponse{}, nil
}
