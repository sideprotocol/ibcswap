package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/gogo/protobuf/proto"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

var (
	StepSend            = 1
	StepReceive         = 2
	StepAcknowledgement = 3
)
var _ types.MsgServer = Keeper{}

// MakeSwap is called when the maker wants to make atomic swap. The method create new order and lock tokens.
// This is the step 1 (Create order & Lock Token) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap.
func (k Keeper) MakeSwap(goCtx context.Context, msg *types.MakeSwapMsg) (*types.MsgMakeSwapResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}

	msgByte, err0 := proto.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	sender, err1 := sdk.AccAddressFromBech32(msg.MakerAddress)
	if err1 != nil {
		return nil, err1
	}

	balance := k.bankKeeper.GetBalance(ctx, sender, msg.SellToken.Denom)
	if balance.Amount.BigInt().Cmp(msg.SellToken.Amount.BigInt()) < 0 {
		return &types.MsgMakeSwapResponse{}, errors.New("insufficient balance")
	}

	escrowAddr := types.GetEscrowAddress(msg.SourcePort, msg.SourceChannel)

	// lock sell token into module
	if err := k.bankKeeper.SendCoins(
		ctx, sender, escrowAddr, sdk.NewCoins(msg.SellToken),
	); err != nil {
		return nil, err
	}

	packet := types.AtomicSwapPacketData{
		Type: types.MAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}

	order := createOrder(ctx, msg, k.channelKeeper)

	if err := k.SendSwapPacket(ctx, msg.SourcePort, msg.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgMakeSwapResponse{OrderId: order.Id}, nil
}

// TakeSwap is the step 5 (Lock Order & Lock Token) of the atomic swap: https://github.com/liangping/ibc/blob/atomic-swap/spec/app/ics-100-atomic-swap/ibcswap.png
// This method lock the order (set a value to the field "Taker") and lock Token
func (k Keeper) TakeSwap(goCtx context.Context, msg *types.TakeSwapMsg) (*types.MsgTakeSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// pre-check
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgByte, err0 := proto.Marshal(msg)
	if err0 != nil {
		return nil, err0
	}

	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return nil, types.ErrOrderDoesNotExists
	}

	sourceChannel := extractSourceChannelForTakerMsg(order.Path)
	sourcePort := extractSourcePortForTakerMsg(order.Path)

	escrowAddr := types.GetEscrowAddress(sourcePort, sourceChannel)

	if order.Status != types.Status_SYNC && order.Status != types.Status_INITIAL {
		return nil, errors.New("order is not in valid state")
	}

	// Make sure the maker's buy token matches the taker's sell token
	if order.Maker.BuyToken.Denom != msg.SellToken.Denom || !order.Maker.BuyToken.Amount.Equal(msg.SellToken.Amount) {
		return &types.MsgTakeSwapResponse{}, errors.New("invalid sell token")
	}

	// Checks if the order has already been taken
	if order.Takers != nil {
		return &types.MsgTakeSwapResponse{}, errors.New("order has already been taken")
	}

	// If `desiredTaker` is set, only the desiredTaker can accept the order.
	if order.Maker.DesiredTaker != "" && order.Maker.DesiredTaker != msg.TakerAddress {
		return &types.MsgTakeSwapResponse{}, errors.New("invalid taker address")
	}

	takerAddr, err := sdk.AccAddressFromBech32(msg.TakerAddress)
	if err != nil {
		return &types.MsgTakeSwapResponse{}, err
	}

	balance := k.bankKeeper.GetBalance(ctx, takerAddr, msg.SellToken.Denom)
	if balance.Amount.BigInt().Cmp(msg.SellToken.Amount.BigInt()) < 0 {
		return &types.MsgTakeSwapResponse{}, errors.New("insufficient balance")
	}

	// Locks the sellToken to the escrow account
	if err = k.bankKeeper.SendCoins(ctx, takerAddr, escrowAddr, sdk.NewCoins(msg.SellToken)); err != nil {
		return &types.MsgTakeSwapResponse{}, err
	}

	packet := types.AtomicSwapPacketData{
		Type: types.TAKE_SWAP,
		Data: msgByte,
		Memo: "",
	}

	if err := k.SendSwapPacket(ctx, sourcePort, sourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	// Update order state
	// Mark that the order has been occupied
	order.Takers = msg
	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgTakeSwapResponse{}, nil
}

