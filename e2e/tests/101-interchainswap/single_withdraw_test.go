package e2e

import (
	"context"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"
	types "github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/strangelove-ventures/ibctest/v6/chain/cosmos"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"
)

func (s *InterchainswapTestSuite) TestSingleWithdrawStatus() {

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

	const initialX = 1000_000 // USDT
	const initialY = 1000_000 // ETH

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

	t.Run("send withdraw message", func(t *testing.T) {

		poolId := types.GetPoolId([]string{chainADenom, chainBDenom})

		poolARes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
		s.Require().NoError(err)
		pool := poolARes.InterchainLiquidityPool

		depositToken := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)}
		amm := types.NewInterchainMarketMaker(&pool)
		out, err := amm.DepositSingleAsset(depositToken)
		s.Require().NoError(err)
		logger.CleanLog("===single-asset deposit===", out)

		msg := types.NewMsgSingleAssetDeposit(
			poolId,
			chainAAddress,
			&sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(initialX)},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		test.WaitForBlocks(ctx, 10, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

		pool = syncPool(s, ctx, chainA, poolId)
		logger.CleanLog("Pool Status after deposit", pool)
		amm = types.NewInterchainMarketMaker(&pool)

		outA, err := amm.SingleWithdraw(*out, chainADenom)
		s.Require().NoError(err)
		logger.CleanLog("====outA====", outA)
		s.Require().Equal(checkSlippage(outA.Amount, depositToken.Amount, 1), true)
	})
}

func syncPool(s *InterchainswapTestSuite, ctx context.Context, chain *cosmos.CosmosChain, poolId string) types.InterchainLiquidityPool {
	poolRes, err := s.QueryInterchainswapPool(ctx, chain, poolId)
	s.Require().NoError(err)
	pool := poolRes.InterchainLiquidityPool
	return pool
}

func checkSlippage(expect, real sdk.Int, slippage int) bool {
	if slippage > 100 || slippage < 0 {
		return false
	}
	differ := sdk.NewInt(0)
	if expect.GTE(real) {
		differ = expect.Sub(real)
	} else {
		differ = real.Sub(expect)
	}
	tolerance := sdk.NewInt(int64(slippage)).Mul(expect)
	return differ.LTE(tolerance)
}
