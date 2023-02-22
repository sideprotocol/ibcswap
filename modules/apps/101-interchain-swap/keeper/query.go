package keeper

import (
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

var _ types.QueryServer = Keeper{}
