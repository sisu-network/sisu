package tssclients

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sisu-network/sisu/utils"
)

type TuktukClient struct {
	client *rpc.Client
}

// DialTuktuk connects a client to the given URL.
func DialTuktuk(rawurl string) (*TuktukClient, error) {
	return dialTuktukContext(context.Background(), rawurl)
}

func dialTuktukContext(ctx context.Context, rawurl string) (*TuktukClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newTuktukClient(c), nil
}

func newTuktukClient(c *rpc.Client) *TuktukClient {
	return &TuktukClient{c}
}

func (c *TuktukClient) CheckHealth() error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "tss_checkHealth")
	if err != nil {
		utils.LogError("Cannot check tuktuk health, err = ", err)
		return err
	}

	return nil
}

func (c *TuktukClient) KeyGen(chainSymbol string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "tss_keyGen", chainSymbol)
	if err != nil {
		utils.LogError("Cannot send keygen request, err = ", err)
		return err
	}

	return nil
}
