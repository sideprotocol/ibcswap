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
		Use:   "create-pool [source-port] [source-channel] [sender] [weight] [tokens] [decimals]",
		Short: "Broadcast message CreatePool",
		Args:  cobra.ExactArgs(6),
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			argSourcePort := args[0]
			argSourceChannel := args[1]
			argSender := args[2]
			argWeight := args[3]
			denomsStr := args[4]
			decimalsStr := args[5]
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			tokens, err := GetTokens(denomsStr)
			if err != nil {
				return err
			}

			decimalList := strings.Split(decimalsStr, ",")
			decimals := []uint32{}
			for _, decimalStrEle := range decimalList {
				decimal, err := strconv.Atoi(decimalStrEle)
				if err != nil {
					return fmt.Errorf("invalid decimal %s", decimalStrEle)
				}
				decimals = append(decimals, uint32(decimal))
			}
			if len(decimals) != 2 {
				return fmt.Errorf("invalid decimals length %s", decimals)
			}

			msg := types.NewMsgCreatePool(
				argSourcePort,
				argSourceChannel,
				argSender,
				argWeight,
				tokens,
				decimals,
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
