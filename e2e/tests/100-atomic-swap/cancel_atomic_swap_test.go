package e2e

import (
	"context"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
	"github.com/stretchr/testify/suite"
)

func TestCancelAtomicSwapTestSuite(t *testing.T) {
	suite.Run(t, new(CancelAtomicSwapTestSuite))
}

type CancelAtomicSwapTestSuite struct {
	testsuite.E2ETestSuite
}

func (s *CancelAtomicSwapTestSuite) TestCancelAtomicSwap() {
	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	relayer, channelA, _ := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())
	chainA, chainB := s.GetChains()

	//create wallets for testing of the maker address on chains A and B.
	chainAMakerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	makerAddressOnChainA := chainAMakerWallet.Bech32Address("cosmos")
	chainBMakerWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	makerReceivingAddressOnChainB := chainBMakerWallet.Bech32Address("cosmos")

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("cancel atomic swap", func(t *testing.T) {
		// get initial balances of the accounts that will be asserted after atomic swap.
		resp, err := s.QueryBalance(ctx, chainA, makerAddressOnChainA, chainA.Config().Denom)
		s.Require().NoError(err)
		initialMakerBalance := resp.Balance.Amount

		// Broadcast Make Swap transaction.
		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msgMake := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressOnChainA, makerReceivingAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		response, err := s.BroadcastMessages(ctx, chainA, chainAMakerWallet, msgMake)
		s.AssertValidTxResponse(response)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// broadcast TAKE SWAP transaction
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := createOrder(msgMake, 1)
		msgCancel := types.NewMsgCancelSwap(makerAddressOnChainA, order.Id, timeoutHeight2, 0)
		resp2, err2 := s.BroadcastMessages(ctx, chainA, chainAMakerWallet, msgCancel)

		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// Assert balances after atomic swap finished
		b1, err := s.QueryBalance(ctx, chainA, makerAddressOnChainA, chainA.Config().Denom)
		s.Require().NoError(err)
		s.Require().Equal(initialMakerBalance.Int64(), b1.Balance.Amount.Int64())
	})

	t.Run("cancel atomic swap that doesn't exist", func(t *testing.T) {
		// Broadcast Make Swap transaction.
		makerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		makerAddressChainA := makerWallet.Bech32Address("cosmos")
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		orderId := "C8EA430579423B692745DA0820948B9C40936E4010BFD2CF21A4768312A513ED"
		msgCancel := types.NewMsgCancelSwap(makerAddressChainA, orderId, timeoutHeight2, 0)

		resp2, err2 := s.BroadcastMessages(ctx, chainA, makerWallet, msgCancel)
		s.Require().NoError(err2)
		s.Require().Equal("failed to execute message; message index: 0: Make Order does not exist", resp2.RawLog)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
	})

	t.Run("cancel atomic swap that was canceled", func(t *testing.T) {
		// Broadcast Make Swap transaction.
		makerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		makerAddressChainA := makerWallet.Bech32Address("cosmos")
		makerWalletChainB := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		makerAddressOnChainB := makerWalletChainB.Bech32Address("cosmos")

		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressChainA, makerAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		response, err := s.BroadcastMessages(ctx, chainA, makerWallet, msg)
		s.AssertValidTxResponse(response)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// broadcast Cancel order
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := createOrder(msg, 3)

		msgCancel := types.NewMsgCancelSwap(makerAddressChainA, order.Id, timeoutHeight2, 0)
		resp2, err2 := s.BroadcastMessages(ctx, chainA, makerWallet, msgCancel)

		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// broadcast Cancel order for the second time
		//msgCancel2 := types.NewMsgCancelSwap(channelA.PortID, channelA.ChannelID, makerAddressChainA, order.Id, timeoutHeight2, 0)
		resp, err := s.BroadcastMessages(ctx, chainA, makerWallet, msgCancel)
		s.Require().NoError(err)
		s.Require().Equal("failed to execute message; message index: 0: order is not in a valid state for cancellation", resp.RawLog)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
	})

	t.Run("cancel atomic swap from sender that is not maker of the order", func(t *testing.T) {
		// Broadcast Make Swap transaction.
		makerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		makerAddressChainA := makerWallet.Bech32Address("cosmos")
		makerWalletChainB := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		makerAddressOnChainB := makerWalletChainB.Bech32Address("cosmos")

		// create wallets for testing of the taker address on chains A and B.
		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressChainA, makerAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		response, err := s.BroadcastMessages(ctx, chainA, makerWallet, msg)
		s.AssertValidTxResponse(response)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// broadcast Cancel order
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := createOrder(msg, 5)

		msgCancel := types.NewMsgCancelSwap(makerAddressChainA, order.Id, timeoutHeight2, 0)
		resp2, err2 := s.BroadcastMessages(ctx, chainA, makerWallet, msgCancel)
		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
	})
}
