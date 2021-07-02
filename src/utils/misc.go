package utils

import (
	"sort"
	"sync"
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
