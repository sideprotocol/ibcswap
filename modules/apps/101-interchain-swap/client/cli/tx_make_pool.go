package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/sideprotocol/ibcswap/v6/modules/apps/101-interchain-swap/types"
	"github.com/spf13/cobra"
)

func CmdMakePool() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make-pool [creator] [counterPartyCreator] [weight] [tokens] [decimals] [swap-fee] [channel]",
		Short: "Broadcast message MakePool",
		Args:  cobra.ExactArgs(7),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			weights, err := parseWeights(args[2])
			if err != nil {
				return err
			}

			tokens, err := GetTokens(args[3])
			if err != nil {
				return err
			}

			decimals, err := parseDecimals(args[4])
			if err != nil {
				return err
			}

			swapFee, err := strconv.Atoi(args[5])
			if err != nil {
				return err
			}

			if swapFee < 0 || swapFee > 10000 {
				return fmt.Errorf("invalid swap value. swapFee has to be in between 0 and 10000")
			}

			msg := types.NewMsgMakePool(
				args[6], // argSourcePort
				args[7], // argSourceChannel
				args[0], // argSender
				args[1], // argCounterPartySender
				types.PoolAsset{
					Side:    types.PoolAssetSide_SOURCE,
					Balance: tokens[0],
					Weight:  weights[0],
					Decimal: decimals[0],
				},
				types.PoolAsset{
					Side:    types.PoolAssetSide_DESTINATION,
					Balance: tokens[1],
					Weight:  weights[1],
					Decimal: decimals[1],
				},
				uint32(swapFee),
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

func parseWeights(weightsStr string) ([]uint32, error) {
	weights := strings.Split(weightsStr, ",")
	if len(weights) != 2 {
		return nil, fmt.Errorf("invalid weights length %v", weights)
	}

	totalWeight := 0
	weightsAsInt := []uint32{}
	for _, weight := range weights {
		weightAsInt, err := strconv.Atoi(weight)
		if err != nil || weightAsInt <= 0 {
			return nil, fmt.Errorf("can't parse weight value %v", err)
		}
		totalWeight += weightAsInt
		weightsAsInt = append(weightsAsInt, uint32(weightAsInt))
	}

	if totalWeight != 100 {
		return nil, fmt.Errorf("weight sum has to be equal to 100 %v", totalWeight)
	}
	return weightsAsInt, nil
}
