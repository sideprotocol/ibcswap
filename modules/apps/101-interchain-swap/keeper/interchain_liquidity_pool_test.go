package keeper_test

import (
	"fmt"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestSetAndGetInterchainLiquidityPool() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().InterchainSwapKeeper

	pool := types.InterchainLiquidityPool{
		Id: "pool1",
	}

	k.SetInterchainLiquidityPool(ctx, pool)

	retrievedPool, found := k.GetInterchainLiquidityPool(ctx, pool.Id)
	suite.Require().True(found, "Expected to find the set pool.")
	suite.Require().Equal(pool, retrievedPool, "The set pool did not match the retrieved pool.")
}

func (suite *KeeperTestSuite) TestGetAllInterchainLiquidityPools() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().InterchainSwapKeeper

	// Add several pools
	expectedPools := make([]types.InterchainLiquidityPool, 5)
	for i := 0; i < 5; i++ {
		pool := types.InterchainLiquidityPool{
			Id: fmt.Sprintf("pool%d", i),
			// ... other fields
		}
		expectedPools[i] = pool
		k.SetInterchainLiquidityPool(ctx, pool)
	}

	pools := k.GetAllInterchainLiquidityPool(ctx)
	suite.Require().Equal(expectedPools, pools, "Retrieved pools did not match the expected pools.")
}
