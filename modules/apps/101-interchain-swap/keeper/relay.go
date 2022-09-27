package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (k Keeper) SendIBCSwapPacket(
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

	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
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

	packet := channeltypes.NewPacket(
		swapPacket.GetBytes(),
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	if err := k.ics4Wrapper.SendPacket(ctx, channelCap, packet); err != nil {
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

func (k Keeper) OnCreatePool(ctx sdk.Context, request *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) OnDeposit(ctx sdk.Context, response *types.MsgCreatePoolResponse) (*types.MsgDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) OnWithdraw(ctx sdk.Context, request *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) OnLeftSwap(ctx sdk.Context, request *types.MsgLeftSwapRequest) (*types.MsgSwapResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) OnRightSwap(ctx sdk.Context, request *types.MsgRightSwapRequest) (*types.MsgSwapResponse, error) {
	//TODO implement me
	panic("implement me")
}
