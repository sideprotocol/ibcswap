package keeper_test

import (
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
