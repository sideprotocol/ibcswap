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

	// if pool.Status != types.PoolStatus_POOL_STATUS_READY {
	// 	return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotReadyForSwap)
	// }

	outToken, err := pool.FindAssetByDenom(msg.DenomOut)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s in pool", types.ErrEmptyDenom)
	}

	// validate asset
	if outToken.Side != types.PoolSide_NATIVE {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotNativeDenom)
	}

	// lock pool token to the swap module
	err = k.LockTokens(
		ctx,
		pool.EncounterPartyPort,
		pool.EncounterPartyChannel,
		sdk.MustAccAddressFromBech32(msg.Sender),
		sdk.NewCoins(*msg.PoolCoin),
	)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", err)
	}

	// construct the IBC data packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)

	k.SendIBCSwapPacket(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), packet)
	return &types.MsgWithdrawResponse{}, nil
}
