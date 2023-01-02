package sisu

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains"
	external "github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/echovl/cardano-go"
)

// A function to make sure that all mocks implement its designated interface.
func checkMock() {
	var _ TxOutputProducer = new(MockTxOutputProducer)
	var _ TxTracker = new(MockTxTracker)
	var _ PostedMessageManager = new(MockPostedMessageManager)
	var _ PartyManager = new(MockPartyManager)
}

///// ManagerContainer

func MockManagerContainer(args ...interface{}) ManagerContainer {
	mc := &DefaultManagerContainer{}

	for _, arg := range args {
		switch t := arg.(type) {
		case PostedMessageManager:
			mc.pmm = t
		case common.GlobalData:
			mc.globalData = t
		case external.DeyesClient:
			mc.deyesClient = t
		case external.DheartClient:
			mc.dheartClient = t
		case config.Config:
			mc.config = t
		case common.TxSubmit:
			mc.txSubmit = t
		case common.AppKeys:
			mc.appKeys = t
		case PartyManager:
			mc.partyManager = t
		case TxTracker:
			mc.txTracker = t
		case keeper.Keeper:
			mc.keeper = t
		case sdk.Context:
			mc.readOnlyContext.Store(t)
		case TxOutputProducer:
			mc.txOutProducer = t
		case ValidatorManager:
			mc.valsManager = t
		case TransferQueue:
			mc.transferOutQueue = t
		case chains.BridgeManager:
			mc.bridgeManager = t
		case keeper.PrivateDb:
			mc.privateDb = t
		}
	}

	return mc
}

///// TxTracker

type MockTxTracker struct {
	AddTransactionFunc          func(txOut *types.TxOutOld)
	UpdateStatusFunc            func(chain string, hash string, status types.TxStatus)
	RemoveTransactionFunc       func(chain string, hash string)
	OnTxFailedFunc              func(chain string, hash string, status types.TxStatus)
	CheckExpiredTransactionFunc func()
}

func (m *MockTxTracker) AddTransaction(txOut *types.TxOutOld) {
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

////// TxOutputProducer

type MockTxOutputProducer struct {
	GetTxOutsFunc                     func(ctx sdk.Context, chain string, transfers []*types.Transfer) ([]*types.TxOutMsg, error)
	PauseContractFunc                 func(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)
	ResumeContractFunc                func(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)
	ContractChangeOwnershipFunc       func(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutMsg, error)
	ContractSetLiquidPoolAddressFunc  func(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutMsg, error)
	ContractEmergencyWithdrawFundFunc func(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutMsg, error)
}

func (m *MockTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string, transfers []*types.Transfer) ([]*types.TxOutMsg, error) {
	if m.GetTxOutsFunc != nil {
		return m.GetTxOutsFunc(ctx, chain, transfers)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error) {
	if m.PauseContractFunc != nil {
		return m.PauseContractFunc(ctx, chain, hash)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error) {
	if m.ResumeContractFunc != nil {
		return m.ResumeContractFunc(ctx, chain, hash)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutMsg, error) {
	if m.ContractChangeOwnershipFunc != nil {
		return m.ContractChangeOwnershipFunc(ctx, chain, contractHash, newOwner)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutMsg, error) {
	if m.ContractSetLiquidPoolAddressFunc != nil {
		return m.ContractSetLiquidPoolAddressFunc(ctx, chain, contractHash, newAddress)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutMsg, error) {
	if m.ContractEmergencyWithdrawFundFunc != nil {
		return m.ContractEmergencyWithdrawFundFunc(ctx, chain, contractHash, tokens, newOwner)
	}

	return nil, nil
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

///// PartyManager

type MockPartyManager struct {
	GetActivePartyPubkeysFunc func() []ctypes.PubKey
}

func (m *MockPartyManager) GetActivePartyPubkeys() []ctypes.PubKey {
	if m.GetActivePartyPubkeysFunc != nil {
		return m.GetActivePartyPubkeysFunc()
	}

	return nil
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
	AddTxOutFunc      func(txOut *types.TxOutOld)
	ProcessTxOutsFunc func(ctx sdk.Context)
}

func (m *MockTxOutQueue) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

func (m *MockTxOutQueue) AddTxOut(txOut *types.TxOutOld) {
	if m.AddTxOutFunc != nil {
		m.AddTxOutFunc(txOut)
	}
}

func (m *MockTxOutQueue) ProcessTxOuts(ctx sdk.Context) {
	if m.ProcessTxOutsFunc != nil {
		m.ProcessTxOutsFunc(ctx)
	}
}
