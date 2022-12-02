package helper

import (
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"os"

	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func GetDevPrivateKey() *ecdsa.PrivateKey {
	// This is the private key for account 0xbeF23B2AC7857748fEA1f499BE8227c5fD07E70c
	bz, err := hex.DecodeString("9f575b88940d452da46a6ceec06a108fcd5863885524aec7fb0bc4906eb63ab1")
	if err != nil {
		panic(err)
	}

	privateKey, err := ethcrypto.ToECDSA(bz)
	if err != nil {
		panic(err)
	}

	return privateKey
}

func GetTokens(file string) []*types.Token {
	tokens := []*types.Token{}

	dat, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &tokens); err != nil {
		panic(err)
	}

	return tokens
}
