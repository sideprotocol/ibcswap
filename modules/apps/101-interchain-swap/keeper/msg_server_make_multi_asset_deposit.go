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
	sourceAsset, err := pool.FindAssetBySide(types.PoolAssetSide_SOURCE)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrNotFoundDenomInPool, "%s", types.ErrFailedMultiAssetDeposit)
	}
	destinationAsset, err := pool.FindAssetBySide(types.PoolAssetSide_DESTINATION)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrNotFoundDenomInPool, "%s:", types.ErrFailedMultiAssetDeposit)
	}

	currentRatio := sourceAsset.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(destinationAsset.Amount)
	inputRatio := msg.Deposits[0].Balance.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(msg.Deposits[1].Balance.Amount)

	if err := types.CheckSlippage(currentRatio, inputRatio, 10); err != nil {
		return nil, errormod.Wrapf(types.ErrInvalidPairRatio, "%s", types.ErrFailedMultiAssetDeposit)
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
	k.AppendMultiDepositOrder(sdkCtx, pool.Id, order)

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
