package e2e

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	clienttypes "github.com/cosmos/ibc-go/v6/modules/core/02-client/types"
	"github.com/gogo/protobuf/proto"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/100-atomic-swap/types"
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

		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

		res, err := s.QueryAtomicswapOrders(ctx, chainA)
		s.Require().NoError(err)
		s.Require().Equal(len(res.Orders), 1)

		res, err = s.QueryAtomicswapByOrders(ctx, chainA, types.OrderType_SellToBuy)
		s.Require().NoError(err)
		s.Require().Equal(len(res.Orders), 1)

		res, err = s.QueryAtomicswapByOrders(ctx, chainB, types.OrderType_BuyToSell)
		s.Require().NoError(err)
		s.Require().Equal(len(res.Orders), 1)

		res, err = s.QuerySubmittedAtomicswap(ctx, chainA, makerAddressOnChainA)
		s.Require().NoError(err)
		s.Require().Equal(len(res.Orders), 1)

		// broadcast TAKE SWAP transaction
		sellToken2 := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight2 := clienttypes.NewHeight(0, 110)
		orderRes, err := s.QueryAtomicswapOrders(ctx, chainA)
		s.Require().NoError(err)
		s.Require().NotZero(len(orderRes.Orders))

		fmt.Println("taker address", orderRes.Orders)
		msgTake := types.NewMsgTakeSwap(orderRes.Orders[0].Id, sellToken2, takerAddressOnChainB, takerReceivingAddressOnChainA, timeoutHeight2, 0, time.Now().UTC().Unix())
		resp2, err2 := s.BroadcastMessages(ctx, chainB, chainBTakerWallet, msgTake)

		s.AssertValidTxResponse(resp2)
		s.Require().NoError(err2)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		// check packet relay status.
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

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

func createOrder(msg *types.MakeSwapMsg, sequence uint64) types.Order {
	path := orderPath(msg.SourcePort, msg.SourceChannel, msg.SourcePort, msg.SourceChannel, sequence)
	return types.Order{
		Id:     generateOrderId(path, msg),
		Status: types.Status_INITIAL,
		Path:   path,
		Maker:  msg,
	}
}

func orderPath(sourcePort, sourceChannel, destPort, destChannel string, sequence uint64) string {
	return fmt.Sprintf("channel/%s/port/%s/channel/%s/port/%s/%d", sourceChannel, sourcePort, destChannel, destPort, sequence)
}

func generateOrderId(orderPath string, msg *types.MakeSwapMsg) string {
	prefix := []byte(orderPath)
	bytes, _ := proto.Marshal(msg)
	hash := sha256.Sum256(append(prefix, bytes...))
	return hex.EncodeToString(hash[:])
}
