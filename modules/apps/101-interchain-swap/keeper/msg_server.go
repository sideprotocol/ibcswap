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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionMakeOrder+"_"+types.EventValueSuffixAcknowledged, poolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)
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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionTakeOrder+"_"+types.EventValueSuffixAcknowledged, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)

	return nil
}

// OnCreatePoolAcknowledged processes the create pool acknowledgement, mints LP tokens, and saves the liquidity pool.
func (k Keeper) OnCancelPoolAcknowledged(ctx sdk.Context, msg *types.MsgCancelPoolRequest) error {
	// Save pool after completing the create operation in the counterparty chain

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}
	sourceAsset, err := pool.FindAssetBySide(types.PoolAssetSide_SOURCE)
	if err != nil {
		return err
	}
	if err = k.UnlockTokens(ctx, msg.SourcePort, msg.SourceChannel, sdk.MustAccAddressFromBech32(msg.Creator), sdk.NewCoins(*sourceAsset)); err != nil {
		return err
	}

	k.RemoveInterchainLiquidityPool(ctx, msg.PoolId)

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionCancelOrder+"_"+types.EventValueSuffixAcknowledged, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)

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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionSingleDeposit+"_"+types.EventValueSuffixAcknowledged, req.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: req.Sender,
		},
		sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: req.Token.String(),
		},
		sdk.Attribute{
			Key:   types.AttributeKeyLpToken,
			Value: res.PoolToken.String(),
		},
	)

	return nil
}

// OnMultiAssetDepositAcknowledged processes a double deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnCancelMultiAssetDepositAcknowledged(ctx sdk.Context, req *types.MsgCancelMultiAssetDepositRequest) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	order, found := k.GetMultiDepositOrder(ctx, req.PoolId, req.OrderId)
	if !found {
		return types.ErrCancelOrder
	}

	// Create escrow module account here

	if err := k.UnlockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(req.Creator), sdk.NewCoins(*order.Deposits[0])); err != nil {
		return types.ErrCancelOrder
	}
	k.RemoveMultiDepositOrder(ctx, req.PoolId, req.OrderId)

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionCancelOrder+"_"+types.EventValueSuffixAcknowledged, req.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: req.Creator,
		},
	)
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
	lpTokenAttr := []sdk.Attribute{}
	for _, poolToken := range stateChange.PoolTokens {
		pool.AddPoolSupply(*poolToken)
		lpTokenAttr = append(lpTokenAttr, sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: poolToken.String(),
		})
	}

	tokenInAttr := []sdk.Attribute{}
	for _, deposit := range order.Deposits {
		pool.AddAsset(*deposit)
		tokenInAttr = append(tokenInAttr, sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: deposit.String(),
		})
	}

	order.Status = types.OrderStatus_COMPLETE

	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	// Update order statuse
	k.SetMultiDepositOrder(ctx, order)

	eventAttr := []sdk.Attribute{
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: req.Sender,
		},
		sdk.Attribute{
			Key:   types.AttributeKeyMultiDepositOrderId,
			Value: req.OrderId,
		},
	}
	eventAttr = append(eventAttr, tokenInAttr...)
	eventAttr = append(eventAttr, lpTokenAttr...)
	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionTakeMultiDeposit+"_"+types.EventValueSuffixReceived, req.PoolId,
		eventAttr...,
	)
	return nil
}

