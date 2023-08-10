package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) EmitEvent(ctx sdk.Context,
	action, poolID, sender string, attr ...sdk.Attribute,
) {
	headerAttr := []sdk.Attribute{
		{
			Key:   types.AttributeKeyAction,
			Value: action,
		},
		{
			Key:   types.AttributeKeyPoolId,
			Value: poolID,
		},
		{
			Key:   types.AttributeKeyName,
			Value: types.EventOwner,
		},
		{
			Key:   types.AttributeKeyMsgSender,
			Value: sender,
		},
	}

	headerAttr = append(headerAttr, attr...)
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.ModuleName,
			headerAttr...,
		),
	)
}
