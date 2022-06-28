package utils

import (
	"fmt"
	"math/big"

	libchain "github.com/sisu-network/lib/chain"
)

var (
	ONE_ETHER_IN_WEI    = big.NewInt(1000000000000000000)
	ONE_ADA_IN_LOVELACE = big.NewInt(1_000_000)
)

func EtherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, ONE_ETHER_IN_WEI)
}

func WeiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, ONE_ETHER_IN_WEI)
}

// LovelaceToWei note: it's temporary conversion to avoid transferring too small token amount
// 1 ADA = 10^18 wei
func LovelaceToWei(lovelace *big.Int) *big.Int {
	return new(big.Int).Mul(lovelace, new(big.Int).Div(ONE_ETHER_IN_WEI, ONE_ADA_IN_LOVELACE))
}

// WeiToLovelace converts ETH wei amount to ADA lovelace amount. 10^18 wei = 10^6 lovelace
func WeiToLovelace(wei *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(wei, ONE_ADA_IN_LOVELACE), ONE_ETHER_IN_WEI)
}

// SourceAmountToLovelace converts an amount from source chain to corresponding amount in lovelace
// in Cardano. For example, 10^18 wei in ETH is equivalent to 10^6 lovelace in Cardano.
func SourceAmountToLovelace(source string, amount *big.Int) (*big.Int, error) {
	if libchain.IsCardanoChain(source) {
		return new(big.Int).SetBytes(amount.Bytes()), nil
	}

	if libchain.IsETHBasedChain(source) {
		return new(big.Int).Div(new(big.Int).Mul(amount, ONE_ADA_IN_LOVELACE), ONE_ETHER_IN_WEI), nil
	}

	return nil, fmt.Errorf("Unknown source %s", source)
}
