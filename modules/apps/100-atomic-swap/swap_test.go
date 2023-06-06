package swap_test

import (
	"testing"

	"github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
	ibctesting "github.com/sideprotocol/ibcswap/v6/testing"
)

type SwapTestSuite struct {
	suite.Suite

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func (suite *SwapTestSuite) SetupTest() {
	suite.coordinator = ibctesting.NewCoordinator(suite.T(), 2)
	suite.chainA = suite.coordinator.GetChain(ibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(ibctesting.GetChainID(2))
}

func NewSwapTestPath(chainA, chainB *ibctesting.TestChain) *ibctesting.Path {
	path := ibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = types.PortID
	path.EndpointB.ChannelConfig.PortID = types.PortID
	path.EndpointA.ChannelConfig.Version = types.Version
	path.EndpointB.ChannelConfig.Version = types.Version

	return path
}

// constructs a send from chainA to chainB on the established channel/connection
// and sends the same coin back from chainB to chainA.
func (suite *SwapTestSuite) TestHandleMsgSwap() {
	// setup between chainA and chainB
	path := NewSwapTestPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	//	originalBalance := suite.chainA.GetSimApp().BankKeeper.GetBalance(suite.chainA.GetContext(), suite.chainA.SenderAccount.GetAddress(), sdk.DefaultBondDenom)
	timeoutHeight := clienttypes.NewHeight(0, 110)

	amount, ok := sdk.NewIntFromString("9223372036854775808") // 2^63 (one above int64)
	suite.Require().True(ok)
	coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, amount)

	// send from chainA to chainB
	msg := types.NewMsgMakeSwap(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID,
		coinToSendToB, coinToSendToB,
		suite.chainA.SenderAccount.GetAddress().String(),
		suite.chainA.SenderAccount.GetAddress().String(), // it's same, since the prefix is same on both chains
		suite.chainB.SenderAccount.GetAddress().String(),
		timeoutHeight, 0, time.Now().UTC().Unix())
	res, err := suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err) // message committed

	packet, err := ibctesting.ParsePacketFromEvents(res.GetEvents())
	suite.Require().NoError(err)

	// relay send
	err = path.RelayPacket(packet)
	suite.Require().NoError(err) // relay committed

	//order := types.NewAtomicOrder(msg, path.EndpointA.ChannelID)
	//suite.chainB.NextBlock()
	//has := suite.chainB.GetSimApp().AtomicSwapKeeper.HasOTCOrder(suite.chainB.GetContext(), order.Id)
	//suite.Require().True(has)

}

func TestSwapTestSuite(t *testing.T) {
	suite.Run(t, new(SwapTestSuite))
}
