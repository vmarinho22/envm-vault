package service

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

const (
	otpLength  = 8
	otpExpiry  = 10 // minutes
	maxOTPCode = 100000000
)

type OTPService struct{}

func NewOTPService() *OTPService {
	return &OTPService{}
}

func (s *OTPService) Generate() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(maxOTPCode))
	if err != nil {
		return "", fmt.Errorf("generate otp: %w", err)
	}
	return fmt.Sprintf("%08d", n.Int64()), nil
}

func (s *OTPService) Hash(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

func (s *OTPService) Verify(code, hash string) bool {
	return s.Hash(code) == hash
}

func (s *OTPService) ExpiryMinutes() int {
	return otpExpiry
}
