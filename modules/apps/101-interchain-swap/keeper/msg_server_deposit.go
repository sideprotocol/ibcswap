package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (k Keeper) Deposit(goCtx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", types.ErrNotFoundPool)
	}

	// Deposit token to Escrow account
	coins := []sdk.Coin{}
	for _, coin := range msg.Tokens {
		accAddress, err := sdk.AccAddressFromBech32(msg.Sender)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", types.ErrInvalidAddress)
		}
		balance := k.bankKeeper.GetBalance(ctx, accAddress, coin.Denom)
		if balance.Amount.Equal(sdk.NewInt(0)) {
			return nil, types.ErrInvalidAmount
		}
		coins = append(coins, *coin)
	}

	if len(coins) == 0 {
		return nil, types.ErrInvalidTokenLength
	}

	// create escrow module account  here
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), coins)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", err)
	}

	// construct ibc packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgDepositResponse{
		PoolToken: pool.Supply,
	}, nil
}
