package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMultiAssetWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi_asset_withdraw [remote sender][demons(aside,bside)][pool coins]",
		Short: "Broadcast message Withdraw",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argRemoteSender := args[0]
			argOutDenoms := args[1]
			argCoin := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denoms := strings.Split(argOutDenoms, ",")
			if len(denoms) != 2 {
				return fmt.Errorf("invalid token length! : %d", len(denoms))
			}

			coins, err := GetTokens(argCoin)
			if err != nil {
				return nil
			}
			if len(coins) != 2 {
				return fmt.Errorf("invalid token length! : %d", len(coins))
			}

			msg := types.NewMsgMultiAssetWithdraw(
				clientCtx.GetFromAddress().String(),
				argRemoteSender,
				denoms[0],
				denoms[1],
				coins[0],
				coins[1],
			)
			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")

			if err1 == nil && err2 == nil {
				timeoutHeight, timeoutTimestamp, err := GetTimeOuts(clientCtx, args[0], args[1], packetTimeoutHeight, uint64(packetTimeoutTimestamp), false)
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

	return cmd
}
