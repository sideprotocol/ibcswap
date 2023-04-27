package keeper

import (
	"context"

	"github.com/btcsuite/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) DoubleDeposit(goCtx context.Context, msg *types.MsgDoubleDepositRequest) (*types.MsgDoubleDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// // validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// check address
	senderPrefix, _, err := bech32.Decode(msg.LocalDeposit.Sender)
	if err != nil {
		return nil, err
	}
	if sdk.GetConfig().GetBech32AccountAddrPrefix() != senderPrefix {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "first address has to be this chain address (%s)", err)
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	}
	_ = pool

	// Deposit token to Escrow account
	coins, err := k.validateDoubleDepositCoins(ctx, &pool, msg.LocalDeposit.Sender, msg.LocalDeposit.Token)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", err)
	}
	if len(coins) == 0 {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "it include invalid coins (%s)")
	}

	// create escrow module account  here
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.LocalDeposit.Sender), coins)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "because of %s", err)
	}

	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	poolTokens, err := amm.DepositDoubleAsset([]*sdk.Coin{
		msg.LocalDeposit.Token,
		msg.RemoteDeposit.Token,
	})

	if err != nil {
		return nil, err
	}

	// construct ibc packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.DOUBLE_DEPOSIT,
		Data:        rawMsgData,
		StateChange: &types.StateChange{PoolTokens: poolTokens},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgDoubleDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}
