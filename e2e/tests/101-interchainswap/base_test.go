package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	//atomicswaptypes "github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
	"github.com/stretchr/testify/suite"
)

func TestInterchainswapTestSuite(t *testing.T) {
	suite.Run(t, new(InterchainswapTestSuite))
}

type InterchainswapTestSuite struct {
	testsuite.E2ETestSuite
}

func (s *InterchainswapTestSuite) TestBasicMsgPacket() {

	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	relayer, channelA := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())

	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	chainBDenom := chainB.Config().Denom

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddress := chainAWallet.Bech32Address("cosmos")

	// allocate tokens to the new account
	initialBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(10000000000)))
	err := s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialBalances)

	s.Require().NoError(err)
	res, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	s.Require().NotEqual(res.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	//chainBAddress = s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)

	// s.Require().NoError(test.WaitForBlocks(ctx, 1, chainA, chainB), "failed to wait for blocks")

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send creat pool message", func(t *testing.T) {
		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"1:2",
			[]string{chainADenom, chainBDenom},
			[]uint32{10, 100},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		fmt.Println("tx:", resp)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)

		// check packet relay status.
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

		poolId := types.GetPoolId(msg.Denoms)
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolInfo := poolRes.InterchainLiquidityPool
		s.Require().EqualValues(msg.SourceChannel, poolInfo.EncounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolInfo.EncounterPartyPort)

	})

	t.Run("send deposit message", func(t *testing.T) {

		beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
		msg := types.NewMsgDeposit(
			poolId,
			chainAAddress,
			[]*sdk.Coin{&depositCoin},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)
		expectedBalance := balanceRes.Balance.Add(depositCoin)
		s.Require().Equal(expectedBalance.Denom, beforeDeposit.Balance.Denom)
		s.Require().Equal(expectedBalance.Amount, beforeDeposit.Balance.Amount)

		// check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolInfo := poolRes.InterchainLiquidityPool
		s.Require().NotEqual(poolInfo.Supply.Amount, sdk.NewInt(0))
	})

	//withdraw
	t.Run("send withdraw message", func(t *testing.T) {

		beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		denomOut := chainADenom
		sender := chainAAddress
		msg := types.NewMsgWithdraw(
			sender,
			poolCoin,
			denomOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)
		s.Require().Equal(balanceRes.Balance.Denom, beforeDeposit.Balance.Denom)
		s.Require().Equal(balanceRes.Balance.Amount, beforeDeposit.Balance.Amount)
	})

	// send swap message
	t.Run("send swap message (don't check ack)", func(t *testing.T) {
		sender := chainAAddress
		tokenIn := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
		tokenOut := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSwap(
			sender,
			10,
			sender,
			&tokenIn,
			&tokenOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)
	})

}

// interchainswapChannelOptions configures both of the chains to have interchainswap enabled.
func interchainswapChannelOptions() func(options *ibc.CreateChannelOptions) {
	return func(opts *ibc.CreateChannelOptions) {
		opts.SourcePortName = types.PortID
		opts.DestPortName = types.PortID
		opts.Order = ibc.Unordered
		opts.Version = types.Version
	}
}

// // interchainswapChannelOptions configures both of the chains to have interchainswap enabled.
// func atomicswapChannelOptions() func(options *ibc.CreateChannelOptions) {
// 	return func(opts *ibc.CreateChannelOptions) {
// 		opts.SourcePortName = atomicswaptypes.PortID
// 		opts.DestPortName = atomicswaptypes.PortID
// 		opts.Order = ibc.Unordered
// 		opts.Version = atomicswaptypes.Version
// 	}
// }
