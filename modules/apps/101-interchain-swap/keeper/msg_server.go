package keeper

import (
	"fmt"

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

	initialLiquidity := sdk.NewInt(int64(msg.InitalLiquidity))
	// when add initial liquidity, we need to update pool token amount.
	pool.AddPoolSupply(sdk.NewCoin(pool.PoolId, initialLiquidity))
	pool.AddAsset(sdk.NewCoin(msg.Denoms[0], initialLiquidity))
	k.SetInterchainLiquidityPool(ctx, *pool)
}

func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgDepositRequest, res *types.MsgDepositResponse) error {
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// pool status update.
	pool.AddPoolSupply(*res.PoolToken)
	for _, token := range req.Tokens {
		pool.AddAsset(*token)
	}

	// mint voucher
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *res.PoolToken)
	if err != nil {
		return err
	}
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnWithdrawAcknowledged(ctx sdk.Context, req *types.MsgWithdrawRequest, res *types.MsgWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolCoin.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	pool.SubPoolSupply(*res.Tokens[0])
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
	pooId := types.GetPoolId([]string{req.TokenOut.Denom, req.TokenIn.Denom})
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

	// assume pool is ready when it is created.
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

	//TODO: Need to implement params module and market maker.
	amm := types.NewInterchainMarketMaker(
		pool,
		types.DefaultMaxFeeRate,
	)

	initialLiquidity := sdk.NewCoin(msg.Denoms[0], sdk.NewInt(int64(msg.InitalLiquidity)))
	lpToken, err := amm.DepositSingleAsset(initialLiquidity)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// pool status update
	pool.AddPoolSupply(*lpToken)
	pool.AddAsset(sdk.NewCoin(msg.Denoms[0], sdk.NewInt(int64(msg.InitalLiquidity))))

	// save pool status
	k.SetInterchainLiquidityPool(ctx, pool)
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
		types.DefaultMaxFeeRate,
	)

	poolToken, err := amm.DepositSingleAsset(*msg.Tokens[0])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// update supply.
	err = pool.AddPoolSupply(*poolToken)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// update pool tokens.
	for _, token := range msg.Tokens {
		err = pool.AddAsset(*token)
		// if err != nil {
		// 	return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
		// }
	}

	// save pool and market.
	k.SetInterchainLiquidityPool(ctx, pool)
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

	// calculate output token.
	amm := types.NewInterchainMarketMaker(
		pool,
		types.DefaultMaxFeeRate,
	)

	outToken, err := amm.Withdraw(*msg.PoolCoin, msg.DenomOut)
	ctx.EventManager().EmitTypedEvents(outToken)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s!", err)
	}

	// save pool and market.
	k.SetInterchainLiquidityPool(ctx, *amm.Pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgWithdrawResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}

func (k Keeper) OnSwapReceived(ctx sdk.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pooId := types.GetPoolId([]string{msg.TokenOut.Denom, msg.TokenIn.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, pooId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	//TODO: need to implement amm part.
	//feeRate := parms.getPoolFeeRate()

	amm := types.NewInterchainMarketMaker(
		pool,
		types.DefaultMaxFeeRate,
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
	fmt.Println(outToken)
	if float64(outToken.Amount.Uint64()) < expected {
		return nil, errorsmod.Wrap(types.ErrFailedOnSwapReceived, "doesn't meet slippage for swap!, %s")
	}

	err = k.UnlockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*outToken))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient!")
	}

	// remove pool token
	pool.SubPoolSupply(*outToken)

	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgSwapResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}
