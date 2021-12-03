package mock

import (
	"github.com/sisu-network/cosmos-sdk/crypto/hd"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	"github.com/sisu-network/cosmos-sdk/types"
)

type KeyringInfo struct {}

func (k *KeyringInfo) GetType() keyring.KeyType {
	return keyring.TypeLocal
}

func (k *KeyringInfo) GetName() string {
	return ""
}

func (k *KeyringInfo) GetPubKey() cryptotypes.PubKey {
	return nil
}

func (k *KeyringInfo) GetAddress() types.AccAddress {
	return nil
}

func (k *KeyringInfo) GetPath() (*hd.BIP44Params, error) {
	return nil, nil
}

func (k *KeyringInfo) GetAlgo() hd.PubKeyType {
	return hd.Ed25519Type
}
