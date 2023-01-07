package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	tcrypto "github.com/tendermint/tendermint/crypto"
)

// GetSortedValidators returns a sorted pubkey based on the hash of each key bytes and a message.
func GetSortedValidators(msg string, keys []tcrypto.PubKey) []tcrypto.PubKey {
	type Wrapper struct {
		key   tcrypto.PubKey
		index int
		hash  string
	}

	wrappers := make([]Wrapper, 0, len(keys))
	for i, key := range keys {
		sha := sha256.New()
		sha.Write(key.Bytes())
		sha.Write([]byte(msg))
		hash := hex.EncodeToString(sha.Sum(nil))

		wrappers = append(wrappers, Wrapper{
			key:   key,
			index: i,
			hash:  hash,
		})
	}

	sort.SliceStable(wrappers, func(i, j int) bool {
		return wrappers[i].hash < wrappers[j].hash
	})

	ret := make([]tcrypto.PubKey, 0, len(keys))
	for _, wrapper := range wrappers {
		ret = append(ret, keys[wrapper.index])
	}

	return ret
}
