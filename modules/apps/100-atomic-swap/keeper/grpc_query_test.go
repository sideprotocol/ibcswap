package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
	ibctesting "github.com/ibcswap/ibcswap/v6/testing"
)

func (suite *KeeperTestSuite) TestQueryParams() {
	ctx := sdk.WrapSDKContext(suite.chainA.GetContext())
	expParams := types.DefaultParams()
	res, _ := suite.queryClient.Params(ctx, &types.QueryParamsRequest{})
	suite.Require().Equal(&expParams, res.Params)
}

func (suite *KeeperTestSuite) TestEscrowAddress() {
	var (
		req *types.QueryEscrowAddressRequest
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				req = &types.QueryEscrowAddressRequest{
					PortId:    ibctesting.SwapPort,
					ChannelId: ibctesting.FirstChannelID,
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.chainA.GetContext())

			res, err := suite.queryClient.EscrowAddress(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				expected := types.GetEscrowAddress(ibctesting.SwapPort, ibctesting.FirstChannelID).String()
				suite.Require().Equal(expected, res.EscrowAddress)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
