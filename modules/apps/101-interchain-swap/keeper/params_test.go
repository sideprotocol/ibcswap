package keeper_test

import "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"

func (suite *KeeperTestSuite) TestParams() {
	expParams := types.DefaultParams()

	params := suite.chainA.GetSimApp().InterchainSwapKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)

	expParams.SwapEnabled = false
	suite.chainA.GetSimApp().InterchainSwapKeeper.SetParams(suite.chainA.GetContext(), expParams)
	params = suite.chainA.GetSimApp().InterchainSwapKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}
