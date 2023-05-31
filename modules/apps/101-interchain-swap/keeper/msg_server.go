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

// OnCreatePoolAcknowledged processes the create pool acknowledgement, mints LP tokens, and saves the liquidity pool.
func (k Keeper) OnCreatePoolAcknowledged(ctx sdk.Context, msg *types.MsgCreatePoolRequest) error {
	// Save pool after completing the create operation in the counterparty chain
	pool := types.NewInterchainLiquidityPool(
		ctx,
		msg.Sender,
		msg.ChainId,
		k.bankKeeper,
		msg.Tokens,
		msg.Decimals,
		msg.Weight,
		msg.SourcePort,
		msg.SourceChannel,
	)

	// When adding initial liquidity, update the pool token amount
	initialLiquidity := sdk.NewCoin(pool.PoolId, msg.Tokens[0].Amount.Add(msg.Tokens[1].Amount))
	pool.AddPoolSupply(initialLiquidity)

	// Mint LP tokens
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(msg.Sender), initialLiquidity)
	if err != nil {
		return err
	}
	// calculate pool price
	// Instantiate an interchain market maker with the default fee rate
	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		pool,
		fee,
	)
	pool.PoolPrice = float32(amm.LpPrice())
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
	if pool.Status == types.PoolStatus_POOL_STATUS_READY {
		// Update pool supply and status
		pool.AddPoolSupply(*res.PoolToken)
		pool.AddAsset(*req.Token)
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
	}
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
	err := k.MintTokens(ctx, sdk.MustAccAddressFromBech32(req.LocalDeposit.Sender), *res.PoolTokens[0])

	if err != nil {
		return err
	}

	// Update pool supply and status
	for _, poolToken := range res.PoolTokens {
		pool.AddPoolSupply(*poolToken)
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_INITIAL {
		pool.AddAsset(*req.LocalDeposit.Token)
		pool.AddAsset(*req.RemoteDeposit.Token)
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
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
	for _, poolAsset := range res.Tokens {
		pool.SubtractAsset(*poolAsset)
	}
	pool.SubtractPoolSupply(*req.PoolCoin)

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

	// save pool
	k.SetInterchainLiquidityPool(ctx, pool)
	return nil
}

func (k Keeper) OnMultiWithdrawAcknowledged(ctx sdk.Context, req *types.MsgMultiAssetWithdrawRequest, res *types.MsgMultiAssetWithdrawResponse) error {

	pool, found := k.GetInterchainLiquidityPool(ctx, req.LocalWithdraw.PoolCoin.Denom)
	if !found {
		return types.ErrNotFoundPool
	}

	// update pool status
	for _, poolAsset := range res.Tokens {
		pool.SubtractAsset(*poolAsset)
	}
	pool.SubtractPoolSupply(*req.LocalWithdraw.PoolCoin)
	pool.SubtractPoolSupply(*req.RemoteWithdraw.PoolCoin)
	//burn voucher token.
	err := k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(req.LocalWithdraw.Sender), *req.LocalWithdraw.PoolCoin)
	if err != nil {
		return err
	}

	// unlock token
	err = k.UnlockTokens(ctx,
		pool.EncounterPartyPort,
		pool.EncounterPartyChannel,
		sdk.MustAccAddressFromBech32(req.LocalWithdraw.Sender),
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

	pooId := types.GetPoolIdWithTokens(msg.Tokens)
	_, found := k.GetInterchainLiquidityPool(ctx, pooId)

	if found {
		return nil, types.ErrAlreadyExistPool
	}

	// assume pool is ready when it is created.
	pool := *types.NewInterchainLiquidityPool(
		ctx,
		msg.Sender,
		msg.ChainId,
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

	// Instantiate an interchain market maker with the default fee rate
	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	lpToken, err := amm.DepositSingleAsset(*msg.Tokens[0])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedOnDepositReceived, "because of %s", err)
	}

	// pool status update
	pool.AddPoolSupply(*lpToken)

	// calculate
	pool.PoolPrice = float32(amm.LpPrice())
	// save pool status
	k.SetInterchainLiquidityPool(ctx, pool)
	return &pooId, nil
}

func (k Keeper) OnSingleAssetDepositReceived(ctx sdk.Context, msg *types.MsgSingleAssetDepositRequest, stateChange *types.StateChange) (*types.MsgSingleAssetDepositResponse, error) {

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrNotFoundPool
	}

	if pool.Status == types.PoolStatus_POOL_STATUS_READY {
		// increase lp token mint amount
		pool.AddPoolSupply(*stateChange.PoolTokens[0])

		// update pool tokens.
		pool.AddAsset(*msg.Token)
	} else {
		// switch pool status to 'READY'
		pool.Status = types.PoolStatus_POOL_STATUS_READY
	}

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
	senderAcc := k.authKeeper.GetAccount(ctx, sdk.MustAccAddressFromBech32(msg.RemoteDeposit.Sender))
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
	escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)

	// Create a deposit message
	sendMsg := banktypes.MsgSend{
		FromAddress: senderAcc.GetAddress().String(),
		ToAddress:   escrowAccount.String(),
		Amount:      sdk.NewCoins(*msg.RemoteDeposit.Token),
	}

	// Recover original signed Tx.
	deposit := types.RemoteDeposit{
		Sequence: senderAcc.GetSequence(),
		Sender:   msg.RemoteDeposit.Sender,
		Token:    msg.RemoteDeposit.Token,
	}
	rawDepositTx, err := types.ModuleCdc.Marshal(&deposit)

	if err != nil {
		return nil, err
	}

	pubKey := senderAcc.GetPubKey()
	isValid := pubKey.VerifySignature(rawDepositTx, msg.RemoteDeposit.Signature)

	if !isValid {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, ":%s", types.ErrInvalidSignature)
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
	if pool.Status == types.PoolStatus_POOL_STATUS_READY {
		pool.AddAsset(*msg.LocalDeposit.Token)
		pool.AddAsset(*msg.RemoteDeposit.Token)
	} else {
		pool.Status = types.PoolStatus_POOL_STATUS_READY
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

// OnSingleAssetWithdrawReceived processes a withdrawal request and returns a response or an error.
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

	// Update pool status by subtracting the supplied pool coin and output token
	for _, poolAsset := range stateChange.Out {
		pool.SubtractAsset(*poolAsset)
	}

	for _, poolToken := range stateChange.PoolTokens {
		pool.SubtractPoolSupply(*poolToken)
	}

	// Save pool
	k.SetInterchainLiquidityPool(ctx, pool)
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
	if !k.bankKeeper.HasSupply(ctx, msg.LocalWithdraw.DenomOut) {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "invalid denom in local withdraw message: %s", msg.LocalWithdraw.DenomOut)
	}

	// Retrieve the liquidity pool
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.RemoteWithdraw.PoolCoin.Denom)
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

	err := k.UnlockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Recipient), sdk.NewCoins(*stateChange.Out[0]))
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
