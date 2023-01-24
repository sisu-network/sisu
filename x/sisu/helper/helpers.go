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
	log.Verbose("totalGas, tokenPrice, nativeTokenPrice, gasCostInToken = ", totalGasCost, tokenPrice,
		nativeTokenPrice, gasCostInToken)
	if err != nil {
		return nil, err
	}

	return gasCostInToken, nil
}

func GasCostInToken(gasCost, tokenPrice, nativeTokenPrice *big.Int) (*big.Int, error) {
	// amount := gasCost * nativeTokenPrice / tokenPrice
	gasInToken := new(big.Int).Mul(gasCost, nativeTokenPrice)
	gasInToken = new(big.Int).Div(gasInToken, tokenPrice)

	return gasInToken, nil
}
