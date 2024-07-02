package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func EncodeWithSHA256(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return strings.ToUpper(hex.EncodeToString(hasher.Sum(nil)))
}
