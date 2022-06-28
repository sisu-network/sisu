package helper

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math/big"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

var BigZero = big.NewInt(0)

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

func GetCardanoTxFeeInToken(adaPrice, tokenPrice, adaForTxFee *big.Int) *big.Int {
	if tokenPrice.Cmp(BigZero) <= 0 {
		return BigZero
	}

	txFeeInAda := new(big.Int).Mul(adaPrice, adaForTxFee)
	return new(big.Int).Div(txFeeInAda, tokenPrice)
}
