package tssclients

import eTypes "github.com/sisu-network/deyes/types"

type MockDeyesClient struct {
	PingFunc              func(source string) error
	DispatchFunc          func(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error)
	AddWatchAddressesFunc func(chain string, addrs []string) error
	GetNonceFunc          func(chain string, address string) int64
	SetSisuReadyFunc      func(isReady bool) error
}

func (c *MockDeyesClient) Ping(source string) error {
	if c.PingFunc != nil {
		return c.PingFunc(source)
	}

	return nil
}

func (c *MockDeyesClient) Dispatch(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error) {
	if c.DispatchFunc != nil {
		return c.DispatchFunc(request)
	}

	return nil, nil
}

func (c *MockDeyesClient) AddWatchAddresses(chain string, addrs []string) error {
	if c.AddWatchAddressesFunc != nil {
		return c.AddWatchAddressesFunc(chain, addrs)
	}

	return nil
}

func (c *MockDeyesClient) GetNonce(chain string, address string) int64 {
	if c.GetNonceFunc != nil {
		return c.GetNonceFunc(chain, address)
	}

	return 0
}

func (c *MockDeyesClient) SetSisuReady(isReady bool) error {
	if c.SetSisuReadyFunc != nil {
		return c.SetSisuReadyFunc(isReady)
	}

	return nil
}