// CancelSwap is the step 10 (Cancel Request) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap.
// It is executed on the Maker chain. Only the maker of the order can cancel the order.
func (k Keeper) CancelSwap(goCtx context.Context, msg *types.CancelSwapMsg) (*types.MsgCancelSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	msgbyte, err := proto.Marshal(msg)
	if err != nil {
		return nil, err
	}

	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return &types.MsgCancelSwapResponse{}, types.ErrOrderDoesNotExists
	}

	// Make sure the sender is the maker of the order.
	if order.Maker.MakerAddress != msg.MakerAddress {
		return &types.MsgCancelSwapResponse{}, fmt.Errorf("sender is not the maker of the order")
	}

	// Make sure the order is in a valid state for cancellation
	if order.Status != types.Status_SYNC && order.Status != types.Status_INITIAL {
		return &types.MsgCancelSwapResponse{}, fmt.Errorf("order is not in a valid state for cancellation")
	}

	packet := types.AtomicSwapPacketData{
		Type: types.CANCEL_SWAP,
		Data: msgbyte,
		Memo: "",
	}

	if err := k.SendSwapPacket(ctx, order.Maker.SourcePort, order.Maker.SourceChannel, msg.TimeoutHeight, msg.TimeoutTimestamp, packet); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitTypedEvents(msg)

	return &types.MsgCancelSwapResponse{}, nil
}

// See createOutgoingPacket in spec:https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay
func (k Keeper) fillAtomicOrder(ctx sdk.Context, escrowAddr sdk.AccAddress, order types.Order, msg *types.TakeSwapMsg, step int) error {

	switch order.Status {
	case types.Status_CANCEL:
		return types.ErrOrderCanceled
	case types.Status_COMPLETE:
		return types.ErrOrderCompleted
	default:
		// continue
	}
	if order.Maker.BuyToken.Denom != msg.SellToken.Denom {
		return types.ErrOrderDenominationMismatched
	}

	if msg.SellToken.IsGTE(order.Maker.BuyToken) {
		order.Takers = msg // types.NewTakerFromMsg(msg)
		order.Status = types.Status_COMPLETE
		order.CompleteTimestamp = msg.CreateTimestamp
	} else {
		return types.ErrOrderInsufficientAmount
	}

	switch step {
	case StepSend:
		// executed on the destination chain of make swap,
		// lock the taker's sell coin into the module
		sender, err1 := sdk.AccAddressFromBech32(msg.TakerAddress)
		if err1 != nil {
			return err1
		}

		// lock sell token into module
		if err := k.bankKeeper.SendCoins(
			ctx, sender, escrowAddr, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	case StepReceive:
		// StepReceive
		// executed on the source chain of make swap,
		// send maker's sell token to taker
		receiver, err1 := sdk.AccAddressFromBech32(msg.TakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// send maker's sell token from module to taker
		if err := k.bankKeeper.SendCoins(
			ctx, escrowAddr, receiver, sdk.NewCoins(order.Maker.SellToken),
		); err != nil {
			return err
		}
		break
	case StepAcknowledgement:
		// executed on the destination chain of the swap,
		// send taker's sell token to maker
		receiver, err1 := sdk.AccAddressFromBech32(order.Maker.MakerReceivingAddress)
		if err1 != nil {
			return err1
		}

		// send taker's sell token from module to maker
		if err := k.bankKeeper.SendCoins(
			ctx, escrowAddr, receiver, sdk.NewCoins(msg.SellToken),
		); err != nil {
			return err
		}
		break
	}

	// save updated status
	k.SetAtomicOrder(ctx, order)

	return nil
}

func (k Keeper) executeCancel(ctx sdk.Context, msg *types.CancelSwapMsg, step int) error {
	// check order status
	if order, ok2 := k.GetAtomicOrder(ctx, msg.OrderId); ok2 {
		if order.Maker.MakerAddress != msg.MakerAddress {
			return types.ErrOrderPermissionIsNotAllowed
		}
		switch order.Status {
		case types.Status_CANCEL:
			return types.ErrOrderCanceled
		case types.Status_COMPLETE:
			return types.ErrOrderCompleted
		default:
			// continue
			if step != StepSend {
				order.CancelTimestamp = msg.CreateTimestamp
				order.Status = types.Status_CANCEL
			}

			// refund
			if step == StepAcknowledgement {
				receiver := sdk.MustAccAddressFromBech32(msg.MakerAddress)
				escrowAddr := types.GetEscrowAddress(order.Maker.SourcePort, order.Maker.SourceChannel)
				k.bankKeeper.SendCoins(ctx, escrowAddr, receiver, sdk.NewCoins(order.Maker.SellToken))
			}
		}
	} else {
		return types.ErrOrderDoesNotExists
	}

	return nil
}

// the following methods are executed On Destination chain.

// OnReceivedMake is the step 3.1 (Save order) from the atomic swap:
// https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap
// The step is executed on the Taker chain.
func (k Keeper) OnReceivedMake(ctx sdk.Context, packet channeltypes.Packet, msg *types.MakeSwapMsg) (string, error) {
	// Check if buyToken is a valid token on the taker chain, could be either native or ibc token
	// Disable it for demo at 2023-3-18
	//supply := k.bankKeeper.GetSupply(ctx, msg.BuyToken.Denom)
	//if supply.Amount.Int64() <= 0 {
	//	return errors.New("buy token does not exist on the taker chain")
	//}

	path := orderPath(msg.SourcePort, msg.SourceChannel, packet.DestinationPort, packet.DestinationChannel, packet.Sequence)
	order := types.Order{
		Id:     generateOrderId(path, msg),
		Status: types.Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}

	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)
	return order.Id, nil
}

// OnReceivedTake is step 7.1 (Transfer Make Token) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap
// The step is executed on the Maker chain.
func (k Keeper) OnReceivedTake(ctx sdk.Context, packet channeltypes.Packet, msg *types.TakeSwapMsg) (string, error) {
	escrowAddr := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())

	// check order status
	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return "", types.ErrOrderDoesNotExists
	}

	//if order.Status != types.Status_SYNC {
	//	return errors.New("invalid order status")
	//}

	if msg.SellToken.Denom != order.Maker.BuyToken.Denom || !msg.SellToken.Amount.Equal(order.Maker.BuyToken.Amount) {
		return "", errors.New("invalid sell token")
	}

	// If `desiredTaker` is set, only the desiredTaker can accept the order.
	if order.Maker.DesiredTaker != "" && order.Maker.DesiredTaker != msg.TakerAddress {
		return "", errors.New("invalid taker address")
	}

	takerReceivingAddr, err := sdk.AccAddressFromBech32(msg.TakerReceivingAddress)
	if err != nil {
		return "", err
	}

	// Send maker.sellToken to taker's receiving address
	if err = k.bankKeeper.SendCoins(ctx, escrowAddr, takerReceivingAddr, sdk.NewCoins(order.Maker.SellToken)); err != nil {
		return "", errors.New("transfer coins failed")
	}

	// Update status of order
	order.Status = types.Status_COMPLETE
	order.Takers = msg
	order.CompleteTimestamp = msg.CreateTimestamp
	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)
	return order.Id, nil
}

