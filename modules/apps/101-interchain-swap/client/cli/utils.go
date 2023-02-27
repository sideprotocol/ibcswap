package cli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GetTokens(argTokens string) ([]*sdk.Coin, error) {
	tokens := strings.Split(argTokens, ",")
	var tokenReg = regexp.MustCompile("([0-9]+)([A-Z]+)")
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid token input %s. Please follow this style `1marscoin,2venuscoin`", argTokens)
	}
	coins := []*sdk.Coin{}
	for _, token := range tokens {
		matches := tokenReg.FindStringSubmatch(token)
		amount, err := strconv.Atoi(matches[0])
		if err != nil {
			return nil, err
		}
		denom := matches[1]
		coins = append(coins, &sdk.Coin{
			Amount: sdk.NewInt(int64(amount)),
			Denom:  denom,
		})
	}
	return coins, nil
}
