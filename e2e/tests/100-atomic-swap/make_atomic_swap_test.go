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
	"github.com/stretchr/testify/suite"
)

func TestMakeAtomicSwapTestSuite(t *testing.T) {
	suite.Run(t, new(MakeAtomicSwapTestSuite))
}

type MakeAtomicSwapTestSuite struct {
	testsuite.E2ETestSuite
}

func (s *MakeAtomicSwapTestSuite) TestMakeAtomicSwap_EdgeCases() {
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

	t.Run("maker try to sell more tokens then he has in his account", func(t *testing.T) {
		// Broadcast Make Swap transaction.
		sellToken := sdk.NewCoin(chainA.Config().Denom, sdk.NewInt(200_000_000))
		buyToken := sdk.NewCoin(chainB.Config().Denom, sdk.NewInt(50))
		timeoutHeight := clienttypes.NewHeight(0, 110)
		msg := types.NewMsgMakeSwap(channelA.PortID, channelA.ChannelID, sellToken, buyToken, makerAddressOnChainA, makerReceivingAddressOnChainB, "", timeoutHeight, 0, time.Now().UTC().Unix())
		resp, err := s.BroadcastMessages(ctx, chainA, chainAMakerWallet, msg)
		s.Require().NoError(err)
		s.Equal("failed to execute message; message index: 0: insufficient balance", resp.RawLog)
	})
}
