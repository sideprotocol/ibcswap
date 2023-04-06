package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	ibctesting "github.com/ibcswap/ibcswap/v6/testing"
	"github.com/ibcswap/ibcswap/v6/testing/testutil/sample"
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
				msg := types.NewMsgCreatePool(
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					suite.chainA.SenderAccount.GetAddress().String(),
					"1:2",
					[]*sdk.Coin{{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(1000),
					}, {
						Denom:  "bside",
						Amount: sdk.NewInt(1000),
					}},
					[]uint32{10, 100},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer deposit request",
			func() {
				suite.coordinator.CreateChannels(path)
				msg := types.NewMsgDeposit(
					"test pool id",
					suite.chainA.SenderAccount.GetAddress().String(),
					[]*sdk.Coin{{Denom: sdk.DefaultBondDenom, Amount: sdk.NewInt(1000)}},
				)

				msgbyte, err = types.ModuleCdc.Marshal(msg)
				suite.Require().NoError(err)
			}, true, true,
		},
		{
			"successful transfer withdraw request",
			func() {
				suite.coordinator.CreateChannels(path)
				msg := types.NewMsgWithdraw(
					suite.chainA.SenderAccount.GetAddress().String(),
					&sdk.Coin{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(10),
					},
					sdk.DefaultBondDenom,
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
				msg := types.NewMsgCreatePool(
					path.EndpointA.ChannelConfig.PortID,
					path.EndpointA.ChannelID,
					suite.chainA.SenderAccount.GetAddress().String(),
					"1:2",
					[]*sdk.Coin{{
						Denom:  sdk.DefaultBondDenom,
						Amount: sdk.NewInt(1000),
					}, {
						Denom:  "bside",
						Amount: sdk.NewInt(1000),
					}},
					[]uint32{10, 100},
				)
				destPort := path.EndpointA.Counterparty.ChannelConfig.PortID
				destChannel := path.EndpointA.ChannelID
				poolId, err := suite.chainA.GetSimApp().InterchainSwapKeeper.OnCreatePoolReceived(
					ctx,
					msg,
					destPort,
					destChannel,
				)
				suite.Require().NoError(err)
				suite.Require().Equal(*poolId, types.GetPoolIdWithTokens(msg.Tokens))
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
