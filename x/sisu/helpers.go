package sisu

import (
	"crypto"
	"encoding/hex"
	"fmt"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

func GetKeygenId(keyType string, block int64, pubKeys []ctypes.PubKey) string {
	// Get hashes of all pubkeys
	digester := crypto.MD5.New()
	for _, pubKey := range pubKeys {
		fmt.Fprint(digester, pubKey.Bytes())
	}
	hash := hex.EncodeToString(digester.Sum(nil))

	return fmt.Sprintf("%s;%d;%s", keyType, block, hash)
}
