package keeper_test

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
)

func (suite *KeeperTestSuite) TestAtomicOrderFiFo() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().AtomicSwapKeeper

	// Number of iterations for the fuzzy test
	const iterations = 100

	// Used to store generated order IDs to later query and verify
	for i := 0; i < iterations; i++ {
		orderID := "test" + strconv.Itoa(rand.Intn(100000)) // Generating a random order ID for uniqueness

		k.AppendAtomicOrder(ctx, types.Order{
			Id:              orderID,
			Side:            types.NATIVE,
			Status:          types.Status_INITIAL,
			Path:            "test" + strconv.Itoa(i),
			Maker:           nil,
			CancelTimestamp: int64(i),
		})
	}
	orders := k.GetAllOrder(ctx)

	total := len(orders)
	suite.Require().Equal(total, iterations)
	for index, order := range orders {
		suite.Require().Equal(order.CancelTimestamp, total-index-1)
	}
	res, err := k.GetAllOrders(ctx, &types.QueryOrdersRequest{})
	suite.Require().NoError(err)
	totalByQuery := len(res.Orders)
	suite.Require().Equal(totalByQuery, iterations)
}

func (suite *KeeperTestSuite) TestMoveOrderToBottom() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().AtomicSwapKeeper

	// Add a few atomic orders
	orderIDs := make([]string, 5)
	for i := 0; i < 5; i++ {
		orderID := fmt.Sprintf("testOrder%d", i)
		orderIDs[i] = orderID
		k.AppendAtomicOrder(ctx, types.Order{
			Id:              orderID,
			Side:            types.NATIVE,
			Status:          types.Status_INITIAL,
			Path:            "path" + strconv.Itoa(i),
			Maker:           nil,
			CancelTimestamp: int64(i),
		})
	}
	// Move the third order to the bottom
	err := k.MoveOrderToBottom(ctx, orderIDs[2])
	suite.Require().NoError(err)

	// Check if the moved order is now the last in the list
	ordersAfterMove := k.GetAllOrder(ctx)
	lastOrder := ordersAfterMove[len(ordersAfterMove)-1]
	suite.Require().Equal(lastOrder.Id, orderIDs[2])

	// Additional check for the updated `orderId -> count` relationship
	movedOrderCount := k.GetAtomicOrderCountByOrderId(ctx, orderIDs[2])
	lastOrderCount := k.GetAtomicOrderCount(ctx) - 1
	suite.Require().Equal(movedOrderCount, lastOrderCount)

	// Verify that all other orders have maintained their original sequence
	for i := 0; i < len(orderIDs)-1; i++ { // Exclude the last order (since it was moved)
		if i < 2 { // For orders before the moved order
			suite.Require().Equal(ordersAfterMove[i].Id, orderIDs[i])
		} else { // For orders after the moved order
			suite.Require().Equal(ordersAfterMove[i].Id, orderIDs[i+1])
		}
	}
}

func (suite *KeeperTestSuite) TestTrimExcessOrders() {
	ctx := suite.chainA.GetContext()
	k := suite.chainA.GetSimApp().AtomicSwapKeeper

	// Add a few more than MaxItems
	for i := uint64(0); i < types.MaxOrderCount+100; i++ {
		orderID := fmt.Sprintf("testOrder%d", i)
		k.AppendAtomicOrder(ctx, types.Order{
			Id:              orderID,
			Side:            types.NATIVE,
			Status:          types.Status_INITIAL,
			Path:            "path" + strconv.FormatUint(i, 10),
			Maker:           nil,
			CancelTimestamp: int64(i),
		})
	}

	// Ensure we have MaxItems + 100 orders
	suite.Require().Equal(uint64(types.MaxOrderCount+100), k.GetAtomicOrderCount(ctx))

	// Trim excess
	k.TrimExcessOrders(ctx)

	// Validate only MaxItems are present now
	suite.Require().Equal(uint64(types.MaxOrderCount), k.GetAtomicOrderCount(ctx))
}
