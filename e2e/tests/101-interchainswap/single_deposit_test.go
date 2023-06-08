package e2e

// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	sdk "github.com/cosmos/cosmos-sdk/types"
// 	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"

// 	test "github.com/strangelove-ventures/ibctest/v6/testutil"

// 	//"github.com/cosmos/ibc-go/e2e/testsuite"
// 	"github.com/cosmos/ibc-go/e2e/testsuite"
// 	"github.com/cosmos/ibc-go/e2e/testvalues"
// )

// func (s *InterchainswapTestSuite) TestSingleDepositStatus() {

// 	t := s.T()
// 	ctx := context.TODO()
// 	logger := testsuite.NewLogger()
// 	// // setup relayers and connection-0 between two chains.
// 	relayer, channelA, channelB := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())
// 	_ = relayer
// 	_ = channelA
// 	chainA, chainB := s.GetChains()

// 	chainADenom := chainA.Config().Denom
// 	chainBDenom := chainB.Config().Denom

// 	logger.CleanLog("get Prefix", chainB.Config().Bech32Prefix)

// 	// // create wallets for testing
// 	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
// 	chainAAddress := chainAWallet.Bech32Address(chainB.Config().Bech32Prefix)

// 	chainBUserMnemonic, err := createNewMnemonic()
// 	s.Require().NoError(err)

// 	chainBWallet := s.CreateUserOnChainBWithMnemonic(ctx, chainBUserMnemonic, testvalues.StartingTokenAmount)
// 	chainBAddress := chainBWallet.Bech32Address(chainB.Config().Bech32Prefix)
// 	priv, _ := getPrivFromNewMnemonic(chainBUserMnemonic)

// 	addr := sdk.AccAddress(priv.PubKey().Address().Bytes())
// 	logger.CleanLog("address:", addr.String())
// 	s.Require().Equal(chainBAddress, addr.String())

// 	chainAWalletForB := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
// 	chainAAddressForB := chainAWalletForB.Bech32Address(chainB.Config().Bech32Prefix)
// 	fmt.Println(chainAAddressForB)

// 	resA, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
// 	s.Require().NotEqual(resA.Balance.Amount, sdk.NewInt(0))
// 	s.Require().NoError(err)

// 	resB, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
// 	s.Require().NotEqual(resB.Balance.Amount, sdk.NewInt(0))
// 	s.Require().NoError(err)

// 	const initialX = 2_000_000 // USDT
// 	const initialY = 1000      // ETH
// 	//const fee = 300

// 	// make force transaction to set pub key
// 	err = s.SendCoins(ctx, chainB, chainBWallet, chainAAddress, sdk.NewCoins(sdk.NewCoin(
// 		chainBDenom, sdk.NewInt(10),
// 	)))
// 	s.Require().NoError(err)

// 	t.Run("start relayer", func(t *testing.T) {
// 		s.StartRelayer(relayer)
// 	})

// 	t.Run("send create pool message", func(t *testing.T) {
// 		depositSignMsg := types.DepositSignature{
// 			Sender:   chainBAddress,
// 			Balance:  &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
// 			Sequence: 1,
// 		}

// 		rawDepositMsg, err := types.ModuleCdc.Marshal(&depositSignMsg)
// 		s.Require().NoError(err)
// 		signature, err := priv.Sign(rawDepositMsg)
// 		s.Require().NoError(err)

// 		msg := types.NewMsgCreatePool(
// 			channelA.PortID,
// 			channelA.ChannelID,
// 			chainAAddress,
// 			chainBAddress,
// 			signature,
// 			types.PoolAsset{
// 				Side:    types.PoolAssetSide_SOURCE,
// 				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
// 				Weight:  50,
// 				Decimal: 6,
// 			},
// 			types.PoolAsset{
// 				Side:    types.PoolAssetSide_TARGET,
// 				Balance: &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
// 				Weight:  50,
// 				Decimal: 6,
// 			},
// 			300,
// 		)
// 		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
// 		s.AssertValidTxResponse(resp)
// 		s.Require().NoError(err)

// 		// wait block when packet relay.
// 		test.WaitForBlocks(ctx, 10, chainA, chainB)
// 		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)

// 		// check pool info in chainA and chainB
// 		poolId := types.GetPoolId(msg.GetLiquidityDenoms())
// 		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
// 		s.Require().NoError(err)
// 		poolAInfo := poolARes.InterchainLiquidityPool

// 		// check pool info sync status.

// 		s.Require().EqualValues(msg.SourceChannel, poolAInfo.CounterPartyChannel)
// 		s.Require().EqualValues(msg.SourcePort, poolAInfo.CounterPartyPort)
// 		//s.Require().EqualValues(msg.Tokens[0].Amount, poolAInfo.Supply.Amount)

