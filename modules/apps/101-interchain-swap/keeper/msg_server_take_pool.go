package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (k msgServer) TakePool(ctx context.Context, msg *types.MsgTakePoolRequest) (*types.MsgTakePoolResponse, error) {

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	_ = sdkCtx
	return nil, nil
}
