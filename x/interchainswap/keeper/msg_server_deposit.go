package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func (k msgServer) Deposit(goCtx context.Context, msg *types.MsgDepositRequest) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_, found := k.GetInterchainLiquidityPool(ctx, msg.PoolId)
	if !found {
		return nil, nil
	}

	for _, denom := range msg.Tokens {
		balance := k.bankViewKeeper.GetBalance(ctx, sdk.AccAddress(msg.Sender), denom.Denom)
		if balance.Amount.Equal(math.NewInt(0)) {
			return nil, nil
		}
	}

	//k.Keeper.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(msg.Sender), types.ModuleName, msg.Tokens)

	return &types.MsgDepositResponse{}, nil
}
