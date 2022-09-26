package keeper_test

import "github.com/ibcswap/ibcswap/v4/modules/apps/100-atomic-swap/types"

func (suite *KeeperTestSuite) TestParams() {
	expParams := types.DefaultParams()

	params := suite.chainA.GetSimApp().IBCSwapKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)

	expParams.SwapEnabled = false
	suite.chainA.GetSimApp().IBCSwapKeeper.SetParams(suite.chainA.GetContext(), expParams)
	params = suite.chainA.GetSimApp().IBCSwapKeeper.GetParams(suite.chainA.GetContext())
	suite.Require().Equal(expParams, params)
}
