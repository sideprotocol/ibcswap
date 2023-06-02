package keeper

import (
	"github.com/btcsuite/btcutil/bech32"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

// OnCreatePoolAcknowledged processes the create pool acknowledgement, mints LP tokens, and saves the liquidity pool.
func (k Keeper) OnCreatePoolAcknowledged(ctx sdk.Context, msg *types.MsgCreatePoolRequest) error {
	// Save pool after completing the create operation in the counterparty chain
	poolId := types.GetPoolId(msg.GetLiquidityDenoms())
	pool := types.NewInterchainLiquidityPool(
		ctx,
		msg.Creator,
		k.bankKeeper,
		poolId,
		msg.Liquidity,
		msg.SourcePort,
		msg.SourceChannel,
	)

	// Mint LP tokens
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(msg.Creator), *pool.Supply)
	if err != nil {
		return err
	}
	// calculate pool price
	// Instantiate an interchain market maker with the default fee rate
	amm := *types.NewInterchainMarketMaker(pool)
	pool.PoolPrice = float32(amm.LpPrice())
	// save initial pool asset amount
	var poolAssets []sdk.Coin
	for _, asset := range msg.Liquidity {
		poolAssets = append(poolAssets, *asset.Balance)
	}
	k.SetInitialPoolAssets(ctx, pool.Id, poolAssets)

	// Save the liquidity pool
	k.SetInterchainLiquidityPool(ctx, *pool)
	return nil
}

// OnSingleDepositAcknowledged processes a single deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, req *types.MsgSingleAssetDepositRequest, res *types.MsgSingleAssetDepositResponse) error {

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
	pool.SubtractAsset(*req.Token)
	pool.SubtractPoolSupply(*res.PoolToken)

	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

// OnMultiAssetDepositAcknowledged processes a double deposit acknowledgement, mints voucher tokens, and updates the liquidity pool.
func (k Keeper) OnMultiAssetDepositAcknowledged(ctx sdk.Context, req *types.MsgMultiAssetDepositRequest, res *types.MsgMultiAssetDepositResponse) error {

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolId)
	if !found {
		return types.ErrNotFoundPool
	}

	// Mint voucher tokens for the sender
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.Deposits[0].Sender), *res.PoolTokens[0])

	if err != nil {
		return err
	}

	// Update pool supply and status
	for _, poolToken := range res.PoolTokens {
		pool.AddPoolSupply(*poolToken)
	}

	for _, deposit := range req.Deposits {
		pool.AddAsset(*deposit.Balance)
	}
	// Save the updated liquidity pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnSingleWithdrawAcknowledged(ctx sdk.Context, req *types.MsgSingleAssetWithdrawRequest, res *types.MsgSingleAssetWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.PoolCoin.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	pool.SubtractAsset(*res.Tokens[0])
	pool.SubtractPoolSupply(*req.PoolCoin)

	//burn voucher token.
	err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *req.PoolCoin)
	if err != nil {
		return err
	}

	// unlock token
	err = k.UnlockTokens(ctx,
		pool.CounterPartyPort,
		pool.CounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.Sender),
		sdk.NewCoins(*res.Tokens[0]),
	)

	if err != nil {
		return err
	}

	// remove pool supply is zero.
	if pool.Supply.Amount.Equal(sdk.NewInt(0)) {
		k.RemoveInterchainLiquidityPool(ctx, pool.Id)
	} else {
		// save pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}
	return nil
}

