package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k Keeper) OnCreatePoolAcknowledged(ctx sdk.Context, msg *types.MsgCreatePoolRequest) {
	//save pool after complete create operation in counter party chain.

	pool := types.NewInterchainLiquidityPool(
		ctx,
		k.bankKeeper,
		msg.Denoms,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)
	k.SetInterchainLiquidityPool(ctx, *pool)
}

func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgDepositRequest, res *types.MsgDepositResponse) error {
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}
	pool.Supply.Amount.Add(res.PoolToken.Amount)
	// mint voucher
	err := k.MintTokens(ctx, sdk.AccAddress(req.Sender), *res.PoolToken)
	if err != nil {
		return err
	}
	k.SetInterchainLiquidityPool(ctx, pool)
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		"Deposit Ack3:",
		sdk.NewAttribute(
			"Deposit",
			pool.Supply.Amount.String(),
		),
	))
	return nil
}

func (k Keeper) OnWithdrawAcknowledged(ctx sdk.Context, req *types.MsgWithdrawRequest, res *types.MsgWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolCoin.Denom)
	if !found {
		return types.ErrNotFoundPool
	}
	pool.Supply.Amount.Sub(res.Tokens[0].Amount)

	k.SetInterchainLiquidityPool(ctx, pool)

	// burn voucher token.
	// err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *req.PoolCoin)
	// if err != nil {
	// 	return err
	// }
	// // unlock token
	// err = k.UnlockTokens(ctx,
	// 	pool.EncounterPartyPort,
	// 	pool.EncounterPartyChannel,
	// 	sdk.MustAccAddressFromBech32(req.Sender),
	// 	sdk.NewCoins(*res.Tokens[0]),
	// )
	// if err != nil {
	// 	return err
	// }
	return nil
}

func (k Keeper) OnSwapAcknowledged(ctx sdk.Context, req *types.MsgSwapRequest, res *types.MsgSwapResponse) error {
	pooId := types.GetPoolId([]string{req.TokenIn.Denom, req.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, pooId)
	if !found {
		return types.ErrNotFoundPool
	}

	assetOut, err := pool.FindAssetByDenom(res.Tokens[0].Denom)
	if err != nil {
		return err
	}
	if assetOut.Balance.Amount.LT(res.Tokens[0].Amount) {
		return types.ErrInvalidAmount
	}

	assetIn, err := pool.FindAssetByDenom(req.TokenIn.Denom)

	if err != nil {
		return err
	}

	assetIn.Balance.Amount.Add(req.TokenIn.Amount)
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnCreatePoolReceived(ctx sdk.Context, msg *types.MsgCreatePoolRequest, destPort, destChannel string) (*string, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pooId := types.GetPoolId(msg.Denoms)
	_, found := k.GetInterchainLiquidityPool(ctx, pooId)

	if found {
		return nil, types.ErrAlreadyExistPool
	}

	pool := *types.NewInterchainLiquidityPool(
		ctx,
		k.bankKeeper,
		msg.Denoms,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)
	//count native tokens
	count := 0
	for _, denom := range msg.Denoms {
		if k.bankKeeper.HasSupply(ctx, denom) {
			count += 1
			pool.UpdateAssetPoolSide(denom, types.PoolSide_NATIVE)
		} else {
			pool.UpdateAssetPoolSide(denom, types.PoolSide_REMOTE)
		}
	}

	if count != 1 {
		return nil, types.ErrInvalidDenomPair
	}

	// save pool status
	k.SetInterchainLiquidityPool(ctx, pool)
	// emit events
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"On Create",
			sdk.NewAttribute(
				"poolId",
				pooId,
			),
		),
	)
	return &pooId, nil
}

func (k Keeper) OnDepositReceived(ctx sdk.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	//TODO: Need to implement params module and market maker.
	amm := types.NewInterchainMarketMaker(
		pool,
		323,
	)

	poolToken, err := amm.DepositSingleAsset(*msg.Tokens[0])
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"mint boucher",
			sdk.NewAttribute(
				"PoolTokenDenom:",
				poolToken.Denom,
			),
			sdk.NewAttribute(
				"PoolTokenDenom:",
				poolToken.Amount.String(),
			),
		),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// is ready for swap.
	//amm.Pool.UpdateAssetPoolSide(msg.Tokens[0].Denom, types.PoolSide(types.PoolStatus_POOL_STATUS_READY)) //= types.PoolStatus_POOL_STATUS_READY

	k.SetInterchainLiquidityPool(ctx, *amm.Pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgDepositResponse{
		PoolToken: poolToken,
	}, nil
}

func (k Keeper) OndWithdrawReceive(ctx sdk.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)

	if !found {
		return nil, types.ErrNotFoundPool
	}

	//TODO: need to implement amm part.
	//feeRate := parms.getPoolFeeRate()

	amm := types.NewInterchainMarketMaker(
		pool,
		323,
	)

	outToken, err := amm.Withdraw(*msg.PoolCoin, msg.DenomOut)
	ctx.EventManager().EmitTypedEvents(outToken)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s!", err)
	}

	k.SetInterchainLiquidityPool(ctx, *amm.Pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	ctx.EventManager().EmitTypedEvents(msg)
	return &types.MsgWithdrawResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}

func (k Keeper) OnSwapReceived(ctx sdk.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pooId := types.GetPoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, pooId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	//TODO: need to implement amm part.
	//feeRate := parms.getPoolFeeRate()

	amm := types.NewInterchainMarketMaker(
		pool,
		323,
	)

	var outToken *sdk.Coin
	var err error
	switch msg.SwapType {
	case types.SwapMsgType_LEFT:
		outToken, err = amm.LeftSwap(*msg.TokenIn, msg.TokenOut.Denom)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrFailedOnSwapReceived, "because of %s", err)
		}
	case types.SwapMsgType_RIGHT:
		outToken, err = amm.RightSwap(*msg.TokenIn, *msg.TokenOut)
		if err != nil {
			return nil, errorsmod.Wrapf(types.ErrFailedOnSwapReceived, "because of %s", err)
		}
	}

	expected := float64(msg.TokenOut.Amount.Uint64()) * (1 - float64(msg.Slippage)/10000)

	if outToken.Amount.LT(sdk.NewInt(int64(expected))) {
		return nil, errorsmod.Wrap(types.ErrFailedOnSwapReceived, "doesn't meet slippage for swap!, %s")
	}

	escrowAddr := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, escrowAddr.String(), sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*outToken))

	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient!")
	}

	k.SetInterchainLiquidityPool(ctx, *amm.Pool)
	ctx.EventManager().EmitTypedEvents(msg)
	return &types.MsgSwapResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}
