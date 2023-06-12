package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

func CmdTakePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "take-pool [pool-id] [creator]",
		Short: "Broadcast message TakePool",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			poolId := args[0]
			creator := args[1]

			if _, err = sdk.AccAddressFromBech32(creator); err != nil {
				return err
			}

			msg := types.NewMsgTakePool(poolId, creator)

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
	cmd.Flags().String("packet-timeout-height", "", "Packet timeout height")
	cmd.Flags().Uint("packet-timeout-timestamp", 0, "Packet timeout timestamp (in nanoseconds)")

	return cmd
}
