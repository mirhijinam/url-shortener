package random

import (
	"crypto/rand"
	"math/big"
)

var alphabet []rune = []rune(
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789",
)

func NewRandomAlias(length int) (string, error) {
	ret := make([]rune, length)

	for i := range ret {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {
			return "", err
		}
		ret[i] = alphabet[randInt.Int64()]
	}

	return string(ret), nil
}
