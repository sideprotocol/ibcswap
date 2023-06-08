package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) MakeMultiAssetDeposit(ctx context.Context, msg *types.MsgMakeMultiAssetDepositRequest) (*types.MsgMultiAssetDepositResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	// check asset owned status
	balance := k.bankKeeper.GetBalance(sdkCtx, sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender), msg.Deposits[0].Balance.Denom)
	if balance.Amount.LT(msg.Deposits[0].Balance.Amount) {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrInEnoughAmount)
	}

	// Check initial deposit condition
	if pool.Status != types.PoolStatus_ACTIVE {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotReadyForSwap)
	}

	// Create escrow module account here
	err = k.LockTokens(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender), sdk.NewCoins(*msg.Deposits[0].Balance))

	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "due to %s", err)
	}

	amm := *types.NewInterchainMarketMaker(
		&pool,
	)

	poolTokens, err := amm.DepositMultiAsset([]*sdk.Coin{
		msg.Deposits[0].Balance,
		msg.Deposits[1].Balance,
	})

	if err != nil {
		return nil, err
	}

	// create order
	orderId := types.GetOrderId(sdkCtx.ChainID())

	order := types.MultiAssetDepositOrder{
		Id:                      orderId,
		PoolId:                  msg.PoolId,
		ChainId:                 sdkCtx.ChainID(),
		SourceMaker:             msg.Deposits[0].Sender,
		DestinationTaker:        msg.Deposits[1].Sender,
		SourceMakerLpToken:      poolTokens[0],
		DestinationTakerLpToken: poolTokens[1],
	}

	// save order in source chain
	k.SetMultiDepositOrder(sdkCtx, order)

	// Construct IBC packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.MAKE_MULTI_DEPOSIT,
		Data:        rawMsgData,
		StateChange: &types.StateChange{PoolTokens: poolTokens},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)
	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	err = k.SendIBCSwapPacket(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}
