package e2e

import (
	"context"
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

func (s *AtomicSwapTestSuite) TestAtomicSwap_HappyPath() {
	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	relayer, channelA := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())
	chainA, chainB := s.GetChains()

	//create wallets for testing of the maker address on chains A and B.
	chainAMakerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	makerAddressOnChainA := chainAMakerWallet.Bech32Address("cosmos")
	chainBMakerWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	makerReceivingAddressOnChainB := chainBMakerWallet.Bech32Address("cosmos")

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("happy path atomic swap", func(t *testing.T) {
		chainADenom := chainA.Config().Denom
		chainBDenom := chainB.Config().Denom

		//create wallets for testing of the taker address on chains A and B.
		chainBTakerWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		takerAddressOnChainB := chainBTakerWallet.Bech32Address("cosmos")
		chainATakerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		takerReceivingAddressOnChainA := chainATakerWallet.Bech32Address("cosmos")

		//get initial balances of the accounts that will be asserted after atomic swap.
		res1, err := s.QueryBalance(ctx, chainA, makerAddressOnChainA, chainADenom)
		s.Require().NoError(err)
		makerInitialBallanceOnChainA := res1.Balance.Amount
		res2, err := s.QueryBalance(ctx, chainB, makerReceivingAddressOnChainB, chainBDenom)
		s.Require().NoError(err)
		makerInitialBallanceOnChainB := res2.Balance.Amount
		res3, err := s.QueryBalance(ctx, chainA, takerReceivingAddressOnChainA, chainADenom)
		s.Require().NoError(err)
		takerInitialBallanceOnChainA := res3.Balance.Amount
		res4, err := s.QueryBalance(ctx, chainB, takerAddressOnChainB, chainBDenom)
		s.Require().NoError(err)
		takerInitialBallanceOnChainB := res4.Balance.Amount

		// Broadcast Make Swap transaction.
		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressOnChainA, makerReceivingAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		resp, err := s.BroadcastMessages(ctx, chainA, chainAMakerWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// broadcast TAKE SWAP transaction
		sellToken2 := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := types.NewAtomicOrder(types.NewMakerFromMsg(msg), msg.SourceChannel)
		msgTake := types.NewMsgTakeSwap(channelA.PortID, channelA.ChannelID, sellToken2, takerAddressOnChainB, takerReceivingAddressOnChainA, timeoutHeight2, 0, time.Now().UTC().Unix())
		msgTake.OrderId = order.Id
		resp2, err2 := s.BroadcastMessages(ctx, chainB, chainBTakerWallet, msgTake)
		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// Assert balances after atomic swap finished
		b1, err := s.QueryBalance(ctx, chainA, makerAddressOnChainA, chainA.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(makerInitialBallanceOnChainA.Int64()-100, b1.Balance.Amount.Int64())

		b2, err := s.QueryBalance(ctx, chainB, makerReceivingAddressOnChainB, chainB.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(makerInitialBallanceOnChainB.Int64()+50, b2.Balance.Amount.Int64())

		b3, err := s.QueryBalance(ctx, chainA, takerReceivingAddressOnChainA, chainA.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(takerInitialBallanceOnChainA.Int64()+100, b3.Balance.Amount.Int64())

		b4, err := s.QueryBalance(ctx, chainB, takerAddressOnChainB, chainB.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(takerInitialBallanceOnChainB.Int64()-50, b4.Balance.Amount.Int64())
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
