package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) MultiAssetWithdraw(goCtx context.Context, msg *types.MsgMultiAssetWithdrawRequest) (*types.MsgMultiAssetWithdrawResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// check out denom
	if !k.bankKeeper.HasSupply(ctx, msg.PoolToken.Denom) {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "invalid denom in local withdraw message:%s", msg.PoolToken.Denom)
	}

	tokenBalance := k.bankKeeper.GetBalance(ctx, sdk.MustAccAddressFromBech32(msg.Receiver), msg.PoolId)
	if tokenBalance.Amount.LT(msg.PoolToken.Amount) {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "sender don't have enough pool token amount:%s", msg.PoolToken.Amount)
	}

	// PoolCoin.Denom is just poolID.
	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolToken.Denom)

	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedWithdraw, "because of %s", types.ErrNotFoundPool)
	}

	amm := *types.NewInterchainMarketMaker(
		&pool,
	)

	outs, err := amm.MultiAssetWithdraw(*msg.PoolToken)
	if err != nil {
		return nil, err
	}

	//burn voucher token.
	err = k.BurnTokens(ctx, sdk.MustAccAddressFromBech32(msg.Receiver), *msg.PoolToken)
	if err != nil {
		return nil, err
	}

	// construct the IBC data packet
	rawMsgData := types.ModuleCdc.MustMarshalJSON(msg)
	rawStateChange := types.ModuleCdc.MustMarshalJSON(&types.StateChange{
		Out:        outs,
		PoolTokens: []*sdk.Coin{msg.PoolToken},
	})

	packet := types.IBCSwapPacketData{
		Type:        types.MULTI_WITHDRAW,
		Data:        rawMsgData,
		StateChange: rawStateChange,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	// Use input timeoutHeight, timeoutStamp
	if msg.TimeoutHeight != nil {
		timeoutHeight = *msg.TimeoutHeight
	}
	if msg.TimeoutTimeStamp != 0 {
		timeoutStamp = msg.TimeoutTimeStamp
	}

	_, err = k.SendIBCSwapPacket(ctx, msg.Port, msg.Channel, timeoutHeight, uint64(timeoutStamp), packet)
	if err != nil {
		return nil, types.ErrFailedWithdraw
	}

	rawOuts := []string{}
	for _, out := range outs {
		rawOuts = append(rawOuts, out.String())
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeLiquidityWithdraw,
			sdk.Attribute{
				Key:   types.AttributeKeyPoolId,
				Value: msg.PoolId,
			},
			sdk.Attribute{
				Key:   types.AttributeKeyLpToken,
				Value: msg.PoolToken.String(),
			},
			sdk.Attribute{
				Key:   types.AttributeKeyTokenOut,
				Value: strings.Join(rawOuts, ":"),
			},
		),
	)

	return &types.MsgMultiAssetWithdrawResponse{}, nil
}
