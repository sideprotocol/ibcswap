package types

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strconv"
	"strings"
)

func ParseWeight(text string, length int) ([]uint32, error) {
	weights := make([]uint32, length)
	raws := strings.Split(text, ":")
	if len(raws) != length {
		return weights, ErrInvalidPairLength
	}

	for i := 0; i < length; i++ {
		if v, err := strconv.ParseUint(raws[i], 10, 8); err != nil {
			return nil, ErrInvalidPairLength
		} else if v > 0 && v < 100 {
			weights[i] = uint32(v)
		} else {
			return nil, ErrInvalidWeightOfPool
		}
	}
	return weights, nil
}

func GeneratePoolId(denoms []string) string {
	sort.Sort(sort.StringSlice(denoms))
	var text string
	for i := 0; i < len(denoms); i++ {
		text += denoms[i]
	}
	hash := sha256.Sum256([]byte(text))
	return hex.EncodeToString(hash[:])
}
