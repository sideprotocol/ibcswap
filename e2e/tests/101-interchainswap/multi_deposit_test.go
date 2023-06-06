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

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	crypto "github.com/cosmos/cosmos-sdk/crypto/types"
	bip39 "github.com/tyler-smith/go-bip39"
)

func (s *InterchainswapTestSuite) TestDoubleDepositStatus() {

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

	// // send swap message
	t.Run("pool status", func(t *testing.T) {

		poolId := types.GetPoolId([]string{
			chainADenom, chainBDenom,
		})

		amountOfChainBUserBeforeTx, err := s.GetBalance(ctx, chainB, chainBAddress, chainBDenom)
		s.Require().NoError(err)
		poolTokenAmountOfChainBUserBeforeTx, err := s.GetBalance(ctx, chainB, chainBAddress, poolId)
		s.Require().NoError(err)

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
					// 	{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
					// 	{Denom: chainBDenom, Amount: sdk.NewInt(1000)},
					// }
					// wallet = *chainAWallet
					// chain = chainA
					// channel = channelA
					// packetSequence = 2
				},
				"double deposit",
				true,
			},
		}

		for _, tc := range testCases {
			switch tc.msgType {
			case "double deposit":

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

				msg := types.NewMsgMultiAssetDeposit(
					poolId,
					[]string{
						chainAAddress,
						chainBAddress,
					},
					depositTokens,
					signedTx,
				)

				txRes, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
				s.Require().NoError(err)
				s.AssertValidTxResponse(txRes)
			}

			test.WaitForBlocks(ctx, 15, chainA, chainB)
			s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

			// pool status log.
			pool, err := getPool(s, ctx, chainA, poolId)
			s.Require().NoError(err)
			amm := types.NewInterchainMarketMaker(pool)

			priceA_B, _ := amm.MarketPrice(chainADenom, chainBDenom)
			priceB_A, _ := amm.MarketPrice(chainBDenom, chainADenom)

			logger.CleanLog("Price: A->B, B->A", *priceA_B, *priceB_A)
			logger.CleanLog("Pool Info", pool)

			amountOfChainBUserAfterTx, err := s.GetBalance(ctx, chainB, chainBAddress, chainBDenom)
			s.Require().NoError(err)
			depositedAmount := amountOfChainBUserBeforeTx.Sub(*amountOfChainBUserAfterTx)
			logger.CleanLog("balance(Before)", amountOfChainBUserBeforeTx)
			logger.CleanLog("balance(After)", amountOfChainBUserAfterTx)
			logger.CleanLog("depositedAmount", depositedAmount)
			s.Require().Equal(depositedAmount.Int64(), int64(initialY))

			poolTokenAmountOfChainBUserAfterTx, err := s.GetBalance(ctx, chainB, chainBAddress, poolId)
			s.Require().NoError(err)
			logger.CleanLog("pool Token(Before)", poolTokenAmountOfChainBUserBeforeTx)
			logger.CleanLog("pool Token(After)", poolTokenAmountOfChainBUserAfterTx)
			//s.Require().NotEqual(poolTokenAmountOfChainBUserBeforeTx.Int64(), poolTokenAmountOfChainBUserAfterTx.Int64())

			// check balance update status
			for _, asset := range pool.Assets {
				if asset.Balance.Denom == chainADenom {
					s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialX*2))
				}
				if asset.Balance.Denom == chainBDenom {
					s.Require().Equal(asset.Balance.Amount, sdk.NewInt(initialY*2))
				}
			}
		}
	})
}

func verifySignedMessage(rawTx []byte, signedMessage []byte, publicKey crypto.PubKey) bool {
	// Replace this with the actual rawTx bytes you want to verify
	return publicKey.VerifySignature(rawTx, signedMessage)
}

func createNewMnemonic() (string, error) {
	// Generate a new random mnemonic
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", fmt.Errorf("failed to generate new entropy: %w", err)
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate new mnemonic: %w", err)
	}
	return mnemonic, nil
}

func getPrivFromNewMnemonic(mnemonic string) (crypto.PrivKey, error) {
	hdPath := hd.CreateHDPath(118, 0, 0).String()
	derivedPriv, _ := hd.Secp256k1.Derive()(mnemonic, "", hdPath)
	priv := hd.Secp256k1.Generate()(derivedPriv)
	return priv, nil
}
