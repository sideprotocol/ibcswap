package keeper

import (
	"context"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) TakeMultiAssetDeposit(ctx context.Context, msg *types.MsgTakeMultiAssetDepositRequest) (*types.MsgMultiAssetDepositResponse, error) {

	//sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// pool, found := k.GetInterchainLiquidityPool(sdkCtx, msg.PoolId)
	// if !found {
	// 	return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	// }

	// // check asset owned status
	// balance := k.bankKeeper.GetBalance(sdkCtx, sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender), msg.Deposits[0].Balance.Denom)

	// // Check initial deposit condition
	// if pool.Status != types.PoolStatus_ACTIVE {
	// 	return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotReadyForSwap)
	// }

	// // Create escrow module account here
	// err = k.LockTokens(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Deposits[0].Sender), sdk.NewCoins(*msg.Deposits[0].Balance))

	// if err != nil {
	// 	return nil, errormod.Wrapf(types.ErrFailedDoubleDeposit, "due to %s", err)
	// }

	// amm := *types.NewInterchainMarketMaker(
	// 	&pool,
	// )

	// if err != nil {
	// 	return nil, err
	// }

	// // Construct IBC packet
	// rawMsgData, err := types.ModuleCdc.Marshal(msg)
	// if err != nil {
	// 	return nil, err
	// }

	// packet := types.IBCSwapPacketData{
	// 	Type:        types.MULTI_DEPOSIT,
	// 	Data:        rawMsgData,
	// 	StateChange: &types.StateChange{PoolTokens: poolTokens},
	// }

	// timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&sdkCtx)

	// // Use input timeoutHeight, timeoutStamp
	// if msg.TimeoutHeight != nil {
	// 	timeoutHeight = *msg.TimeoutHeight
	// }
	// if msg.TimeoutTimeStamp != 0 {
	// 	timeoutStamp = msg.TimeoutTimeStamp
	// }

	// err = k.SendIBCSwapPacket(sdkCtx, pool.CounterPartyPort, pool.CounterPartyChannel, timeoutHeight, timeoutStamp, packet)
	// if err != nil {
	// 	return nil, err
	// }

	// return &types.MsgMultiAssetDepositResponse{
	// 	PoolTokens: poolTokens,
	// }, nil
	return nil, nil
}
