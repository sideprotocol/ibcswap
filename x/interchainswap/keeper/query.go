package keeper

import (
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

var _ types.QueryServer = Keeper{}
