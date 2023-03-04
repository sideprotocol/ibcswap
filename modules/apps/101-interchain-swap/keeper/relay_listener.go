package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) CreateIBCSwapAMM(ctx sdk.Context, poolId string) (types.BalancerAMM, error) {
	pool, ok := k.GetBalancerPool(ctx, poolId)
	if !ok {
		return types.BalancerAMM{}, types.ErrInvalidPoolId
	}

	params := k.GetParams(ctx)

	amm := types.NewBalanceAMM(&pool, int64(params.MaxFeeRate))
	return amm, nil

}

func (k Keeper) SendIBCSwapDelegationDataPacket(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
	swapPacket types.IBCSwapPacketData,
) error {

	if err := swapPacket.ValidateBasic(); err != nil {
		return err
	}

	if !k.GetSwapEnabled(ctx) {
		return types.ErrSendDisabled
	}

	_, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	//destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	//destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	_, found2 := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found2 {
		return sdkerrors.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	// begin createOutgoingPacket logic
	// See spec for this logic: https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	//packet := channeltypes.NewPacket(
	//	swapPacket.GetBytes(),
	//	sequence,
	//	sourcePort,
	//	sourceChannel,
	//	destinationPort,
	//	destinationChannel,
	//	timeoutHeight,
	//	timeoutTimestamp,
	//)

	_, err := k.ics4Wrapper.SendPacket(ctx, channelCap, sourcePort, sourceChannel, timeoutHeight, timeoutTimestamp, swapPacket.GetBytes())
	if err != nil {
		return err
	}

	defer func() {
		//if sendingCoin.Amount.IsInt64() {
		//	telemetry.SetGaugeWithLabels(
		//		[]string{"tx", "msg", "ibc", "swap"},
		//		float32(sendingCoin.Amount.Int64()),
		//		[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, "fullDenomPath")},
		//	)
		//}
	}()

	return nil
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data *types.IBCSwapPacketData) error {

	ctx.Logger().Debug("refundPacketToken: %s", data)

	return nil
}

func (k Keeper) OnCreatePoolReceived(ctx sdk.Context, msg *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {

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

	k.SetBalancerPool(ctx, pool)

	return &types.MsgCreatePoolResponse{
		PoolId: pool.Id,
	}, nil

}

func (k Keeper) OnSingleDepositReceived(ctx sdk.Context, msg *types.MsgSingleDepositRequest) (*types.MsgSingleDepositResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	amm, err := k.CreateIBCSwapAMM(ctx, msg.PoolId)
	if err != nil {
		return nil, err
	}

	poolToken, err := amm.Deposit(msg.Tokens)
	if err != nil {
		return nil, err
	}

	k.SetBalancerPool(ctx, *amm.Pool) // update pool states

	return &types.MsgSingleDepositResponse{
		PoolToken: &poolToken,
	}, nil
}

func (k Keeper) OnWithdrawReceived(ctx sdk.Context, msg *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	amm, err := k.CreateIBCSwapAMM(ctx, msg.PoolToken.Denom) // Pool Token denomination is the pool Id
	if err != nil {
		return nil, err
	}

	outToken, err := amm.Withdraw(msg.PoolToken, msg.DenomOut)
	if err != nil {
		return nil, err
	}

	k.SetBalancerPool(ctx, *amm.Pool) // update pool states

	// only output one asset in the pool
	return &types.MsgWithdrawResponse{Token: &outToken}, nil
}

func (k Keeper) OnLeftSwapReceived(ctx sdk.Context, msg *types.MsgLeftSwapRequest) (*types.MsgSwapResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	poolId := types.GeneratePoolId([]string{msg.TokenIn.Denom, msg.DenomOut})

	amm, err := k.CreateIBCSwapAMM(ctx, poolId) // Pool Token denomination is the pool Id
	if err != nil {
		return nil, err
	}

	outToken, err := amm.LeftSwap(msg.TokenIn, msg.DenomOut)
	if err != nil {
		return nil, err
	}

	k.SetBalancerPool(ctx, *amm.Pool) // update pool states

	// only output one asset in the pool
	return &types.MsgSwapResponse{
		TokenOut: &outToken,
	}, nil
}

func (k Keeper) OnRightSwapReceived(ctx sdk.Context, msg *types.MsgRightSwapRequest) (*types.MsgSwapResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	poolId := types.GeneratePoolId([]string{msg.TokenIn.Denom, msg.TokenOut.Denom})

	amm, err := k.CreateIBCSwapAMM(ctx, poolId) // Pool Token denomination is the pool Id
	if err != nil {
		return nil, err
	}

	outToken, err := amm.RightSwap(msg.TokenIn, msg.TokenOut)
	if err != nil {
		return nil, err
	}

	k.SetBalancerPool(ctx, *amm.Pool) // update pool states

	// only output one asset in the pool
	return &types.MsgSwapResponse{
		TokenOut: &outToken,
	}, nil
}

func (k Keeper) OnCreatePoolAcknowledged(ctx sdk.Context, request *types.MsgCreatePoolRequest, response *types.MsgCreatePoolResponse) error {
	//TODO implement me
	// do nothing
	return nil
}

func (k Keeper) OnSingleDepositAcknowledged(ctx sdk.Context, request *types.MsgSingleDepositRequest, response *types.MsgSingleDepositResponse) error {
	if err := request.ValidateBasic(); err != nil {
		return err
	}

	if pool, ok := k.GetBalancerPool(ctx, request.PoolId); !ok {
		return types.ErrInvalidPoolId
	} else {
		pool.PoolToken.Add(*response.PoolToken)
		length := len(request.Tokens)
		for i := 0; i < length; i++ {
			pool.UpdateBalance(*request.Tokens[i], true)
		}
		k.SetBalancerPool(ctx, pool) // update pool states
	}

	receiver := sdk.MustAccAddressFromBech32(request.Sender)
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(*response.PoolToken)); err != nil {
		return err
	}
	return k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, receiver, sdk.NewCoins(*response.PoolToken))
}

func (k Keeper) OnWithdrawAcknowledged(ctx sdk.Context, request *types.MsgWithdrawRequest, response *types.MsgWithdrawResponse) error {
	sender := sdk.MustAccAddressFromBech32(request.Sender)
	escrowAddr := types.GetEscrowAddress(request.SourcePort, request.SourceChannel)
	if pool, ok := k.GetBalancerPool(ctx, request.PoolToken.Denom); !ok {
		return types.ErrInvalidPoolId
	} else {
		pool.PoolToken.Sub(*request.PoolToken)
		pool.UpdateBalance(*response.Token, false)
		k.SetBalancerPool(ctx, pool) // update pool states
	}
	// currently, only one asset is allowed to withdraw, which is what you deposited, what you can withdraw
	if err := k.bankKeeper.SendCoins(ctx, escrowAddr, sender, sdk.NewCoins(*response.Token)); err != nil {
		return err
	}
	// burn pool token deposited of the request message
	// TODO check if we need transfer from escrowAddress first before burn
	return k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(*request.PoolToken))
}

func (k Keeper) OnLeftSwapAcknowledged(ctx sdk.Context,
	packet channeltypes.Packet, request *types.MsgLeftSwapRequest, response *types.MsgSwapResponse) error {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) OnRightSwapAcknowledged(ctx sdk.Context,
	packet channeltypes.Packet, request *types.MsgRightSwapRequest, response *types.MsgSwapResponse) error {
	//TODO implement me
	panic("implement me")
}
