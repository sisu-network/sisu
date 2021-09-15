package tssclients

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sisu-network/sisu/utils"
)

type DeyesClient struct {
	client *rpc.Client
}

func DialDeyes(rawurl string) (*DeyesClient, error) {
	return dialDeyesContext(context.Background(), rawurl)
}

func dialDeyesContext(ctx context.Context, rawurl string) (*DeyesClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDeyesClient(c), nil
}

func newDeyesClient(c *rpc.Client) *DeyesClient {
	return &DeyesClient{c}
}

func (c *DeyesClient) CheckHealth() error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "deyes_checkHealth")
	if err != nil {
		utils.LogError("Cannot check deyes health, err = ", err)
		return err
	}

	return nil
}

// Informs the deyes that Sisu server is ready to accept transaction.
func (c *DeyesClient) SetSisuReady(chain string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setSisuReady", chain)
	if err != nil {
		utils.LogError("Cannot Set readiness for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

// Adds a list of addresses to watch on a specific chain
func (c *DeyesClient) AddWatchAddresses(chain string, addrs []string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_addWatchAddresses", chain, addrs)
	if err != nil {
		utils.LogError("Cannot Set readiness for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

func (c *DeyesClient) Dispatch(chain string, tx []byte) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_dispatchTx", chain, tx)
	if err != nil {
		utils.LogError("Cannot Dispatch tx to the chain", chain, "err =", err)
		return err
	}

	return nil
}
