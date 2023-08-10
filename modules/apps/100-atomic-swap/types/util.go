package types

import (
	"crypto/rand"
	"encoding/base64"
)

// generateRandomString generates a random string of the length n.
func GenerateRandomString(chainID string, n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return chainID + base64.URLEncoding.EncodeToString(b)
}

func GetEventValueWithSuffix(value, suffix string) string {
	return value + "_" + suffix
}
