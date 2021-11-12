package tssclients

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	eTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
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
		log.Error("Cannot check deyes health, err = ", err)
		return err
	}

	return nil
}

// Informs the deyes that Sisu server is ready to accept transaction.
func (c *DeyesClient) SetSisuReady(chain string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setSisuReady", chain)
	if err != nil {
		log.Error("Cannot Set readiness for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

// Adds a list of addresses to watch on a specific chain
func (c *DeyesClient) AddWatchAddresses(chain string, addrs []string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_addWatchAddresses", chain, addrs)
	if err != nil {
		log.Error("Cannot Set readiness for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

func (c *DeyesClient) Dispatch(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error) {
	var result = &eTypes.DispatchedTxResult{}
	err := c.client.CallContext(context.Background(), &result, "deyes_dispatchTx", request)
	if err != nil {
		log.Error("Cannot Dispatch tx to the chain", request.Chain, "err =", err)
		return result, err
	}

	log.Verbose("Tx has been dispatched")

	return result, nil
}

func (c *DeyesClient) GetNonce(chain string, address string) int64 {
	var result int64
	err := c.client.CallContext(context.Background(), &result, "deyes_getNonce", chain, address)
	if err != nil {
		log.Error("Cannot get nonce for chain and address", chain, address, "err =", err)
	}

	return result
}
