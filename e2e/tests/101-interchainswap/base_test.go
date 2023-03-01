package e2e

import (
	//"context"
	"context"
	"fmt"
	"testing"

	//"time"

	//"github.com/strangelove-ventures/ibctest/broadcast"
	//"github.com/strangelove-ventures/ibctest/chain/cosmos"
	//"github.com/strangelove-ventures/ibctest/ibc"
	//"github.com/strangelove-ventures/ibctest/test"
	//"github.com/stretchr/testify/suite"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	"github.com/strangelove-ventures/ibctest/ibc"
	"github.com/stretchr/testify/suite"

	types "github.com/sideprotocol/ibcswap/v4/modules/apps/101-interchain-swap/types"
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
	relayer, channelA, _ := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())

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
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)
		fmt.Println("OnCreate RES:", resp)
		fmt.Println("OnCreate ERR:", err)

		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 0)
		//s.Require().NoError(test.WaitForBlocks(ctx, 30, chainA, chainB), "failed to wait for blocks")
		// poolId := types.GetPoolId(msg.Denoms)
		// poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		// s.Require().NoError(err)
		// poolInfo := poolRes.InterchainLiquidityPool
		// s.Require().EqualValues(msg.SourceChannel, poolInfo.EncounterPartyChannel)
		// s.Require().EqualValues(msg.SourcePort, poolInfo.EncounterPartyPort)

	})

	t.Run("send deposit message", func(t *testing.T) {

		//beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		//s.Require().NoError(err)

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(100000)}
		msg := types.NewMsgDeposit(
			poolId,
			chainAAddress,
			[]*sdk.Coin{&depositCoin},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		fmt.Println("Deposit Res:", res)
		fmt.Println("Deposit ERR:", err)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		// s.Require().NoError(err)
		// expectedBalance := balanceRes.Balance.Add(depositCoin)
		// s.Require().Equal(expectedBalance.Denom, beforeDeposit.Balance.Denom)
		// s.Require().Equal(expectedBalance.Amount, beforeDeposit.Balance.Amount)
	})

	// //withdraw
	t.Run("send withdraw message", func(t *testing.T) {
		//beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		fmt.Println("PoolCoin:", poolCoin.Amount)

		denomOut := chainADenom
		sender := chainAAddress

		//coin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(100)}
		msg := types.NewMsgWithdraw(
			sender,
			poolCoin,
			denomOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		fmt.Println("Res:", resp)
		fmt.Println("ERR:", err)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		// s.Require().NoError(err)
		// s.Require().Equal(balanceRes.Balance.Denom, beforeDeposit.Balance.Denom)
		// s.Require().Equal(balanceRes.Balance.Amount, beforeDeposit.Balance.Amount)
	})

	// // // send swap message
	// t.Run("send swap message", func(t *testing.T) {
	// 	sender := chainAAddress
	// 	tokenIn := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}
	// 	tokenOut := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
	// 	msg := types.NewMsgSwap(
	// 		sender,
	// 		10,
	// 		sender,
	// 		&tokenIn,
	// 		&tokenOut,
	// 	)
	// 	resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
	// 	fmt.Println(resp)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)
	// })
}

// interchainswapChannelOptions configures both of the chains to have interchainswap enabled.
func interchainswapChannelOptions() func(options *ibc.CreateChannelOptions) {
	return func(opts *ibc.CreateChannelOptions) {
		opts.Version = "ics101-1"
		opts.Order = ibc.Unordered
		opts.DestPortName = types.ModuleName
		opts.SourcePortName = types.ModuleName
	}
}
