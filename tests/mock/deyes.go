package mock

import (
	eTypes "github.com/sisu-network/deyes/types"
	"math/rand"
)

type DeyesClient struct {}

func (c *DeyesClient) CheckHealth() error {
	return nil
}

// Informs the deyes that Sisu server is ready to accept transaction.
func (c *DeyesClient) SetSisuReady(chain string) error {
	return nil
}

// Adds a list of addresses to watch on a specific chain
func (c *DeyesClient) AddWatchAddresses(chain string, addrs []string) error {
	return nil
}

func (c *DeyesClient) Dispatch(_ *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error) {
	var result = &eTypes.DispatchedTxResult{}
	return result, nil
}

func (c *DeyesClient) GetNonce(_ string, _ string) int64 {
	return rand.Int63()
}
