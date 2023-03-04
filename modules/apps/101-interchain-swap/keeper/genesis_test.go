package keeper_test

import "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"

func (suite *KeeperTestSuite) TestGenesis() {
	genesis := suite.chainA.GetSimApp().IBCInterchainSwapKeeper.ExportGenesis(suite.chainA.GetContext())
	suite.Require().Equal(types.PortID, genesis.PortId)
	suite.Require().NotPanics(func() {
		suite.chainA.GetSimApp().IBCInterchainSwapKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