// OnReceivedCancel is the step 12 (Cancel Order) of the atomic swap: https://github.com/liangping/ibc/tree/atomic-swap/spec/app/ics-100-atomic-swap.
// This step is executed on the Taker chain.
func (k Keeper) OnReceivedCancel(ctx sdk.Context, packet channeltypes.Packet, msg *types.CancelSwapMsg) (string, error) {
	if err := msg.ValidateBasic(); err != nil {
		return "", err
	}
	// check order status

	order, ok := k.GetAtomicOrder(ctx, msg.OrderId)
	if !ok {
		return "", errors.New("order not found")
	}

	if order.Status != types.Status_SYNC && order.Status != types.Status_INITIAL {
		return "", errors.New("invalid order status")
	}

	if order.Takers != nil {
		return "", errors.New("the maker order has already been occupied")
	}

	// Update status of order
	order.Status = types.Status_CANCEL
	order.CancelTimestamp = msg.CreateTimestamp
	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)
	return order.Id, nil
}

func createOrder(ctx sdk.Context, msg *types.MakeSwapMsg, channelKeeper types.ChannelKeeper) types.Order {
	channel, _ := channelKeeper.GetChannel(ctx, msg.SourcePort, msg.SourceChannel)
	sequence, _ := channelKeeper.GetNextSequenceSend(ctx, msg.SourcePort, msg.SourceChannel)
	path := orderPath(msg.SourcePort, msg.SourceChannel, channel.Counterparty.PortId, channel.Counterparty.ChannelId, sequence)
	return types.Order{
		Id:     generateOrderId(path, msg),
		Status: types.Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}
}

func orderPath(sourcePort, sourceChannel, destPort, destChannel string, sequence uint64) string {
	return fmt.Sprintf("channel/%s/port/%s/channel/%s/port/%s/%d", sourceChannel, sourcePort, destChannel, destPort, sequence)
}

func generateOrderId(orderPath string, msg *types.MakeSwapMsg) string {
	prefix := []byte(orderPath)
	bytes, _ := proto.Marshal(msg)
	hash := sha256.Sum256(append(prefix, bytes...))
	return hex.EncodeToString(hash[:])
}

func extractSourceChannelForTakerMsg(path string) string {
	parts := strings.Split(path, "/")
	return parts[5]
}

func extractSourcePortForTakerMsg(path string) string {
	parts := strings.Split(path, "/")
	return parts[7]
}
