package keeper_test

import (
	"fmt"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

// func (suite *KeeperTestSuite) TestAppendAndGetInterchainLiquidityPool() {
// 	ctx := suite.chainA.GetContext()
// 	k := suite.chainA.GetSimApp().InterchainSwapKeeper

// 	// Set a few pools
// 	poolCount := 5 // Let's say you want to set 5 pools for this test. Adjust as needed.
// 	createdPools := make([]types.InterchainLiquidityPool, poolCount)

// 	for i := 1; i <= poolCount; i++ {
// 		pool := types.InterchainLiquidityPool{
// 			Id: fmt.Sprintf("pool%d", i),
// 		}
// 		createdPools[i-1] = pool
// 		k.AppendInterchainLiquidityPool(ctx, pool)
// 	}

// 	// Get and verify each set pool
// 	for _, createdPool := range createdPools {
// 		retrievedPool, found := k.GetInterchainLiquidityPool(ctx, createdPool.Id)
// 		suite.Require().True(found, "Expected to find the set pool.")
// 		suite.Require().Equal(createdPool, retrievedPool, "The set pool did not match the retrieved pool.")
// 	}

//		// Now, let's test exceeding the max pool count (if applicable)
//		// You'd need to set pools greater than the `types.MaxPoolCount`, and then ensure that the oldest pools are no longer retrievable.
//		if poolCount < types.MaxPoolCount {
//			for i := poolCount + 1; i <= types.MaxPoolCount+1; i++ { // +1 to exceed max count by one
//				pool := types.InterchainLiquidityPool{
//					Id: fmt.Sprintf("pool%d", i),
//				}
//				k.AppendInterchainLiquidityPool(ctx, pool)
//			}
//			// Now, the oldest pool (pool1) should not be found
//			_, found := k.GetInterchainLiquidityPool(ctx, "pool1")
//			suite.Require().False(found, "Expected not to find the oldest pool.")
//		}
//	}
func (suite *KeeperTestSuite) TestSetAndGetInterchainLiquidityPool() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().InterchainSwapKeeper

	// Set a few pools
	poolCount := 5 // Let's say you want to set 5 pools for this test. Adjust as needed.
	createdPools := make([]types.InterchainLiquidityPool, poolCount)

	for i := 1; i <= poolCount; i++ {
		pool := types.InterchainLiquidityPool{
			Id: fmt.Sprintf("pool%d", i),
		}
		createdPools[i-1] = pool
		k.AppendInterchainLiquidityPool(ctx, pool)
	}
	for i := 1; i <= poolCount; i++ {
		pool, found := k.GetInterchainLiquidityPool(ctx, fmt.Sprintf("pool%d", i))
		suite.Require().True(found)
		pool.Status = types.PoolStatus_ACTIVE
		k.SetInterchainLiquidityPool(ctx, pool)
		pool, found = k.GetInterchainLiquidityPool(ctx, fmt.Sprintf("pool%d", i))
		suite.Require().True(found)
		suite.Require().Equal(types.PoolStatus_ACTIVE, pool.Status)
	}
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
		k.AppendInterchainLiquidityPool(ctx, pool)
	}

	pools := k.GetAllInterchainLiquidityPool(ctx)
	suite.Require().Equal(expectedPools, pools, "Retrieved pools did not match the expected pools.")
}

func (suite *KeeperTestSuite) TestRemoveInterchainLiquidityPools() {
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
		k.AppendInterchainLiquidityPool(ctx, pool)
	}

	for i := 0; i < 5; i++ {
		pool := types.InterchainLiquidityPool{
			Id: fmt.Sprintf("pool%d", i),
			// ... other fields
		}
		expectedPools[i] = pool
		k.RemoveInterchainLiquidityPool(ctx, pool.Id)
	}
	pools := k.GetAllInterchainLiquidityPool(ctx)
	suite.Require().Equal(len(pools), 0)
}
