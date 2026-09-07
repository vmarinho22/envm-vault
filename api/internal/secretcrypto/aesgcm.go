package secretcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"github.com/google/uuid"
)

const (
	KeySize   = 32
	NonceSize = 12
)

var (
	ErrInvalidKey      = errors.New("encryption key must be 32 bytes")
	ErrInvalidBlob     = errors.New("invalid ciphertext blob")
	ErrWrongKeyVersion = errors.New("unknown encryption key version")
	ErrDecryptFailed   = errors.New("decrypt failed")
)

// Crypter encrypts and decrypts secret values at rest.
type Crypter interface {
	Encrypt(plaintext, aad []byte) (blob []byte, keyVersion int, err error)
	Decrypt(blob, aad []byte, keyVersion int) ([]byte, error)
}

type AESGCM struct {
	gcm        cipher.AEAD
	keyVersion int
}

// NewAESGCM builds a Crypter from a 32-byte key and version tag.
func NewAESGCM(key []byte, keyVersion int) (*AESGCM, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKey
	}
	if keyVersion < 1 {
		return nil, fmt.Errorf("key version must be >= 1")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("gcm: %w", err)
	}
	return &AESGCM{gcm: gcm, keyVersion: keyVersion}, nil
}

func (c *AESGCM) Encrypt(plaintext, aad []byte) ([]byte, int, error) {
	nonce := make([]byte, NonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, 0, fmt.Errorf("nonce: %w", err)
	}
	ciphertext := c.gcm.Seal(nil, nonce, plaintext, aad)
	blob := make([]byte, 0, NonceSize+len(ciphertext))
	blob = append(blob, nonce...)
	blob = append(blob, ciphertext...)
	return blob, c.keyVersion, nil
}

func (c *AESGCM) Decrypt(blob, aad []byte, keyVersion int) ([]byte, error) {
	if keyVersion != c.keyVersion {
		return nil, ErrWrongKeyVersion
	}
	if len(blob) < NonceSize+c.gcm.Overhead() {
		return nil, ErrInvalidBlob
	}
	nonce := blob[:NonceSize]
	ciphertext := blob[NonceSize:]
	plaintext, err := c.gcm.Open(nil, nonce, ciphertext, aad)
	if err != nil {
		return nil, ErrDecryptFailed
	}
	return plaintext, nil
}

// AAD builds authenticated associated data for a secret value row.
func AAD(secretID, environmentID uuid.UUID) []byte {
	out := make([]byte, 32)
	copy(out[:16], secretID[:])
	copy(out[16:], environmentID[:])
	return out
}
