package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/ibcswap/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSingleAssetWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "single_asset_withdraw [demon] [pool coins]",
		Short: "Broadcast message Withdraw",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argOutDenoms := args[1]
			argCoin := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			coins, err := GetTokens(argCoin)
			if err != nil {
				return nil
			}
			if len(coins) != 2 {
				return fmt.Errorf("invalid token length! : %d", len(coins))
			}

			msg := types.NewMsgSingleAssetWithdraw(
				clientCtx.GetFromAddress().String(),
				argOutDenoms,
				coins[0],
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
