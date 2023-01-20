package utils

import (
	"encoding/binary"
	"math/big"
)

var (
	OneEtherInWei    = int64(1000000000000000000)
	OneAdaInLoveLace = int64(1_000_000)
)

func EtherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(OneEtherInWei))
}

func WeiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(OneEtherInWei))
}

// LovelaceToWei note: it's temporary conversion to avoid transferring too small token amount
// 1 ADA = 10^18 wei
func LovelaceToWei(lovelace *big.Int) *big.Int {
	return new(big.Int).Mul(lovelace, new(big.Int).Div(big.NewInt(OneEtherInWei),
		big.NewInt(OneAdaInLoveLace)))
}

// WeiToLovelace converts ETH wei amount to ADA lovelace amount. 10^18 wei = 10^6 lovelace
func WeiToLovelace(wei *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(wei, big.NewInt(OneAdaInLoveLace)),
		big.NewInt(OneEtherInWei))
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
