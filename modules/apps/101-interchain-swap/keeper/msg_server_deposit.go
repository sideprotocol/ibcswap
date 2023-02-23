package keeper

import (
	"context"
	"encoding/json"

	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "because of %s", err)
	}

	// Deposit token to Escrow account
	var coins sdk.Coins
	for _, denom := range msg.Tokens {
		balance := k.bankViewKeeper.GetBalance(ctx, sdk.AccAddress(msg.Sender), denom.Denom)
		if balance.Amount.Equal(sdk.NewInt(0)) {
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

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)

	// construct ibc packet
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

	return &types.MsgDepositResponse{
		PoolToken: pool.Supply,
	}, nil
}
