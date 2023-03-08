package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// GetSwapEnabled retrieves the send enabled boolean from the paramstore
func (k Keeper) GetSwapEnabled(ctx sdk.Context) bool {
	var res bool
	k.paramstore.Get(ctx, types.KeySwapEnabled, &res)
	return res
}

func (k Keeper) GetSwapMaxFeeRate(ctx sdk.Context) uint32 {
	var res uint32
	k.paramstore.Get(ctx, types.KeySwapMaxFeeRate, &res)
	return res
}

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams(k.GetSwapEnabled(ctx), k.GetSwapMaxFeeRate(ctx))
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
