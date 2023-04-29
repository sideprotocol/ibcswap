package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) MultiAssetWithdraw(goCtx context.Context, msg *types.MsgMultiAssetWithdrawRequest) (*types.MsgMultiAssetWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// check out denom
	if !k.bankKeeper.HasSupply(ctx, msg.LocalWithdraw.DenomOut) {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "invalid denom in local withdraw message:%s", msg.LocalWithdraw.DenomOut)
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.LocalWithdraw.PoolCoin.Denom)

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

	localOut, err := amm.Withdraw(*msg.LocalWithdraw.PoolCoin, msg.LocalWithdraw.DenomOut)

	if err != nil {
		return nil, err
	}

	remoteOut, err := amm.Withdraw(*msg.LocalWithdraw.PoolCoin, msg.RemoteWithdraw.DenomOut)

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
			Out:        []*sdk.Coin{localOut, remoteOut},
			PoolTokens: []*sdk.Coin{msg.LocalWithdraw.PoolCoin, msg.RemoteWithdraw.PoolCoin},
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)

	err = k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}
	return &types.MsgMultiAssetWithdrawResponse{}, nil
}