func (k Keeper) OnMultiWithdrawAcknowledged(ctx sdk.Context, req *types.MsgMultiAssetWithdrawRequest, res *types.MsgMultiAssetWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.Withdraws[0].Balance.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	for _, poolAsset := range res.Tokens {
		pool.SubtractAsset(*poolAsset)
	}
	for _, poolToken := range req.Withdraws {
		pool.SubtractPoolSupply(*poolToken.Balance)
	}
	//burn voucher token.
	err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.Sender), *req.Withdraws[0].Balance)
	if err != nil {
		return err
	}

	// unlock token
	err = k.UnlockTokens(ctx,
		pool.CounterPartyPort,
		pool.CounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.Sender),
		sdk.NewCoins(*res.Tokens[0]),
	)

	if err != nil {
		return err
	}

	// save pool
	k.SetInterchainLiquidityPool(ctx, pool)
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
	pool.SubtractAsset(*res.Tokens[0])
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnCreatePoolReceived(ctx sdk.Context, msg *types.MsgCreatePoolRequest, destPort, destChannel string) (*string, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	poolID := types.GetPoolId(msg.GetLiquidityDenoms())
	_, found := k.GetInterchainLiquidityPool(ctx, poolID)

	if found {
		return nil, types.ErrAlreadyExistPool
	}

	// assume pool is ready when it is created.
	pool := *types.NewInterchainLiquidityPool(
		ctx,
		msg.Creator,
		k.bankKeeper,
		poolID,
		msg.Liquidity,
		msg.SourcePort,
		msg.SourceChannel,
	)

	if !k.bankKeeper.HasSupply(ctx, msg.Liquidity[1].Balance.Denom) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "due to %s", types.ErrInvalidDecimalPair)
	}

	creatorAddr := sdk.MustAccAddressFromBech32(msg.CounterPartyCreator)
	liquidity := k.bankKeeper.GetBalance(ctx, creatorAddr, msg.Liquidity[1].Balance.Denom)
	if liquidity.Amount.LTE(sdk.NewInt(0)) {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "due to %s", types.ErrInEnoughAmount)
	}
	creatorAcc := k.authKeeper.GetAccount(ctx, creatorAddr)
	escrowAccount := types.GetPoolId(msg.GetLiquidityDenoms())
	//

	sendMsg := banktypes.MsgSend{
		FromAddress: creatorAcc.GetAddress().String(),
		ToAddress:   escrowAccount,
		Amount:      sdk.NewCoins(*msg.Liquidity[1].Balance),
	}

	err := k.VerifySignature(ctx, creatorAcc, *msg.Liquidity[1].Balance, msg.CounterPartySig)
	if err != nil {
		return nil, err
	}
	_, err = k.executeDepositTx(ctx, &sendMsg)
	if err != nil {
		return nil, err
	}
	
	// Instantiate an interchain market maker with the default fee rate
	amm := *types.NewInterchainMarketMaker(&pool)
	// calculate
	pool.PoolPrice = float32(amm.LpPrice())
	// save pool status
	k.SetInterchainLiquidityPool(ctx, pool)

	// save initial pool asset amount
	var poolAssets []sdk.Coin
	for _, coin := range msg.Liquidity {
		poolAssets = append(poolAssets, *coin.Balance)
	}
	k.SetInitialPoolAssets(ctx, poolID, poolAssets)

	return &poolID, nil
}

func (k Keeper) OnSingleAssetDepositReceived(ctx sdk.Context, msg *types.MsgSingleAssetDepositRequest, stateChange *types.StateChange) (*types.MsgSingleAssetDepositResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// update pool status
	pool.AddPoolSupply(*stateChange.PoolTokens[0])
	pool.AddAsset(*msg.Token)

	// save pool.
	k.SetInterchainLiquidityPool(ctx, pool)
	return &types.MsgSingleAssetDepositResponse{
		PoolToken: stateChange.PoolTokens[0],
	}, nil
}

