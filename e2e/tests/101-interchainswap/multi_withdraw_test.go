package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
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

	t.Run("send multi asset withdraw message", func(t *testing.T) {

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
