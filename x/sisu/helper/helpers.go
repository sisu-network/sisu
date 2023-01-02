package helper

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
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

func GetChainGasCostInToken(ctx sdk.Context, k keeper.Keeper, tokenId, chainId string,
	totalGasCost *big.Int) (*big.Int, error) {
	chain := k.GetChain(ctx, chainId)

	tokens := k.GetTokens(ctx, []string{tokenId, chain.NativeToken})
	token := tokens[tokenId]
	nativeToken := tokens[chain.NativeToken]
	if token == nil {
		return nil, fmt.Errorf("GetChainGasCostInToken: cannot find token %s", tokenId)
	}
	if nativeToken == nil {
		return nil, fmt.Errorf("GetChainGasCostInToken: cannot find token %s", chain.NativeToken)
	}

	tokenPrice, ok := new(big.Int).SetString(token.Price, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid token price: %s", token.Price)
	}

	if cmp := tokenPrice.Cmp(utils.ZeroBigInt); cmp <= 0 {
		return nil, fmt.Errorf("Token price must be positive: %s", tokenPrice)
	}

	nativeTokenPrice, ok := new(big.Int).SetString(nativeToken.Price, 10)
	if !ok {
		return nil, fmt.Errorf("Invalid native token price %s, token = %s", nativeToken.Price, nativeToken.Id)
	}

	gasCostInToken, err := GetGasCostInToken(totalGasCost, tokenPrice, nativeTokenPrice)
	log.Verbose("totalGas, tokenPrice, nativeTokenPrice, gasCost = ", totalGasCost, tokenPrice,
		nativeTokenPrice, gasCostInToken)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return gasCostInToken, nil
}

func GetGasCostInToken(gasCost, tokenPrice, nativeTokenPrice *big.Int) (*big.Int, error) {
	// amount := gasCost * nativeTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, nativeTokenPrice)
	gasInToken = new(big.Int).Div(gasInToken, tokenPrice)

	return gasInToken, nil
}
