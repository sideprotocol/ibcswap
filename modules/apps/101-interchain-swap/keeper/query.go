package keeper

import (
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

var _ types.QueryServer = Keeper{}
