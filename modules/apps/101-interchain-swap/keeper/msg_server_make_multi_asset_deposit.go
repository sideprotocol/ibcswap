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

	// Check initial deposit condition
	if pool.Status != types.PoolStatus_ACTIVE {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotReadyForSwap)
	}

	// Check input ration of tokens
	sourceAsset, err := pool.FindAssetByDenom(msg.Deposits[0].Balance.Denom)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrNotFoundDenomInPool, "%s", types.ErrFailedMultiAssetDeposit)
	}

	destinationAsset, err := pool.FindAssetByDenom(msg.Deposits[1].Balance.Denom)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrNotFoundDenomInPool, "%s:", types.ErrFailedMultiAssetDeposit)
	}

	currentRatio := sdk.NewDecFromInt(sourceAsset.Balance.Amount).Quo(sdk.NewDecFromInt(destinationAsset.Balance.Amount))
	inputRatio := sdk.NewDecFromInt(msg.Deposits[0].Balance.Amount).Quo(sdk.NewDecFromInt(msg.Deposits[1].Balance.Amount))

	if err := types.CheckSlippage(currentRatio, inputRatio, 10); err != nil {
		return nil, errormod.Wrapf(types.ErrInvalidPairRatio, "%d:%d:%s", currentRatio, inputRatio, types.ErrFailedMultiAssetDeposit)
	}

	// Create escrow module account here
	err = k.LockTokens(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender), sdk.NewCoins(*msg.Deposits[0].Balance))

	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedMultiAssetDeposit, "due to %s", err)
	}

	amm := *types.NewInterchainMarketMaker(
		&pool,
	)

	poolTokens, err := amm.DepositMultiAsset(sdk.Coins{
		*msg.Deposits[0].Balance,
		*msg.Deposits[1].Balance,
	})

	if err != nil {
		return nil, err
	}

	// create order
	order := types.MultiAssetDepositOrder{
		PoolId:           msg.PoolId,
		ChainId:          sdkCtx.ChainID(),
		SourceMaker:      msg.Deposits[0].Sender,
		DestinationTaker: msg.Deposits[1].Sender,
		Deposits:         types.GetCoinsFromDepositAssets(msg.Deposits),
		Status:           types.OrderStatus_PENDING,
		CreatedAt:        sdkCtx.BlockHeight(),
	}

	// save order in source chain
	//set orderId
	creator := sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender)
	acc := k.authKeeper.GetAccount(sdkCtx, creator)
	orderId := types.GetOrderId(order.SourceMaker, acc.GetSequence())
	order.Id = orderId

	k.SetMultiDepositOrder(sdkCtx, order)
	// Construct IBC packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.MAKE_MULTI_DEPOSIT,
		Data: rawMsgData,
		StateChange: &types.StateChange{
			MultiDepositOrderId: order.Id,
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)
	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	_, err = k.SendIBCSwapPacket(sdkCtx, msg.Port, msg.Channel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMakeMultiDepositOrder,
			sdk.Attribute{
				Key:   types.AttributeKeyPoolId,
				Value: msg.PoolId,
			},
			sdk.Attribute{
				Key:   types.AttributeKeyMultiDepositOrderId,
				Value: orderId,
			},
		))

	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}
