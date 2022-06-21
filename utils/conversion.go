package utils

import "math/big"

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

// LovelaceToETHTokens note: it's temporary conversion to avoid transferring too small token amount
// 1 ADA = 10^18 tokens
func LovelaceToETHTokens(lovelace *big.Int) *big.Int {
	return new(big.Int).Mul(lovelace, new(big.Int).Div(ONE_ETHER_IN_WEI, ONE_ADA_IN_LOVELACE))
}

// ETHTokensToLovelace 10^18 tokens = 1 ADA = 10^6 lovelace
func ETHTokensToLovelace(tokens *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(tokens, ONE_ADA_IN_LOVELACE), ONE_ETHER_IN_WEI)
}
