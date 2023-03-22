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

func (k Keeper) OnCreatePoolAcknowledged(ctx sdk.Context, msg *types.MsgCreatePoolRequest) error {
	//save pool after complete create operation in counter party chain.

	pool := types.NewInterchainLiquidityPool(
		ctx,
		msg.Sender,
		k.bankKeeper,
		msg.Tokens,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)

	// when add initial liquidity, we need to update pool token amount.
	initialLiquidity := sdk.NewCoin(pool.PoolId, msg.Tokens[0].Amount)
	pool.AddPoolSupply(initialLiquidity)

	// mint lpToken.
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(msg.Sender), initialLiquidity)
	if err != nil {
		return err
	}

	k.SetInterchainLiquidityPool(ctx, *pool)
	return nil
}

func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgDepositRequest, res *types.MsgDepositResponse) error {
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// mint voucher
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *res.PoolToken)
	if err != nil {
		return err
	}

	// pool status update.
	pool.AddPoolSupply(*res.PoolToken)

	if pool.Status != types.PoolStatus_POOL_STATUS_INITIAL {
		for _, token := range req.Tokens {
			pool.AddAsset(*token)
		}
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
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

	//burn voucher token.
	err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *req.PoolCoin)
	if err != nil {
		return err
	}
	// unlock token
	err = k.UnlockTokens(ctx,
		pool.EncounterPartyPort,
		pool.EncounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.Sender),
		sdk.NewCoins(*res.Tokens[0]),
	)
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) OnSwapAcknowledged(ctx sdk.Context, req *types.MsgSwapRequest, res *types.MsgSwapResponse) error {
	pooId := types.GetPoolId([]string{req.TokenIn.Denom, req.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, pooId)
	if !found {
		return types.ErrNotFoundPool
	}

	// pool status update
	pool.AddAsset(*req.TokenIn)
	pool.SubAsset(*res.Tokens[0])
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnCreatePoolReceived(ctx sdk.Context, msg *types.MsgCreatePoolRequest, destPort, destChannel string) (*string, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pooId := types.GetPoolIdWithTokens(msg.Tokens)
	_, found := k.GetInterchainLiquidityPool(ctx, pooId)

	if found {
		return nil, types.ErrAlreadyExistPool
	}

	// assume pool is ready when it is created.
	pool := *types.NewInterchainLiquidityPool(
		ctx,
		msg.Sender,
		k.bankKeeper,
		msg.Tokens,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)
	//count native tokens
	count := 0
	for _, token := range msg.Tokens {
		if k.bankKeeper.HasSupply(ctx, token.Denom) {
			count += 1
			pool.UpdateAssetPoolSide(token.Denom, types.PoolSide_NATIVE)
		} else {
			pool.UpdateAssetPoolSide(token.Denom, types.PoolSide_REMOTE)
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

	lpToken, err := amm.DepositSingleAsset(*msg.Tokens[0])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// pool status update
	pool.AddPoolSupply(*lpToken)

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
		return nil, err
	}

	// increase lp token mint amount
	pool.AddPoolSupply(*poolToken)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	if pool.Status == types.PoolStatus_POOL_STATUS_READY {
		// update pool tokens.
		for _, token := range msg.Tokens {
			pool.AddAsset(*token)
		}
	} else {
		// switch pool status to 'READY'
		pool.Status = types.PoolStatus_POOL_STATUS_READY
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

	pooId := types.GetPoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})
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

	if float64(outToken.Amount.Uint64()) < expected {
		return nil, errorsmod.Wrap(types.ErrFailedOnSwapReceived, "doesn't meet slippage for swap!, %s")
	}

	err = k.UnlockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*outToken))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient!")
	}

	// update pool status
	pool.SubAsset(*outToken)
	pool.AddAsset(*msg.TokenIn)

	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgSwapResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}
