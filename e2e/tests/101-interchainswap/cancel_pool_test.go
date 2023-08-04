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

func (s *InterchainswapTestSuite) TestCancelPoolMsgPacket() {

	t := s.T()
	ctx := context.TODO()
	logger := testsuite.NewLogger()
	// setup relayers and connection-0 between two chains.
	relayer, channelA, _ := s.SetupChainsRelayerAndChannel(ctx, interchainswapChannelOptions())

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

	addr := sdk.AccAddress(priv.PubKey().Address().Bytes())
	logger.CleanLog("address:", addr.String())
	s.Require().Equal(chainBAddress, addr.String())

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

	t.Run("send cancel pool message", func(t *testing.T) {

		pool := getFirstPool(s, ctx, chainA)
		msg := types.NewMsgCancelPool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			pool.Id,
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)

		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		// wait block when packet relay.
		test.WaitForBlocks(ctx, 5, chainA, chainB)
		s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 3)

		// check pool info in chainA and chainB
		// check pool info in chainA and chainB
		poolsResInChainA, err := s.QueryInterchainswapPools(ctx, chainA)
		s.Require().NoError(err)
		s.Require().Equal(len(poolsResInChainA.InterchainLiquidityPool), 0)
		poolsResInChainB, err := s.QueryInterchainswapPools(ctx, chainB)
		s.Require().NoError(err)
		s.Require().Equal(len(poolsResInChainB.InterchainLiquidityPool), 0)

	})

}
