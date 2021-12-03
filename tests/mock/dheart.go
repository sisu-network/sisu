package mock

import (
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	htypes "github.com/sisu-network/dheart/types"
)

type DheartClient struct {
}

func (c *DheartClient) SetPrivKey(encodedKey string, keyType string) error {
	return nil
}

func (c *DheartClient) CheckHealth() error {
	return nil
}

func (c *DheartClient) KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error {
	return nil
}

func (c *DheartClient) KeySign(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error {
	return nil
}
