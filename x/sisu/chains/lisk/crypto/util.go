package crypto

import (
	"crypto/sha256"
	"math/big"
	"strings"
)

var (
	GENERATOR = []int{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}
	CHARSET   = "zxvcpmbn3465o978uyrtkqew2adsjhfg"
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

func ConvertUInt5ToBase32(uint5Array []byte) string {
	result := ""
	charsets := strings.Split(CHARSET, "")
	for _, value := range uint5Array {
		result += charsets[value]
	}
	return result
}

// GetBigNumberStringFromBytes returns the BigNumber representation of the bytes as string
func GetBigNumberStringFromBytes(data []byte) string {
	numericAddress := new(big.Int)
	numericAddress.SetBytes(data)

	return numericAddress.Text(10)
}

func ConvertUIntArray(uintArray []byte, fromBits int, toBits int) []byte {
	maxValue := (1 << toBits) - 1
	accumulator := 0
	bits := 0

	var result []byte
	for _, p := range uintArray {
		byteValue := p
		if byteValue < 0 || byteValue>>fromBits != 0 {
			return make([]byte, 0)
		}
		accumulator = (accumulator << fromBits) | int(byteValue)
		bits += fromBits
		for bits >= toBits {
			bits -= toBits
			result = append(result, byte((accumulator>>bits)&maxValue))
		}
	}
	return result
}

func CreateChecksum(uint5Array []byte) []byte {
	values := append(uint5Array, []byte{0, 0, 0, 0, 0, 0}...)
	mod := Polymod(values) ^ 1
	var result []byte
	for p := 0; p < 6; p += 1 {
		result = append(result, byte((mod>>(5*(5-p)))&31))
	}
	return result
}

func Polymod(uint5Array []byte) int {
	chk := 1
	for _, value := range uint5Array {

		top := chk >> 25
		chk = ((chk & 0x1ffffff) << 5) ^ int(value)
		for i := 0; i < 5; i += 1 {
			if ((top >> i) & 1) > 0 {
				chk ^= GENERATOR[i]
			}
		}
	}
	return chk
}
