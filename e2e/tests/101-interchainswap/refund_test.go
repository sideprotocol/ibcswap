package e2e

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	test "github.com/strangelove-ventures/ibctest/v6/testutil"

	// "github.com/cosmos/ibc-go/e2e/testsuite"
	"github.com/cosmos/ibc-go/e2e/testvalues"

	//atomicswaptypes "github.com/ibcswap/ibcswap/v6/modules/apps/100-atomic-swap/types"
	types "github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	// "github.com/strangelove-ventures/ibctest/v6/ibc"
	// "github.com/stretchr/testify/suite"
)

func (s *InterchainswapTestSuite) TestRefundMsgOnTimeoutPacket() {

	t := s.T()
	ctx := context.TODO()

	// setup relayers and connection-0 between two chains.
	_, channelA, _ := s.SetupChainsRelayerAndChannel(
		ctx,
		interchainswapChannelOptions(),
	)

	chainA, chainB := s.GetChains()

	chainADenom := chainA.Config().Denom
	chainBDenom := chainB.Config().Denom

	// // create wallets for testing
	chainAWallet := s.CreateUserOnChainA(ctx, testvalues.StartingTokenAmount)
	chainAAddress := chainAWallet.Bech32Address("cosmos")

	// allocate tokens to the new account
	initialBalances := sdk.NewCoins(sdk.NewCoin(chainADenom, sdk.NewInt(10000000000)))
	err := s.SendCoinsFromModuleToAccount(ctx, chainA, chainAWallet, initialBalances)

	s.Require().NoError(err)
	res, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	s.Require().NotEqual(res.Balance.Amount, sdk.NewInt(0))
	s.Require().NoError(err)

	//chainBAddress = s.CreateUserOnChainB(ctx, testvalues.StartingTokenAmount)

	// s.Require().NoError(test.WaitForBlocks(ctx, 1, chainA, chainB), "failed to wait for blocks")

	// t.Run("stop relayer", func(t *testing.T) {
	// 	s.StopRelayer(ctx, relayer)
	// })

	t.Run("send create pool message", func(t *testing.T) {

		res, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		s.Require().NoError(err)
		balanceBefore := res.Balance.Amount

		msg := types.NewMsgCreatePool(
			channelA.PortID,
			channelA.ChannelID,
			chainAAddress,
			"1:1",
			[]*sdk.Coin{
				{Denom: chainADenom, Amount: sdk.NewInt(100000)},
				{Denom: chainBDenom, Amount: sdk.NewInt(10000)},
			},
			[]uint32{6, 6},
		)

		resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
		s.AssertValidTxResponse(resp)
		s.Require().NoError(err)

		height, err := chainA.Height(ctx)
		s.Require().NoError(err)
		expireHeight := height + 100

		// 150 block later than current block
		// Step 4: Wait for the expected timeout duration to pas
		// wait block when packet relay. Force timeout to refund
		test.WaitForBlocks(ctx, int(expireHeight), chainA, chainB)
		// check packet relay status.
		//s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 1)
		//s.Require().NoError(err)

		res, err = s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
		balanceAfter := res.Balance.Amount
		s.Require().NoError(err)
		s.Require().Equal(balanceBefore, balanceAfter)

	})

	// t.Run("send deposit message, with refund", func(t *testing.T) {

	// 	beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)

	// 	poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
	// 	depositCoin := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
	// 	msg := types.NewMsgDeposit(
	// 		poolId,
	// 		chainAAddress,
	// 		[]*sdk.Coin{&depositCoin},
	// 	)

	// 	resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)

	// 	s.StopRelayer(ctx, relayer)

	// 	// check packet relayed or not. - Force to timeout
	// 	test.WaitForBlocks(ctx, 250, chainA, chainB)
	// 	// s.AssertPacketRelayed(ctx, chainA, channelA.PortID, channelA.ChannelID, 2)

	// 	refundedDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)
	// 	s.Require().Equal(beforeDeposit, refundedDeposit)

	// })

	// // withdraw
	// t.Run("send withdraw message", func(t *testing.T) {

	// 	beforeDeposit, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)
	// 	poolId := types.GetPoolId([]string{chainADenom, chainBDenom})
	// 	poolRes, err := s.QueryInterchainswapPool(ctx, chainA, poolId)
	// 	s.Require().NoError(err)
	// 	poolCoin := poolRes.InterchainLiquidityPool.Supply
	// 	s.Require().NotEqual(poolCoin.Amount, sdk.NewInt(0))

	// 	denomOut := chainADenom
	// 	sender := chainAAddress
	// 	msg := types.NewMsgWithdraw(
	// 		sender,
	// 		poolCoin,
	// 		denomOut,
	// 	)
	// 	resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)

	// 	balanceRes, err := s.QueryBalance(ctx, chainA, chainAAddress, chainADenom)
	// 	s.Require().NoError(err)
	// 	s.Require().Equal(balanceRes.Balance.Denom, beforeDeposit.Balance.Denom)
	// 	s.Require().Equal(balanceRes.Balance.Amount, beforeDeposit.Balance.Amount)
	// })

	// // send swap message
	// t.Run("send swap message (don't check ack)", func(t *testing.T) {
	// 	sender := chainAAddress
	// 	tokenIn := sdk.Coin{Denom: chainADenom, Amount: sdk.NewInt(1000)}
	// 	tokenOut := sdk.Coin{Denom: chainBDenom, Amount: sdk.NewInt(1000)}
	// 	msg := types.NewMsgSwap(
	// 		sender,
	// 		10,
	// 		sender,
	// 		&tokenIn,
	// 		&tokenOut,
	// 	)
	// 	resp, err := s.BroadcastMessages(ctx, chainA, chainAWallet, msg)
	// 	s.AssertValidTxResponse(resp)
	// 	s.Require().NoError(err)

	// })

}
