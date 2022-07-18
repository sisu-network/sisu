package helper

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
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

func GetChainGasCostInToken(ctx sdk.Context, k keeper.Keeper, tokenId, chainId string) (*big.Int, error) {
	chain := k.GetChain(ctx, chainId)
	gasPrice := chain.GasPrice
	token := k.GetTokens(ctx, []string{tokenId})[tokenId]
	if token == nil {
		return nil, fmt.Errorf("GetChainGasCostInToken: cannot find token ", tokenId)
	}

	price, ok := new(big.Int).SetString(token.Price, 10)
}

func GetGasCostInToken(gas, gasPrice, tokenPrice, nativeTokenPrice *big.Int) (*big.Int, error) {
	// Get total gas cost
	gasCost := new(big.Int).Mul(gas, gasPrice)

	// amount := gasCost * nativeTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, nativeTokenPrice)
	gasInToken = new(big.Int).Div(gasInToken, tokenPrice)

	return gasInToken, nil
}
