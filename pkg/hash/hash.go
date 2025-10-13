package hash

import (
	"crypto/sha256"
	"encoding/hex"
)

func Hash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func Verify(data, hashed string) bool {
	return Hash(data) == hashed
}
