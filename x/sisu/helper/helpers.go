package helper

import (
	"crypto"
	"encoding/hex"
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/external"
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

func GetChainGasCostInToken(ctx sdk.Context, k keeper.Keeper, deyesClient external.DeyesClient,
	tokenId string, chainId string, totalGasCost *big.Int) (*big.Int, error) {
	// 1. Get native token price
	chain := k.GetChain(ctx, chainId)
	if chain == nil {
		return nil, fmt.Errorf("Invalid chain %s", chainId)
	}

	nativeTokenPrice, err := deyesClient.GetTokenPrice(chain.NativeToken)
	if err != nil {
		return nil, err
	}

	// 2. Get token price
	tokenPrice, err := deyesClient.GetTokenPrice(tokenId)
	if err != nil {
		return nil, err
	}

	// 3. Calculate how many token needed to use to cover the gas cost.
	gasCostInToken, err := GasCostInToken(totalGasCost, tokenPrice, nativeTokenPrice)
	if err != nil {
		return nil, err
	}

	return gasCostInToken, nil
}

func GasCostInToken(gasCost, tokenPrice, nativeTokenPrice *big.Int) (*big.Int, error) {
	// amount := gasCost * nativeTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, nativeTokenPrice)
	gasInToken = new(big.Int).Div(gasInToken, tokenPrice)

	log.Verbose("totalGas, tokenPrice, nativeTokenPrice, gasCostInToken = ", gasCost, tokenPrice,
		nativeTokenPrice, gasInToken)

	return gasInToken, nil
}

func Difference(a, b *big.Int) float64 {
	x := a
	if a.Cmp(b) == -1 {
		x = b
	}

	diff := new(big.Int).Abs(new(big.Int).Sub(a, b))
	diffPercentage := new(big.Float).Quo(new(big.Float).SetInt(diff), new(big.Float).SetInt(x))

	p, _ := diffPercentage.Float64()
	return p
}
