package keeper_test

import (
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

func (suite *KeeperTestSuite) TestGenesis() {

	genesis := suite.chainA.GetSimApp().AtomicSwapKeeper.ExportGenesis(suite.chainA.GetContext())

	suite.Require().Equal(types.PortID, genesis.PortId)

	suite.Require().NotPanics(func() {
		suite.chainA.GetSimApp().AtomicSwapKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
