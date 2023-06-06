package cli

import (
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdSingleAssetWithdraw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "single_asset_withdraw [demon] [pool coins]",
		Short: "Broadcast message Withdraw",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argOutDenoms := args[0]
			argCoin := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coinArgs := strings.Split(argCoin, "pool")
			amount, err := strconv.Atoi(coinArgs[0])
			if err != nil {
				return nil
			}
			poolCoin := sdk.Coin{
				Denom:  "pool" + coinArgs[1],
				Amount: sdk.NewInt(int64(amount)),
			}
			if err != nil {
				return nil
			}

			msg := types.NewMsgSingleAssetWithdraw(
				clientCtx.GetFromAddress().String(),
				argOutDenoms,
				&poolCoin,
			)

			packetTimeoutHeight, err1 := cmd.Flags().GetString("packet-timeout-height")
			packetTimeoutTimestamp, err2 := cmd.Flags().GetUint("packet-timeout-timestamp")

			pool, err := QueryPool(clientCtx, poolCoin.Denom)
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
