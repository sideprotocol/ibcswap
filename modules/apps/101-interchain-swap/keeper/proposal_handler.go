package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// HandleMarketFeeUpdateProposal is a handler for executing a passed community spend proposal
func HandleMarketFeeUpdateProposal(ctx sdk.Context, k Keeper, p *types.MarketFeeUpdateProposal) error {

	k.SetSwapFeeRate(ctx, p.FeeRate)
	logger := k.Logger(ctx)
	logger.Info("updated market pool: %s", p.PoolId, "feeRate:", p.FeeRate)
	return nil
}
