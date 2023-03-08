package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
)

func (suite *KeeperTestSuite) TestMsgWithdraw() {
	var (
		msg    *types.MsgWithdrawRequest
		poolId *string
		err    error
	)

	testCases := []struct {
		name     string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				// deposit first of all.
				depositMsg := types.NewMsgDeposit(
					*poolId,
					suite.chainA.SenderAccount.GetAddress().String(),
					[]*sdk.Coin{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)}},
				)

				err := suite.chainA.GetSimApp().InterchainSwapKeeper.OnSingleDepositAcknowledged(
					suite.chainA.GetContext(),
					depositMsg,
					&types.MsgDepositResponse{
						PoolToken: &sdk.Coin{
							Denom:  *poolId,
							Amount: math.NewInt(1000),
						},
					},
				)
				suite.NoError(err)
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
				msg.Sender = sample.AccAddress()
			},
			false,
		},
	}

	for _, tc := range testCases {
		// create pool first of all.
		poolId, err = suite.SetupPool()
		suite.Require().NoError(err)
		fmt.Println(poolId)

		//
		coin := sdk.NewCoin(*poolId, sdk.NewInt(10))
		msg = types.NewMsgWithdraw(
			suite.chainA.SenderAccount.GetAddress().String(),
			&coin,
			sdk.DefaultBondDenom,
		)

		tc.malleate()

		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().InterchainSwapKeeper)
		res, err := msgSrv.Withdraw(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
