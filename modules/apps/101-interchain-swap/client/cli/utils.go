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
	fmt.Println("Tokens", tokens)
	var tokenReg = regexp.MustCompile(`^(\d+)([a-zA-Z]+)$`)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("invalid token input %s. Please follow this style `1marscoin,2venuscoin`", argTokens)
	}
	coins := []*sdk.Coin{}
	for _, token := range tokens {
		matches := tokenReg.FindStringSubmatch(token)
		fmt.Println("token", matches)
		amount, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}
		denom := matches[2]
		coins = append(coins, &sdk.Coin{
			Amount: sdk.NewInt(int64(amount)),
			Denom:  denom,
		})
	}
	return coins, nil
}
