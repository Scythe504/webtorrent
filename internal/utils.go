package internal

import (
	"crypto/rand"
	"math/big"
	"encoding/base64"
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

func GenerateSecureToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}
