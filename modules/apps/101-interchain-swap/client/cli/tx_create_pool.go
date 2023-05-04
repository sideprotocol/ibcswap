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

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create_pool [source-port] [source-channel] [sender] [weight] [tokens] [decimals]",
		Short: "Broadcast message CreatePool",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens, err := GetTokens(args[4])
			if err != nil {
				return err
			}

			decimals, err := parseDecimals(args[5])
			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(
				args[0], // argSourcePort
				args[1], // argSourceChannel
				args[2], // argSender
				args[3], // argWeight
				tokens,
				decimals,
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
	cmd.Flags().String("packet-timeout-height", "", "Packet timeout height")
	cmd.Flags().Uint("packet-timeout-timestamp", 0, "Packet timeout timestamp (in nanoseconds)")

	return cmd
}

func parseDecimals(decimalsStr string) ([]uint32, error) {
	decimalList := strings.Split(decimalsStr, ",")
	decimals := make([]uint32, 0, len(decimalList))

	for _, decimalStr := range decimalList {
		decimal, err := strconv.Atoi(decimalStr)
		if err != nil {
			return nil, fmt.Errorf("invalid decimal %s", decimalStr)
		}
		decimals = append(decimals, uint32(decimal))
	}

	if len(decimals) != 2 {
		return nil, fmt.Errorf("invalid decimals length %v", decimals)
	}

	return decimals, nil
}
