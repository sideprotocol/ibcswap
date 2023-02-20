package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/sideprotocol/ibcswap/v4/testutil/keeper"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/keeper"
	"github.com/sideprotocol/ibcswap/v4/x/interchainswap/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.InterchainswapKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
