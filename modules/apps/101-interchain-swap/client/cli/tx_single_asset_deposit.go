package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSingleAssetDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "single_asset_deposit [pool-id] [sender] [pool-token] [port] [channel]",
		Short: "Broadcast message Deposit",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argPoolId := args[0]
			argSender := args[1]
			argTokens := args[2]
			argPort := args[3]
			argChannel := args[4]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens, err := GetTokens(argTokens)
			if err != nil {
				return err
			}

			msg := types.NewMsgSingleAssetDeposit(
				argPoolId,
				argSender,
				tokens[0],
				argPort,
				argChannel,
			)

			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")
			pool, err := QueryPool(clientCtx, argPoolId)
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

// QueryPool fetches the pool information from the chain using the gRPC client
