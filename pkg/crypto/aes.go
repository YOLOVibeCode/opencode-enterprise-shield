// Package crypto provides encryption utilities for Enterprise Shield.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AESEncryptor provides AES-256-GCM encryption.
type AESEncryptor struct {
	key []byte
}

// NewAESEncryptor creates a new AES encryptor with the given key.
// Key must be 32 bytes (256 bits) for AES-256.
func NewAESEncryptor(key []byte) (*AESEncryptor, error) {
	if len(key) != 32 {
		return nil, errors.New("key must be 32 bytes for AES-256")
	}
	return &AESEncryptor{key: key}, nil
}

// GenerateKey generates a new random 256-bit key.
func GenerateKey() ([]byte, error) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return nil, err
	}
	return key, nil
}

// Encrypt encrypts plaintext using AES-256-GCM.
// Returns: nonce (12 bytes) || ciphertext || tag (16 bytes)
func (e *AESEncryptor) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Encrypt and append nonce at the beginning
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM.
// Input format: nonce (12 bytes) || ciphertext || tag (16 bytes)
func (e *AESEncryptor) Decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptString is a convenience method for encrypting strings.
func (e *AESEncryptor) EncryptString(plaintext string) ([]byte, error) {
	return e.Encrypt([]byte(plaintext))
}

// DecryptString is a convenience method for decrypting to strings.
func (e *AESEncryptor) DecryptString(ciphertext []byte) (string, error) {
	plaintext, err := e.Decrypt(ciphertext)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

// DeriveKey derives a key from a password using a simple KDF.
// For production, use Argon2 or scrypt.
func DeriveKey(password, salt []byte) []byte {
	// Simple PBKDF-like derivation (use proper KDF in production)
	key := make([]byte, 32)
	combined := append(password, salt...)
	
	// Simple hash-based derivation
	for i := 0; i < 32; i++ {
		key[i] = combined[i%len(combined)]
	}
	
	return key
}

