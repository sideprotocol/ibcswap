package keeper

import (
	"context"


	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) MultiAssetDeposit(goCtx context.Context, msg *types.MsgMultiAssetDepositRequest) (*types.MsgMultiAssetDepositResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	// validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	}

	balance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.LocalDeposit.Sender), msg.LocalDeposit.Token.Denom)

	if balance.Amount.LT(msg.LocalDeposit.Token.Amount) {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrInEnoughAmount)
	}

	// check initial deposit condition
	if pool.Status == types.PoolStatus_POOL_STATUS_INITIAL {
		for _, asset := range sdk.NewCoins(*msg.LocalDeposit.Token, *msg.RemoteDeposit.Token) {
			orderedAmount, err := pool.FindAssetByDenom(asset.Denom)
			if err != nil {
				return nil, err
			}
			if !orderedAmount.Balance.Amount.Equal(asset.Amount) {
				return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotAllowedAmount)
			}
		}
	} else {
		//check the ratio of local amount and remote amount
		ratioOfTokensIn := msg.LocalDeposit.Token.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(msg.RemoteDeposit.Token.Amount)
		localAssetInPool, _ := pool.FindAssetByDenom(msg.LocalDeposit.Token.Denom)
		remoteAssetAmountInPool, _ := pool.FindAssetByDenom(msg.LocalDeposit.Token.Denom)
		ratioOfAssetsInPool := localAssetInPool.Balance.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(remoteAssetAmountInPool.Balance.Amount)
		if !ratioOfTokensIn.Equal(ratioOfAssetsInPool) {
			return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrInvalidPairRatio)
		}
	}

	// create escrow module account  here
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.LocalDeposit.Sender), sdk.NewCoins(*msg.LocalDeposit.Token))

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "because of %s", err)
	}

	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	poolTokens, err := amm.DepositDoubleAsset([]*sdk.Coin{
		msg.LocalDeposit.Token,
		msg.RemoteDeposit.Token,
	})

	if err != nil {
		return nil, err
	}

	// construct ibc packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.DOUBLE_DEPOSIT,
		Data:        rawMsgData,
		StateChange: &types.StateChange{PoolTokens: poolTokens},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgMultiAssetDepositResponse{
		PoolTokens: poolTokens,
	}, nil
}
