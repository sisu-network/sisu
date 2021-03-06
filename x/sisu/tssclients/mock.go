package tssclients

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	eTypes "github.com/sisu-network/deyes/types"
	htypes "github.com/sisu-network/dheart/types"
)

func check() {
	var _ DeyesClient = new(MockDeyesClient)
	var _ DheartClient = new(MockDheartClient)
}

///// DeyesClient

type MockDeyesClient struct {
	PingFunc              func(source string) error
	DispatchFunc          func(request *eTypes.DispatchedTxRequest) (*eTypes.DispatchedTxResult, error)
	SetChainAccountFunc   func(chain string, addr string) error
	SetGatewayAddressFunc func(chain string, addr string) error
	GetNonceFunc          func(chain string, address string) int64
	SetSisuReadyFunc      func(isReady bool) error
	GetGasPricesFunc      func(chains []string) ([]int64, error)
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

func (c *MockDeyesClient) SetChainAccount(chain string, addr string) error {
	if c.SetChainAccountFunc != nil {
		return c.SetChainAccountFunc(chain, addr)
	}

	return nil
}

func (c *MockDeyesClient) SetGatewayAddress(chain string, addr string) error {
	if c.SetGatewayAddressFunc != nil {
		return c.SetGatewayAddressFunc(chain, addr)
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

func (c *MockDeyesClient) GetGasPrices(chains []string) ([]int64, error) {
	if c.GetGasPricesFunc != nil {
		return c.GetGasPricesFunc(chains)
	}

	return nil, nil
}

///// DheartClient

type MockDheartClient struct {
	SetPrivKeyFunc   func(encodedKey string, keyType string) error
	PingFunc         func(string) error
	KeyGenFunc       func(keygenId string, chain string, pubKeys []ctypes.PubKey) error
	KeySignFunc      func(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error
	BlockEndFunc     func(blockHeight int64) error
	SetSisuReadyFunc func(isReady bool) error
}

func (m *MockDheartClient) SetPrivKey(encodedKey string, keyType string) error {
	if m.SetPrivKeyFunc != nil {
		return m.SetPrivKeyFunc(encodedKey, keyType)
	}

	return nil
}

func (m *MockDheartClient) Ping(s string) error {
	if m.PingFunc != nil {
		return m.PingFunc(s)
	}

	return nil
}

func (m *MockDheartClient) KeyGen(keygenId string, chain string, pubKeys []ctypes.PubKey) error {
	if m.KeyGenFunc != nil {
		return m.KeyGenFunc(keygenId, chain, pubKeys)
	}

	return nil
}

func (m *MockDheartClient) KeySign(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error {
	if m.KeySignFunc != nil {
		return m.KeySignFunc(req, pubKeys)
	}

	return nil
}

func (m *MockDheartClient) BlockEnd(blockHeight int64) error {
	if m.BlockEndFunc != nil {
		return m.BlockEndFunc(blockHeight)
	}

	return nil
}

func (m *MockDheartClient) SetSisuReady(isReady bool) error {
	if m.SetSisuReadyFunc != nil {
		return m.SetSisuReadyFunc(isReady)
	}

	return nil
}
