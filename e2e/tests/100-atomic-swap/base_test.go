package e2e

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	clienttypes "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
	"github.com/sideprotocol/ibcswap/v4/modules/apps/31-atomic-swap/types"
	"github.com/strangelove-ventures/ibctest/ibc"
	"github.com/strangelove-ventures/ibctest/test"
	"github.com/stretchr/testify/suite"
)

func TestAtomicSwapTestSuite(t *testing.T) {
	suite.Run(t, new(AtomicSwapTestSuite))
}

type AtomicSwapTestSuite struct {
	testsuite.E2ETestSuite
}

func (s *AtomicSwapTestSuite) TestMakeSwap() {
	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	relayer, channelA, _ := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())
	_, channelB, _ := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())

	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	//chainBDenom := chainB.Config().Denom

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	senderAddress := chainAWallet.Bech32Address("cosmos")
	chainBWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	senderReceivingAddress := chainBWallet.Bech32Address("token")

	// allocate tokens to the new account
	initialBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(10000000000)))
	err := s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialBalances)

	s.Require().NoError(err)
	res, err := s.QueryBalance(ctx, chainA, senderAddress, chainADenom)
	s.Require().NotEqual(res.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send make swap message", func(t *testing.T) {
		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)

		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, senderAddress, senderReceivingAddress, "", timeoutHeight, 0, time.Now().UTC().Unix())
		msg.Packet = channeltypes.NewPacket(
			[]byte{},
			0,
			channelA.PortID,
			channelA.ChannelID,
			channelB.PortID,
			channelB.ChannelID,
			clienttypes.NewHeight(0, 100),
			0,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 1, chainA, chainB)

		// check packet relay status.
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

	})
}

// atomicSwapChannelOptions configures both of the chains to have atomic swap enabled.
func atomicSwapChannelOptions() func(options *ibc.CreateChannelOptions) {
	return func(opts *ibc.CreateChannelOptions) {
		opts.Version = "ics100-1"
		opts.Order = ibc.Unordered
		opts.DestPortName = types.ModuleName
		opts.SourcePortName = types.ModuleName
	}
}
