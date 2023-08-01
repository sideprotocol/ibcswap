package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
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

// OnCreatePoolAcknowledged processes the create pool acknowledgement, mints LP tokens, and saves the liquidity pool.
func (k Keeper) OnMakePoolAcknowledged(ctx sdk.Context, msg *types.MsgMakePoolRequest, poolId string) error {
	// Save pool after completing the create operation in the counterparty chain

	pool := types.NewInterchainLiquidityPool(
		ctx,
		msg.Creator,
		msg.CounterPartyCreator,
		k.bankKeeper,
		poolId,
		msg.Liquidity,
		msg.SwapFee,
		msg.SourcePort,
		msg.SourceChannel,
	)

	pool.SourceChainId = ctx.ChainID()

	// Mint LP tokens
	totalAmount := sdk.NewInt(0)
	for _, asset := range msg.Liquidity {
		totalAmount = totalAmount.Add(asset.Balance.Amount)
	}
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(msg.Creator), sdk.Coin{
		Denom: pool.Supply.Denom, Amount: totalAmount.Mul(sdk.NewInt(int64(msg.Liquidity[0].Weight))).Quo(sdk.NewInt((100))),
	})

	if err != nil {
		return err
	}

	k.SetInterchainLiquidityPool(ctx, *pool)
	return nil
}

// OnCreatePoolAcknowledged processes the create pool acknowledgement, mints LP tokens, and saves the liquidity pool.
func (k Keeper) OnTakePoolAcknowledged(ctx sdk.Context, msg *types.MsgTakePoolRequest) error {
	// Save pool after completing the create operation in the counterparty chain

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)

	if !found {
		return types.ErrNotFoundPool
	}

	pool.Status = types.PoolStatus_ACTIVE

	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

// OnSingleAssetDepositAcknowledged processes a single deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnSingleAssetDepositAcknowledged(ctx sdk.Context, req *types.MsgSingleAssetDepositRequest, res *types.MsgSingleAssetDepositResponse) error {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// Mint voucher tokens for the sender
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *res.PoolToken)
	if err != nil {
		return err
	}

	// update pool status
	pool.AddAsset(*req.Token)
	pool.AddPoolSupply(*res.PoolToken)

	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

// OnMultiAssetDepositAcknowledged processes a double deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnMakeMultiAssetDepositAcknowledged(ctx sdk.Context, req *types.MsgMakeMultiAssetDepositRequest) error {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

// OnMultiAssetDepositAcknowledged processes a double deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnTakeMultiAssetDepositAcknowledged(ctx sdk.Context, req *types.MsgTakeMultiAssetDepositRequest, stateChange types.StateChange) error {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	order, found := k.GetMultiDepositOrder(ctx, req.PoolId, req.OrderId)
	if !found {
		return types.ErrNotFoundMultiDepositOrder
	}
	// Update pool supply and status
	for _, poolToken := range stateChange.PoolTokens {
		pool.AddPoolSupply(*poolToken)
	}
	for _, deposit := range order.Deposits {
		pool.AddAsset(*deposit)
	}

	order.Status = types.OrderStatus_COMPLETE

	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	// Update order statuse
	k.SetMultiDepositOrder(ctx, order)
	return nil
}

func (k Keeper) OnMultiAssetWithdrawAcknowledged(ctx sdk.Context, req *types.MsgMultiAssetWithdrawRequest, res *types.MsgMultiAssetWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolToken.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	for _, poolAsset := range res.Tokens {
		pool.SubtractAsset(*poolAsset)
	}
	pool.SubtractPoolSupply(*req.PoolToken)

	//burn voucher token.
	err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.Receiver), *req.PoolToken)
	if err != nil {
		return err
	}

	// unlock token
	err = k.UnlockTokens(ctx,
		pool.CounterPartyPort,
		pool.CounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.Receiver),
		sdk.NewCoins(*res.Tokens[0]),
	)

	if err != nil {
		return err
	}

	if pool.Supply.Amount.Equal(sdk.NewInt(0)) {
		k.RemoveInterchainLiquidityPool(ctx, req.PoolId)
	} else {
		// Save pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}
	return nil
}

func (k Keeper) OnSwapAcknowledged(ctx sdk.Context, req *types.MsgSwapRequest, res *types.MsgSwapResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// pool status update
	pool.AddAsset(*req.TokenIn)
	pool.SubtractAsset(*req.TokenOut)
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

// onReceive
func (k Keeper) OnMakePoolReceived(ctx sdk.Context, msg *types.MsgMakePoolRequest, poolID, sourceChainId string) (*string, error) {

	_, found := k.GetInterchainLiquidityPool(ctx, poolID)

	if found {
		return nil, types.ErrAlreadyExistPool
	}

	// assume pool is ready when it is created.
	pool := *types.NewInterchainLiquidityPool(
		ctx,
		msg.Creator,
		msg.CounterPartyCreator,
		k.bankKeeper,
		poolID,
		msg.Liquidity,
		msg.SwapFee,
		msg.SourcePort,
		msg.SourceChannel,
	)

	pool.SourceChainId = sourceChainId

	if !k.bankKeeper.HasSupply(ctx, msg.Liquidity[1].Balance.Denom) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "due to %s", types.ErrInvalidDecimalPair)
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	return &poolID, nil
}

