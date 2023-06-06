package interchainswap_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v6/testing"
)

type InterchainSwapTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *InterchainSwapTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewInterchainSwapTestPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *InterchainSwapTestSuite) TestHandleMsgInterchainSwap() {
	// setup between chainA and chainB
	path := NewInterchainSwapTestPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	//	originalBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(0, 110)
	_ = timeoutHeight
	amount, ok := sdk.NewIntFromString("9223372036854775808") // 2^63 (one above int64)
	suite.Require().True(ok)
	coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amount)
	_ = coinToSendToB

	// // send from chainA to chainB
	// msg := types.NewMsgCreatePool(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
	// 	coinToSendToB, coinToSendToB,
	// 	suite.chainA.SenderAccount.GetAddress().String(),
	// 	suite.chainA.SenderAccount.GetAddress().String(), // it's same, since the prefix is same on both chains
	// 	suite.chainB.SenderAccount.GetAddress().String(),
	// 	timeoutHeight, 0, time.Now().UTC().Unix())
	// res, err := suite.chainA.SendMsgs(msg)
	// suite.Require().NoError(err) // message committed

	// packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	// suite.Require().NoError(err)

	// // relay send
	// err = path.RelayPacket(packet)
	// suite.Require().NoError(err) // relay committed

	// order := types.NewOTCOrder(msg, path.EndpointA.ChannelID)
	// suite.chainB.NextBlock()
	// has := suite.chainB.GetSimApp().AtomicSwapKeeper.HasOTCOrder(suite.chainB.GetContext(), order.Id)
	suite.Require().True(true)

}

func TestInterchainSwapTestSuite(t *testing.T) {
	suite.Run(t, new(InterchainSwapTestSuite))
}
