package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashPass(s string) string {

	data := []byte(s)
	hash := sha256.Sum256(data)
	hashString := hex.EncodeToString(hash[:])

	return hashString
}
