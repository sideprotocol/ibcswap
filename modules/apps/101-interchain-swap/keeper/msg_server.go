package keeper

import (
	"github.com/btcsuite/btcutil/bech32"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	initialLiquidity := sdk.NewCoin(pool.PoolId, msg.Tokens[0].Amount.Add(msg.Tokens[1].Amount))
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

// OnDoubleDepositAcknowledged processes a double deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnDoubleDepositAcknowledged(ctx sdk.Context, req *types.MsgDoubleDepositRequest, res *types.MsgDoubleDepositResponse) error {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// Mint voucher tokens for the sender
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.Senders[0]), *res.PoolTokens[0])

	if err != nil {
		return err
	}

	// Update pool supply and status
	for _, poolToken := range res.PoolTokens {
		pool.AddPoolSupply(*poolToken)
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_INITIAL {
		for _, token := range req.Tokens {
			pool.AddAsset(*token)
		}
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
	}

	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnWithdrawAcknowledged(ctx sdk.Context, req *types.MsgWithdrawRequest, res *types.MsgWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolCoin.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	pool.SubPoolSupply(*res.Tokens[0])
	for _, token := range res.Tokens {
		pool.SubAsset(*token)
	}
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
		&pool,
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
		&pool,
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

// OnDoubleDepositReceived processes a double deposit request and returns a response or an error.
func (k Keeper) OnDoubleDepositReceived(ctx sdk.Context, msg *types.MsgDoubleDepositRequest) (*types.MsgDoubleDepositResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// // Verify the sender's address
	secondSenderAcc := k.authKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(msg.Senders[1]))
	senderPrefix, _, err := bech32.Decode(secondSenderAcc.GetAddress().String())
	if err != nil {
		return nil, err
	}
	if sdk.GetConfig().GetBech32AccountAddrPrefix() != senderPrefix {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "first address has to be this chain address (%s)", err)
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	}

	// Lock assets from senders to escrow account
	escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	// Create a deposit message
	sendMsg := banktypes.MsgSend{
		FromAddress: secondSenderAcc.GetAddress().String(),
		ToAddress:   escrowAccount.String(),
		Amount:      sdk.NewCoins(*msg.Tokens[1]),
	}

	deposit := types.EncounterPartyDepositTx{
		AccountSequence: secondSenderAcc.GetSequence(),
		Sender:          secondSenderAcc.GetAddress().String(),
		Token:           msg.Tokens[1],
	}

	rawDepositTx, err := types.ModuleCdc.Marshal(&deposit)
	if err != nil {
		return nil, err
	}

	pubKey := secondSenderAcc.GetPubKey()

	isValid := pubKey.VerifySignature(rawDepositTx, msg.EncounterPartySignature)
	if !isValid {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, ":%s", types.ErrInvalidSignature)
	}

	_, err = k.executeDepositTx(ctx, &sendMsg)
	if err != nil {
		return nil, err
	}

	// Instantiate an interchain market maker
	amm := types.NewInterchainMarketMaker(
		&pool,
		types.DefaultMaxFeeRate,
	)

	// Process double asset deposit
	poolTokens, err := amm.DepositDoubleAsset(msg.Tokens)
	if err != nil {
		return nil, err
	}

	// Increase LP token mint amount
	for _, token := range poolTokens {
		pool.AddPoolSupply(*token)
	}

	// Update pool tokens or switch pool status to 'READY'
	if pool.Status == types.PoolStatus_POOL_STATUS_READY {
		for _, token := range msg.Tokens {
			pool.AddAsset(*token)
		}
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
	}

	// Mint voucher tokens for the sender
	err = k.MintTokens(ctx, secondSenderAcc.GetAddress(), *poolTokens[1])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, ":%s", err)
	}

	// Save pool and market
	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgDoubleDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}

