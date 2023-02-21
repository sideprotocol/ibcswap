package keeper

import (
	"context"
	"encoding/json"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, fmt.Errorf("tokenId: %s  %s", &pool.PoolId, types.ErrNotFoundPool)
	}

	// Deposit token to Escrow account
	var coins sdk.Coins
	for _, denom := range msg.Tokens {
		balance := k.bankViewKeeper.GetBalance(ctx, sdk.AccAddress(msg.Sender), denom.Denom)
		if balance.Amount.Equal(math.NewInt(0)) {
			return nil, types.ErrInvalidAmount
		}
		coin := sdk.Coin{
			Denom:  denom.Denom,
			Amount: denom.Amount,
		}
		coins.Add(coin)
	}

	escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, escrowAccount, types.ModuleName, coins)

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut()

	//Send packet
	rawMsgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}
	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{}, nil
}
