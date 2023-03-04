package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	//sdk "github.com/cosmos/cosmos-sdk/types"
	//channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

var (
	StepSend            = 1
	StepReceive         = 2
	StepAcknowledgement = 3
)

var _ types.MsgServer = Keeper{}

func (k Keeper) DelegateCreatePool(goctx context.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	_, err1 := sdk.AccAddressFromBech32(msg.Sender)
	if err1 != nil {
		return nil, err1
	}

	pool := types.NewBalancerLiquidityPool(msg.Denoms, msg.Decimals, msg.Weight)
	if err := pool.Validate(); err != nil {
		return nil, err
	}

	// count native tokens
	count := 0
	for _, denom := range msg.Denoms {
		if k.bankKeeper.HasSupply(ctx, denom) {
			count += 1
			pool.UpdateAssetPoolSide(denom, types.PoolSide_POOL_SIDE_NATIVE_ASSET)
		} else {
			pool.UpdateAssetPoolSide(denom, types.PoolSide_POOL_SIDE_REMOTE_ASSET)
		}
	}
	if count == 0 {
		return nil, types.ErrNoNativeTokenInPool
	}

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	packet := types.NewIBCSwapPacketData(types.CREATE_POOL, msgByte, nil)
	if err := k.SendIBCSwapDelegationDataPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgCreatePoolResponse{}, nil
}

func (k Keeper) DelegateSingleDeposit(goctx context.Context, msg *types.MsgSingleDepositRequest) (*types.MsgSingleDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.Sender)
	if err1 != nil {
		return nil, err1
	}

	// deposit assets to the swap module
	length := len(msg.Tokens)
	var coins = make([]sdk.Coin, length)
	for i := 0; i < length; i++ {
		coins[i] = *msg.Tokens[i]
	}
	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
	k.bankKeeper.SendCoins(ctx, sender, escrowAddr, sdk.NewCoins(coins...))

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	packet := types.NewIBCSwapPacketData(types.SINGLE_DEPOSIT, msgByte, nil)
	if err := k.SendIBCSwapDelegationDataPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	// k.bankKeeper.MintCoins()
	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgSingleDepositResponse{}, nil
}

func (k Keeper) DelegateWithdraw(ctx2 context.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	ctx := sdk.UnwrapSDKContext(ctx2)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.Sender)
	if err1 != nil {
		return nil, err1
	}

	// deposit assets to the swap module
	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
	k.bankKeeper.SendCoins(ctx, sender, escrowAddr, sdk.NewCoins(*msg.PoolToken))

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	packet := types.NewIBCSwapPacketData(types.WITHDRAW, msgByte, nil)
	if err := k.SendIBCSwapDelegationDataPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgWithdrawResponse{}, nil
}

func (k Keeper) DelegateLeftSwap(goctx context.Context, msg *types.MsgLeftSwapRequest) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.Sender)
	if err1 != nil {
		return nil, err1
	}

	// deposit assets to the swap module
	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
	k.bankKeeper.SendCoins(ctx, sender, escrowAddr, sdk.NewCoins(*msg.TokenIn))

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	packet := types.NewIBCSwapPacketData(types.LEFT_SWAP, msgByte, nil)
	if err := k.SendIBCSwapDelegationDataPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgSwapResponse{}, nil
}

func (k Keeper) DelegateRightSwap(goctx context.Context, msg *types.MsgRightSwapRequest) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goctx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.Sender)
	if err1 != nil {
		return nil, err1
	}

	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)
	k.bankKeeper.SendCoins(ctx, sender, escrowAddr, sdk.NewCoins(*msg.TokenIn))

	msgByte, err0 := types.ModuleCdc.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	packet := types.NewIBCSwapPacketData(types.RIGHT_SWAP, msgByte, nil)
	if err := k.SendIBCSwapDelegationDataPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgSwapResponse{}, nil
}
