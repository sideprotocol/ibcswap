package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"

	test "github.com/strangelove-ventures/ibctest/v6/testutil"

	//"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
)

func (s *InterchainswapTestSuite) TestSingleDepositStatus() {

	t := s.T()
	ctx := context.TODO()
	logger := testsuite.NewLogger()
	// // setup relayers and connection-0 between two chains.
	relayer, channelA, channelB := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())

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

	resA, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	s.Require().NotEqual(resA.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	resB, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
	s.Require().NotEqual(resB.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	// make force transaction to set pub key
	err = s.SendCoins(ctx, chainB, chainBWallet, chainAAddress, sdk.NewCoins(sdk.NewCoin(
		chainBDenom, sdk.NewInt(10),
	)))
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
				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
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

		// check pool info sync status.
		s.Require().EqualValues(msg.SourceChannel, poolA.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolA.CounterPartyPort)

		poolB := getFirstPool(s, ctx, chainB)

		s.Require().EqualValues(msg.SourceChannel, poolB.CounterPartyChannel)
		s.Require().EqualValues(msg.SourcePort, poolB.CounterPartyPort)
		//s.Require().EqualValues(msg.Tokens[1].Amount, poolBInfo.Supply.Amount)

		fmt.Println(poolB)
		logger.CleanLog("Create Pool: PoolA", poolA)
		fmt.Println("===================")
		logger.CleanLog("Create Pool: PoolB", poolB)

		// compare pool info sync status
		s.Require().EqualValues(poolA.Supply, poolB.Supply)
		s.Require().EqualValues(poolA.Assets[0].Balance.Amount, poolB.Assets[0].Balance.Amount)
		s.Require().EqualValues(poolA.Assets[1].Balance.Amount, poolB.Assets[1].Balance.Amount)
		s.Require().Equal(poolA.Status, types.PoolStatus_INITIALIZED)

		// check liquidity status in escrow account and my wallet.
		escrowAccount := types.GetEscrowAddress(poolA.CounterPartyPort, poolA.CounterPartyChannel)
		resA, err := s.QueryBalance(ctx, chainA, escrowAccount.String(), chainADenom)
		s.Require().NoError(err)
		for _, asset := range poolA.Assets {
			if asset.Balance.Denom == chainADenom {
				s.Require().Equal(asset.Balance, resA.Balance)
			}
		}
		s.Require().Equal(poolA.Supply.Amount, sdk.NewInt(initialX+initialY))
	})

	t.Run("send take pool message", func(t *testing.T) {

		pool := getFirstPool(s, ctx, chainA)
		msg := types.NewMsgTakePool(
			chainBAddress,
			pool.Id,
		)

		resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

		// check pool info in chainA and chainB
		poolA := getFirstPool(s, ctx, chainA)
		poolB := getFirstPool(s, ctx, chainB)

		fmt.Println(poolB)
		logger.CleanLog("Take Pool: PoolA", poolA)
		fmt.Println("===================")
		logger.CleanLog("Take Pool: PoolB", poolB)

		// compare pool info sync status
		s.Require().EqualValues(poolA.Supply, poolB.Supply)
		s.Require().EqualValues(poolA.Assets[0].Balance.Amount, poolB.Assets[0].Balance.Amount)
		s.Require().EqualValues(poolA.Assets[1].Balance.Amount, poolB.Assets[1].Balance.Amount)
		s.Require().Equal(poolA.Status, types.PoolStatus_ACTIVE)

		// check liquidity status in escrow account and my wallet.
		escrowAccount := types.GetEscrowAddress(poolA.CounterPartyPort, poolA.CounterPartyChannel)
		resA, err := s.QueryBalance(ctx, chainA, escrowAccount.String(), chainADenom)
		s.Require().NoError(err)

		resB, err := s.QueryBalance(ctx, chainB, escrowAccount.String(), chainBDenom)
		s.Require().NoError(err)

		for _, asset := range poolA.Assets {
			if asset.Balance.Denom == chainADenom {
				s.Require().Equal(asset.Balance, resA.Balance)

			}
			if asset.Balance.Denom == chainBDenom {
				s.Require().Equal(asset.Balance, resB.Balance)
			}
		}
		s.Require().Equal(poolA.Supply.Amount, sdk.NewInt(initialX+initialY))
	})

	// // send swap message
	t.Run("pool status", func(t *testing.T) {

		pool := getFirstPool(s, ctx, chainA)
		poolId := pool.Id

		depositTokens := []*sdk.Coin{
			{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
			{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
		}

		wallet := *chainAWallet
		chain := chainA
		channel := channelA
		packetSequence := 2

		testCases := []struct {
			name     string
			malleate func()
			msgType  string
			expPass  bool
		}{

			{
				"make deposit Assets (initialX,initialY)",
				func() {
				},
				"make multi-deposit",
				true,
			},

			{
				"take multi-deposit Assets (initialX,initialY)",
				func() {
					wallet = *chainBWallet
					chain = chainB
					channel = channelB
					packetSequence = 2
				},
				"take multi-deposit",
				true,
			},
		}

		for _, tc := range testCases {
			tc.malleate()
			switch tc.msgType {
			case "make multi-deposit":

				sourceAsset, err := pool.FindAssetBySide(types.PoolAssetSide_SOURCE)
				s.Require().NoError(err)
				destinationAsset, err := pool.FindAssetBySide(types.PoolAssetSide_DESTINATION)
				s.Require().NoError(err)

				currentRatio := sourceAsset.Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(destinationAsset.Amount)
				inputRatio := depositTokens[0].Amount.Mul(sdk.NewInt(types.Multiplier)).Quo(depositTokens[1].Amount)

				logger.CleanLog("=current ratio=", currentRatio)
				logger.CleanLog("=input ratio=", currentRatio)

				err = types.CheckSlippage(currentRatio, inputRatio, 10)
				s.NoError(err)
				msg := types.NewMsgMakeMultiAssetDeposit(
					poolId,
					[]string{
						chainAAddress,
						chainBAddress,
					},
					depositTokens,
				)

				txRes, err := s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)
			case "take multi-deposit":
				msg := types.NewMsgTakeMultiAssetDeposit(
					chainBAddress,
					poolId,
					0,
				)

				txRes, err := s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)
			}

			test.WaitForBlocks(ctx, 15, chainA, chainB)
			s.AssertPacketRelayed(ctx, chain, channel.PortID, channel.ChannelID, uint64(packetSequence))

			// pool status log.
			if tc.msgType == "take multi-deposit" {
				poolA := getFirstPool(s, ctx, chainA)
				poolB := getFirstPool(s, ctx, chainB)
				s.Require().NoError(err)
				amm := types.NewInterchainMarketMaker(&poolA)

				priceA_B, _ := amm.MarketPrice(chainADenom, chainBDenom)
				priceB_A, _ := amm.MarketPrice(chainBDenom, chainADenom)

				logger.CleanLog("Price: A->B, B->A", *priceA_B, *priceB_A)
				logger.CleanLog("PoolA", poolA)
				logger.CleanLog("PoolB", poolB)

				s.Require().EqualValues(poolA.Id, poolB.Id)

				for i := 0; i < len(poolA.Assets); i++ {
					s.Require().Equal(poolA.Assets[i].Balance.Amount, poolB.Assets[i].Balance.Amount)
				}

				// check balance update status
				for _, asset := range poolA.Assets {
					if asset.Balance.Denom == chainADenom {
						s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialX*2))
					}
					if asset.Balance.Denom == chainBDenom {
						s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialY*2))
					}
				}

			}

		}
	})

	t.Run("send withdraw asset message", func(t *testing.T) {

		pool := getFirstPool(s, ctx, chainA)
		sourceMakerPoolToken, err := s.QueryBalance(ctx, chainA, chainAAddress, pool.Id)
		s.Require().NoError(err)
		// destinationTakerPoolToken, err := s.QueryBalance(ctx, chainB, chainBAddress, pool.Id)
		// s.Require().NoError(err)

		msg := types.NewMsgMultiAssetWithdraw(
			pool.Id,
			chainAAddress,
			chainBAddress,
			sourceMakerPoolToken.Balance,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

		// check pool info in chainA and chainB
		poolA := getFirstPool(s, ctx, chainA)
		poolB := getFirstPool(s, ctx, chainB)

		fmt.Println(poolB)
		logger.CleanLog("Take Pool: PoolA", poolA)
		fmt.Println("===================")
		logger.CleanLog("Take Pool: PoolB", poolB)
	})

}
