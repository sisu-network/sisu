package tssclients

import (
	"context"

	"github.com/ethereum/go-ethereum/rpc"
	eTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/lib/log"
)

type DeyesClient interface {
	Ping(source string) error
	Dispatch(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error)
	SetGatewayAddress(chain string, addr string) error
	GetNonce(chain string, address string) int64
	SetSisuReady(isReady bool) error
	GetGasPrices(chains []string) ([]int64, error)
}

type defaultDeyesClient struct {
	client *rpc.Client
}

func DialDeyes(rawurl string) (DeyesClient, error) {
	return dialDeyesContext(context.Background(), rawurl)
}

func dialDeyesContext(ctx context.Context, rawurl string) (DeyesClient, error) {
	c, err := rpc.DialContext(ctx, rawurl)
	if err != nil {
		return nil, err
	}
	return newDeyesClient(c), nil
}

func newDeyesClient(c *rpc.Client) DeyesClient {
	return &defaultDeyesClient{c}
}

func (c *defaultDeyesClient) Ping(source string) error {
	var result interface{}
	err := c.client.CallContext(context.Background(), &result, "deyes_ping", source)
	if err != nil {
		log.Error("Cannot ping deyes, err = ", err)
		return err
	}

	return nil
}

// Informs the deyes that Sisu server is ready to accept transaction.
func (c *defaultDeyesClient) SetSisuReady(isReady bool) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setSisuReady", isReady)
	if err != nil {
		log.Error("Cannot Set readiness for deyes, err = ", err)
		return err
	}

	return nil
}

func (c *defaultDeyesClient) SetGatewayAddress(chain string, addr string) error {
	var result string
	err := c.client.CallContext(context.Background(), &result, "deyes_setGatewayAddress", chain, addr)
	if err != nil {
		log.Error("Cannot set gateway address for deyes, chain = ", chain, "err = ", err)
		return err
	}

	return nil
}

func (c *defaultDeyesClient) Dispatch(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error) {
	var result = &eTypes.DispatchedTxResult{}
	err := c.client.CallContext(context.Background(), &result, "deyes_dispatchTx", request)
	if err != nil {
		log.Error("Cannot Dispatch tx to the chain", request.Chain, "err =", err)
		return result, err
	}

	log.Verbose("Tx has been dispatched")

	return result, nil
}

func (c *defaultDeyesClient) GetNonce(chain string, address string) int64 {
	var result int64
	err := c.client.CallContext(context.Background(), &result, "deyes_getNonce", chain, address)
	if err != nil {
		log.Error("Cannot get nonce for chain and address", chain, address, "err =", err)
	}

	return result
}

func (c *defaultDeyesClient) GetGasPrices(chains []string) ([]int64, error) {
	result := make([]int64, 0)
	err := c.client.CallContext(context.Background(), &result, "deyes_getGasPrices", chains)
	if err != nil {
		log.Error("Cannot get gas price for chains = ", chains, "err = ", err)
		return nil, err
	}

	return result, nil
}
