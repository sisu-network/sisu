package types

import (
	"fmt"
	"math/big"

	"github.com/sisu-network/deyes/utils"
)

// GetAddressForChain returns the address for this token on a particular chain. Return empty string
// if not found.
func (t *Token) GetAddressForChain(c string) string {
	for i, chain := range t.Chains {
		if chain == c {
			return t.Addresses[i]
		}
	}

	return ""
}

// GetDeciamls returns the decimal of this token on a particular chain. Return 0 if not found.
func (t *Token) GetDeciamls(c string) byte {
	for i, chain := range t.Chains {
		if chain == c {
			return t.Decimals[i]
		}
	}

	return 0
}

// ConvertAmountToSisuAmount converts an amount on a chain to Sisu amount (18 decimals).
func (t *Token) ConvertAmountToSisuAmount(chain string, amount *big.Int) (*big.Int, error) {
	if amount == nil {
		return nil, fmt.Errorf("ConvertAmountToSisuAmount: Amount is nil")
	}

	var decimal byte
	for i, c := range t.Chains {
		if chain == c {
			decimal = t.Decimals[i]
			break
		}
	}

	if decimal == 0 {
		return nil, fmt.Errorf("Cannot find chain %s", chain)
	}

	pow := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	ret := new(big.Int).Mul(amount, utils.ONE_ETHER_IN_WEI)
	ret = ret.Quo(ret, pow)

	return ret, nil
}
