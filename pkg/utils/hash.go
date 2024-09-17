package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sum256(str string) string {
	hash := sha256.Sum256([]byte(str))
	return hex.EncodeToString(hash[:])
}