// 		poolBRes, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
// 		s.Require().NoError(err)
// 		poolBInfo := poolBRes.InterchainLiquidityPool
// 		s.Require().EqualValues(msg.SourceChannel, poolBInfo.CounterPartyChannel)
// 		s.Require().EqualValues(msg.SourcePort, poolBInfo.CounterPartyPort)
// 		//s.Require().EqualValues(msg.Tokens[1].Amount, poolBInfo.Supply.Amount)

// 		fmt.Println(poolAInfo)
// 		logger.CleanLog("Create Pool: PoolA", poolAInfo)
// 		fmt.Println("===================")
// 		logger.CleanLog("Create Pool: PoolB", poolBInfo)

// 		// compare pool info sync status
// 		s.Require().EqualValues(poolAInfo.Supply, poolBInfo.Supply)
// 		s.Require().EqualValues(poolAInfo.Assets[0].Balance.Amount, poolBInfo.Assets[0].Balance.Amount)
// 		s.Require().EqualValues(poolAInfo.Assets[1].Balance.Amount, poolBInfo.Assets[1].Balance.Amount)
// 	})

// 	t.Run("send deposit message", func(t *testing.T) {

// 		// check the balance of the chainA account.
// 		beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
// 		s.Require().NoError(err)
// 		s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

// 		// prepare deposit message.
// 		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
// 		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)}

// 		msg := types.NewMsgSingleAssetDeposit(
// 			poolId,
// 			chainBAddress,
// 			&depositCoin,
// 		)
// 		resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)
// 		s.AssertValidTxResponse(resp)
// 		s.Require().NoError(err)

// 		balanceRes, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
// 		s.Require().NoError(err)
// 		expectedBalance := balanceRes.Balance.Add(depositCoin)
// 		s.Require().Equal(expectedBalance.Denom, beforeDeposit.Balance.Denom)
// 		s.Require().Equal(expectedBalance.Amount, beforeDeposit.Balance.Amount)

// 		// check packet relayed or not.
// 		test.WaitForBlocks(ctx, 15, chainA, chainB)
// 		s.AssertPacketRelayed(ctx, chainB, channelB.PortID, channelB.ChannelID, 2)

// 		poolResInChainA, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
// 		s.Require().NoError(err)
// 		poolInChainA := poolResInChainA.InterchainLiquidityPool
// 		s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainA.Status)

// 		poolResInChainB, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
// 		s.Require().NoError(err)
// 		poolInChainB := poolResInChainB.InterchainLiquidityPool
// 		s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainB.Status)

// 		logger.CleanLog("Send Deposit(After):PoolA", poolInChainA)
// 		logger.CleanLog("Send Deposit(After):PoolB", poolInChainB)

// 	})

// 	depositX := initialX * 30 / 100
// 	//depositY := initialY * 30 / 100

// 	t.Run("send deposit messages", func(t *testing.T) {

// 		sender := chainAAddress
// 		denom := chainADenom
// 		chain := chainA
// 		//signer := chainAWallet
// 		depositAmount := sdk.NewInt(int64(depositX))

// 		testCases := []struct {
// 			name     string
// 			malleate func()
// 			expPass  bool
// 		}{

// 			{
// 				"single deposit Asset A(30%)",
// 				func() {
// 				},
// 				true,
// 			},
// 			// {
// 			// 	"single deposit Asset B(30%)",
// 			// 	func() {
// 			// 		sender = chainBAddress
// 			// 		denom = chainBDenom
// 			// 		chain = chainB
// 			// 		signer = chainBWallet
// 			// 		depositAmount = sdk.NewInt(int64(depositY))
// 			// 	},
// 			// 	true,
// 			// },
// 		}

// 		for _, tc := range testCases {
// 			logger.CleanLog(tc.name, "start")
// 			// check the balance of the chainA account.
// 			beforeDeposit, err := s.QueryBalance(ctx, chain, sender, denom)
// 			s.Require().NoError(err)
// 			s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

// 			// prepare deposit message.
// 			poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
// 			depositCoin := sdk.Coin{Denom: denom, Amount: depositAmount}

// 			poolRes, err := s.QueryInterchainswapPool(ctx, chain, poolId)
// 			s.Require().NoError(err)

// 			amm := types.NewInterchainMarketMaker(&poolRes.InterchainLiquidityPool)

// 			lp, err := amm.DepositSingleAsset(depositCoin)
// 			logger.CleanLog("lp price:", poolRes.InterchainLiquidityPool.PoolPrice)
// 			logger.CleanLog("deposit error", err)
// 			logger.CleanLog("lp amount", lp)
// 		}

// 	})

// }