func (k Keeper) OnMultiAssetWithdrawAcknowledged(ctx sdk.Context, req *types.MsgMultiAssetWithdrawRequest, stateChange types.StateChange) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	for _, poolAsset := range stateChange.Out {
		pool.SubtractAsset(*poolAsset)
	}
	pool.SubtractPoolSupply(*req.PoolToken)

	nativeToken, err := pool.FindDenomBySide(types.PoolAssetSide_SOURCE)
	if err != nil {
		return err
	}

	out, err := stateChange.FindOutByDenom(*nativeToken)
	if err != nil {
		return err
	}

	// unlock token
	err = k.UnlockTokens(ctx,
		pool.CounterPartyPort,
		pool.CounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.Receiver),
		sdk.NewCoins(*out),
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

	// emit events
	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolCreator,
			Value: req.Receiver,
		},
		{
			Key:   types.AttributeKeyLpToken,
			Value: req.PoolToken.String(),
		},
	}

	k.EmitEvent(
		ctx, types.EventValueActionWithdrawMultiDeposit+"_"+types.EventValueSuffixAcknowledged, req.PoolId,
		eventAttr...,
	)

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

	// Emit events
	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolId,
			Value: req.PoolId,
		},
		{
			Key:   types.AttributeKeyTokenIn,
			Value: req.TokenIn.String(),
		},
		{
			Key:   types.AttributeKeyTokenOut,
			Value: req.TokenOut.String(),
		},
	}

	k.EmitEvent(
		ctx, types.EventValueActionSwap+"_"+types.EventValueSuffixAcknowledged, req.PoolId,
		eventAttr...,
	)
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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionMakeOrder+"_"+types.EventValueSuffixReceived, poolID,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)

	return &poolID, nil
}

func (k Keeper) OnTakePoolReceived(ctx sdk.Context, msg *types.MsgTakePoolRequest) (*types.MsgTakePoolResponse, error) {

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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionTakeOrder+"_"+types.EventValueSuffixReceived, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)

	return &types.MsgTakePoolResponse{
		PoolId: pool.Id,
	}, nil
}

func (k Keeper) OnCancelPoolReceived(ctx sdk.Context, msg *types.MsgCancelPoolRequest) (*types.MsgCancelPoolResponse, error) {

	if _, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId); !found {
		return nil, types.ErrNotFoundPool
	}
	// remove pool status
	k.RemoveInterchainLiquidityPool(ctx, msg.PoolId)

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionCancelOrder+"_"+types.EventValueSuffixReceived, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
	)
	return &types.MsgCancelPoolResponse{
		PoolId: msg.PoolId,
	}, nil
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

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionSingleDeposit+"_"+types.EventValueSuffixReceived, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Sender,
		},
		sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: msg.Token.String(),
		},
		sdk.Attribute{
			Key:   types.AttributeKeyLpToken,
			Value: stateChange.Out[0].String(),
		},
	)

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

	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Deposits[0].Sender,
		},
		{
			Key:   types.AttributeKeyMultiDepositOrderId,
			Value: order.Id,
		},
	}

	eventAttr = append(eventAttr, types.GetEventAttrOfAsset(msg.Deposits)...)

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionMakeMultiDeposit+"_"+types.EventValueSuffixReceived, msg.PoolId,
		eventAttr...,
	)
	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: []*sdk.Coin{},
	}, nil
}

