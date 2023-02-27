package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/echovl/cardano-go"
)

// A function to make sure that all mocks implement its designated interface.
func checkMock() {
	var _ chains.TxOutputProducer = new(background.MockTxOutputProducer)
	var _ components.TxTracker = new(MockTxTracker)
	var _ components.PostedMessageManager = new(MockPostedMessageManager)
}

///// TxTracker

type MockTxTracker struct {
	AddTransactionFunc          func(txOut *types.TxOut)
	UpdateStatusFunc            func(chain string, hash string, status types.TxStatus)
	RemoveTransactionFunc       func(chain string, hash string)
	OnTxFailedFunc              func(chain string, hash string, status types.TxStatus)
	CheckExpiredTransactionFunc func()
}

func (m *MockTxTracker) AddTransaction(txOut *types.TxOut) {
	if m.AddTransactionFunc != nil {
		m.AddTransactionFunc(txOut)
	}
}

func (m *MockTxTracker) UpdateStatus(chain string, hash string, status types.TxStatus) {
	if m.UpdateStatusFunc != nil {
		m.UpdateStatusFunc(chain, hash, status)
	}
}

func (m *MockTxTracker) RemoveTransaction(chain string, hash string) {
	if m.RemoveTransactionFunc != nil {
		m.RemoveTransactionFunc(chain, hash)
	}
}

func (m *MockTxTracker) OnTxFailed(chain string, hash string, status types.TxStatus) {
	if m.OnTxFailedFunc != nil {
		m.OnTxFailedFunc(chain, hash, status)
	}
}

func (m *MockTxTracker) CheckExpiredTransaction() {
	if m.CheckExpiredTransactionFunc != nil {
		m.CheckExpiredTransactionFunc()
	}
}

///// PostedMessageManager

type MockPostedMessageManager struct {
	ShouldProcessMsgFunc func(ctx sdk.Context, msg sdk.Msg) (bool, []byte)
}

func (m *MockPostedMessageManager) ShouldProcessMsg(ctx sdk.Context, msg sdk.Msg) (bool, []byte) {
	if m.ShouldProcessMsgFunc != nil {
		return m.ShouldProcessMsgFunc(ctx, msg)
	}

	return false, nil
}

///// CardanoNode

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

///// TxInQueue

type MockTransferQueue struct {
	StartFunc                         func(ctx sdk.Context)
	ProcessTransfersFunc              func(ctx sdk.Context)
	ClearInMemoryPendingTransfersFunc func(chain string)
}

func (m *MockTransferQueue) Start(ctx sdk.Context) {
	if m.StartFunc != nil {
		m.StartFunc(ctx)
	}
}
func (m *MockTransferQueue) ProcessTransfers(ctx sdk.Context) {
	if m.ProcessTransfersFunc != nil {
		m.ProcessTransfersFunc(ctx)
	}
}
func (m *MockTransferQueue) ClearInMemoryPendingTransfers(chain string) {
	if m.ClearInMemoryPendingTransfersFunc != nil {
		m.ClearInMemoryPendingTransfersFunc(chain)
	}
}

///// TxOutQueue

type MockTxOutQueue struct {
	StartFunc         func()
	AddTxOutFunc      func(txOut *types.TxOut)
	ProcessTxOutsFunc func(ctx sdk.Context)
}

func (m *MockTxOutQueue) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

func (m *MockTxOutQueue) AddTxOut(txOut *types.TxOut) {
	if m.AddTxOutFunc != nil {
		m.AddTxOutFunc(txOut)
	}
}

func (m *MockTxOutQueue) ProcessTxOuts(ctx sdk.Context) {
	if m.ProcessTxOutsFunc != nil {
		m.ProcessTxOutsFunc(ctx)
	}
}
