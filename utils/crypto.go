package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"io"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/ethereum/go-ethereum/accounts"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/go-bip39"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func AESDEncrypt(message []byte, key []byte) ([]byte, error) {
	// generate a new aes cipher using our 32 byte long key
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// here we encrypt our text using the Seal function
	// Seal encrypts and authenticates plaintext, authenticates the
	// additional data and appends the result to dst, returning the updated
	// slice. The nonce must be NonceSize() bytes long and unique for all
	// time, for a given key.
	return gcm.Seal(nonce, nonce, message, nil), nil
}

func GetCosmosPubKey(keyType string, bz []byte) (cryptotypes.PubKey, error) {
	switch keyType {
	case "ed25519":
		edKey := &ed25519.PubKey{}
		err := edKey.UnmarshalAmino(bz)
		return edKey, err

	case "secp256k1":
		secpKey := &secp256k1.PubKey{}
		err := secpKey.UnmarshalAmino(bz)
		return secpKey, err
	}

	return nil, fmt.Errorf("Cannot find public key")
}

func GetEcdsaPrivateKey(mnemonic string) (*ecdsa.PrivateKey, common.Address) {
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, "")
	if err != nil {
		panic(err)
	}

	dpath, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		panic(err)
	}

	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)

	key := masterKey
	for _, n := range dpath {
		key, err = key.Derive(n)
		if err != nil {
			panic(err)
		}
	}

	privateKey, err := key.ECPrivKey()
	if err != nil {
		panic(err)
	}

	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.PublicKey
	addr := crypto.PubkeyToAddress(publicKey)

	return privateKeyECDSA, addr
}
