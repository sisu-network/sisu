package types

import (
	"fmt"

	cryptosdk "github.com/cosmos/cosmos-sdk/crypto/types"

	sdked25519 "github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func (valPubkey *ValPubkey) GetCosmosPubkey() (cryptosdk.PubKey, error) {
	switch valPubkey.Type {
	case "ed25519":
		return &sdked25519.PubKey{Key: valPubkey.Bytes}, nil

	case "secp256k1":
		return &secp256k1.PubKey{Key: valPubkey.Bytes}, nil

	default:
		return nil, fmt.Errorf("Invalid key type: %s", valPubkey.Type)
	}
}
