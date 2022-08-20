package keeper

import (
	//"strings"

	metrics "github.com/armon/go-metrics"
	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v4/modules/core/24-host"
	coretypes "github.com/cosmos/ibc-go/v4/modules/core/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/31-atomic-swap/types"
)

func (k Keeper) SendSwap(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	sendingCoin, receivingCoin sdk.Coin,
	sender sdk.AccAddress,
	senderReceivingAddress string,
	expectedCounterPartyAddress string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) error {
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

	packetData := types.NewAtomicSwapPacketData(
		sendingCoin.Denom,
		sendingCoin.Amount.String(),
		receivingCoin.Denom,
		receivingCoin.Amount.String(),
		expectedCounterPartyAddress,
		sender.String(),
		senderReceivingAddress,
	)

	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
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
		if sendingCoin.Amount.IsInt64() {
			telemetry.SetGaugeWithLabels(
				[]string{"tx", "msg", "ibc", "swap"},
				float32(sendingCoin.Amount.Int64()),
				[]metrics.Label{telemetry.NewLabel(coretypes.LabelDenom, "fullDenomPath")},
			)
		}
	}()

	return nil
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SwapPacketData) error {
	// validate packet data upon receiving
	if err := data.ValidateBasic(); err != nil {
		return err
	}

	if !k.GetSwapEnabled(ctx) {
		return types.ErrReceiveDisabled
	}

	// decode the receiver address
	receiver, err := sdk.AccAddressFromBech32(data.SenderReceivingAddress)
	if err != nil {
		return err
	}

	// parse the transfer amount
	transferAmount, ok := sdk.NewIntFromString(data.ReceivingTokenAmount)
	if !ok {
		return sdkerrors.Wrapf(types.ErrInvalidAmount, "unable to parse transfer amount (%s) into math.Int", data.ReceivingTokenAmount)
	}

	labels := []metrics.Label{
		telemetry.NewLabel(coretypes.LabelSourcePort, packet.GetSourcePort()),
		telemetry.NewLabel(coretypes.LabelSourceChannel, packet.GetSourceChannel()),
	}

	ctx.Logger().Info("Received Message: %s, %s, %s", receiver, transferAmount, labels)

	return nil
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SwapPacketData, ack channeltypes.Acknowledgement) error {
	switch ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return k.refundPacketToken(ctx, packet, data)
	default:
		// the acknowledgement succeeded on the receiving chain so nothing
		// needs to be executed and no error needs to be returned
		return nil
	}
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.SwapPacketData) error {
	return k.refundPacketToken(ctx, packet, data)
}

func (k Keeper) refundPacketToken(ctx sdk.Context, packet channeltypes.Packet, data types.SwapPacketData) error {

	ctx.Logger().Debug("refundPacketToken: %s, %s, %s", data.ReceivingTokenDenom, data.SendingTokenDenom, data.Sender)

	return nil
}
