package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"
)

// GetSwapEnabled retrieves the send enabled boolean from the paramstore
func (k Keeper) GetSwapEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramSpace.Get(ctx, types.KeySwapEnabled, &res)
	return res
}

func (k Keeper) GetSwapMaxFeeRate(ctx sdk.Context) uint32 {
	var res uint32
	k.paramSpace.Get(ctx, types.KeySwapMaxFeeRate, &res)
	return res
}

// GetParams returns the total set of ibc-transfer parameters.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.GetSwapEnabled(ctx), k.GetSwapMaxFeeRate(ctx))
}

// SetParams sets the total set of ibc-transfer parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
