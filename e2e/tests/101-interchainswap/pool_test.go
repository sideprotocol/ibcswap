package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v6/ibc"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
)

func (s *InterchainswapTestSuite) TestPoolStatus() {

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
	fmt.Println(chainAAddressForB)

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

	const initialX = 2_000_000 // USDT
	const initialY = 1000      // ETH
	const fee = 300

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send create pool message", func(t *testing.T) {
		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"20:80",
			[]*sdk.Coin{
				{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
				{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
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
		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)}

		msg := types.NewMsgDeposit(
			poolId,
			chainBAddress,
			[]*sdk.Coin{&depositCoin},
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

	})

	// send swap message
	t.Run("pool status", func(t *testing.T) {

		var depositCoin sdk.Coin
		var chain *cosmos.CosmosChain
		var wallet ibc.Wallet

		poolId := types.GetPoolId([]string{
			chainADenom, chainBDenom,
		})

		var sender string
		var recipient string
		var tokenIn sdk.Coin
		var tokenOut *sdk.Coin
		var poolCoin sdk.Coin
		var denomOut string

		var channel ibc.ChannelOutput
		var txRes sdk.TxResponse
		var err error
		var packetSequence int
		fmt.Println(packetSequence)

		testCases := []struct {
			name     string
			malleate func()
			msgType  string
			expPass  bool
		}{
			{
				"deposit Asset A (initialX)",
				func() {
					depositCoin = sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)}
					wallet = *chainAWallet
					chain = chainA
					packetSequence = 2
				},
				"deposit",
				true,
			},
			{
				"deposit Asset B (initialY)",
				func() {
					depositCoin = sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)}
					wallet = *chainBWallet
					chain = chainB
					channel = channelB
					packetSequence = 2
				},
				"deposit",
				true,
			},
			{
				"swap Asset A(10_0000) to Asset B",
				func() {
					tokenIn = sdk.NewCoin(chainADenom, sdk.NewInt(10_0000))
					wallet = *chainAWallet
					chain = chainA
					sender = wallet.Bech32Address("cosmos")
					recipient = chainAAddressForB
					channel = channelA
					packetSequence = 3
				},
				"swap",
				true,
			},
			{
				"swap Asset A(100_0000) to Asset B",
				func() {
					tokenIn = sdk.NewCoin(chainADenom, sdk.NewInt(100_0000))
					wallet = *chainAWallet
					chain = chainA
					sender = wallet.Bech32Address("cosmos")
					recipient = chainAAddressForB
					channel = channelA
					packetSequence = 4
				},
				"swap",
				true,
			},
			{
				"deposit Asset A (100_000)",
				func() {
					depositCoin = sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(100_000)}
					wallet = *chainAWallet
					chain = chainA
					channel = channelA
					packetSequence = 5
				},
				"deposit",
				true,
			},
			{
				"withdraw Asset A (50%)",
				func() {
					// initial liquidity token remove:
					lpAmount, err := s.QueryBalance(ctx, chainA, chainAAddress, poolId)
					s.Require().NoError(err)
					halfAmount := lpAmount.Balance.Amount.Sub(sdk.NewInt(initialX)).Quo(sdk.NewInt(2))
					poolCoin = sdk.NewCoin(poolId, halfAmount)
					s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

					denomOut = chainADenom
					sender = chainAAddress
					wallet = *chainAWallet
					chain = chainA
					channel = channelA
					packetSequence = 6
				},
				"withdraw",
				true,
			},
			{
				"withdraw Asset B (20%)",
				func() {

					lpAmount, err := s.QueryBalance(ctx, chainB, chainBAddress, poolId)
					s.Require().NoError(err)
					halfAmount := lpAmount.Balance.Amount.Sub(sdk.NewInt(initialY)).Quo(sdk.NewInt(10))
					poolCoin = sdk.NewCoin(poolId, halfAmount)
					s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

					denomOut = chainBDenom
					wallet = *chainBWallet
					chain = chainB
					channel = channelB
					sender = wallet.Bech32Address("cosmos")
					packetSequence = 3
				},
				"withdraw",
				true,
			},
		}

		for _, tc := range testCases {

			fmt.Printf("=======%s=======\n\n", tc.name)
			tc.malleate()

			switch tc.msgType {
			case "deposit":
				msg := types.NewMsgDeposit(
					poolId,
					wallet.Bech32Address("cosmos"),
					[]*sdk.Coin{&depositCoin},
				)
				txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)

			case "swap":
				// pool status log.
				poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
				s.Require().NoError(err)
				pool := poolRes.InterchainLiquidityPool

				am := types.NewInterchainMarketMaker(
					&pool,
					fee,
				)

				tokenOut, _ = am.LeftSwap(tokenIn, chainBDenom)
				msg := types.NewMsgSwap(
					types.SwapMsgType_LEFT,
					sender,
					100,
					recipient,
					&tokenIn,
					tokenOut,
				)

				txRes, err := s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)

			case "withdraw":

				pool, err := getPool(
					s, ctx, chain, poolId,
				)
				s.Require().NoError(err)

				amm := types.NewInterchainMarketMaker(
					pool,
					fee,
				)

				outToken, err := amm.Withdraw(poolCoin, denomOut)
				s.Require().NoError(err)
				s.Require().NotEqual(outToken.Amount, sdk.NewInt(0))
				s.Require().NotNil(outToken)

				msg := types.NewMsgWithdraw(
					sender,
					&poolCoin,
					denomOut,
				)

				err = msg.ValidateBasic()
				s.Require().NoError(err)

				txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				//fmt.Println(txRes)
				s.AssertValidTxResponse(txRes)
				s.Require().NoError(err)
			}

			test.WaitForBlocks(ctx, 15, chainA, chainB)
			s.AssertPacketRelayed(ctx, chain, channel.PortID, channel.ChannelID, uint64(packetSequence))

			// pool status log.
			pool, err := getPool(s, ctx, chain, poolId)

			s.Require().NoError(err)
			amm := types.NewInterchainMarketMaker(
				pool,
				fee,
			)

			priceA_B, _ := amm.MarketPrice(chainADenom, chainBDenom)
			priceB_A, _ := amm.MarketPrice(chainADenom, chainBDenom)

			logger.CleanLog("Price: A->B, B->A", *priceA_B, *priceB_A)
			logger.CleanLog("Pool Info", pool)
		}

	})

}

func getPool(suit *InterchainswapTestSuite, ctx context.Context, chain *cosmos.CosmosChain, poolId string) (*types.InterchainLiquidityPool, error) {
	poolRes, err := suit.QueryInterchainswapPool(ctx, chain, poolId)
	if err != nil {
		return nil, err
	}
	pool := poolRes.InterchainLiquidityPool
	return &pool, nil
}
