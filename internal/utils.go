package internal

import (
	"crypto/rand"
	"math/big"
)

func RandomId() string {
	charset := "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ_-"

	var shortCode = ""
	charsetLen := int64(len(charset))

	for range 8 {
		rd, _ := rand.Int(rand.Reader, big.NewInt(charsetLen))
		shortCode += string(charset[int(rd.Int64())])
	}

	return shortCode
}
