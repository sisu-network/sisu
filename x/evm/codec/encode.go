package codec

import (
	"bytes"
	"encoding/binary"
	"io"
)

func EncodePrefixedLength(data []byte) ([]byte, error) {
	var buf = new(bytes.Buffer)

	size := len(data)
	err := EncodeUvarint(buf, uint64(size))
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Copy from go-amino@v0.15.1
// EncodeUvarint is used to encode golang's int, int32, int64 by default. unless specified differently by the
// `binary:"fixed32"`, `binary:"fixed64"`, or `binary:"zigzag32"` `binary:"zigzag64"` tags.
// It matches protobufs varint encoding.
func EncodeUvarint(w io.Writer, u uint64) (err error) {
	var buf [10]byte
	n := binary.PutUvarint(buf[:], u)
	_, err = w.Write(buf[0:n])
	return
}
