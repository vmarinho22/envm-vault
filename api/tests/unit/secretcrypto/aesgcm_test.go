package secretcrypto_test

import (
	"bytes"
	"testing"

	"github.com/envm-vault/envm-vault/api/internal/secretcrypto"
	"github.com/google/uuid"
)

func testKey(t *testing.T) []byte {
	t.Helper()
	key := make([]byte, secretcrypto.KeySize)
	for i := range key {
		key[i] = byte(i + 1)
	}
	return key
}

func TestAESGCM_RoundTrip(t *testing.T) {
	c, err := secretcrypto.NewAESGCM(testKey(t), 1)
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	secretID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	envID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	aad := secretcrypto.AAD(secretID, envID)

	blob, ver, err := c.Encrypt([]byte("super-secret"), aad)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if ver != 1 {
		t.Fatalf("keyVersion = %d, want 1", ver)
	}

	got, err := c.Decrypt(blob, aad, ver)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if !bytes.Equal(got, []byte("super-secret")) {
		t.Fatalf("plaintext = %q", got)
	}
}

func TestAESGCM_TamperFails(t *testing.T) {
	c, err := secretcrypto.NewAESGCM(testKey(t), 1)
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	aad := secretcrypto.AAD(uuid.Nil, uuid.Nil)
	blob, ver, err := c.Encrypt([]byte("x"), aad)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	blob[len(blob)-1] ^= 0xff
	if _, err := c.Decrypt(blob, aad, ver); err == nil {
		t.Fatal("expected decrypt failure on tampered blob")
	}
}

func TestAESGCM_WrongAADFails(t *testing.T) {
	c, err := secretcrypto.NewAESGCM(testKey(t), 1)
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	secretID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	envID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	aad := secretcrypto.AAD(secretID, envID)
	blob, ver, err := c.Encrypt([]byte("x"), aad)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	wrong := secretcrypto.AAD(secretID, uuid.Nil)
	if _, err := c.Decrypt(blob, wrong, ver); err == nil {
		t.Fatal("expected decrypt failure with wrong AAD")
	}
}

func TestAESGCM_WrongKeyVersion(t *testing.T) {
	c, err := secretcrypto.NewAESGCM(testKey(t), 1)
	if err != nil {
		t.Fatalf("NewAESGCM: %v", err)
	}
	aad := secretcrypto.AAD(uuid.Nil, uuid.Nil)
	blob, _, err := c.Encrypt([]byte("x"), aad)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if _, err := c.Decrypt(blob, aad, 2); err == nil {
		t.Fatal("expected wrong key version error")
	}
}

func TestNewAESGCM_InvalidKey(t *testing.T) {
	if _, err := secretcrypto.NewAESGCM([]byte("short"), 1); err == nil {
		t.Fatal("expected invalid key error")
	}
}
