package tssclients

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	ctypes "github.com/sisu-network/cosmos-sdk/crypto/types"
	"github.com/sisu-network/lib/log"
)

//go:generate mockgen -source=x/tss/tssclients/dheart_client.go -destination=tests/mock/tss/tssclients/dheart_client.go -package=mock

type DheartClient interface {
	SetPrivKey(encodedKey string, keyType string) error
	CheckHealth() error
	KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error
	BlockEnd(blockHeight int64) error
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

func (c *DefaultDheartClient) CheckHealth() error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "tss_checkHealth")
	if err != nil {
		log.Error("Cannot check Dheart health, err = ", err)
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
