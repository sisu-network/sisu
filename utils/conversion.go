package utils

import (
	"encoding/binary"
	"math/big"
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

func Uint64ToBytes(num uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(num))

	return b
}

func BytesToUint64(bz []byte) uint64 {
	return binary.BigEndian.Uint64(bz)
}

func Uint32ToBytes(num uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, num)

	return b
}

func BytesToUint32(bz []byte) uint32 {
	return binary.BigEndian.Uint32(bz)
}
