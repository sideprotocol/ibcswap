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

func CmdCreatePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-pool [source-port] [source-channel] [sender] [weight] [denoms] [decimals] [initial liquidity amount]",
		Short: "Broadcast message CreatePool",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argSourcePort := args[0]
			argSourceChannel := args[1]
			argSender := args[2]
			argWeight := args[3]
			denomsStr := args[4]
			decimalsStr := args[5]
			initialLiquidity := args[6]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			denoms := strings.Split(denomsStr, ",")
			decimalList := strings.Split(decimalsStr, ",")
			decimals := []uint32{}
			for _, decimalStrEle := range decimalList {
				decimal, err := strconv.Atoi(decimalStrEle)
				if err != nil {
					return fmt.Errorf("invalid decimal %s", decimalStrEle)
				}
				decimals = append(decimals, uint32(decimal))
			}

			initialLiquidityAmount, err := strconv.ParseUint(initialLiquidity, 10, 64)

			if err != nil {
				return err
			}

			msg := types.NewMsgCreatePool(
				argSourcePort,
				argSourceChannel,
				argSender,
				argWeight,
				denoms,
				decimals,
				initialLiquidityAmount,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
