package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) CreatePool(goCtx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// validate message
	err := host.PortIdentifierValidator(msg.SourcePort)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	err = host.ChannelIdentifierValidator(msg.SourceChannel)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	err = msg.ValidateBasic()
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedSwap, "because of %s", err)
	}

	localAssetCount := 0
	for _, denom := range msg.Denoms {
		if k.bankKeeper.HasSupply(ctx, denom) {
			localAssetCount += 1
		}
	}

	// should have 1 native asset on the chain
	if localAssetCount < 1 {
		return nil, types.ErrNumberOfLocalAsset
	}

	// check user owned initial liquidity or not
	holdingNativeCoin := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Denoms[0])
	if holdingNativeCoin.Amount.LT(sdk.NewInt(int64(msg.InitalLiquidity))) {
		return nil, types.ErrEmptyInitialLiquidity
	}
	lockedNativeCoin := sdk.NewCoin(msg.Denoms[0], sdk.NewInt(int64(msg.InitalLiquidity)))

	// move initial fund to liquidity pool
	err = k.LockTokens(ctx, msg.SourcePort, msg.SourceChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(lockedNativeCoin))

	if err != nil {
		return nil, err
	}

	poolData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// construct IBC data packet
	packet := types.IBCSwapDataPacket{
		Type: types.MessageType_CREATE,
		Data: poolData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}
	poolId := types.GetPoolId(msg.Denoms)
	return &types.MsgCreatePoolResponse{
		PoolId: poolId,
	}, nil
}
