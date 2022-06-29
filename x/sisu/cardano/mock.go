package cardano

import "github.com/echovl/cardano-go"

type MockCardanoNode struct {
	// UTxOs returns a list of unspent transaction outputs for a given address
	UTxOsFunc func(cardano.Address) ([]cardano.UTxO, error)

	// Tip returns the node's current tip
	TipFunc func() (*cardano.NodeTip, error)

	// SubmitTx submits a transaction to the node using cbor encoding
	SubmitTxFunc func(*cardano.Tx) (*cardano.Hash32, error)

	// ProtocolParams returns the Node's Protocol Parameters
	ProtocolParamsFunc func() (*cardano.ProtocolParams, error)

	// Network returns the node's current network type
	NetworkFunc func() cardano.Network
}

// UTxOs returns a list of unspent transaction outputs for a given address
func (m *MockCardanoNode) UTxOs(addr cardano.Address) ([]cardano.UTxO, error) {
	if m.UTxOsFunc != nil {
		return m.UTxOsFunc(addr)
	}

	return nil, nil
}

// Tip returns the node's current tip
func (m *MockCardanoNode) Tip() (*cardano.NodeTip, error) {
	if m.TipFunc != nil {
		return m.TipFunc()
	}

	return nil, nil
}

// SubmitTx submits a transaction to the node using cbor encoding
func (m *MockCardanoNode) SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error) {
	if m.SubmitTxFunc != nil {
		return m.SubmitTxFunc(tx)
	}

	return nil, nil
}

// ProtocolParams returns the Node's Protocol Parameters
func (m *MockCardanoNode) ProtocolParams() (*cardano.ProtocolParams, error) {
	if m.ProtocolParamsFunc != nil {
		return m.ProtocolParamsFunc()
	}

	return nil, nil
}

// Network returns the node's current network type
func (m *MockCardanoNode) Network() cardano.Network {
	if m.NetworkFunc != nil {
		return m.NetworkFunc()
	}

	return cardano.Testnet
}
