package crypto

import (
	"golang.org/x/crypto/ed25519"
)

// SignMessageWithPrivateKey takes a message and a privateKey and returns a signature as hex string
func SignMessageWithPrivateKey(message string, privateKey []byte) []byte {
	rawMessage := []byte(message)

	signedMessage := ed25519.Sign(privateKey, rawMessage)

	return signedMessage
}

// SignDataWithPrivateKey takes data and a privateKey and returns a signature
func SignDataWithPrivateKey(data []byte, privateKey []byte) []byte {
	signedMessage := ed25519.Sign(privateKey, data)

	return signedMessage
}

// VerifyMessageWithPublicKey takes a message, signature and publicKey and verifies it
func VerifyMessageWithPublicKey(message string, signature []byte, publicKey []byte) bool {
	isValid := ed25519.Verify(publicKey, []byte(message), signature)

	return isValid
}

// VerifyDataWithPublicKey takes data, a signature and a publicKey and verifies it
func VerifyDataWithPublicKey(data []byte, signature []byte, publicKey []byte) bool {
	isValid := ed25519.Verify(publicKey, data, signature)

	return isValid
}
