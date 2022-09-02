package keeper_test

import (
	"fmt"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/31-atomic-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v4/testing"
)

// test sending from chainA to chainB using both coin that orignate on
// chainA and coin that orignate on chainB
func (suite *KeeperTestSuite) TestSendSwap() {
	var (
		amount sdk.Coin
		path   *ibctesting.Path
		err    error
	)

	testCases := []struct {
		msg            string
		malleate       func()
		sendFromSource bool
		expPass        bool
	}{
		{
			"successful transfer from source chain",
			func() {
				suite.coordinator.CreateSwapChannels(path)
				amount = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
			}, true, true,
		},
		{
			"successful transfer with coin from counterparty chain",
			func() {
				// send coin from chainA back to chainB
				suite.coordinator.CreateSwapChannels(path)
				amount = types.GetTransferCoin(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, sdk.DefaultBondDenom, sdk.NewInt(100))
			}, false, true,
		},
		{
			"source channel not found",
			func() {
				// channel references wrong ID
				suite.coordinator.CreateSwapChannels(path)
				path.EndpointA.ChannelID = ibctesting.InvalidID
				amount = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
			}, true, false,
		},
		{
			"next seq send not found",
			func() {
				path.EndpointA.ChannelID = "channel-0"
				path.EndpointB.ChannelID = "channel-0"
				// manually create channel so next seq send is never set
				suite.chainA.App.GetIBCKeeper().ChannelKeeper.SetChannel(
					suite.chainA.GetContext(),
					path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
					channeltypes.NewChannel(channeltypes.OPEN, channeltypes.ORDERED, channeltypes.NewCounterparty(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID), []string{path.EndpointA.ConnectionID}, ibctesting.DefaultChannelVersion),
				)
				suite.chainA.CreateChannelCapability(suite.chainA.GetSimApp().ScopedIBCMockKeeper, path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
				amount = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
			}, true, false,
		},

		// createOutgoingPacket tests
		// - source chain
		{
			"send coin failed",
			func() {
				suite.coordinator.CreateSwapChannels(path)
				amount = sdk.NewCoin("randomdenom", sdk.NewInt(100))
			}, true, false,
		},
		{
			"channel capability not found",
			func() {
				suite.coordinator.CreateSwapChannels(path)
				cap := suite.chainA.GetChannelCapability(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)

				// Release channel capability
				suite.chainA.GetSimApp().ScopedTransferKeeper.ReleaseCapability(suite.chainA.GetContext(), cap)
				amount = sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100))
			}, true, false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			path = NewSwapPath(suite.chainA, suite.chainB)
			suite.coordinator.SetupConnections(path)

			tc.malleate()

			msg := types.NewMsgMakeSwap(
				path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
				amount, amount,
				suite.chainA.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				clienttypes.NewHeight(0, 110), 0,
				time.Now().UTC().Unix(),
			)

			msgbyte, err1 := types.ModuleCdc.Marshal(msg)
			suite.Require().NoError(err1)

			packet := types.NewAtomicSwapPacketData(types.MAKE_SWAP, msgbyte, "")

			err = suite.chainA.GetSimApp().IBCSwapKeeper.SendSwapPacket(
				suite.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
				clienttypes.NewHeight(0, 110), 0,
				packet,
			)

			if tc.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}
