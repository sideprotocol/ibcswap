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
	for _, token := range msg.Tokens {
		if k.bankKeeper.HasSupply(ctx, token.Denom) {
			localAssetCount += 1
		}
	}

	// should have 1 native asset on the chain
	if localAssetCount < 1 {
		return nil, types.ErrNumberOfLocalAsset
	}

	// check user owned initial liquidity or not
	holdingNativeCoin := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.Sender), msg.Tokens[0].Denom)
	if holdingNativeCoin.Amount.LT(msg.Tokens[0].Amount) {
		return nil, types.ErrEmptyInitialLiquidity
	}

	// move initial fund to liquidity pool
	err = k.LockTokens(ctx, msg.SourcePort, msg.SourceChannel, sdk.MustAccAddressFromBech32(msg.Sender), sdk.NewCoins(*msg.Tokens[0]))

	if err != nil {
		return nil, err
	}

	poolData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	// construct IBC data packet
	packet := types.IBCSwapPacketData{
		Type: types.CREATE_POOL,
		Data: poolData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}
	poolId := types.GetPoolIdWithTokens(msg.Tokens)
	return &types.MsgCreatePoolResponse{
		PoolId: poolId,
	}, nil
}
