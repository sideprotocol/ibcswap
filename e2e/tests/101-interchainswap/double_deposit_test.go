package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"

	test "github.com/strangelove-ventures/ibctest/v6/testutil"

	//"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	//"github.com/strangelove-ventures/ibctest/v6/ibc"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	bip39 "github.com/tyler-smith/go-bip39"
)

func (s *InterchainswapTestSuite) TestDoubleDepositStatus() {

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

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddress := chainAWallet.Bech32Address("cosmos")

	chainBWallet := s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)
	chainBAddress := chainBWallet.Bech32Address("cosmos")

	// // create wallets for testing
	chainBUserWallet := s.CreateUserOnChainA(ctx, 10)
	chainBUserAddress := chainBUserWallet.Bech32Address("cosmos")
	_ = chainBUserAddress

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

	// wallet create
	add, mnemonic, priv, err := createNewWallet()
	logger.CleanLog("wallet", add.String(), mnemonic, priv, err)
	funds := sdk.NewCoins(sdk.NewCoin(chainBDenom, sdk.NewInt(1000)))
	err = s.SendCoins(ctx, chainB, chainBWallet, add.String(), funds)
	s.Require().NoError(err)

	res, err := s.QueryBalance(ctx, chainB, add.String(), chainBDenom)
	s.Require().NoError(err)
	logger.CleanLog("Account status", res.Balance.Amount)

	//sig, err := priv.Sign(rawTx)

	//logger.CleanLog("signed", sig, err)
	//verifySignedMessage(sig,&add)

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

		//var depositCoins []*sdk.Coin
		// var chain *cosmos.CosmosChain
		// var wallet ibc.Wallet

		// poolId := types.GetPoolId([]string{
		// 	chainADenom, chainBDenom,
		// })

		// var channel ibc.ChannelOutput
		// var txRes sdk.TxResponse

		// var packetSequence int

		poolId := types.GetPoolId([]string{
			chainADenom, chainBDenom,
		})
		pool, err := getPool(s, ctx, chainA, poolId)
		s.Require().NoError(err)
		escrowAccount := types.GetEscrowAddress(pool.EncounterPartyPort, pool.EncounterPartyChannel)

		testCases := []struct {
			name     string
			malleate func()
			msgType  string
			expPass  bool
		}{

			{
				"double deposit Assets (initialX,initialY)",
				func() {
					// depositCoins = []*sdk.Coin{
					// 	&sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
					// 	&sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(initialY)},
					// }
					// wallet = *chainAWallet
					// chain = chainA
					// channel = channelA
					// packetSequence = 5
				},
				"double deposit",
				true,
			},
		}

		for _, tc := range testCases {

			switch tc.msgType {
			case "double deposit":

				logger.CleanLog("escrow Account:", escrowAccount.String())

				sendTxMsg := banktypes.NewMsgSend(
					sdk.MustAccAddressFromBech32(chainBAddress),
					escrowAccount,
					sdk.NewCoins(sdk.NewCoin(chainBDenom, sdk.NewInt(1000))),
				)

				// depositTx := types.EncounterPartyDepositTx{
				// 	AccountSequence: nonce,
				// 	Sender:          suite.chainB.SenderAccount.GetAddress().String(),
				// 	Tokens:          []*sdk.Coin{{Denom: denomPair[0], Amount: sdk.NewInt(1000)}, {Denom: denomPair[1], Amount: sdk.NewInt(1000)}},
				// }

				rawTx, err := types.ModuleCdc.Marshal(sendTxMsg)
				if err != nil {
					fmt.Println(err)
					return
				}

				signedTx, err := priv.Sign(rawTx)
				pubKey := priv.PubKey()
				s.Require().Equal(verifySignedMessage(rawTx, signedTx, pubKey), true)

				// msg := types.NewMsgDoubleDeposit(
				// 	poolId,
				// 	chainAWallet.Bech32Address("cosmos"),
				// 	chainBWallet.Bech32Address("cosmos"),
				// 	[]*sdk.Coin{
				// 		{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
				// 		{Denom: chainBDenom, Amount: sdk.NewInt(1000)},
				// 	},
				// 	rawTx,
				// )

				// // fmt.Println("Deposit Message", msg)

				// txRes, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
				// fmt.Println("Res:", txRes)
				// logger.CleanLog("DoubleDeposit:", txRes)
				// s.Require().NoError(err)
				// s.AssertValidTxResponse(txRes)

			}

			// test.WaitForBlocks(ctx, 15, chainA, chainB)
			// s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, uint64(2))

			// res, err := s.QueryBalance(ctx, chainB, escrowAccount.String(), chainBDenom)
			// s.Require().NoError(err)

			// logger.CleanLog("My Balance", res)
			// // pool status log.
			// pool, err := getPool(s, ctx, chain, poolId)

			// s.Require().NoError(err)
			// amm := types.NewInterchainMarketMaker(
			// 	pool,
			// 	fee,
			// )

			// priceA_B, _ := amm.MarketPrice(chainADenom, chainBDenom)
			// priceB_A, _ := amm.MarketPrice(chainADenom, chainBDenom)

			// logger.CleanLog("Price: A->B, B->A", *priceA_B, *priceB_A)
			// logger.CleanLog("Pool Info", pool)
		}

	})
}

func createNewWallet() (sdk.AccAddress, string, crypto.PrivKey, error) {
	// Generate a new random mnemonic
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to generate new entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, "", nil, fmt.Errorf("failed to generate new mnemonic: %w", err)
	}

	// Create a BIP39 seed from the mnemonic
	seed := bip39.NewSeed(mnemonic, "")

	// Create a BIP32 master key from the seed
	masterKey, _ := hd.ComputeMastersFromSeed(seed)

	// Create a secp256k1 private key from the master key

	privKey := secp256k1.GenPrivKeyFromSecret(masterKey[:])

	// Get the address associated with the private key
	addr := sdk.AccAddress(privKey.PubKey().Address().Bytes())

	return addr, mnemonic, privKey, nil
}

func verifySignedMessage(rawTx []byte, signedMessage []byte, publicKey crypto.PubKey) bool {
	// Replace this with the actual rawTx bytes you want to verify
	return publicKey.VerifySignature(rawTx, signedMessage)
}
