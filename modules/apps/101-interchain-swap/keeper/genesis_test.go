package keeper_test

import (
	"fmt"

	host "github.com/cosmos/ibc-go/v6/modules/core/24-host"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestGenesis() {
	genesis := suite.chainA.GetSimApp().InterchainSwapKeeper.ExportGenesis(suite.chainA.GetContext())

	_, _, err := suite.chainA.GetSimApp().ScopedInterchainSwapKeeper.LookupModules(
		suite.chainA.GetContext(), host.PortPath(genesis.PortId),
	)
	fmt.Println(err)
	suite.Require().NoError(err)
	fmt.Println(genesis)
	suite.Require().Equal(types.PortID, genesis.PortId)
	suite.Require().NotPanics(func() {
		suite.chainA.GetSimApp().InterchainSwapKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	})
}
