package utils

import "sync"

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
