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

func FromByteToInt(i interface{}) int {
	var buf bytes.Buffer
	var result int
	err := binary.Read(&buf, binary.BigEndian, &result)
	if err != nil {
		log.Errorf("Failed to read binary float, err = %s", err)
	}

	return result
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
	amount = amount.Mul(amount, big.NewInt(int64(10000-rate)))
	amount = amount.Div(amount, big.NewInt(10_000))

	return amount
}
