package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errormod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) SingleAssetDeposit(ctx context.Context, msg *types.MsgSingleAssetDepositRequest) (*types.MsgSingleAssetDepositResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errormod.Wrapf(types.ErrFailedDeposit, "%s", types.ErrNotFoundPool)
	}

	// Deposit token to Escrow account
	accAddress, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedDeposit, "%s", types.ErrInvalidAddress)
	}
	balance := k.bankKeeper.GetBalance(sdkCtx, accAddress, msg.Token.Denom)
	if balance.Amount.Equal(sdk.NewInt(0)) {
		return nil, types.ErrInvalidAmount
	}

	if pool.Status == types.PoolStatus_POOL_STATUS_INITIAL {
		poolAsset, err := pool.FindAssetByDenom(msg.Token.Denom)
		if err == nil {
			if !poolAsset.Balance.Amount.Equal(msg.Token.Amount) {
				return nil, types.ErrInvalidInitialDeposit
			}
		}
	}

	// Deposit assets to the escrowed account
	err = k.LockTokens(sdkCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*msg.Token))
	if err != nil {
		return nil, errormod.Wrapf(types.ErrFailedDeposit, "%s", err)
	}

	fee := k.GetSwapFeeRate(sdkCtx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	poolToken, err := amm.DepositSingleAsset(*msg.Token)
	if err != nil {
		return nil, err
	}

	// Construct IBC packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.SINGLE_DEPOSIT,
		Data:        rawMsgData,
		StateChange: &types.StateChange{PoolTokens: []*sdk.Coin{poolToken}},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)
	err = k.SendIBCSwapPacket(sdkCtx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgSingleAssetDepositResponse{
		PoolToken: pool.Supply,
	}, nil
}
