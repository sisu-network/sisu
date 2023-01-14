package crypto

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
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

// GetAddressFromPublicKey takes a Lisk public key and returns the associated address
func GetAddressFromPublicKey(publicKey []byte) string {
	publicKeyHash := sha256.Sum256(publicKey)
	return hex.EncodeToString(publicKeyHash[:20])
}

func AddressToLisk32(address []byte) string {
	var byteSequence []byte
	for _, b := range address {
		byteSequence = append(byteSequence, b)
	}
	uint5Address := ConvertUIntArray(byteSequence, 8, 5)
	uint5Checksum := CreateChecksum(uint5Address)

	return "lsk"+ConvertUInt5ToBase32(append(uint5Address, uint5Checksum...))
}