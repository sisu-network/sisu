package cardano

import (
	"github.com/echovl/cardano-go"
)

// MockCardanoClient implements CardanoClient
type MockCardanoClient struct {
	BalanceFunc        func(address cardano.Address) (*cardano.Value, error)
	UTxOsFunc          func(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error)
	TipFunc            func() (*cardano.NodeTip, error)
	ProtocolParamsFunc func() (*cardano.ProtocolParams, error)
	SubmitTxFunc       func(tx *cardano.Tx) (*cardano.Hash32, error)
}

func (m *MockCardanoClient) Balance(address cardano.Address) (*cardano.Value, error) {
	if m.BalanceFunc != nil {
		return m.BalanceFunc(address)
	}

	return nil, nil
}

func (m *MockCardanoClient) UTxOs(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error) {
	if m.UTxOsFunc != nil {
		return m.UTxOsFunc(addr, maxBlock)
	}

	return nil, nil
}

func (m *MockCardanoClient) Tip() (*cardano.NodeTip, error) {
	if m.TipFunc != nil {
		return m.TipFunc()
	}

	return nil, nil
}

func (m *MockCardanoClient) ProtocolParams() (*cardano.ProtocolParams, error) {
	if m.ProtocolParamsFunc != nil {
		return m.ProtocolParamsFunc()
	}

	return nil, nil
}

func (m *MockCardanoClient) SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error) {
	if m.SubmitTxFunc != nil {
		return m.SubmitTxFunc(tx)
	}

	return nil, nil
}
