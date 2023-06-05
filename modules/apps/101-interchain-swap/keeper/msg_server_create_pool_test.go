package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/keeper"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
)

func (suite *KeeperTestSuite) TestMsgCreatePool() {
	var msg *types.MsgCreatePoolRequest

	testCases := []struct {
		name    string
		expPass bool
	}{
		{
			"success",
			true,
		},
		// {
		// 	"channel does not exist",
		// 	false,
		// },
	}

	for _, tc := range testCases {
		suite.SetupTest()

		path := NewInterchainSwapPath(suite.chainA, suite.chainB)
		suite.coordinator.Setup(path)

		msg = types.NewMsgCreatePool(
			path.EndpointA.ChannelConfig.PortID,
			path.EndpointA.ChannelID,
			suite.chainA.SenderAccount.GetAddress().String(),
			suite.chainB.SenderAccount.GetAddress().String(),
			[]byte("0"),
			types.PoolAsset{
				Side: types.PoolAssetSide_SOURCE,
				Balance: &sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				},
				Weight:  50,
				Decimal: 6,
			},

			types.PoolAsset{
				Side: types.PoolAssetSide_SOURCE,
				Balance: &sdk.Coin{
					Denom:  sdk.DefaultBondDenom,
					Amount: sdk.NewInt(1000),
				},
				Weight:  50,
				Decimal: 6,
			},
			300,
		)
		msgSrv := keeper.NewMsgServerImpl(suite.chainA.GetSimApp().InterchainSwapKeeper)
		res, err := msgSrv.CreatePool(sdk.WrapSDKContext(suite.chainA.GetContext()), msg)

		if tc.expPass {
			suite.Require().NoError(err)
			suite.Require().NotNil(res)
		} else {
			suite.Require().Error(err)
			suite.Require().Nil(res)
		}
	}
}
