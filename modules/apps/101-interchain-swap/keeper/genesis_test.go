package keeper_test

import (
	"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestGenesis() {

	genesis := suite.chainA.GetSimApp().IBCSwapKeeper.ExportGenesis(suite.chainA.GetContext())

	suite.Require().Equal(types.PortID, genesis.PortId)

	suite.Require().NotPanics(func() {
		suite.chainA.GetSimApp().IBCSwapKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
