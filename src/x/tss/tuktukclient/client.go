package tuktukclient

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sisu-network/sisu/utils"
)

type Client struct {
	client *rpc.Client
}

// Dial connects a client to the given URL.
func Dial(rawurl string) (*Client, error) {
	return DialContext(context.Background(), rawurl)
}

func DialContext(ctx context.Context, rawurl string) (*Client, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return NewClient(c), nil
}

func NewClient(c *rpc.Client) *Client {
	return &Client{c}
}

func (tc *Client) GetVersion() (string, error) {
	var result string
	err := tc.client.CallContext(context.Background(), &result, "tss_version")
	if err != nil {
		utils.LogError("Cannot get TSS version, err = ", err)
		return "", err
	}

	return result, nil
}

func (tc *Client) KeyGen(chainSymbol string) error {
	var result string
	err := tc.client.CallContext(context.Background(), &result, "tss_keyGen", chainSymbol)
	if err != nil {
		utils.LogError("Cannot send keygen request, err = ", err)
		return err
	}

	return nil
}
