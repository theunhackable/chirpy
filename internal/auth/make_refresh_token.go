package auth

import (
	"crypto/rand"
	"encoding/hex"
)

func MakeRefreshToken() string {
	key := make([]byte, 256)
	rand.Read(key)
	return hex.EncodeToString(key)
}
