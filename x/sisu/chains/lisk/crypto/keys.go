package crypto

import (
	"bytes"
	"crypto/sha256"

	"golang.org/x/crypto/ed25519"
)

// GetPrivateKeyFromSecret takes a Lisk secret and returns the associated private key
func GetPrivateKeyFromSecret(secret string) []byte {
	secretHash := GetSHA256Hash(secret)
	_, prKey, _ := ed25519.GenerateKey(bytes.NewReader(secretHash[:sha256.Size]))

	return prKey
}

// GetPublicKeyFromSecret takes a Lisk secret and returns the associated public key
func GetPublicKeyFromSecret(secret string) []byte {
	secretHash := GetSHA256Hash(secret)
	pKey, _, _ := ed25519.GenerateKey(bytes.NewReader(secretHash[:sha256.Size]))

	return pKey
}
