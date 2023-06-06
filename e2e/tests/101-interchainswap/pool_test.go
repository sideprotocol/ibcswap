package e2e

import (
	"context"
	"fmt"
	"math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	//"github.com/cosmos/cosmos-sdk/x/auth/tx"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v6/ibc"

	//"github.com/strangelove-ventures/ibctest/v6/ibc"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
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

	chainBUserMnemonic, err := createNewMnemonic()
	s.Require().NoError(err)
	chainBWallet := s.CreateUserOnChainBWithMnemonic(ctx, chainBUserMnemonic, testvalues.StartingTokenAmount)
	chainBAddress := chainBWallet.Bech32Address(chainB.Config().Bech32Prefix)
	priv, _ := getPrivFromNewMnemonic(chainBUserMnemonic)

	chainAWalletForB := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddressForB := chainAWalletForB.Bech32Address("cosmos")
	fmt.Println(chainAAddressForB)

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

	t.Run("send create pool message", func(t *testing.T) {

		depositSignMsg := types.DepositSignature{
			Sender:   chainBAddress,
			Balance:  &sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)},
			Sequence: 0,
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
				Balance: &sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(20000)},
				Weight:  50,
				Decimal: 6,
			},
			types.PoolAsset{
				Side:    types.PoolAssetSide_TARGET,
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
		poolId := types.GetPoolId(msg.GetLiquidityDenoms())
		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		poolAInfo := poolARes.InterchainLiquidityPool

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

	t.Run("send deposit message (enable pool) with invalid amount", func(t *testing.T) {

		// check the balance of the chainA account.
		beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

		// prepare deposit message.
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY * 0.1)}

		msg := types.NewMsgSingleAssetDeposit(
			poolId,
			chainBAddress,
			&depositCoin,
		)
		resp, err := s.BroadcastMessages(ctx, chainB, chainBWallet, msg)
		s.Require().NoError(err)
		s.Require().Equal(resp.Code, uint32(1538))
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
		s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainA.Status)

		poolResInChainB, err := s.QueryInterchainswapPool(ctx, chainB, poolId)
		s.Require().NoError(err)
		poolInChainB := poolResInChainB.InterchainLiquidityPool
		s.Require().Equal(types.PoolStatus_ACTIVE, poolInChainB.Status)

		logger.CleanLog("Send Deposit(After):PoolA", poolInChainA)
		logger.CleanLog("Send Deposit(After):PoolB", poolInChainB)

	})

	// // send swap message
	t.Run("pool status", func(t *testing.T) {

		var depositCoin sdk.Coin
		var depositCoins []*sdk.Coin
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

		testCases := []struct {
			name     string
			malleate func()
			msgType  string
			expPass  bool
		}{
			{
				"deposit Asset A (initialX)",
				func() {

					depositCoin = sdk.NewCoin(
						chainADenom,
						sdk.NewInt(initialX*0.015),
					)
					wallet = *chainAWallet
					chain = chainA
					channel = channelB
					packetSequence = 3

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
					packetSequence = 3
				},
				"deposit",
				true,
			},

			{
				"swap Asset A(100_000) to Asset B",
				func() {
					tokenIn = sdk.NewCoin(chainADenom, sdk.NewInt(100_000))
					wallet = *chainAWallet
					chain = chainA
					denomOut = chainBDenom
					sender = wallet.Bech32Address("cosmos")
					recipient = chainAAddressForB
					channel = channelA
					packetSequence = 4
				},
				"swap",
				true,
			},
			{
				"swap Asset B(100) to Asset A",
				func() {
					tokenIn = sdk.NewCoin(chainBDenom, sdk.NewInt(100))
					wallet = *chainBWallet
					chain = chainB
					sender = wallet.Bech32Address("cosmos")
					recipient = chainAAddressForB
					denomOut = chainADenom
					channel = channelB
					packetSequence = 4
				},
				"swap",
				true,
			},
			{
				"double deposit Assets (initialX,initialY)",
				func() {
					depositCoins = []*sdk.Coin{
						&sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
						&sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
					}
					wallet = *chainAWallet
					chain = chainA
					channel = channelA
					packetSequence = 5
				},
				"double deposit",
				true,
			},
			{
				"withdraw Asset A(50%)",
				func() {
					// initial liquidity token remove:
					lpAmount, err := s.QueryBalance(ctx, chainA, chainAAddress, poolId)
					s.Require().NoError(err)
					halfAmount := lpAmount.Balance.Amount.Quo(sdk.NewInt(2))
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
					halfAmount := lpAmount.Balance.Amount.Quo(sdk.NewInt(5))
					poolCoin = sdk.NewCoin(poolId, halfAmount)
					pool, err := getPool(
						s, ctx, chain, poolId,
					)
					asset, err := pool.FindAssetByDenom(chainBDenom)
					s.Require().NoError(err)
					ratio := 1 - float64(poolCoin.Amount.Int64())/float64(pool.Supply.Amount.Int64())
					exponent := 1 / float64(asset.Weight)
					factor := (1 - math.Pow(ratio, exponent)) * 1e18
					amountOut := asset.Balance.Amount.Mul(sdk.NewInt(int64(factor))).Quo(sdk.NewInt(1e18))

					logger.CleanLog("Withdrow params:", asset, ratio)
					fmt.Println("Amount Out:", amountOut)

					s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

					denomOut = chainBDenom
					wallet = *chainBWallet
					chain = chainB
					channel = channelB
					sender = wallet.Bech32Address("cosmos")
					packetSequence = 5
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
				msg := types.NewMsgSingleAssetDeposit(
					poolId,
					wallet.Bech32Address("cosmos"),
					&depositCoin,
				)
				txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)

			case "double deposit":

				pool, err := getPool(s, ctx, chain, poolId)
				escrowAcc := types.GetEscrowAddress(pool.CounterPartyPort, pool.CounterPartyChannel)

				encodingConfig := chain.Config().EncodingConfig.TxConfig
				//signMode := tx.DefaultSignModes
				txBuilder := encodingConfig.NewTxBuilder()
				//txBuilder := tx.NewTxConfig(encodingConfig, signMode).NewTxBuilder()

				sendTxMsg := banktypes.NewMsgSend(
					sdk.MustAccAddressFromBech32(chainBAddress),
					escrowAcc,
					sdk.NewCoins(*depositCoins[1]),
				)

				txBuilder.SetMsgs(sendTxMsg)
				txBuilder.SetMemo("ITF signed transaction")
				txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewInt64Coin("stake", 10)))
				txBuilder.SetGasLimit(200000)

				txBuilder.SetSignatures()

				hdPath := hd.CreateHDPath(118, 0, 0).String()
				derivedPriv, err := hd.Secp256k1.Derive()(wallet.Mnemonic, "", hdPath)
				if err != nil {
					println("err", err)
					continue
				}
				priv := hd.Secp256k1.Generate()(derivedPriv)

				sigV2 := signing.SignatureV2{
					PubKey: priv.PubKey(),
					Data: &signing.SingleSignatureData{
						SignMode:  chain.Config().EncodingConfig.TxConfig.SignModeHandler().DefaultMode(),
						Signature: nil,
					},
					Sequence: 0,
				}

				txBuilder.SetSignatures(sigV2)

				//s.Require().Equal(ok, true)

				// msg := types.NewMsgDoubleDeposit(
				// 	poolId,
				// 	wallet.Bech32Address("cosmos"),
				// 	depositCoins,
				// 	&types.CPDepositTx{
				// 		Tx:        tx,
				// 		Signature: "",
				// 	},
				// )

				// txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				// s.Require().NoError(err)
				// s.AssertValidTxResponse(txRes)
			case "swap":
				// pool status log.
				poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
				s.Require().NoError(err)
				pool := poolRes.InterchainLiquidityPool

				am := types.NewInterchainMarketMaker(&pool)

				tokenOut, _ = am.LeftSwap(tokenIn, denomOut)
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

				amm := types.NewInterchainMarketMaker(pool)

				outToken, err := amm.SingleWithdraw(poolCoin, denomOut)
				fmt.Println("pool coin:", poolCoin)
				fmt.Println("OutToken:", outToken)
				fmt.Println("ERR", err)
				s.Require().NoError(err)
				s.Require().NotEqual(outToken.Amount, sdk.NewInt(0))
				s.Require().NotNil(outToken)

				msg := types.NewMsgMultiAssetWithdraw(
					poolId,
					chainAAddress,
					chainBAddress,
					&poolCoin,
					&poolCoin,
				)

				err = msg.ValidateBasic()
				s.Require().NoError(err)

				txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.AssertValidTxResponse(txRes)
				s.Require().NoError(err)
			}

			test.WaitForBlocks(ctx, 15, chainA, chainB)
			s.AssertPacketRelayed(ctx, chain, channel.PortID, channel.ChannelID, uint64(packetSequence))

			// pool status log.
			pool, err := getPool(s, ctx, chain, poolId)

			s.Require().NoError(err)
			amm := types.NewInterchainMarketMaker(pool)

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
