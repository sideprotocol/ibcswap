package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) SingleDeposit(goCtx context.Context, msg *types.MsgSingleDepositRequest) (*types.MsgSingleDepositResponse, error) {
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
	accAddress, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", types.ErrInvalidAddress)
	}
	balance := k.bankKeeper.GetBalance(ctx, accAddress, msg.Token.Denom)
	if balance.Amount.Equal(sdk.NewInt(0)) {
		return nil, types.ErrInvalidAmount
	}

	if pool.Status == types.PoolStatus_POOL_STATUS_INITIAL {
		poolAsset, err := pool.FindAssetByDenom(msg.Token.Denom)
		if err == nil {
			if !poolAsset.Balance.Amount.Equal(msg.Token.Amount) {
				return nil, types.ErrInvalidInitialDeposit
			}
		}
	}

	// deposit assets to the escrowed account
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*msg.Token))
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", err)
	}

	
	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	poolToken, err := amm.DepositSingleAsset(*msg.Token)
	if err != nil {
		return nil, err
	}

	// construct ibc packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.SINGLE_DEPOSIT,
		Data:        rawMsgData,
		StateChange: &types.StateChange{PoolTokens: []*sdk.Coin{poolToken}},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgSingleDepositResponse{
		PoolToken: pool.Supply,
	}, nil
}
