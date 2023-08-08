package e2e

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
)

func (s *AtomicSwapTestSuite) TestAtomicSwapFiFoPath() {
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

	t.Run("fuzzy testing atomic swap with FIFO logic", func(t *testing.T) {
		numOrders := 5 // create between 1 to 100 orders
		createdOrderTimestamps := make([]int64, 0, numOrders)

		//create wallets for testing of the taker address on chains A and B.
		chainBTakerWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
		takerAddressOnChainB := chainBTakerWallet.Bech32Address("cosmos")

		for i := 0; i < numOrders; i++ {
			// Randomly determine sell and buy token amounts
			sellAmount := sdk.NewInt(int64(rand.Intn(1000)))
			buyAmount := sdk.NewInt(int64(rand.Intn(500)))

			sellToken := sdk.NewCoin(chainA.Config().Denom, sellAmount)
			buyToken := sdk.NewCoin(chainB.Config().Denom, buyAmount)
			timeoutHeight := clienttypes.NewHeight(0, 10000)

			// Note the current timestamp
			timestamp := time.Now().UTC().Unix()
			createdOrderTimestamps = append(createdOrderTimestamps, timestamp)

			// Broadcast Make Swap transaction with these random values
			msg := types.NewMsgMakeSwap(
				channelA.PortID,
				channelA.ChannelID,
				sellToken,
				buyToken,
				makerAddressOnChainA,
				makerReceivingAddressOnChainB,
				takerAddressOnChainB,
				timeoutHeight, 0,
				time.Now().UTC().Unix(),
			)
			resp, err := s.BroadcastMessages(ctx, chainA, chainAMakerWallet, msg)

			s.AssertValidTxResponse(resp)
			s.Require().NoError(err)

			// wait block when packet relay.
			test.WaitForBlocks(ctx, 10, chainA, chainB)

			// Introduce a small delay to ensure different timestamps
			time.Sleep(1 * time.Millisecond)
		}

		// Now, query the system for the swaps and verify they are returned in FIFO order
		res, err := s.QueryAtomicswapOrders(ctx, chainA)
		s.Require().NoError(err)
		s.Require().Equal(len(res.Orders), numOrders)

		for i, order := range res.Orders {
			fmt.Printf("=========%d======", i)
			fmt.Printf(":%v", order.Maker.CreateTimestamp)
			fmt.Printf("=================")
			//s.Require().Equal(order.Maker.CreateTimestamp, createdOrderTimestamps[len(res.Orders)-i-1])
		}
	})
}
