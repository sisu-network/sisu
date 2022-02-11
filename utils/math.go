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

func Float64ToByte(f float64) []byte {
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.BigEndian, f)
	if err != nil {
		log.Error("Failed to write binary float, err =", err)
	}
	return buf.Bytes()
}
