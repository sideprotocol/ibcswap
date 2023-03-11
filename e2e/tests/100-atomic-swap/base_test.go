package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
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
	relayer, channelA := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())

	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	//chainBDenom := chainB.Config().Denom

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	senderAddress := chainAWallet.Bech32Address("cosmos")
	chainBWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	senderReceivingAddress := chainBWallet.Bech32Address("cosmos")

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
		msg := types.NewMsgMakeSwap(
			channelA.PortID,
			channelA.ChannelID,
			sellToken,
			buyToken,
			senderAddress,
			senderReceivingAddress,
			"",
			timeoutHeight,
			0,
			time.Now().UTC().Unix(),
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// TAKE SWAP

		sellToken2 := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		senderAddressB := chainBWallet.Bech32Address("cosmos")
		senderReceivingAddressA := chainAWallet.Bech32Address("cosmos")

		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := types.NewAtomicOrder(types.NewMakerFromMsg(msg), msg.SourceChannel)
		msgTake := types.NewMsgTakeSwap(
			channelA.PortID,
			channelA.ChannelID,
			sellToken2,
			senderAddressB,
			senderReceivingAddressA,
			timeoutHeight2,
			0,
			time.Now().UTC().Unix(),
		)
		msgTake.OrderId = order.Id

		resp2, err2 := s.BroadcastMessages(ctx, chainB, chainBWallet, msgTake)

		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// check packet relay status.
		//s.AssertPacketRelayed(ctx, chainB, channelA.PortID, channelA.ChannelID, 1)

		b1, err := s.QueryBalance(ctx, chainA, senderAddress, chainA.Config().Denom)
		s.Require().NoError(err)
		fmt.Println("FINALLLLL 1")
		fmt.Println(b1.Balance.String())

		b2, err := s.QueryBalance(ctx, chainA, senderReceivingAddressA, chainA.Config().Denom)
		s.Require().NoError(err)
		fmt.Println("FINALLLLL 2")
		fmt.Println(b2.Balance.String())

		b3, err := s.QueryBalance(ctx, chainB, senderAddressB, chainB.Config().Denom)
		s.Require().NoError(err)
		fmt.Println("FINALLLLL 3")
		fmt.Println(b3.Balance.String())

		b4, err := s.QueryBalance(ctx, chainB, senderReceivingAddress, chainB.Config().Denom)
		s.Require().NoError(err)
		fmt.Println("FINALLLLL 4")
		fmt.Println(b4.Balance.String())

	})

	//t.Run("send take swap message", func(t *testing.T) {
	//	sellToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
	//	senderAddressB := chainBWallet.Bech32Address("cosmos")
	//	senderReceivingAddressA := chainAWallet.Bech32Address("CosmosA")
	//
	//	timeoutHeight := clienttypes.NewHeight(0, 110)
	//	msgTake := types.NewMsgTakeSwap(
	//		channelA.PortID,
	//		channelA.ChannelID,
	//		sellToken,
	//		senderAddressB,
	//		senderReceivingAddressA,
	//		timeoutHeight,
	//		0,
	//		time.Now().UTC().Unix(),
	//	)
	//
	//	resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msgTake)
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Printf("Response from TakeSwap: %#v\n", resp)
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//	fmt.Println("-------------------")
	//
	//	s.AssertValidTxResponse(resp)
	//	s.Require().NoError(err)
	//
	//	// wait block when packet relay.
	//	test.WaitForBlocks(ctx, 10, chainA, chainB)
	//
	//	// check packet relay status.
	//	s.AssertPacketRelayed(ctx, chainB, channelA.PortID, channelA.ChannelID, 1)
	//
	//})
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
