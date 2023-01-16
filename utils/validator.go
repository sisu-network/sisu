package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	cryptosdk "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// GetSortedValidators returns a sorted pubkey based on the hash of each key bytes and a message.
func GetSortedValidators(msg string, nodes []*types.Node) []*types.Node {
	type Wrapper struct {
		key   cryptosdk.PubKey
		index int
		hash  string
	}

	wrappers := make([]Wrapper, 0, len(nodes))
	for i, node := range nodes {
		sha := sha256.New()
		sha.Write(node.ValPubkey.Bytes)
		sha.Write([]byte(msg))
		hash := hex.EncodeToString(sha.Sum(nil))

		sdkKey, err := node.ValPubkey.GetCosmosPubkey()
		if err != nil {
			return nil
		}

		wrappers = append(wrappers, Wrapper{
			key:   sdkKey,
			index: i,
			hash:  hash,
		})
	}

	sort.SliceStable(wrappers, func(i, j int) bool {
		return wrappers[i].hash < wrappers[j].hash
	})

	ret := make([]*types.Node, 0, len(nodes))
	for _, wrapper := range wrappers {
		ret = append(ret, nodes[wrapper.index])
	}

	return ret
}
