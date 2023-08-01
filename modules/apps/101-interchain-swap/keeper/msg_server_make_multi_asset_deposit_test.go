package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suit *KeeperTestSuite) TestMakeMultiAssetDeposit() {
	ctx := suit.chainA.GetContext()

	// create mock pool
	const poolId = "test-pool"
	const aDenom = "denom-a"
	const bDenom = "denom-b"

	k := suit.chainA.GetSimApp().InterchainSwapKeeper
	k.SetInterchainLiquidityPool(ctx, types.InterchainLiquidityPool{
		Id:                  poolId,
		PoolPrice:           0,
		SourceCreator:       suit.chainA.SenderAccount.GetAddress().String(),
		DestinationCreator:  suit.chainB.SenderAccount.GetAddress().String(),
		CounterPartyPort:    types.ModuleName,
		CounterPartyChannel: "channel-0",
		Assets: []*types.PoolAsset{
			{Side: types.PoolAssetSide_SOURCE, Decimal: 6, Weight: 50, Balance: &sdk.Coin{
				Denom:  aDenom,
				Amount: sdk.NewInt(100),
			}},
			{Side: types.PoolAssetSide_DESTINATION, Decimal: 6, Weight: 50, Balance: &sdk.Coin{
				Denom:  bDenom,
				Amount: sdk.NewInt(100),
			}},
		},
		Status:  types.PoolStatus_ACTIVE,
		SwapFee: 300,
		Supply: &sdk.Coin{
			Amount: sdk.NewInt(200),
			Denom:  poolId,
		},
	})

	// send multi deposit message
	deposits := sdk.Coins{
		{
			Denom:  aDenom,
			Amount: sdk.NewInt(100),
		},
		{
			Denom:  bDenom,
			Amount: sdk.NewInt(100),
		},
	}
	k.
		MakeMultiAssetDeposit(
			ctx, types.NewMsgMakeMultiAssetDeposit(
				poolId,
				[]string{
					suit.chainA.SenderAccount.GetAddress().String(),
					suit.chainB.SenderAccount.GetAddress().String(),
				},
				deposits,
				"interchainswap",
				"channel-0",
			))
	orderId := types.GetOrderId(suit.chainA.SenderAccount.GetAddress().String(), 0)
	order, found := k.GetMultiDepositOrder(ctx, poolId, orderId)
	suit.Require().Equal(ctx.ChainID(), order.ChainId)
	suit.Require().Equal(found, true)
	orders := k.GetAllMultiDepositOrder(ctx, poolId)
	suit.Require().Equal(len(orders), 1)
	fmt.Println("Orders:", orders)
}
