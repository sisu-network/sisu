package tssclients

import (
	"context"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/ethereum/go-ethereum/rpc"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/lib/log"
)

//go:generate mockgen -source=x/sisu/tssclients/dheart_client.go -destination=tests/mock/x/sisu/tssclients/dheart_client.go -package=mock

type DheartClient interface {
	SetPrivKey(encodedKey string, keyType string) error
	Ping(string) error
	KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error
	KeySign(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error
	Reshare(oldPubKeys, newPubKeys []ctypes.PubKey) error
	BlockEnd(blockHeight int64) error
	SetSisuReady(isReady bool) error
}

type DefaultDheartClient struct {
	client *rpc.Client
}

// DialDheart connects a client to the given URL.
func DialDheart(rawurl string) (*DefaultDheartClient, error) {
	return dialDheartContext(context.Background(), rawurl)
}

func dialDheartContext(ctx context.Context, rawurl string) (*DefaultDheartClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDefaultDheartClient(c), nil
}

func newDefaultDheartClient(c *rpc.Client) *DefaultDheartClient {
	return &DefaultDheartClient{c}
}

func (c *DefaultDheartClient) SetPrivKey(encodedKey string, keyType string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_setPrivKey", encodedKey, keyType)
	if err != nil {
		log.Error("Cannot do handshare with dheart, err = ", err)
		return err
	}

	return nil
}

func (c *DefaultDheartClient) Ping(source string) error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "tss_ping", source)
	if err != nil {
		log.Error("Cannot ping dheart, err = ", err)
		return err
	}

	return nil
}

func (c *DefaultDheartClient) KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error {
	// Wrap pubkeys
	wrappers := c.getPubkeyWrapper(pubKeys)

	log.Info("Broadcasting keygen to Dheart")

	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_keyGen", keygenId, chain, wrappers)
	if err != nil {
		log.Error("Cannot send keygen request, err = ", err)
		return err
	}

	return nil
}

func (c *DefaultDheartClient) getPubkeyWrapper(pubKeys []ctypes.PubKey) []htypes.PubKeyWrapper {
	wrappers := make([]htypes.PubKeyWrapper, len(pubKeys))
	for i, pubKey := range pubKeys {
		switch pubKey.Type() {
		case "ed25519":
			wrappers[i] = htypes.PubKeyWrapper{
				KeyType: pubKey.Type(),
				Key:     pubKey.Bytes(),
			}
		case "secp256k1":
			wrappers[i] = htypes.PubKeyWrapper{
				KeyType: pubKey.Type(),
				Key:     pubKey.Bytes(),
			}
		}
	}

	return wrappers
}

func (c *DefaultDheartClient) KeySign(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error {
	log.Verbose("Broadcasting key signing to Dheart")

	wrappers := c.getPubkeyWrapper(pubKeys)

	var r interface{}
	err := c.client.CallContext(context.Background(), &r, "tss_keySign", req, wrappers)
	if err != nil {
		log.Error("Cannot send KeySign request, err = ", err)
		return err
	}

	return nil
}

func (c *DefaultDheartClient) BlockEnd(blockHeight int64) error {
	var r interface{}
	err := c.client.CallContext(context.Background(), &r, "tss_blockEnd", blockHeight)
	if err != nil {
		log.Error("Cannot call block end, err = ", err)
		return err
	}

	return nil
}

func (c *DefaultDheartClient) SetSisuReady(isReady bool) error {
	var r interface{}
	err := c.client.CallContext(context.Background(), &r, "tss_setSisuReady", isReady)
	if err != nil {
		log.Error("Cannot call SetSisuReady, err = ", err)
		return err
	}

	return nil
}