func (k Keeper) OnTakePoolReceived(ctx sdk.Context, msg *types.MsgTakePoolRequest) (*string, error) {

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}
	pool.Status = types.PoolStatus_ACTIVE

	// mint voucher token
	asset, err := pool.FindPoolAssetBySide(types.PoolAssetSide_DESTINATION)
	if err != nil {
		return nil, err
	}

	totalAmount := pool.SumOfPoolAssets()
	if err = k.MintTokens(ctx, sdk.MustAccAddressFromBech32(pool.SourceCreator), sdk.Coin{
		Denom: pool.Supply.Denom, Amount: totalAmount.Mul(sdk.NewInt(int64(asset.Weight))).Quo(sdk.NewInt((100))),
	}); err != nil {
		return nil, err
	}

	// save pool status
	k.SetInterchainLiquidityPool(ctx, pool)
	return &pool.Id, nil
}

func (k Keeper) OnSingleAssetDepositReceived(ctx sdk.Context, msg *types.MsgSingleAssetDepositRequest, stateChange *types.StateChange) (*types.MsgSingleAssetDepositResponse, error) {

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// update pool status
	pool.AddPoolSupply(*stateChange.PoolTokens[0])
	pool.AddAsset(*msg.Token)

	k.SetInterchainLiquidityPool(ctx, pool)
	return &types.MsgSingleAssetDepositResponse{
		PoolToken: stateChange.PoolTokens[0],
	}, nil
}

// OnMultiAssetDepositReceived processes a double deposit request and returns a response or an error.
func (k Keeper) OnMakeMultiAssetDepositReceived(ctx sdk.Context, msg *types.MsgMakeMultiAssetDepositRequest, stateChange *types.StateChange) (*types.MsgMultiAssetDepositResponse, error) {

	// Verify the sender's address
	_, err := sdk.AccAddressFromBech32(msg.Deposits[1].Sender)
	if err != nil {
		return nil, err
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	// create order
	order := types.MultiAssetDepositOrder{
		Id:               stateChange.MultiDepositOrderId,
		PoolId:           msg.PoolId,
		ChainId:          pool.SourceChainId,
		SourceMaker:      msg.Deposits[0].Sender,
		DestinationTaker: msg.Deposits[1].Sender,
		Deposits:         types.GetCoinsFromDepositAssets(msg.Deposits),
		Status:           types.OrderStatus_PENDING,
		CreatedAt:        ctx.BlockHeight(),
	}

	k.SetMultiDepositOrder(ctx, order)
	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: []*sdk.Coin{},
	}, nil
}

// OnMultiAssetDepositReceived processes a double deposit request and returns a response or an error.
func (k Keeper) OnTakeMultiAssetDepositReceived(ctx sdk.Context, msg *types.MsgTakeMultiAssetDepositRequest, stateChange *types.StateChange) (*types.MsgMultiAssetDepositResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	order, found := k.GetMultiDepositOrder(ctx, msg.PoolId, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFoundPool, "%s", types.ErrFailedMultiAssetDeposit)
	}

	order.Status = types.OrderStatus_COMPLETE

	// pool status update
	for _, supply := range stateChange.PoolTokens {
		pool.AddPoolSupply(*supply)
	}
	for _, asset := range order.Deposits {
		pool.AddAsset(*asset)
	}

	// Mint voucher tokens for the sender
	totalPoolToken := sdk.NewCoin(msg.PoolId, sdk.NewInt(0))
	for _, poolToken := range stateChange.PoolTokens {
		totalPoolToken = totalPoolToken.Add(*poolToken)
	}
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(order.SourceMaker), totalPoolToken)

	if err != nil {
		return nil, err
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetMultiDepositOrder(ctx, order)
	return &types.MsgMultiAssetDepositResponse{}, nil
}

// OnMultiAssetWithdrawReceived processes a withdrawal request and returns a response or an error.
func (k Keeper) OnMultiAssetWithdrawReceived(ctx sdk.Context, msg *types.MsgMultiAssetWithdrawRequest, stateChange *types.StateChange) (*types.MsgMultiAssetWithdrawResponse, error) {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolToken.Denom)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// Update pool status by subtracting the supplied pool coin and output token
	for _, poolAsset := range stateChange.Out {
		pool.SubtractAsset(*poolAsset)
	}
	for _, poolToken := range stateChange.PoolTokens {
		pool.SubtractPoolSupply(*poolToken)
	}

	// escrow operation
	err := k.UnlockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.CounterPartyReceiver), sdk.NewCoins(*stateChange.Out[1]))

	if err != nil {
		return nil, err
	}

	if pool.Supply.Amount.LTE(sdk.NewInt(0)) {
		k.RemoveInterchainLiquidityPool(ctx, msg.PoolId)
	} else {
		// Save pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}

	return &types.MsgMultiAssetWithdrawResponse{
		Tokens: stateChange.Out,
	}, nil
}

// OnSwapReceived processes a swap request and returns a response or an error.
func (k Keeper) OnSwapReceived(ctx sdk.Context, msg *types.MsgSwapRequest, stateChange *types.StateChange) (*types.MsgSwapResponse, error) {

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)

	if !found {
		return nil, types.ErrNotFoundPool
	}

	_, err := sdk.AccAddressFromBech32(msg.Recipient)
	if err != nil {
		return nil, err
	}

	err = k.UnlockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*stateChange.Out[0]))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient")
	}

	// Update pool status by subtracting output token and adding input token
	pool.SubtractAsset(*stateChange.Out[0])
	pool.AddAsset(*msg.TokenIn)

	// Save pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return &types.MsgSwapResponse{
		Tokens: stateChange.Out,
	}, nil
}
