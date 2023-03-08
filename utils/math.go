package utils

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/sisu-network/lib/log"
)

func Float32ToByte(f float32) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		log.Error("Failed to write binary float, err =", err)
	}
	return buf.Bytes()
}

func ToByte(i interface{}) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, i)
	if err != nil {
		log.Errorf("Failed to write binary float, err = %s", err)
	}
	return buf.Bytes()
}

func FromByteToInt(bz []byte) int {
	value := binary.BigEndian.Uint32(bz)
	return int(value)
}

func FromByteToInt64(bz []byte) int64 {
	value := binary.BigEndian.Uint64(bz)
	return int64(value)
}

func MaxUint64(a, b uint64) uint64 {
	if a < b {
		return b
	}

	return a
}

// SubtractCommissionRate returns an amount after substracting commission rate. 1 commission rate
// unit is 0.01%
func SubtractCommissionRate(amount *big.Int, rate int32) *big.Int {
	amount = new(big.Int).Mul(amount, big.NewInt(int64(10_000-rate)))
	amount = new(big.Int).Div(amount, big.NewInt(10_000))
	return amount
}
