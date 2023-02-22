package keeper

import (
	"context"
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Withdraw(goCtx context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	err := msg.ValidateBasic()

	if err != nil {
		return nil, err
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolCoin.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrNotFoundPool, "failed to withdraw because of %s")
	}

	if pool.Status != types.PoolStatus_POOL_STATUS_READY {
		return nil, errorsmod.Wrapf(types.ErrNotReadyForSwap, "failed to withdraw because of %s")
	}

	asset, err := pool.FindAssetByDenom(msg.DenomOut)

	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrEmptyDenom, "failed to withdraw because of %s in pool")
	}

	// validate asset
	if asset.Side != types.PoolSide_NATIVE {
		return nil, errorsmod.Wrapf(types.ErrNotNativeDenom, "failed to withdraw because of %s")
	}

	// lock pool token to the swap module
	escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)
	k.Keeper.bankKeeper.SendCoinsFromModuleToAccount(ctx, string(escrowAccount), sdk.AccAddress(msg.Sender), sdk.NewCoins(*msg.PoolCoin))

	// construct the IBC data packet
	rawMsgData, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	_ = types.IBCSwapDataPacket{
		Type: types.MessageType_DEPOSIT,
		Data: rawMsgData,
	}
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(pool.EncounterPartyPort, pool.EncounterPartyChannel))
	if !ok {
		return nil, nil
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut()

	k.channelKeeper.SendPacket(ctx, channelCap, pool.EncounterPartyPort, pool.EncounterPartyChannel, timeoutHeight, uint64(timeoutStamp), rawMsgData)
	return &types.MsgWithdrawResponse{}, nil
}
