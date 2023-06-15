package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMultiAssetWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi_asset_withdraw [poolId] [receiver] [remote sender][pool coin]",
		Short: "Broadcast message Withdraw",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			poolId := args[0]
			argSender := args[1]
			argRemoteSender := args[2]
			argCoin := args[3]

			sdk.MustAccAddressFromBech32(argSender)
			sdk.MustAccAddressFromBech32(argRemoteSender)

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := GetTokens(argCoin)
			if err != nil {
				return nil
			}
			if len(coins) != 0 {
				return fmt.Errorf("invalid token length! : %d", len(coins))
			}

			msg := types.NewMsgMultiAssetWithdraw(
				poolId,
				argSender,
				argRemoteSender,
				coins[0],
			)
			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")

			pool, err := QueryPool(clientCtx, coins[0].Denom)
			if err != nil {
				return err
			}

			if err1 == nil && err2 == nil {
				timeoutHeight, timeoutTimestamp, err := GetTimeOuts(clientCtx, pool.CounterPartyPort, pool.CounterPartyChannel, packetTimeoutHeight, uint64(packetTimeoutTimestamp), false)

				if err == nil {
					msg.TimeoutHeight = timeoutHeight
					msg.TimeoutTimeStamp = *timeoutTimestamp
				}
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	cmd.Flags().String("packet-timeout-height", "", "Packet timeout height")
	cmd.Flags().Uint("packet-timeout-timestamp", 0, "Packet timeout timestamp (in nanoseconds)")

	return cmd
}
