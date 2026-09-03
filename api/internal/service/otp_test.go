package service

import (
	"testing"
)

func TestOTPService_GenerateAndVerify(t *testing.T) {
	otp := NewOTPService()

	code, err := otp.Generate()
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(code) != 8 {
		t.Errorf("Generate() length = %d, want 8", len(code))
	}

	hash := otp.Hash(code)
	if hash == "" {
		t.Fatal("Hash() returned empty string")
	}

	if !otp.Verify(code, hash) {
		t.Error("Verify() should return true for valid code")
	}

	if otp.Verify("00000000", hash) {
		t.Error("Verify() should return false for invalid code")
	}
}

func TestOTPService_HashIsDeterministic(t *testing.T) {
	otp := NewOTPService()

	h1 := otp.Hash("12345678")
	h2 := otp.Hash("12345678")

	if h1 != h2 {
		t.Error("Hash() should be deterministic")
	}
}

func TestOTPService_ExpiryMinutes(t *testing.T) {
	otp := NewOTPService()
	if otp.ExpiryMinutes() != 10 {
		t.Errorf("ExpiryMinutes() = %d, want 10", otp.ExpiryMinutes())
	}
}
