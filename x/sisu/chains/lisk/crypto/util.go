package crypto

import (
	"crypto/sha256"
	"math/big"
)

// GetSHA256Hash returns the SHA256 hash of a string as byte slice
func GetSHA256Hash(stringToHash string) [sha256.Size]byte {
	return sha256.Sum256([]byte(stringToHash))
}

// GetFirstEightBytesReversed returns the first 8 bytes of a byte slice in reversed order.
func GetFirstEightBytesReversed(bytes []byte) []byte {
	if len(bytes) < 8 {
		return nil
	}

	result := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		result[7-i] = bytes[i]
	}

	return result
}

// GetBigNumberStringFromBytes returns the BigNumber representation of the bytes as string
func GetBigNumberStringFromBytes(data []byte) string {
	numericAddress := new(big.Int)
	numericAddress.SetBytes(data)

	return numericAddress.Text(10)
}
