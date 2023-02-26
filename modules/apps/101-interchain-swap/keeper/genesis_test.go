package keeper_test

import (
	"fmt"

	//"github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestGenesis() {

	fmt.Println("Get start!")
	genesis := suite.chainA.GetSimApp().IBCInterchainSwapKeeper.ExportGenesis(suite.chainA.GetContext())

	fmt.Println("Genesis:", genesis)
	//suite.Require().Equal(types.PortID, genesis.PortId)

	// suite.Require().NotPanics(func() {
	// 	suite.chainA.GetSimApp().IBCInterchainSwapKeeper.InitGenesis(suite.chainA.GetContext(), *genesis)
	// })
}
