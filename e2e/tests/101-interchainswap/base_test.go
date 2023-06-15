package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
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

	chainBUserMnemonic, err := createNewMnemonic()
	s.Require().NoError(err)
	chainBWallet := s.CreateUserOnChainBWithMnemonic(ctx, chainBUserMnemonic, testvalues.StartingTokenAmount)
	chainBAddress := chainBWallet.Bech32Address(chainB.Config().Bech32Prefix)
	//priv, _ := getPrivFromNewMnemonic(chainBUserMnemonic)

	chainAWalletForB := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddressForB := chainAWalletForB.Bech32Address("cosmos")

	// allocate tokens to the new account
	initialATokenBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(1000000000000)))
	err = s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialATokenBalances)
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

	const initialX = 2_000_000 // USDT
	const initialY = 1000      // ETH

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send make pool message", func(t *testing.T) {

		msg := types.NewMsgMakePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			chainBAddress,
			types.PoolAsset{
				Side:    types.PoolAssetSide_SOURCE,
				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(20000)},
				Weight:  50,
				Decimal: 6,
			},
			types.PoolAsset{
				Side:    types.PoolAssetSide_DESTINATION,
				Balance: &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)},
				Weight:  50,
				Decimal: 6,
			},
			300,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

		// check pool info in chainA and chainB
		poolA := getFirstPool(s, ctx, chainA)

		s.Require().EqualValues(msg.SourceChannel, poolA.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolA.CounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[0].Amount, poolAInfo.Supply.Amount)

		poolB := getFirstPool(s, ctx, chainB)
		s.Require().EqualValues(msg.SourceChannel, poolB.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolB.CounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[1].Amount, poolBInfo.Supply.Amount)

		logger.CleanLog("Create Pool: PoolA", poolA)
		fmt.Println("===================")
		logger.CleanLog("Create Pool: PoolB", poolB)

		// compare pool info sync status
		s.Require().EqualValues(poolA.Supply, poolB.Supply)
		s.Require().EqualValues(poolA.Assets[0].Balance.Amount, poolB.Assets[0].Balance.Amount)
		s.Require().EqualValues(poolA.Assets[1].Balance.Amount, poolB.Assets[1].Balance.Amount)
	})

	t.Run("send take pool message (enable pool)", func(t *testing.T) {

		// check the balance of the chainA account.
		beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

		pool := getFirstPool(s, ctx, chainA)
		// prepare deposit message.
		//poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}

		msg := types.NewMsgTakePool(chainBAddress, pool.Id)
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

		poolA := getFirstPool(s, ctx, chainA)
		s.Require().Equal(poolA.Status, types.PoolStatus_ACTIVE)

		poolB := getFirstPool(s, ctx, chainB)
		s.Require().NoError(err)

		s.Require().Equal(poolB.Status, types.PoolStatus_ACTIVE)

		logger.CleanLog("Send Deposit(After):PoolA", poolA)
		logger.CleanLog("Send Deposit(After):PoolB", poolB)
		fmt.Print(chainAAddressForB)
	})

	// send swap message
	t.Run("send swap message", func(t *testing.T) {
		sender := chainBAddress
		recipient := chainAAddressForB
		// check chain A pool status
		poolA := getFirstPool(s, ctx, chainA)

		tokenIn := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}

		// calculate estimate output.
		am := types.NewInterchainMarketMaker(
			&poolA,
		)

		outToken, err := am.LeftSwap(tokenIn, chainADenom)
		logger.CleanLog("Swap:(Estimated output token)", outToken)

		s.Require().NoError(err)
		s.Require().Greater(outToken.Amount.Uint64(), uint64(0))

		tokenOut := outToken

		senderBefore, err := s.QueryBalance(ctx, chainB, chainBWallet.Bech32Address("cosmos"), chainBDenom)
		s.Require().NoError(err)
		// check chain B pool status
		poolB := getFirstPool(s, ctx, chainB)
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
			poolA.Id,
			300,
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

		// check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

		// check withdraw status of a recipient.
		recipientATokenRes, err := s.QueryBalance(ctx, chainA, recipient, chainADenom)
		s.Require().NoError(err)
		s.Require().Greater(recipientATokenRes.Balance.Amount.Int64(), testvalues.StartingTokenAmount)

		// check chain A pool status

		poolA = getFirstPool(s, ctx, chainA)

		// check chain B pool status
		poolB = getFirstPool(s, ctx, chainB)

		assetBAfter, err := poolB.FindAssetByDenom(chainBDenom)

		s.Require().NoError(err)
		s.Require().Greater(assetBAfter.Balance.Amount.Uint64(), assetBBefore.Balance.Amount.Uint64())

		logger.CleanLog("Swap(After): poolA:", poolA)
		fmt.Println("------------------------------------")
		logger.CleanLog("Swap(After): poolB:", poolB)
	})

	// single withdraw
	t.Run("send multi asset withdraw message", func(t *testing.T) {
		pool := getFirstPool(s, ctx, chainA)
		amm := types.NewInterchainMarketMaker(&pool)

		chainAPoolToken, err := s.QueryBalance(ctx, chainA, chainAAddress, pool.Id)
		s.Require().NoError(err)
		chainBPoolToken, err := s.QueryBalance(ctx, chainB, chainBAddress, pool.Id)
		s.Require().NoError(err)

		msg := types.NewMsgMultiAssetWithdraw(
			pool.Id,
			chainAAddress,
			chainBAddress,
			chainAPoolToken.Balance,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		pool = getFirstPool(s, ctx, chainA)
		amm = types.NewInterchainMarketMaker(&pool)
		outA, err := amm.MultiAssetWithdraw(*chainAPoolToken.Balance, chainADenom)
		s.Require().NoError(err)
		outB, err := amm.MultiAssetWithdraw(*chainBPoolToken.Balance, chainBDenom)
		s.Require().NoError(err)
		logger.CleanLog("====outA====", outA)
		logger.CleanLog("====outB====", outB)

		withdrawMsg := types.NewMsgMultiAssetWithdraw(
			pool.Id,
			chainAAddress,
			chainBAddress,
			chainAPoolToken.Balance,
		)

		resp, err = s.BroadcastMessages(ctx, chainA, chainAWallet, withdrawMsg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

		pool = getFirstPool(s, ctx, chainA)
		logger.CleanLog("pool information", pool)

		// withdraw
		for _, asset := range pool.Assets {
			if asset.Balance.Denom == chainADenom {
				s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialX))
			}
			if asset.Balance.Denom == chainBDenom {
				s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialY))
			}
		}
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
	chainBInvalidAddress := "(invalidAddress)"

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
		msg := types.NewMsgMakePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"",
			types.PoolAsset{
				Side:    types.PoolAssetSide_SOURCE,
				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(20000)},
				Weight:  50,
				Decimal: 6,
			},
			types.PoolAsset{
				Side:    types.PoolAssetSide_DESTINATION,
				Balance: &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)},
				Weight:  50,
				Decimal: 6,
			},
			300,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)
	})

	t.Run("send deposit message with invalid address", func(t *testing.T) {

		pool := getFirstPool(s, ctx, chainA)
		depositCoin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSingleAssetDeposit(
			pool.Id,
			chainAInvalidAddress,
			&depositCoin,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("invalid address", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send withdraw message with invalid address", func(t *testing.T) {
		pool := getFirstPool(s, ctx, chainA)
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, pool.Id)
		s.Require().NoError(err)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		msg := types.NewMsgMultiAssetWithdraw(
			pool.Id,
			chainAInvalidAddress,
			chainBInvalidAddress,
			poolCoin,
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
			"test",
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

		pool := getFirstPool(s, ctx, chainA)
		depositCoin := sdk.Coin{Denom: "astros", Amount: sdk.NewInt(1000)}
		msg := types.NewMsgSingleAssetDeposit(
			pool.Id,
			chainAAddress,
			&depositCoin,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("failed to execute message; message index: 0: Invalid token amount", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send withdraw message with invalid denom", func(t *testing.T) {
		pool := getFirstPool(s, ctx, chainA)
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, pool.Id)
		s.Require().NoError(err)
		poolCoin := poolRes.InterchainLiquidityPool.Supply
		s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

		msg := types.NewMsgMultiAssetWithdraw(
			pool.Id,
			chainAInvalidAddress,
			chainBInvalidAddress,
			poolCoin,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.Require().Equal("failed to execute message; message index: 0: Invalid token amount", resp.RawLog)
		s.Require().NoError(err)

	})

	t.Run("send swap message (don't check ack) with invalid amount", func(t *testing.T) {
		sender := chainAAddress
		tokenIn := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000000000000)}
		tokenOut := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}
		pool := getFirstPool(s, ctx, chainA)
		msg := types.NewMsgSwap(
			types.SwapMsgType_LEFT,
			pool.Id,
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

func getFirstPool(s *InterchainswapTestSuite, ctx context.Context, chain *cosmos.CosmosChain) types.InterchainLiquidityPool {
	// check pool info in chainA and chainB
	poolsRes, err := s.QueryInterchainswapPools(ctx, chain)
	s.Require().NoError(err)
	pools := poolsRes.InterchainLiquidityPool
	s.Require().Greater(len(pools), 0)
	return pools[0]
}
