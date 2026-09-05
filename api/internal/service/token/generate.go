package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	tokenLiteralPrefix = "ev_"
	tokenRandomBytes   = 32
	tokenDisplayLen    = 11
)

func GeneratePlaintext() (string, error) {
	buf := make([]byte, tokenRandomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}
	return tokenLiteralPrefix + hex.EncodeToString(buf), nil
}

func HashToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

func DisplayPrefix(plaintext string) string {
	if len(plaintext) < tokenDisplayLen {
		return plaintext
	}
	return plaintext[:tokenDisplayLen]
}
