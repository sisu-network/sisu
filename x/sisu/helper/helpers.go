package helper

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func GetKeygenId(keyType string, block int64, pubKeys []ctypes.PubKey) string {
	// Get hashes of all pubkeys
	digester := crypto.MD5.New()
	for _, pubKey := range pubKeys {
		fmt.Fprint(digester, pubKey.Bytes())
	}
	hash := hex.EncodeToString(digester.Sum(nil))

	return fmt.Sprintf("%s;%d;%s", keyType, block, hash)
}

func GetGasCostInToken(gas, gasPrice, tokenPrice, nativeTokenPrice *big.Int) (*big.Int, error) {
	// Get total gas cost
	gasCost := new(big.Int).Mul(gas, gasPrice)

	// amount := gasCost * nativeTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, nativeTokenPrice)
	gasInToken = new(big.Int).Div(gasInToken, tokenPrice)

	return gasInToken, nil
}

// BytesToValPubKey converts from byte array to validator public key
// Cosmos sdk use ed25519 format to present validator public key
func BytesToValPubKey(bz []byte) ctypes.PubKey {
	return &ed25519.PubKey{Key: bz}
}
