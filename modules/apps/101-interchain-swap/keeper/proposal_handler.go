package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	errorsmod "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// HandleMarketFeeUpdateProposal is a handler for executing a passed community spend proposal
func HandleMarketFeeUpdateProposal(ctx sdk.Context, k Keeper, p *types.MarketFeeUpdateProposal) error {
	_, found := k.GetInterchainLiquidityPool(ctx, p.PoolId)
	if !found {
		return errorsmod.ErrNotFound
	}

	market, found := k.GetInterchainMarketMaker(ctx, p.PoolId)
	if !found {
		return errorsmod.ErrNotFound
	}

	market.FeeRate = uint64(p.FeeRate)
	k.SetInterchainMarketMaker(ctx, market)

	logger := k.Logger(ctx)
	logger.Info("updated market pool: %s", p.PoolId, "feeRate:", p.FeeRate)
	return nil
}
