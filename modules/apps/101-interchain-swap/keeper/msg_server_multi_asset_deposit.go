package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) MultiAssetDeposit(ctx context.Context, msg *types.MsgMultiAssetDepositRequest) (*types.MsgMultiAssetDepositResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	}

	balance := k.bankKeeper.GetBalance(sdkCtx, sdk.MustAccAddressFromBech32(msg.LocalDeposit.Sender), msg.LocalDeposit.Token.Denom)

	if balance.Amount.LT(msg.LocalDeposit.Token.Amount) {
		return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrInEnoughAmount)
	}

	// Check initial deposit condition
	if pool.Status == types.PoolStatus_POOL_STATUS_INITIAL {
		for _, asset := range sdk.NewCoins(*msg.LocalDeposit.Token, *msg.RemoteDeposit.Token) {
			orderedAmount, err := pool.FindAssetByDenom(asset.Denom)
			if err != nil {
				return nil, err
			}
			if !orderedAmount.Balance.Amount.Equal(asset.Amount) {
				return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotAllowedAmount)
			}
		}
	} else {
		// Check the ratio of local amount and remote amount
		ratioOfTokensIn := msg.LocalDeposit.Token.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(msg.RemoteDeposit.Token.Amount)
		localAssetInPool, _ := pool.FindAssetByDenom(msg.LocalDeposit.Token.Denom)
		remoteAssetAmountInPool, _ := pool.FindAssetByDenom(msg.LocalDeposit.Token.Denom)
		ratioOfAssetsInPool := localAssetInPool.Balance.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(remoteAssetAmountInPool.Balance.Amount)
		if !ratioOfTokensIn.Equal(ratioOfAssetsInPool) {
			return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrInvalidPairRatio)
		}
	}

	// Create escrow module account here
	err = k.LockTokens(sdkCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.LocalDeposit.Sender), sdk.NewCoins(*msg.LocalDeposit.Token))

	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "due to %s", err)
	}

	fee := k.GetSwapFeeRate(sdkCtx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	poolTokens, err := amm.DepositMultiAsset([]*sdk.Coin{
		msg.LocalDeposit.Token,
		msg.RemoteDeposit.Token,
	})

	if err != nil {
		return nil, err
	}

	// Construct IBC packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.MULTI_DEPOSIT,
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

	err = k.SendIBCSwapPacket(sdkCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}
