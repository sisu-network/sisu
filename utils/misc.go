package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"sort"
	"strconv"
	"sync"

	"golang.org/x/crypto/sha3"
)

func WaitInfinitely() {
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func CopyBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	cb := make([]byte, len(b))
	copy(cb, b)
	return cb
}

func SortInt64(arr []int64) []int64 {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })

	return arr
}

// Max returns the larger of x or y.
func MaxInt(x, y int) int {
	if x < y {
		return y
	}
	return x
}

// Min returns the smaller of x or y.
func MinInt(x, y int) int {
	if x > y {
		return y
	}
	return x
}

func GetTxInHash(blockHeight int64, chain string, txBytes []byte) string {
	bz := []byte(chain + strconv.FormatInt(blockHeight, 10))
	bz = append(txBytes, bz...)
	return KeccakHash32(string(bz))
}

// Hash a string and return the first 32 bytes of the hash.
func KeccakHash32(s string) string {
	hash := sha3.NewLegacyKeccak256()

	var buf []byte
	hash.Write([]byte(s))
	buf = hash.Sum(nil)

	encoded := hex.EncodeToString(buf)
	if len(encoded) > 32 {
		encoded = encoded[:32]
	}

	return encoded
}

// Returns a random index in the range [0..size-1] from an integer and a hash. This is a
// determistic "random" value.
func GetRandomIndex(blockHeight int64, hash string, size int) int {
	sum := sha256.Sum256([]byte(hash + strconv.FormatInt(blockHeight, 10)))
	z := new(big.Int)
	z.SetBytes(sum[:])

	// z = z % size
	z = z.Rem(z, big.NewInt(int64(size)))

	return int(z.Int64())
}
