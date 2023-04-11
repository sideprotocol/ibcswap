package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	"github.com/strangelove-ventures/ibctest/v6/ibc"

	//"github.com/strangelove-ventures/ibctest/v6/ibc"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"

	tx "github.com/cosmos/cosmos-sdk/types/tx"
)

func (s *InterchainswapTestSuite) TestDoubleDepositStatus() {

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

	t.Run("send deposit message (enable pool) with invalid amount", func(t *testing.T) {

		// check the balance of the chainA account.
		beforeDeposit, err := s.QueryBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		s.Require().NotEqual(beforeDeposit.Balance.Amount, sdk.NewInt(0))

		// prepare deposit message.
		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
		depositCoin := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY * 0.1)}

		msg := types.NewMsgDeposit(
			poolId,
			chainBAddress,
			[]*sdk.Coin{&depositCoin},
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

	// // send swap message
	t.Run("pool status", func(t *testing.T) {

		var depositCoins []*sdk.Coin
		var chain *cosmos.CosmosChain
		var wallet ibc.Wallet

		poolId := types.GetPoolId([]string{
			chainADenom, chainBDenom,
		})

		var channel ibc.ChannelOutput
		var txRes sdk.TxResponse

		var packetSequence int

		testCases := []struct {
			name     string
			malleate func()
			msgType  string
			expPass  bool
		}{

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
		}

		for _, tc := range testCases {

			switch tc.msgType {
			case "double deposit":

				pool, err := getPool(s, ctx, chain, poolId)
				escrowAcc := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)

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
				txS := txBuilder.GetTx()

				rawBytes, err := encodingConfig.TxEncoder()(txS)
				cosmosTx := &tx.Tx{}
				err = cosmosTx.Unmarshal(rawBytes)

				msg := types.NewMsgDoubleDeposit(
					poolId,
					wallet.Bech32Address("cosmos"),
					depositCoins,
					&types.CPDepositTx{
						Tx:        cosmosTx,
						Signature: "",
					},
				)

				txRes, err = s.BroadcastMessages(ctx, chain, &wallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)

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
