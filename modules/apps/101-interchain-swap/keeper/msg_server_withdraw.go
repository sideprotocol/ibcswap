package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()

	if err != nil {
		return nil, err
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotFoundPool)
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_READY {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotReadyForSwap)
	}


	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", err)
	}

	fee := k.GetSwapFeeRate(ctx)
	amm := *types.NewInterchainMarketMaker(
		&pool,
		fee,
	)

	out, err := amm.Withdraw(*msg.PoolCoin)
	if err != nil {
		return nil, err
	}

	// construct the IBC data packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.WITHDRAW,
		Data: rawMsgData,
		StateChange: &types.StateChange{
			Out:        out,
			PoolTokens: []*sdk.Coin{msg.PoolCoin},
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)

	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}
	return &types.MsgWithdrawResponse{}, nil
}
