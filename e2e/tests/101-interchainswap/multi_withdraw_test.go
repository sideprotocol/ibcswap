package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
)

func (s *InterchainswapTestSuite) TestMultiWithdrawStatus() {

	t := s.T()
	ctx := context.TODO()
	logger := testsuite.NewLogger()
	// // setup relayers and connection-0 between two chains.
	relayer, channelA, _ := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())
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

	// make force transaction to set pub key
	err = s.SendCoins(ctx, chainB, chainBWallet, chainAAddress, sdk.NewCoins(sdk.NewCoin(
		chainBDenom, sdk.NewInt(10),
	)))
	s.Require().NoError(err)

	t.Run("start relayer", func(t *testing.T) {
		s.StartRelayer(relayer)
	})

	t.Run("send create pool message", func(t *testing.T) {
		depositSignMsg := types.DepositSignature{
			Sender:   chainBAddress,
			Balance:  &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
			Sequence: 1,
		}

		rawDepositMsg, err := types.ModuleCdc.Marshal(&depositSignMsg)
		s.Require().NoError(err)
		signature, err := priv.Sign(rawDepositMsg)
		s.Require().NoError(err)

		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			chainBAddress,
			signature,
			types.PoolAsset{
				Side:    types.PoolAssetSide_SOURCE,
				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
				Weight:  50,
				Decimal: 6,
			},
			types.PoolAsset{
				Side:    types.PoolAssetSide_TARGET,
				Balance: &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
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
		poolId := types.GetPoolId(msg.GetLiquidityDenoms())
		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolAInfo := poolARes.InterchainLiquidityPool

		escrowAccount := types.GetEscrowAddress(poolAInfo.CounterPartyPort, poolAInfo.CounterPartyChannel)

		// check escrow status
		escrowBalance, err := s.QueryBalance(ctx, chainA, escrowAccount.String(), chainADenom)
		s.Require().NoError(err)
		s.Require().Equal(sdk.NewInt(initialX), escrowBalance.Balance.Amount)

		// check pool info sync status.

		s.Require().EqualValues(msg.SourceChannel, poolAInfo.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolAInfo.CounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[0].Amount, poolAInfo.Supply.Amount)

		poolBRes, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolBInfo := poolBRes.InterchainLiquidityPool
		s.Require().EqualValues(msg.SourceChannel, poolBInfo.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolBInfo.CounterPartyPort)
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

	// t.Run("send deposit message (enable pool)", func(t *testing.T) {

	// 	// check the balance of the chainA account.
	// 	beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
	// 	s.Require().NoError(err)
	// 	s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

	// 	// prepare deposit message.
	// 	poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
	// 	depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)}

	// 	msg := types.NewMsgSingleAssetDeposit(
	// 		poolId,
	// 		chainBAddress,
	// 		&depositCoin,
	// 	)
	// 	resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)

	// 	balanceRes, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
	// 	s.Require().NoError(err)
	// 	expectedBalance := balanceRes.Balance.Add(depositCoin)
	// 	s.Require().Equal(expectedBalance.Denom, beforeDeposit.Balance.Denom)
	// 	s.Require().Equal(expectedBalance.Amount, beforeDeposit.Balance.Amount)

	// 	// check packet relayed or not.
	// 	test.WaitForBlocks(ctx, 15, chainA, chainB)
	// 	s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

	// 	poolResInChainA, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
	// 	s.Require().NoError(err)
	// 	poolInChainA := poolResInChainA.InterchainLiquidityPool
	// 	s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainA.Status)

	// 	poolResInChainB, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
	// 	s.Require().NoError(err)
	// 	poolInChainB := poolResInChainB.InterchainLiquidityPool
	// 	s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainB.Status)

	// 	logger.CleanLog("Send Deposit(After):PoolA", poolInChainA)
	// 	logger.CleanLog("Send Deposit(After):PoolB", poolInChainB)

	// })

	// // single withdraw
	// t.Run("send withdraw message", func(t *testing.T) {

	// 	testCases := []struct {
	// 		name    string
	// 		chain   *cosmos.CosmosChain
	// 		wallet  *ibc.Wallet       // Assuming Wallet type exists
	// 		channel ibc.ChannelOutput // Assuming Channel type exists
	// 		address string
	// 		denom   string
	// 	}{
	// 		{
	// 			name:    "Chain A Test",
	// 			chain:   chainA,
	// 			wallet:  chainAWallet,
	// 			channel: channelA,
	// 			address: chainAAddress,
	// 			denom:   chainADenom,
	// 		},
	// 		{
	// 			name:    "Chain B Test",
	// 			chain:   chainB,
	// 			wallet:  chainBWallet,
	// 			channel: channelB,
	// 			address: chainBAddress,
	// 			denom:   chainBDenom,
	// 		},
	// 	}

	// 	for _, tc := range testCases {
	// 		t.Run(tc.name, func(t *testing.T) {
	// 			beforeWithdraw, err := s.QueryBalance(ctx, tc.chain, tc.address, tc.denom)
	// 			s.Require().NoError(err)
	// 			poolId := types.GetPoolId([]string{chainADenom, chainBDenom})

	// 			lpBalanceOfSender, err := s.QueryBalance(ctx, tc.chain, tc.address, poolId)
	// 			s.Require().NoError(err)

	// 			poolRes, err := s.QueryInterchainswapPool(ctx, tc.chain, poolId)
	// 			s.Require().NoError(err)
	// 			pool := poolRes.InterchainLiquidityPool

	// 			amm := types.NewInterchainMarketMaker(&pool)
	// 			out, err := amm.SingleWithdraw(*lpBalanceOfSender.Balance, tc.denom)
	// 			s.Require().NoError(err)
	// 			logger.CleanLog(tc.denom, out)
	// 			logger.CleanLog("Owned Pool token", *lpBalanceOfSender.Balance)
	// 			if tc.chain == chainA {
	// 				s.Require().Equal(out.Amount, sdk.NewInt(initialX))
	// 			} else {
	// 				//s.Require().Equal(out.Amount, sdk.NewInt(initialY))
	// 			}

	// 			msg := types.NewMsgSingleAssetWithdraw(
	// 				tc.address,
	// 				tc.denom,
	// 				lpBalanceOfSender.Balance,
	// 			)

	// 			resp, err := s.BroadcastMessages(ctx, tc.chain, tc.wallet, msg)
	// 			s.AssertValidTxResponse(resp)
	// 			s.Require().NoError(err)

	// 			// check packet relayed or not.
	// 			test.WaitForBlocks(ctx, 10, chainA, chainB)
	// 			s.AssertPacketRelayed(ctx, tc.chain, tc.channel.PortID, tc.channel.ChannelID, 2)

	// 			// check token is withdrawn or not
	// 			balanceRes, err := s.QueryBalance(ctx, tc.chain, tc.address, tc.denom)
	// 			s.Require().NoError(err)
	// 			s.Require().Equal(balanceRes.Balance.Denom, beforeWithdraw.Balance.Denom)
	// 			logger.CleanLog("Withdraw Res", balanceRes.Balance, beforeWithdraw.Balance)
	// 			s.Require().Equal(balanceRes.Balance.Amount.GT(beforeWithdraw.Balance.Amount), true)
	// 			poolRes, err = s.QueryInterchainswapPool(ctx, tc.chain, poolId)
	// 			logger.CleanLog("Pool Status after withdraw", poolRes)
	// 			s.Require().NoError(err)
	// 			if tc.name == "Chain A Test" {
	// 				s.Require().NoError(err)
	// 				asset, err := poolRes.InterchainLiquidityPool.FindAssetByDenom(tc.denom)
	// 				s.Require().Error(err)
	// 				s.Require().Equal(asset.Balance.Amount, sdk.NewInt(0))
	// 				s.Require().Equal(pool.Supply.Amount, sdk.NewInt(initialX+initialY))
	// 			} else {
	// 				s.Require().Error(err)
	// 			}
	// 		})
	// 	}
	// })

	t.Run("send multi deposit message", func(t *testing.T) {

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})

		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		pool := poolARes.InterchainLiquidityPool

		depositTokens := []*sdk.Coin{
			{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
			{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
		}

		remoteDepositTx := &types.DepositSignature{
			Sequence: 1,
			Sender:   chainBAddress,
			Balance:  &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
		}

		rawTx := types.ModuleCdc.MustMarshal(remoteDepositTx)
		if err != nil {
			fmt.Println(err)
			return
		}

		s.Require().NoError(err)
		signedTx, err := priv.Sign(rawTx)
		s.Require().NoError(err)
		pubKey := priv.PubKey()
		s.Require().Equal(verifySignedMessage(rawTx, signedTx, pubKey), true)

		amm := types.NewInterchainMarketMaker(&pool)

		outs, err := amm.DepositMultiAsset(depositTokens)
		s.Require().NoError(err)
		logger.CleanLog("=== multi-asset deposit===", outs)

		msg := types.NewMsgMultiAssetDeposit(
			poolId,
			[]string{
				chainAAddress,
				chainBAddress,
			},
			depositTokens,
			signedTx,
		)
		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		poolARes, err = s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		pool = poolARes.InterchainLiquidityPool
		amm = types.NewInterchainMarketMaker(&pool)
		outA, err := amm.MultiAssetWithdraw(*outs[0], chainADenom)
		s.Require().NoError(err)
		outB, err := amm.MultiAssetWithdraw(*outs[1], chainBDenom)
		s.Require().NoError(err)
		logger.CleanLog("====outA====", outA)
		logger.CleanLog("====outB====", outB)

		withdrawMsg := types.NewMsgMultiAssetWithdraw(
			poolId,
			chainAAddress,
			chainBAddress,
			outs[0],
			outs[1],
		)

		resp, err = s.BroadcastMessages(ctx, chainA, chainAWallet, withdrawMsg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// check packet relayed or not.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)
		poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		pool = poolRes.InterchainLiquidityPool

		for _, asset := range pool.Assets {
			if asset.Balance.Denom == chainADenom {
				s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialX))
			}
			if asset.Balance.Denom == chainBDenom {
				s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialY))
			}
		}
	})

	// Multi withdraw
	// t.Run("send withdraw message", func(t *testing.T) {
	// 	beforeWithdraw, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)
	// 	poolId := types.GetPoolId([]string{chainADenom, chainBDenom})

	// 	lpBalanceOfSenderInChainA, err := s.QueryBalance(ctx, chainA, chainAAddress, poolId)
	// 	s.Require().NoError(err)

	// 	lpBalanceOfSenderInChainB, err := s.QueryBalance(ctx, chainB, chainBAddress, poolId)
	// 	s.Require().NoError(err)

	// 	poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
	// 	s.Require().NoError(err)
	// 	pool := poolRes.InterchainLiquidityPool

	// 	amm := types.NewInterchainMarketMaker(&pool)
	// 	outOfChainA, err := amm.MultiAssetWithdraw(*lpBalanceOfSenderInChainA.Balance, chainADenom)
	// 	s.Require().NoError(err)
	// 	logger.CleanLog(chainADenom, outOfChainA)
	// 	logger.CleanLog("Owned Pool token", *lpBalanceOfSenderInChainA.Balance)

	// 	outOfChainB, err := amm.MultiAssetWithdraw(*lpBalanceOfSenderInChainB.Balance, chainBDenom)
	// 	s.Require().NoError(err)
	// 	logger.CleanLog(chainADenom, outOfChainA)
	// 	logger.CleanLog(chainBDenom, outOfChainB)
	// 	logger.CleanLog("Owned Pool token in chainA", *lpBalanceOfSenderInChainA.Balance)
	// 	logger.CleanLog("Owned Pool token in chainB ", *lpBalanceOfSenderInChainB.Balance)

	// 	// if tc.chain == chainA {
	// 	// 	s.Require().Equal(out.Amount, sdk.NewInt(initialX))
	// 	// } else {
	// 	// 	//s.Require().Equal(out.Amount, sdk.NewInt(initialY))
	// 	// }

	// 	msg := types.NewMsgMultiAssetWithdraw(
	// 		poolId,
	// 		chainAAddress,
	// 		chainBAddress,
	// 		lpBalanceOfSenderInChainA.Balance,
	// 		lpBalanceOfSenderInChainB.Balance,
	// 	)

	// 	resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)

	// 	// check packet relayed or not.
	// 	test.WaitForBlocks(ctx, 10, chainA, chainB)
	// 	s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

	// 	// check token is withdrawn or not
	// 	balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)
	// 	s.Require().Equal(balanceRes.Balance.Denom, beforeWithdraw.Balance.Denom)
	// 	logger.CleanLog("Withdraw Res", balanceRes.Balance, beforeWithdraw.Balance)
	// 	s.Require().Equal(balanceRes.Balance.Amount.GT(beforeWithdraw.Balance.Amount), true)
	// 	poolRes, err = s.QueryInterchainswapPool(ctx, chainA, poolId)
	// 	logger.CleanLog("Pool Status after withdraw", poolRes)
	// 	s.Require().NoError(err)
	// 	// if tc.name == "Chain A Test" {
	// 	// 	s.Require().NoError(err)
	// 	// 	asset, err := poolRes.InterchainLiquidityPool.FindAssetByDenom(tc.denom)
	// 	// 	s.Require().Error(err)
	// 	// 	s.Require().Equal(asset.Balance.Amount, sdk.NewInt(0))
	// 	// 	s.Require().Equal(pool.Supply.Amount, sdk.NewInt(initialX+initialY))
	// 	// } else {
	// 	// 	s.Require().Error(err)
	// 	// }
	// })
}
