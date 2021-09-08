package utils

import "math/big"

var (
	ONE_ETHER_IN_WEI = big.NewInt(1000000000000000000)
)

func EtherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, ONE_ETHER_IN_WEI)
}

func WeiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, ONE_ETHER_IN_WEI)
}
