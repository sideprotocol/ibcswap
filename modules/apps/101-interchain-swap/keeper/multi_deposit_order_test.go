package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestSetMultiDepositOrder() {
	k := suite.chainA.GetSimApp().InterchainSwapKeeper
	ctx := suite.chainA.GetContext()

	// create order
	order := types.MultiAssetDepositOrder{
		PoolId:           "test",
		ChainId:          "test-a",
		SourceMaker:      "test",
		DestinationTaker: "test",
		Deposits:         []*sdk.Coin{{Denom: "test", Amount: sdk.NewInt(100)}, {Denom: "test1", Amount: sdk.NewInt(100)}},
		Status:           types.OrderStatus_PENDING,
		CreatedAt:        123,
	}

	orderId := types.GetOrderId(order.SourceMaker, 0)
	order.Id = orderId

	k.SetMultiDepositOrder(
		ctx,
		order,
	)
	suite.Require().Equal(orderId, uint64(0))

	storedOrder, found := k.GetMultiDepositOrder(
		ctx,
		order.PoolId,
		orderId,
	)
	suite.Require().Equal(found, true)
	suite.Require().Equal(storedOrder, order)

	storedOrders := k.GetAllMultiDepositOrder(
		ctx,
		order.PoolId,
	)
	suite.Require().Equal(1, len(storedOrders))
}

func (suite *KeeperTestSuite) TestGetLatestMultiDepositOrder() {
	k := suite.chainA.GetSimApp().InterchainSwapKeeper
	ctx := suite.chainA.GetContext()

	const poolId = "test"

	// create order
	order := types.MultiAssetDepositOrder{
		PoolId:           poolId,
		ChainId:          "test-a",
		SourceMaker:      "test",
		DestinationTaker: "test",
		Deposits:         []*sdk.Coin{{Denom: "test", Amount: sdk.NewInt(100)}, {Denom: "test1", Amount: sdk.NewInt(100)}},
		Status:           types.OrderStatus_PENDING,
		CreatedAt:        123,
	}

	orderId := types.GetOrderId(order.SourceMaker, 0)
	order.Id = orderId

	k.SetMultiDepositOrder(
		ctx,
		order,
	)
	suite.Require().Equal(orderId, uint64(0))

	storedOrder, found := k.GetMultiDepositOrder(
		ctx,
		order.PoolId,
		orderId,
	)
	suite.Require().Equal(found, true)
	suite.Require().Equal(storedOrder, order)
	suite.Require().Equal(order.Status, types.OrderStatus_PENDING)

}
