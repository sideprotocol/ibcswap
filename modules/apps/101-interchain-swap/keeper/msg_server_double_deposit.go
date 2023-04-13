package keeper

import (
	"context"

	"github.com/btcsuite/btcutil/bech32"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k Keeper) DoubleDeposit(goCtx context.Context, msg *types.MsgDoubleDepositRequest) (*types.MsgDoubleDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// // validate message
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}

	// check address
	senderPrefix, _, err := bech32.Decode(msg.Senders[0])
	if err != nil {
		return nil, err
	}
	if sdk.GetConfig().GetBech32AccountAddrPrefix() != senderPrefix {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "first address has to be this chain address (%s)", err)
	}

	pool, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", types.ErrNotFoundPool)
	}

	// Deposit token to Escrow account
	coins, err := k.validateCoins(ctx, &pool, msg.Senders[0], msg.Tokens)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDeposit, "%s", err)
	}
	if len(coins) == 0 {
		return nil, types.ErrInvalidSignature
	}

	// create escrow module account  here
	err = k.LockTokens(ctx, pool.EncounterPartyPort, pool.EncounterPartyChannel, sdk.MustAccAddressFromBech32(msg.Senders[0]), coins)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrFailedDoubleDeposit, "%s", err)
	}

	// construct ibc packet
	rawMsgData, err := types.ModuleCdc.Marshal(msg)
	if err != nil {
		return nil, err
	}

	packet := types.IBCSwapPacketData{
		Type: types.DOUBLE_DEPOSIT,
		Data: rawMsgData,
	}

	timeoutHeight, timeoutStamp := types.GetDefaultTimeOut(&ctx)
	err = k.SendIBCSwapPacket(ctx, types.PortID, "channel-0", timeoutHeight, timeoutStamp, packet)
	if err != nil {
		return nil, err
	}

	return &types.MsgDoubleDepositResponse{
		PoolTokens: nil,
	}, nil
}
