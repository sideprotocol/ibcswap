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
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
	"github.com/stretchr/testify/suite"
)

func TestTakeAtomicSwapTestSuite(t *testing.T) {
	suite.Run(t, new(TakeAtomicSwapTestSuite))
}

type TakeAtomicSwapTestSuite struct {
	testsuite.E2ETestSuite
}

func (s *TakeAtomicSwapTestSuite) TestTakeAtomicSwap() {
	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	relayer, channelA := s.SetupChainsRelayerAndChannel(ctx, atomicSwapChannelOptions())
	chainA, chainB := s.GetChains()

	//create wallets for testing of the maker address on chains A and B.
	//chainAMakerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	//makerAddressOnChainA := chainAMakerWallet.Bech32Address("cosmos")
	//chainBMakerWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	//makerReceivingAddressOnChainB := chainBMakerWallet.Bech32Address("cosmos")

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("try to take order that is canceled", func(t *testing.T) {
		// Broadcast Make Swap transaction.
		makerWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		makerAddressChainA := makerWallet.Bech32Address("cosmos")
		makerWalletChainB := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		makerAddressOnChainB := makerWalletChainB.Bech32Address("cosmos")

		// create wallets for testing of the taker address on chains A and B.
		takerWalletChainB := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		chainBTakerAddress := takerWalletChainB.Bech32Address("cosmos")
		takerWalletChainA := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
		chainATakerReceivingAddress := takerWalletChainA.Bech32Address("cosmos")

		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(100))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressChainA, makerAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		response, err := s.BroadcastMessages(ctx, chainA, makerWallet, msg)
		s.AssertValidTxResponse(response)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		fmt.Println("Address; ", makerAddressChainA)

		// broadcast Cancel order
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		order := types.NewAtomicOrder(types.NewMakerFromMsg(msg), msg.SourceChannel)

		msgCancel := types.NewMsgCancelSwap(channelA.PortID, channelA.ChannelID, makerAddressChainA, order.Id, timeoutHeight2, 0)
		msgCancel.OrderId = order.Id
		fmt.Println("ORDER ID IN E2E E2E E2E: ", order.Id)
		resp2, err2 := s.BroadcastMessages(ctx, chainA, makerWallet, msgCancel)

		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// try to TAKE canceled order
		sellToken2 := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		msgTake := types.NewMsgTakeSwap(channelA.PortID, channelA.ChannelID, sellToken2, chainBTakerAddress, chainATakerReceivingAddress, timeoutHeight2, 0, time.Now().UTC().Unix())
		msgTake.OrderId = order.Id
		respTake, err := s.BroadcastMessages(ctx, chainB, takerWalletChainB, msgTake)
		s.Require().NoError(err)
		s.AssertValidTxResponse(resp2)
		fmt.Println("TAKE ORDER THAT IS CANCELED:----------------", respTake.RawLog)
		s.Require().Equal("failed to execute message; message index: 0: order is not in valid state", respTake.RawLog)
	})

}
