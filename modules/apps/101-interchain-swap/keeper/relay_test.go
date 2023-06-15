package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v6/testing"
	"github.com/sideprotocol/ibcswap/v6/testing/testutil/sample"
)

// test sending from chainA to chainB using both coin that orignate on
// chainA and coin that orignate on chainB
func (suite *KeeperTestSuite) TestSendSwap() {
	var (
		//amount sdk.Coin
		msgbyte []byte
		path    *ibctesting.Path
		err     error
	)

	testCases := []struct {
		msg            string
		malleate       func()
		sendFromSource bool
		expPass        bool
	}{
		{
			"successful transfer swap request",
			func() {
				suite.coordinator.CreateChannels(path) //CreateInterchainSwapChannels(path)
				msg := &types.MsgSwapRequest{
					SwapType: types.SwapMsgType_LEFT,
					Sender:   sample.AccAddress(),
					TokenIn: &sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(100),
					},
					TokenOut: &sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(100),
					},
				}

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)

			}, true, true,
		},
		{
			"successful transfer creat pool request",
			func() {
				suite.coordinator.CreateChannels(path)

				msg := types.NewMsgMakePool(
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
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

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer deposit request",
			func() {
				suite.coordinator.CreateChannels(path)
				msg := types.NewMsgSingleAssetDeposit(
					"test pool id",
					suite.chainA.SenderAccount.GetAddress().String(),
					&sdk.Coin{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer withdraw request",
			func() {
				suite.coordinator.CreateChannels(path)
				msg := types.NewMsgMultiAssetWithdraw(
					types.GetPoolId(suite.chainA.GetContext().ChainID(), []string{sdk.DefaultBondDenom, sdk.DefaultBondDenom}),
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
					&sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(10),
					},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewInterchainSwapPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate()
			packet := types.IBCSwapPacketData{
				Type: types.LEFT_SWAP,
				Data: msgbyte,
			}

			err = suite.chainA.GetSimApp().InterchainSwapKeeper.SendIBCSwapPacket(
				suite.chainA.GetContext(),
				path.EndpointA.ChannelConfig.PortID,
				path.EndpointA.ChannelID,
				clienttypes.NewHeight(20, 110), 0,
				packet,
			)
			fmt.Println(err)
			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

// msg := &types.MsgSwapRequest{
// 	SwapType: types.SwapMsgType_LEFT,
// 	Sender:   sample.AccAddress(),
// 	TokenIn: &sdk.Coin{
// 		Denom:  sdk.DefaultBondDenom,
// 		Amount: sdk.NewInt(100),
// 	},
// 	TokenOut: &sdk.Coin{
// 		Denom:  sdk.DefaultBondDenom,
// 		Amount: sdk.NewInt(100),
// 	},
// }

// msgbyte, err = types.ModuleCdc.Marshal(msg)
// suite.Require().NoError(err)
func (suite *KeeperTestSuite) TestOnReceived() {
	var (
		//amount sdk.Coin
		path *ibctesting.Path
		err  error
	)
	// create pool first of all.
	pooId, err := suite.SetupPool()
	suite.Require().NoError(err)
	fmt.Println(pooId)

	testCases := []struct {
		msg            string
		malleate       func()
		sendFromSource bool
		expPass        bool
	}{

		{
			"successful on create received.",
			func() {
				ctx := suite.chainA.GetContext()
				msg := types.NewMsgMakePool(
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					suite.chainA.SenderAccount.GetAddress().String(),
					suite.chainB.SenderAccount.GetAddress().String(),
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

				poolId := types.GetPoolId(suite.chainA.ChainID, []string{
					sdk.DefaultBondDenom, sdk.DefaultBondDenom,
				})
				_, err := suite.chainA.GetSimApp().InterchainSwapKeeper.OnMakePoolReceived(
					ctx,
					msg,
					poolId,
				)
				suite.Require().NoError(err)
			}, true, true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewInterchainSwapPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate()

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