// OnWithdrawReceive processes a withdrawal request and returns a response or an error.
func (k Keeper) OnWithdrawReceived(ctx sdk.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// Instantiate an interchain market maker with the default fee rate
	amm := types.NewInterchainMarketMaker(
		&pool,
		types.DefaultMaxFeeRate,
	)

	// Calculate output token
	outToken, err := amm.Withdraw(*msg.PoolCoin, msg.DenomOut)

	// Check for errors in the withdrawal process
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnWithdrawReceived, "because of %s!", err)
	}

	// Ensure output token amount is greater than zero
	if outToken.Amount.LTE(sdk.NewInt(0)) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnWithdrawReceived, "because of %s!", "zero amount")
	}

	// Update pool status by subtracting the supplied pool coin and output token
	pool.SubPoolSupply(*msg.PoolCoin)
	pool.SubAsset(*outToken)

	// Save pool and market
	k.SetInterchainLiquidityPool(ctx, *amm.Pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgWithdrawResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}

// OnSwapReceived processes a swap request and returns a response or an error.
func (k Keeper) OnSwapReceived(ctx sdk.Context, msg *types.MsgSwapRequest) (*types.MsgSwapResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	poolID := types.GetPoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, poolID)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// Instantiate an interchain market maker with the default fee rate
	amm := types.NewInterchainMarketMaker(
		&pool,
		types.DefaultMaxFeeRate,
	)

	var outToken *sdk.Coin
	var err error
	switch msg.SwapType {
	case types.SwapMsgType_LEFT:
		outToken, err = amm.LeftSwap(*msg.TokenIn, msg.TokenOut.Denom)
	case types.SwapMsgType_RIGHT:
		outToken, err = amm.RightSwap(*msg.TokenIn, *msg.TokenOut)
	}

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnSwapReceived, "because of %s", err)
	}

	expected := float64(msg.TokenOut.Amount.Uint64()) * (1 - float64(msg.Slippage)/10000)

	if float64(outToken.Amount.Uint64()) < expected {
		return nil, errorsmod.Wrap(types.ErrFailedOnSwapReceived, "doesn't meet slippage for swap!, %s")
	}

	err = k.UnlockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*outToken))
	if err != nil {
		return nil, errorsmod.Wrap(err, "failed to move assets from escrow address to recipient!")
	}

	// Update pool status by subtracting output token and adding input token
	pool.SubAsset(*outToken)
	pool.AddAsset(*msg.TokenIn)

	// Save pool and market
	k.SetInterchainLiquidityPool(ctx, pool)
	k.SetInterchainMarketMaker(ctx, *amm)
	return &types.MsgSwapResponse{
		Tokens: []*sdk.Coin{outToken},
	}, nil
}

func (k Keeper) executeDepositTx(ctx sdk.Context, msg sdk.Msg) ([]byte, error) {

	txMsgData := &sdk.TxMsgData{
		MsgResponses: make([]*codectypes.Any, 1),
	}

	// CacheContext returns a new context with the multi-store branched into a cached storage object
	// writeCache is called only if all msgs succeed, performing state transitions atomically
	cacheCtx, writeCache := ctx.CacheContext()
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	any, err := k.executeMsg(cacheCtx, msg)
	if err != nil {
		return nil, err
	}
	writeCache()

	txMsgData.MsgResponses[0] = any
	txResponse, err := k.cdc.Marshal(txMsgData)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to marshal tx data")
	}

	return txResponse, nil
}

// Attempts to get the message handler from the router and if found will then execute the message.
// If the message execution is successful, the proto marshaled message response will be returned.
func (k Keeper) executeMsg(ctx sdk.Context, msg sdk.Msg) (*codectypes.Any, error) {
	handler := k.msgRouter.Handler(msg)
	if handler == nil {
		return nil, types.ErrInvalidMsgRouter
	}

	res, err := handler(ctx, msg)
	if err != nil {
		return nil, err
	}

	// NOTE: The sdk msg handler creates a new EventManager, so events must be correctly propagated back to the current context
	ctx.EventManager().EmitEvents(res.GetEvents())

	// Each individual sdk.Result has exactly one Msg response. We aggregate here.
	msgResponse := res.MsgResponses[0]
	if msgResponse == nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidMsg, "got nil Msg response for msg %s", sdk.MsgTypeURL(msg))
	}

	return msgResponse, nil
}
