package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) MultiAssetWithdraw(goCtx context.Context, msg *types.MsgMultiAssetWithdrawRequest) (*types.MsgMultiAssetWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// check out denom
	if !k.bankKeeper.HasSupply(ctx, msg.PoolToken.Denom) {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "invalid denom in local withdraw message:%s", msg.PoolToken.Denom)
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolToken.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotFoundPool)
	}

	amm := *types.NewInterchainMarketMaker(
		&pool,
	)

	srcDenom, _ := pool.FindDenomBySide(types.PoolAssetSide_SOURCE)
	srcOut, err := amm.MultiAssetWithdraw(sdk.Coin{
		Denom: pool.Id, Amount: msg.PoolToken.Amount.Quo(sdk.NewInt(2)),
	}, *srcDenom)
	if err != nil {
		return nil, err
	}
	targetDenom, _ := pool.FindDenomBySide(types.PoolAssetSide_DESTINATION)
	targetOut, err := amm.MultiAssetWithdraw(sdk.Coin{
		Denom: pool.Id, Amount: msg.PoolToken.Amount.Quo(sdk.NewInt(2)),
	}, *targetDenom)

	if err != nil {
		return nil, err
	}

	// construct the IBC data packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.MULTI_WITHDRAW,
		Data: rawMsgData,
		StateChange: &types.StateChange{
			Out:        []*sdk.Coin{srcOut, targetOut},
			PoolTokens: []*sdk.Coin{msg.PoolToken},
		},
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	err = k.SendIBCSwapPacket(ctx, pool.CounterPartyPort, pool.CounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}
	return &types.MsgMultiAssetWithdrawResponse{}, nil
}