// OnMultiAssetDepositReceived processes a double deposit request and returns a response or an error.
func (k Keeper) OnCancelMultiAssetDepositReceived(ctx sdk.Context, msg *types.MsgCancelMultiAssetDepositRequest) (*types.MsgCancelMultiAssetDepositResponse, error) {

	// Retrieve the liquidity pool
	if _, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId); !found {
		return nil, errorsmod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	if _, found := k.GetMultiDepositOrder(ctx, msg.PoolId, msg.OrderId); !found {
		return nil, errorsmod.Wrapf(types.ErrNotFoundPool, "%s", types.ErrFailedMultiAssetDeposit)
	}

	k.RemoveMultiDepositOrder(ctx, msg.PoolId, msg.OrderId)

	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionCancelMultiDeposit+"_"+types.EventValueSuffixReceived, msg.PoolId,
		sdk.Attribute{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Creator,
		},
		sdk.Attribute{
			Key:   types.AttributeKeyMultiDepositOrderId,
			Value: msg.OrderId,
		},
	)
	return &types.MsgCancelMultiAssetDepositResponse{
		PoolId:  msg.PoolId,
		OrderId: msg.PoolId,
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

	tokenInAttr := []sdk.Attribute{}
	for _, asset := range order.Deposits {
		pool.AddAsset(*asset)
		tokenInAttr = append(tokenInAttr, sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: asset.String(),
		})
	}

	// Mint voucher tokens for the sender
	totalPoolToken := sdk.NewCoin(msg.PoolId, sdk.NewInt(0))
	lpTokenAttr := []sdk.Attribute{}

	for _, poolToken := range stateChange.PoolTokens {
		totalPoolToken = totalPoolToken.Add(*poolToken)
		lpTokenAttr = append(lpTokenAttr, sdk.NewAttribute(
			types.AttributeKeyLpToken, poolToken.String(),
		))
		lpTokenAttr = append(lpTokenAttr, sdk.Attribute{
			Key:   types.AttributeKeyTokenIn,
			Value: poolToken.String(),
		})
	}
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(order.SourceMaker), totalPoolToken)

	if err != nil {
		return nil, err
	}

	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetMultiDepositOrder(ctx, order)

	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Sender,
		},
		{
			Key:   types.AttributeKeyMultiDepositOrderId,
			Value: msg.OrderId,
		},
	}
	eventAttr = append(eventAttr, tokenInAttr...)
	eventAttr = append(eventAttr, lpTokenAttr...)
	// emit events
	k.EmitEvent(
		ctx, types.EventValueActionTakeMultiDeposit+"_"+types.EventValueSuffixAcknowledged, msg.PoolId,
		eventAttr...,
	)

	return &types.MsgMultiAssetDepositResponse{}, nil
}

// OnMultiAssetWithdrawReceived processes a withdrawal request and returns a response or an error.
func (k Keeper) OnMultiAssetWithdrawReceived(ctx sdk.Context, msg *types.MsgMultiAssetWithdrawRequest, stateChange types.StateChange) (*types.MsgMultiAssetWithdrawResponse, error) {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// Update pool status by subtracting the supplied pool coin and output token

	rawOuts := []string{}
	for _, poolAsset := range stateChange.Out {
		pool.SubtractAsset(*poolAsset)
		rawOuts = append(rawOuts, poolAsset.String())
	}
	pool.SubtractPoolSupply(*msg.PoolToken)

	nativeDenom, err := pool.FindDenomBySide(types.PoolAssetSide_SOURCE)
	if err != nil {
		return nil, err
	}

	out, err := stateChange.FindOutByDenom(*nativeDenom)
	if err != nil {
		return nil, err
	}
	// escrow operation
	err = k.UnlockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.CounterPartyReceiver), sdk.NewCoins(*out))

	if err != nil {
		return nil, err
	}

	if pool.Supply.Amount.LTE(sdk.NewInt(0)) {
		k.RemoveInterchainLiquidityPool(ctx, msg.PoolId)
	} else {
		// Save pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}

	// emit events
	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolCreator,
			Value: msg.Receiver,
		},
		{
			Key:   types.AttributeKeyLpToken,
			Value: msg.PoolToken.String(),
		},
		{
			Key:   types.TypeMsgWithdraw,
			Value: msg.PoolToken.String(),
		},
	}

	k.EmitEvent(
		ctx, types.EventValueActionWithdrawMultiDeposit+"_"+types.EventValueSuffixReceived, msg.PoolId,
		eventAttr...,
	)

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

	// Emit events
	eventAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyPoolId,
			Value: msg.PoolId,
		},
		{
			Key:   types.AttributeKeyTokenIn,
			Value: msg.TokenIn.String(),
		},
		{
			Key:   types.AttributeKeyTokenOut,
			Value: msg.TokenOut.String(),
		},
	}

	k.EmitEvent(
		ctx, types.EventValueActionSwap+"_"+types.EventValueSuffixReceived, msg.PoolId,
		eventAttr...,
	)
	return &types.MsgSwapResponse{
		Tokens: stateChange.Out,
	}, nil
}
