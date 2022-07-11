package cardano

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cardanogo "github.com/echovl/cardano-go"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/echovl/cardano-go"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// GetCardanoMultiAsset gets Cardano policy address and asset name from transfer out data.
// Cardano multi-asset = policyID + asset nam
func GetCardanoMultiAsset(chain string, token *types.Token, assetAmount uint64) (*cardanogo.MultiAsset, error) {
	// Find the address of the cardano token
	for i, addr := range token.Addresses {
		if libchain.IsCardanoChain(token.Chains[i]) && chain == token.Chains[i] {
			index := strings.Index(addr, ":")
			if index < 0 {
				return nil, fmt.Errorf("cannot find policy id and asset for token: %s", token.Id)
			}

			policyID, err := cardano.NewHash28(addr[:index])
			if err != nil {
				return nil, fmt.Errorf("Invalid policy id, s = %s", addr)
			}
			assetName := cardano.NewAssetName(addr[index+1:])

			fmt.Println("GetCardanoMultiAsset: assetAmount = ", assetAmount)

			asset := cardanogo.NewAssets().Set(assetName, cardano.BigNum(assetAmount))
			multiAsset := cardanogo.NewMultiAsset().Set(cardano.NewPolicyIDFromHash(policyID), asset)

			return multiAsset, nil
		}
	}

	return nil, fmt.Errorf("Cannot find cardano token for %s", token.Id)
}

func GetTokenFromCardanoAsset(ctx sdk.Context, k keeper.Keeper, assetFullName string, cardanoChain string) *types.Token {
	tokens := k.GetAllTokens(ctx)
	for _, token := range tokens {
		for i, chain := range token.Chains {
			if chain == cardanoChain {
				addr := token.Addresses[i]
				index := strings.Index(addr, ":")
				if index < 0 {
					continue
				}

				policyID, err := cardano.NewHash28(addr[:index])
				if err != nil {
					continue
				}

				assetName := cardano.NewAssetName(addr[index+1:])
				if assetFullName == policyID.String()+assetName.String() {
					return token
				}
			}
		}
	}

	return nil
}
