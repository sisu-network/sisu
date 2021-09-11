package tssclients

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/rpc"
	dTypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/utils"
)

type DheartClient struct {
	client *rpc.Client
}

// DialDheart connects a client to the given URL.
func DialDheart(rawurl string) (*DheartClient, error) {
	return dialDheartContext(context.Background(), rawurl)
}

func dialDheartContext(ctx context.Context, rawurl string) (*DheartClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDheartClient(c), nil
}

func newDheartClient(c *rpc.Client) *DheartClient {
	return &DheartClient{c}
}

func (c *DheartClient) SetPrivKey(encodedKey string, keyType string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_setPrivKey", encodedKey, keyType)
	if err != nil {
		utils.LogError("Cannot do handshare with dheart, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) CheckHealth() error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "tss_checkHealth")
	if err != nil {
		utils.LogError("Cannot check Dheart health, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error {
	// Wrap pubkeys
	wrappers := make([]dTypes.PubKeyWrapper, len(pubKeys))
	for i, pubKey := range pubKeys {
		switch pubKey.Type() {
		case "ed25519":
			wrappers[i] = dTypes.PubKeyWrapper{
				KeyType: pubKey.Type(),
				Ed25519: &ed25519.PubKey{
					Key: pubKey.Bytes(),
				},
			}
		case "secp256k1":
			wrappers[i] = dTypes.PubKeyWrapper{
				KeyType: pubKey.Type(),
				Secp256k1: &secp256k1.PubKey{
					Key: pubKey.Bytes(),
				},
			}
		}
	}

	utils.LogInfo("Broadcasting keygen to Dheart")

	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_keyGen", keygenId, chain, wrappers)
	if err != nil {
		utils.LogError("Cannot send keygen request, err = ", err)
		return err
	}

	return nil
}

func (c *DheartClient) KeySign(req *dTypes.KeysignRequest) error {
	utils.LogVerbose("Broadcasting key signing to Dheart")

	var r interface{}
	err := c.client.CallContext(context.Background(), &r, "tss_keySign", req)
	if err != nil {
		utils.LogError("Cannot send KeySign request, err = ", err)
		return err
	}

	return nil
}