// OnMultiAssetDepositReceived processes a double deposit request and returns a response or an error.
func (k Keeper) OnMultiAssetDepositReceived(ctx sdk.Context, msg *types.MsgMultiAssetDepositRequest, stateChange *types.StateChange) (*types.MsgMultiAssetDepositResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Verify the sender's address
	senderAcc := k.authKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(msg.Deposits[1].Sender))
	senderPrefix, _, err := bech32.Decode(senderAcc.GetAddress().String())
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
	escrowAccount := types.GetEscrowAddress(pool.CounterPartyPort, pool.CounterPartyChannel)

	// Create a deposit message
	sendMsg := banktypes.MsgSend{
		FromAddress: senderAcc.GetAddress().String(),
		ToAddress:   escrowAccount.String(),
		Amount:      sdk.NewCoins(*msg.Deposits[1].Balance),
	}

	err = k.VerifySignature(ctx, senderAcc, *msg.Deposits[1].Balance, msg.Deposits[1].Signature)
	if err != nil {
		return nil, err
	}

	_, err = k.executeDepositTx(ctx, &sendMsg)
	if err != nil {
		return nil, err
	}

	// Increase LP token mint amount
	for _, token := range stateChange.PoolTokens {
		pool.AddPoolSupply(*token)
	}

	// Update pool tokens or switch pool status to 'READY'
	for _, deposit := range msg.Deposits {
		pool.AddAsset(*deposit.Balance)
	}

	// Mint voucher tokens for the sender
	err = k.MintTokens(ctx, senderAcc.GetAddress(), *stateChange.PoolTokens[1])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, ":%s", err)
	}
	// Save pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: stateChange.PoolTokens,
	}, nil
}

func (k Keeper) OnSingleAssetWithdrawReceived(ctx sdk.Context, msg *types.MsgSingleAssetWithdrawRequest, stateChange *types.StateChange) (*types.MsgSingleAssetWithdrawResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	// update pool status
	pool.SubtractAsset(*stateChange.Out[0])
	for _, poolToken := range stateChange.PoolTokens {
		pool.SubtractPoolSupply(*poolToken)
	}

	// remove pool supply is zero.
	if pool.Supply.Amount.Equal(sdk.NewInt(0)) {
		k.RemoveInterchainLiquidityPool(ctx, pool.Id)
	} else {
		// save pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}
	return &types.MsgSingleAssetWithdrawResponse{
		Tokens: stateChange.Out,
	}, nil
}

// OnMultiAssetWithdrawReceived processes a withdrawal request and returns a response or an error.
func (k Keeper) OnMultiAssetWithdrawReceived(ctx sdk.Context, msg *types.MsgMultiAssetWithdrawRequest, stateChange *types.StateChange) (*types.MsgMultiAssetWithdrawResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Validate remote denom
	if !k.bankKeeper.HasSupply(ctx, msg.Withdraws[1].Balance.Denom) {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "invalid denom in local withdraw message: %s", msg.Withdraws[1].Balance.Denom)
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.Withdraws[1].Balance.Denom)
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

	// Save pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return &types.MsgMultiAssetWithdrawResponse{
		Tokens: stateChange.Out,
	}, nil
}

// OnSwapReceived processes a swap request and returns a response or an error.
func (k Keeper) OnSwapReceived(ctx sdk.Context, msg *types.MsgSwapRequest, stateChange *types.StateChange) (*types.MsgSwapResponse, error) {

	// Validate the message
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	poolID := types.GetPoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})
	pool, found := k.GetInterchainLiquidityPool(ctx, poolID)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	err := k.UnlockTokens(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*stateChange.Out[0]))
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

func (k Keeper) VerifySignature(ctx sdk.Context, sender authtypes.AccountI, deposit sdk.Coin, signature []byte) error {

	// Recover original signed Tx.
	depositMsg := types.DepositSignature{
		Sequence: sender.GetSequence(),
		Sender:   sender.GetAddress().String(),
		Balance:  &deposit,
	}

	rawDepositTx, err := types.ModuleCdc.Marshal(&depositMsg)

	if err != nil {
		return err
	}
	pubKey := sender.GetPubKey()
	isValid := pubKey.VerifySignature(rawDepositTx, signature)

	if !isValid {
		return errorsmod.Wrapf(types.ErrFailedDoubleDeposit, ":%s", types.ErrInvalidSignature)
	}
	return nil
}
