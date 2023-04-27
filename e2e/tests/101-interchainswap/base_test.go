package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

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
	logger := testsuite.NewLogger()
	// setup relayers and connection-0 between two chains.
	relayer, channelA, channelB := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())

	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	chainBDenom := chainB.Config().Denom

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddress := chainAWallet.Bech32Address("cosmos")
	chainBWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	chainBAddress := chainBWallet.Bech32Address("cosmos")

	chainAWalletForB := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddressForB := chainAWalletForB.Bech32Address("cosmos")

	// allocate tokens to the new account
	initialATokenBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(1000000000000)))
	err := s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialATokenBalances)
	s.Require().NoError(err)
	// allocate tokens to the new account
	initialBTokenBalances := sdk.NewCoins(sdk.NewCoin(chainBDenom, sdk.NewInt(1000000000000)))
	err = s.SendCoinsFromModuleToAccount(ctx, chainB, chainBWallet, initialBTokenBalances)
	s.Require().NoError(err)

	resA, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	s.Require().NotEqual(resA.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	resB, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
	s.Require().NotEqual(resB.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	//chainBAddress = s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)

	// s.Require().NoError(test.WaitForBlocks(ctx, 1, chainA, chainB), "failed to wait for blocks")

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send create pool message", func(t *testing.T) {
		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"50:50",
			[]*sdk.Coin{
				{Denom: chainADenom, Amount: sdk.NewInt(100000)},
				{Denom: chainBDenom, Amount: sdk.NewInt(10000)},
			},
			[]uint32{6, 6},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

		// check pool info in chainA and chainB
		poolId := types.GetPoolIdWithTokens(msg.Tokens)
		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolAInfo := poolARes.InterchainLiquidityPool

		// check pool info sync status.

		s.Require().EqualValues(msg.SourceChannel, poolAInfo.EncounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolAInfo.EncounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[0].Amount, poolAInfo.Supply.Amount)

		poolBRes, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolBInfo := poolBRes.InterchainLiquidityPool
		s.Require().EqualValues(msg.SourceChannel, poolBInfo.EncounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolBInfo.EncounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[1].Amount, poolBInfo.Supply.Amount)

		fmt.Println(poolAInfo)
		logger.CleanLog("Create Pool: PoolA", poolAInfo)
		fmt.Println("===================")
		logger.CleanLog("Create Pool: PoolB", poolBInfo)

		// compare pool info sync status
		s.Require().EqualValues(poolAInfo.Supply, poolBInfo.Supply)
		s.Require().EqualValues(poolAInfo.Assets[0].Balance.Amount, poolBInfo.Assets[0].Balance.Amount)
		s.Require().EqualValues(poolAInfo.Assets[1].Balance.Amount, poolBInfo.Assets[1].Balance.Amount)
	})

	t.Run("send deposit message (enable pool)", func(t *testing.T) {

		// check the balance of the chainA account.
		beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

		// prepare deposit message.
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(10000)}

		msg := types.NewMsgSingleDeposit(
			poolId,
			chainBAddress,
			&depositCoin,
		)
		resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		balanceRes, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		expectedBalance := balanceRes.Balance.Add(depositCoin)
		s.Require().Equal(expectedBalance.Denom, beforeDeposit.Balance.Denom)
		s.Require().Equal(expectedBalance.Amount, beforeDeposit.Balance.Amount)

		// check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 1)

		poolResInChainA, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolInChainA := poolResInChainA.InterchainLiquidityPool
		s.Require().Equal(poolInChainA.Status, types.PoolStatus_POOL_STATUS_READY)

		poolResInChainB, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolInChainB := poolResInChainB.InterchainLiquidityPool
		s.Require().Equal(poolInChainB.Status, types.PoolStatus_POOL_STATUS_READY)

		logger.CleanLog("Send Deposit(After):PoolA", poolInChainA)
		logger.CleanLog("Send Deposit(After):PoolB", poolInChainB)
		fmt.Print(chainAAddressForB)

	})

	// send swap message
	t.Run("send swap message", func(t *testing.T) {
		sender := chainBAddress
		recipient := chainAAddressForB
		// check chain A pool status
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolA := poolARes.InterchainLiquidityPool

		tokenIn := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}

		// calculate estimate output.
		am := types.NewInterchainMarketMaker(
			&poolA,
			types.DefaultMaxFeeRate,
		)

		outToken, err := am.LeftSwap(tokenIn, chainADenom)
		logger.CleanLog("Swap:(Estimated output token)", outToken)

		s.Require().NoError(err)
		s.Require().Greater(outToken.Amount.Uint64(), uint64(0))

		tokenOut := outToken

		senderBefore, err := s.QueryBalance(ctx, chainB, chainBWallet.Bech32Address("cosmos"), chainBDenom)
		s.Require().NoError(err)
		// check chain B pool status
		poolBRes, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolB := poolBRes.InterchainLiquidityPool
		assetBBefore, err := poolB.FindAssetByDenom(chainBDenom)
		s.Require().NoError(err)

		logger.CleanLog("Swap(Before): poolA:", poolA)
		fmt.Println("------------------------------------")
		logger.CleanLog("Swap(Before): poolB:", poolB)

		s.Require().NotEqual(poolA.Supply.Amount, sdk.NewInt(0))
		s.Require().NotEqual(poolB.Supply.Amount, sdk.NewInt(0))

		// swap
		msg := types.NewMsgSwap(
			types.SwapMsgType_LEFT,
			sender,
			1,
			recipient,
			&tokenIn,
			tokenOut,
		)

		resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)
		//fmt.Println("Swap Tx:", resp)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// check token deposit state from wallet
		senderAfter, err := s.QueryBalance(ctx, chainB, chainBWallet.Bech32Address("cosmos"), chainBDenom)
		s.Require().NoError(err)

		differ := senderBefore.Balance.Amount.Sub(senderAfter.Balance.Amount)
		logger.CleanLog("SendDeposit:", differ.Int64())
		s.Require().Equal(differ.Int64(), int64(1000))

		// // check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

		// check withdraw status of a recipient.
		recipientATokenRes, err := s.QueryBalance(ctx, chainA, recipient, chainADenom)
		s.Require().NoError(err)
		s.Require().Greater(recipientATokenRes.Balance.Amount.Int64(), testvalues.StartingTokenAmount)

		// check chain A pool status
		poolARes, err = s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolA = poolARes.InterchainLiquidityPool

		// check chain B pool status
		poolBRes, err = s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolB = poolBRes.InterchainLiquidityPool

		assetBAfter, err := poolB.FindAssetByDenom(chainBDenom)

		s.Require().NoError(err)
		s.Require().Greater(assetBAfter.Balance.Amount.Uint64(), assetBBefore.Balance.Amount.Uint64())

		logger.CleanLog("Swap(After): poolA:", poolA)
		fmt.Println("------------------------------------")
		logger.CleanLog("Swap(After): poolB:", poolB)
	})

	// full withdraw
	t.Run("send withdraw message", func(t *testing.T) {

		beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolCoin := sdk.NewCoin(poolId, sdk.NewInt(1000))
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		denomOut := chainADenom
		sender := chainAAddress
		msg := types.NewMsgWithdraw(
			sender,
			&poolCoin,
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

}

func (s *InterchainswapTestSuite) TestBasicMsgPacketErrors() {
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
	chainAInvalidAddress := "(invalidAddress)"

	// allocate tokens to the new account
	initialBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(10000000000)))
	err := s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialBalances)

	s.Require().NoError(err)
	res, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	s.Require().NotEqual(res.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send create pool message with invalid address", func(t *testing.T) {
		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"50:50",
			[]*sdk.Coin{
				{Denom: chainADenom, Amount: sdk.NewInt(100000)},
				{Denom: chainBDenom, Amount: sdk.NewInt(10000)},
			},
			[]uint32{6, 6},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)
	})

	t.Run("send deposit message with invalid address", func(t *testing.T) {

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSingleDeposit(
			poolId,
			chainAInvalidAddress,
			&depositCoin,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send withdraw message with invalid address", func(t *testing.T) {
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		denomOut := chainADenom
		sender := chainAInvalidAddress
		msg := types.NewMsgWithdraw(
			sender,
			poolCoin,
			denomOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send swap message (don't check ack) with invalid address", func(t *testing.T) {
		sender := chainAInvalidAddress
		tokenIn := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}
		tokenOut := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSwap(
			types.SwapMsgType_LEFT,
			sender,
			10,
			sender,
			&tokenIn,
			&tokenOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)
	})

	t.Run("send deposit message with invalid denon", func(t *testing.T) {

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: "astros", Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSingleDeposit(
			poolId,
			chainAAddress,
			&depositCoin,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("failed to execute message; message index: 0: Invalid token amount", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send withdraw message with invalid denom", func(t *testing.T) {
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		denomOut := "astros"
		sender := chainAInvalidAddress
		msg := types.NewMsgWithdraw(
			sender,
			poolCoin,
			denomOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("failed to execute message; message index: 0: Invalid token amount", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send swap message (don't check ack) with invalid amount", func(t *testing.T) {
		sender := chainAAddress
		tokenIn := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000000000000)}
		tokenOut := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}

		msg := types.NewMsgSwap(
			types.SwapMsgType_LEFT,
			sender,
			10,
			sender,
			&tokenIn,
			&tokenOut,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("failed to execute message; message index: 0: 99998000atoma is smaller than 1000000000000atoma: insufficient funds", resp.RawLog)
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
