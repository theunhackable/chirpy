package auth

import (
	"crypto/rand"
	"encoding/hex"
	"log"
)

func MakeRefreshToken() string {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Printf("error generating refresh token: %v", err)
	}
	return hex.EncodeToString(key)
}
