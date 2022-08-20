package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/sideprotocol/ibcswap/v4/modules/apps/31-atomic-swap/types"
)

var _ types.MsgServer = Keeper{}

// See createOutgoingPacket in spec:https://github.com/cosmos/ibc/tree/master/spec/app/ics-020-fungible-token-transfer#packet-relay

// Swap defines a rpc handler method for MsgTransfer.
func (k Keeper) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender, err := sdk.AccAddressFromBech32(msg.SenderAddress)
	if err != nil {
		return nil, err
	}
	if err := k.SendSwap(
		ctx, msg.SourcePort, msg.SourceChannel,
		msg.SendingToken, msg.SendingToken,
		sender, msg.SenderReceivingAddress, msg.ExpectedCounterpartyAddress,
		msg.TimeoutHeight, msg.TimeoutTimestamp,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC fungible token transfer", "token",
		msg.SendingToken.Denom, "amount", msg.SendingToken.Amount.String(), "sender",
		msg.SenderAddress, "receiver", msg.SenderReceivingAddress, "counterparty", msg.ExpectedCounterpartyAddress)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSwap,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.SenderAddress),
			sdk.NewAttribute(types.AttributeKeyReceiver, msg.SenderReceivingAddress),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		),
	})

	return &types.MsgSwapResponse{}, nil
}
