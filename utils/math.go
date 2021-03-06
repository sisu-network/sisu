package utils

import (
	"bytes"
	"encoding/binary"

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
		log.Error("Failed to write binary float, err =", err)
	}
	return buf.Bytes()
}

func MaxUint64(a, b uint64) uint64 {
	if a < b {
		return b
	}

	return a
}
