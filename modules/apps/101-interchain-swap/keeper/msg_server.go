package keeper

import (
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
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

func (k Keeper) onCreatePoolAcknowledged(ctx sdk.Context, msg *types.MsgCreatePoolRequest) {
	//TODO:
}

func (k Keeper) onSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgDepositRequest, res *types.MsgDepositResponse) error {
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}
	pool.Supply.Amount.Add(res.PoolToken.Amount)
	k.SetInterchainLiquidityPool(ctx, pool)
	voucherAmount := sdk.NewCoins(*res.PoolToken)
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, voucherAmount)
	if err != nil {
		return err
	}
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(req.Sender), voucherAmount)
	return nil
}

func (k Keeper) onWithdrawAcknowledged(ctx sdk.Context, req *types.MsgWithdrawRequest, res *types.MsgWithdrawResponse) error {
	pool, found := k.GetInterchainLiquidityPool(ctx, "")
	if !found {
		return types.ErrNotFoundPool
	}
	pool.Supply.Amount.Sub(res.Tokens[0].Amount)
	k.SetInterchainLiquidityPool(ctx, pool)

	voucherAmount := sdk.NewCoins(*res.Tokens[0])
	err := k.bankKeeper.MintCoins(ctx, types.ModuleName, voucherAmount)
	if err != nil {
		return err
	}
	k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(req.Sender), types.ModuleName, voucherAmount)
	return nil
}

func (k Keeper) onSwapAcknowledged(ctx sdk.Context, req *types.MsgSwapRequest, res *types.MsgSwapResponse) error {
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
	pool, found := k.GetInterchainLiquidityPool(ctx, pooId)

	if !found {
		return nil, types.ErrNotFoundPool
	}

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

	if count == 1 {
		return nil, types.ErrInvalidDenomPair
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	ctx.EventManager().EmitTypedEvents(msg)
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
	_ = pool

	//TODO: Need to implement params module and market maker.

	amm := types.NewInterchainMarketMaker(
		pool,
		323,
	)

	poolToken, err := amm.DepositSingleAsset(*msg.Tokens[0])

	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to deposit because %s")
	}

	k.SetInterchainLiquidityPool(ctx, pool)

	ctx.EventManager().EmitTypedEvents(msg)
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

	_ = pool
	//TODO: need to implement amm part.
	//feeRate := parms.getPoolFeeRate()

	amm := types.NewInterchainMarketMaker(
		pool,
		323,
	)

	outToken, err := amm.Withdraw(*msg.PoolCoin, msg.DenomOut)

	if err != nil {
		return nil, errorsmod.Wrapf(err, "failed to withdraw because of %s!")
	}

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
			return nil, errorsmod.Wrapf(err, "failed to swap because of %s")
		}
	case types.SwapMsgType_RIGHT:
		outToken, err = amm.RightSwap(*msg.TokenIn, *msg.TokenOut)
		if err != nil {
			return nil, errorsmod.Wrapf(err, "failed to swap because of %s")
		}
	}

	expected := float64(msg.TokenOut.Amount.Uint64()) * (1 - float64(msg.Slippage)/10000)

	if outToken.Amount.LT(sdk.NewInt(int64(expected))) {
		return nil, errorsmod.Wrap(err, "doesn't meet slippage for swap!")
	}

	escrowAddr := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, escrowAddr.String(), sdk.AccAddress(msg.Recipient), sdk.NewCoins(*outToken))

	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient!")
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	ctx.EventManager().EmitTypedEvents(msg)
	return &types.MsgSwapResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}
