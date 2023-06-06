package keeper

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

var (
	StepSend            = 1
	StepReceive         = 2
	StepAcknowledgement = 3
)
var _ types.MsgServer = Keeper{}

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
// https://github.com/cosmos/ibc/tree/main/spec/app/ics-100-atomic-swap
// The step is executed on the Taker chain.
func (k Keeper) OnReceivedMake(ctx sdk.Context, packet channeltypes.Packet, orderId, path string, msg *types.MakeSwapMsg) (string, error) {
	// Check if buyToken is a valid token on the taker chain, could be either native or ibc token
	// Disable it for demo at 2023-3-18
	//supply := k.bankKeeper.GetSupply(ctx, msg.BuyToken.Denom)
	//if supply.Amount.Int64() <= 0 {
	//	return errors.New("buy token does not exist on the taker chain")
	//}
	order := types.Order{
		Id:     orderId,
		Side:   types.REMOTE,
		Status: types.Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}

	k.SetAtomicOrder(ctx, order)

	ctx.EventManager().EmitTypedEvents(msg)
	return order.Id, nil
}

// OnReceivedTake is step 7.1 (Transfer Make Token) of the atomic swap: https://github.com/cosmos/ibc/tree/main/spec/app/ics-100-atomic-swap
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

// OnReceivedCancel is the step 12 (Cancel Order) of the atomic swap: https://github.com/cosmos/ibc/tree/main/spec/app/ics-100-atomic-swap.
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

func createOrder(ctx sdk.Context, msg *types.MakeSwapMsg, channelKeeper types.ChannelKeeper) (*types.Order, error) {
	channel, found := channelKeeper.GetChannel(ctx, msg.SourcePort, msg.SourceChannel)
	if !found {
		return nil, types.ErrNotFoundChannel
	}
	sequence := types.GenerateRandomString(ctx.ChainID(), 10)
	path := orderPath(msg.SourcePort, msg.SourceChannel, channel.Counterparty.PortId, channel.Counterparty.ChannelId, sequence)
	return &types.Order{
		Id:     generateOrderId(path, msg),
		Side:   types.NATIVE,
		Status: types.Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}, nil
}

func orderPath(sourcePort, sourceChannel, destPort, destChannel, sequence string) string {
	return fmt.Sprintf("channel/%s/port/%s/channel/%s/port/%s/%d", sourceChannel, sourcePort, destChannel, destPort, sequence)
}

func generateOrderId(orderPath string, msg *types.MakeSwapMsg) string {
	prefix := []byte(orderPath)
	bytes, _ := types.ModuleCdc.MarshalJSON(msg)
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
