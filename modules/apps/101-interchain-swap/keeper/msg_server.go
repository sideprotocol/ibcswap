package keeper

import (
	"context"
	//sdk "github.com/cosmos/cosmos-sdk/types"
	//channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/ibcswap/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

var (
	StepSend            = 1
	StepReceive         = 2
	StepAcknowledgement = 3
)
var _ types.MsgServer = Keeper{}

func (k Keeper) CreatePool(ctx context.Context, request *types.MsgCreatePoolRequest) (*types.MsgCreatePoolResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Deposit(ctx context.Context, response *types.MsgCreatePoolResponse) (*types.MsgDepositResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) Withdraw(ctx context.Context, request *types.MsgWithdrawRequest) (*types.MsgWithdrawResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) LeftSwap(ctx context.Context, request *types.MsgLeftSwapRequest) (*types.MsgSwapResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (k Keeper) RightSwap(ctx context.Context, request *types.MsgRightSwapRequest) (*types.MsgSwapResponse, error) {
	//TODO implement me
	panic("implement me")
}
