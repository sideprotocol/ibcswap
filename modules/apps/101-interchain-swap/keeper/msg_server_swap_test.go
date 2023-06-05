package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
)

func (suite *KeeperTestSuite) TestMsgSwap() {
	var msg *types.MsgSwapRequest
	var poolId *string
	var err error
	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				ctx := suite.chainA.GetContext()
				pool, found := suite.chainA.GetSimApp().InterchainSwapKeeper.GetInterchainLiquidityPool(ctx, *poolId)
				suite.Require().Equal(found, true)
				pool.Status = types.PoolStatus_ACTIVE
				suite.chainA.GetSimApp().InterchainSwapKeeper.SetInterchainLiquidityPool(ctx, pool)
			},
			true,
		},
		{
			"invalid address",
			func() {
				msg.Sender = "invalid address"
			},
			false,
		},
		{
			"invalid amount",
			func() {
				msg.TokenIn = &sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(0),
				}
			},
			false,
		},
	}

	for _, tc := range testCases {
		// create pool first of all.
		poolId, err = suite.SetupPool()
		suite.Require().NoError(err)
		fmt.Println(poolId)

		sender := suite.chainA.SenderAccount
		//
		msg = &types.MsgSwapRequest{
			SwapType:  types.SwapMsgType_LEFT,
			Sender:    sender.GetAddress().String(),
			Recipient: sample.AccAddress(),
			Slippage:  10,
			TokenIn: &sdk.Coin{
				Denom:  sdk.DefaultBondDenom,
				Amount: sdk.NewInt(100),
			},
			TokenOut: &sdk.Coin{
				Denom:  "bside",
				Amount: sdk.NewInt(100),
			},
		}

		tc.malleate()
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().InterchainSwapKeeper)

		res, err := msgSrv.Swap(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
