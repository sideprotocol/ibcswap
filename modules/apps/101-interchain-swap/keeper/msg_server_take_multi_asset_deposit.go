package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) TakeMultiAssetDeposit(ctx context.Context, msg *types.MsgTakeMultiAssetDepositRequest) (*types.MsgMultiAssetDepositResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedMultiAssetDeposit, "%s", types.ErrNotFoundPool)
	}

	order, found := k.GetMultiDepositOrder(sdkCtx, msg.PoolId, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFoundMultiDepositOrder, "%s", types.ErrFailedMultiAssetDeposit)
	}

	if msg.Sender != order.DestinationTaker {
		return nil, errorsmod.Wrapf(types.ErrMultipleAssetDepositNotAllowed, "due to %s of other's", types.ErrFailedMultiAssetDeposit)
	}

	if order.Status == types.OrderStatus_COMPLETE {
		return nil, errorsmod.Wrapf(types.ErrAlreadyCompletedOrder, "due to %s of other's", types.ErrFailedMultiAssetDeposit)
	}

	// estimate pool token
	amm := types.NewInterchainMarketMaker(&pool)
	poolTokens, err := amm.DepositMultiAsset(sdk.Coins{
		*order.Deposits[0],
		*order.Deposits[1],
	})

	// check asset owned status
	asset := order.Deposits[1]
	if err != nil {
		return nil, errorsmod.Wrapf(err, "due to %s of other's", types.ErrFailedMultiAssetDeposit)
	}

	balance := k.bankKeeper.GetBalance(sdkCtx, sdk.MustAccAddressFromBech32(msg.Sender), asset.Denom)

	if balance.Amount.LT(asset.Amount) {
		return nil, errorsmod.Wrapf(types.ErrInEnoughAmount, "due to %s of Lp", types.ErrFailedMultiAssetDeposit)
	}

	// Create escrow module account here
	err = k.LockTokens(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*asset))

	if err != nil {
		return nil, errorsmod.Wrapf(err, "due to %s", types.ErrFailedMultiAssetDeposit)
	}

	// Construct IBC packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type:        types.TAKE_MULTI_DEPOSIT,
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

	_, err = k.SendIBCSwapPacket(sdkCtx, msg.Port, msg.Channel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	// Emit events
	poolTokenEventValues := []string{}
	for _, poolToken := range poolTokens {
		poolTokenEventValues = append(poolTokenEventValues, poolToken.String())
	}
	mintedPoolTokens := strings.Join(poolTokenEventValues, ":")
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeTakeMultiDepositOrder,
			sdk.Attribute{
				Key:   types.AttributeKeyPoolId,
				Value: msg.PoolId,
			},
			sdk.Attribute{
				Key:   types.AttributeKeyMultiDepositOrderId,
				Value: order.Id,
			},
			sdk.Attribute{
				Key:   types.AttributeKeyLpToken,
				Value: mintedPoolTokens,
			},
		))

	return &types.MsgMultiAssetDepositResponse{}, nil
}
