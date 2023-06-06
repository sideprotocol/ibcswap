package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdMultiAssetDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multi_asset_deposit [pool-id] [local sender] [remote sender] [pool-tokens(1000aside,1000bside)] [remote sender signature]",
		Short: "Broadcast message Deposit",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId := args[0]
			argLocalSender := args[1]
			argRemoteSender := args[2]
			argTokens := args[3]
			argSignature := args[4]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens, err := GetTokens(argTokens)
			if err != nil {
				return err
			}

			msg := types.NewMsgMultiAssetDeposit(
				argPoolId,
				[]string{argLocalSender, argRemoteSender},
				tokens,
				[]byte(argSignature),
			)

			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")

			pool, err := QueryPool(clientCtx, argPoolId)
			if err != nil {
				return err
			}

			if err1 == nil && err2 == nil {
				timeoutHeight, timeoutTimestamp, err := GetTimeOuts(clientCtx, pool.CounterPartyPort, pool.CounterPartyChannel, packetTimeoutHeight, uint64(packetTimeoutTimestamp), false)
				fmt.Println("Timeout Height:", timeoutHeight)
				fmt.Println("Timeout Timestamp:", timeoutTimestamp)
				fmt.Println("Timeouts Err:", err)
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
