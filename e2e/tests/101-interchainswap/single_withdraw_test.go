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

func (s *InterchainswapTestSuite) TestSingleWithdrawStatus() {

	t := s.T()
	ctx := context.TODO()
	logger := testsuite.NewLogger()
	// // setup relayers and connection-0 between two chains.
	relayer, channelA, channelB := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())
	_ = relayer
	_ = channelA
	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	chainBDenom := chainB.Config().Denom

	logger.CleanLog("get Prefix", chainB.Config().Bech32Prefix)

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddress := chainAWallet.Bech32Address(chainB.Config().Bech32Prefix)

	chainBUserMnemonic, err := createNewMnemonic()
	s.Require().NoError(err)

	chainBWallet := s.CreateUserOnChainBWithMnemonic(ctx, chainBUserMnemonic, testvalues.StartingTokenAmount)
	chainBAddress := chainBWallet.Bech32Address(chainB.Config().Bech32Prefix)
	priv, _ := getPrivFromNewMnemonic(chainBUserMnemonic)

	addr := sdk.AccAddress(priv.PubKey().Address().Bytes())
	logger.CleanLog("address:", addr.String())
	s.Require().Equal(chainBAddress, addr.String())

	chainAWalletForB := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddressForB := chainAWalletForB.Bech32Address(chainB.Config().Bech32Prefix)
	fmt.Println(chainAAddressForB)

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
			"50:50",
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

		escrowAccount := types.GetEscrowAddress(poolAInfo.EncounterPartyPort, poolAInfo.EncounterPartyChannel)

		// check escrow status
		escrowBalance, err := s.QueryBalance(ctx, chainA, escrowAccount.String(), chainADenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.NewInt(initialX), escrowBalance.Balance.Amount)

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

		msg := types.NewMsgSingleAssetDeposit(
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
		test.WaitForBlocks(ctx, 15, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

		poolResInChainA, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolInChainA := poolResInChainA.InterchainLiquidityPool
		s.Require().Equal(types.PoolStatus_POOL_STATUS_READY, poolInChainA.Status)

		poolResInChainB, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolInChainB := poolResInChainB.InterchainLiquidityPool
		s.Require().Equal(types.PoolStatus_POOL_STATUS_READY, poolInChainB.Status)

		logger.CleanLog("Send Deposit(After):PoolA", poolInChainA)
		logger.CleanLog("Send Deposit(After):PoolB", poolInChainB)

	})

	// single withdraw
	t.Run("send withdraw message", func(t *testing.T) {

		testCases := []struct {
			name    string
			chain   *cosmos.CosmosChain
			wallet  *ibc.Wallet       // Assuming Wallet type exists
			channel ibc.ChannelOutput // Assuming Channel type exists
			address string
			denom   string
		}{
			{
				name:    "Chain A Test",
				chain:   chainA,
				wallet:  chainAWallet,
				channel: channelA,
				address: chainAAddress,
				denom:   chainADenom,
			},
			{
				name:    "Chain B Test",
				chain:   chainB,
				wallet:  chainBWallet,
				channel: channelB,
				address: chainBAddress,
				denom:   chainBDenom,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				beforeWithdraw, err := s.QueryBalance(ctx, tc.chain, tc.address, tc.denom)
				s.Require().NoError(err)
				poolId := types.GetPoolId([]string{chainADenom, chainBDenom})

				lpBalanceOfSender, err := s.QueryBalance(ctx, tc.chain, tc.address, poolId)
				s.Require().NoError(err)

				poolRes, err := s.QueryInterchainswapPool(ctx, tc.chain, poolId)
				s.Require().NoError(err)
				pool := poolRes.InterchainLiquidityPool

				amm := types.NewInterchainMarketMaker(&pool, fee)
				out, err := amm.SingleWithdraw(*lpBalanceOfSender.Balance, tc.denom)
				s.Require().NoError(err)

				if tc.chain == chainA {
					s.Require().Equal(out.Amount, sdk.NewInt(initialX))
				} else {
					s.Require().Equal(out.Amount, sdk.NewInt(initialY))
				}

				msg := types.NewMsgSingleAssetWithdraw(
					tc.address,
					tc.denom,
					lpBalanceOfSender.Balance,
				)

				resp, err := s.BroadcastMessages(ctx, tc.chain, tc.wallet, msg)
				s.AssertValidTxResponse(resp)
				s.Require().NoError(err)

				// check packet relayed or not.
				test.WaitForBlocks(ctx, 10, chainA, chainB)
				s.AssertPacketRelayed(ctx, tc.chain, tc.channel.PortID, tc.channel.ChannelID, 2)

				// check token is withdrawn or not
				balanceRes, err := s.QueryBalance(ctx, tc.chain, tc.address, tc.denom)
				s.Require().NoError(err)
				s.Require().Equal(balanceRes.Balance.Denom, beforeWithdraw.Balance.Denom)
				logger.CleanLog("Withdraw Res", balanceRes.Balance, beforeWithdraw.Balance)
				s.Require().Equal(balanceRes.Balance.Amount.GT(beforeWithdraw.Balance.Amount), true)
				poolRes, err = s.QueryInterchainswapPool(ctx, tc.chain, poolId)
				if tc.name == "Chain A Test" {
					s.Require().NoError(err)
					asset, err := poolRes.InterchainLiquidityPool.FindAssetByDenom(tc.denom)
					s.Require().NoError(err)
					s.Require().Equal(asset.Balance.Amount, sdk.NewInt(0))
					s.Require().Equal(pool.Supply.Amount, sdk.NewInt(initialX+initialY))
				} else {
					s.Require().Error(err)
				}
			})
		}
	})
}
