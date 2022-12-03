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

// GetDecimalsForChain returns the decimal of this token on a particular chain. Return 0 if not found.
func (t *Token) GetDecimalsForChain(c string) uint32 {
	for i, chain := range t.Chains {
		if chain == c {
			return t.Decimals[i]
		}
	}

	return 0
}

// GetUnits returns an absolute value of a `value` unit. For example, 2 ETH on Ethereum with decimal
// 18 will return 2 * 10 ^ 18 units.
func (t *Token) GetUnits(chain string, value int) (*big.Int, error) {
	decimal := t.GetDecimalsForChain(chain)
	if decimal == 0 {
		return nil, fmt.Errorf("Cannot find decimal for chain %s", chain)
	}

	bigValue := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimal)), nil)
	bigValue = bigValue.Mul(bigValue, big.NewInt(int64(value)))

	return bigValue, nil
}

// ConvertAmountToSisuAmount converts an amount on a chain to Sisu amount (18 decimals).
func (t *Token) ConvertAmountToSisuAmount(chain string, amount *big.Int) (*big.Int, error) {
	if amount == nil {
		return nil, fmt.Errorf("ConvertAmountToSisuAmount: Amount is nil")
	}

	var decimal uint32
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
