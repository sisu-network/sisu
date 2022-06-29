package sisu

import (
	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
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
		case tssclients.DeyesClient:
			mc.deyesClient = t
		case tssclients.DheartClient:
			mc.dheartClient = t
		case config.TssConfig:
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
		case world.WorldState:
			mc.worldState = t
		case ValidatorManager:
			mc.valsManager = t
		case TxInQueue:
			mc.txInQueue = t
		case TxOutQueue:
			mc.txOutQueue = t
		}
	}

	return mc
}

///// TxTracker

type MockTxTracker struct {
	AddTransactionFunc          func(txOut *types.TxOut, txIn *types.TxIn)
	UpdateStatusFunc            func(chain string, hash string, status types.TxStatus)
	RemoveTransactionFunc       func(chain string, hash string)
	OnTxFailedFunc              func(chain string, hash string, status types.TxStatus)
	CheckExpiredTransactionFunc func()
}

func (m *MockTxTracker) AddTransaction(txOut *types.TxOut, txIn *types.TxIn) {
	if m.AddTransactionFunc != nil {
		m.AddTransactionFunc(txOut, txIn)
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
	GetTxOutsFunc                     func(ctx sdk.Context, height int64, tx []*types.TxIn) []*types.TxOutWithSigner
	PauseContractFunc                 func(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)
	ResumeContractFunc                func(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error)
	ContractChangeOwnershipFunc       func(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error)
	ContractSetLiquidPoolAddressFunc  func(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error)
	ContractEmergencyWithdrawFundFunc func(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error)
}

func (m *MockTxOutputProducer) GetTxOuts(ctx sdk.Context, height int64, tx []*types.TxIn) []*types.TxOutWithSigner {
	if m.GetTxOutsFunc != nil {
		return m.GetTxOutsFunc(ctx, height, tx)
	}

	return nil
}

func (m *MockTxOutputProducer) PauseContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error) {
	if m.PauseContractFunc != nil {
		return m.PauseContractFunc(ctx, chain, hash)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ResumeContract(ctx sdk.Context, chain string, hash string) (*types.TxOutWithSigner, error) {
	if m.ResumeContractFunc != nil {
		return m.ResumeContractFunc(ctx, chain, hash)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractChangeOwnership(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error) {
	if m.ContractChangeOwnershipFunc != nil {
		return m.ContractChangeOwnershipFunc(ctx, chain, contractHash, newOwner)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractSetLiquidPoolAddress(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error) {
	if m.ContractSetLiquidPoolAddressFunc != nil {
		return m.ContractSetLiquidPoolAddressFunc(ctx, chain, contractHash, newAddress)
	}

	return nil, nil
}

func (m *MockTxOutputProducer) ContractEmergencyWithdrawFund(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error) {
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

///// TxInQueue

type MockTxInQueue struct {
	StartFunc        func()
	AddTxInFunc      func(txIn *types.TxIn)
	ProcessTxInsFunc func()
}

func (m *MockTxInQueue) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

func (m *MockTxInQueue) AddTxIn(txIn *types.TxIn) {
	if m.AddTxInFunc != nil {
		m.AddTxInFunc(txIn)
	}
}

func (m *MockTxInQueue) ProcessTxIns() {
	if m.ProcessTxInsFunc != nil {
		m.ProcessTxInsFunc()
	}
}

///// TxOutQueue

type MockTxOutQueue struct {
	StartFunc         func()
	AddTxOutFunc      func(txOut *types.TxOut)
	ProcessTxOutsFunc func()
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

func (m *MockTxOutQueue) ProcessTxOuts() {
	if m.ProcessTxOutsFunc != nil {
		m.ProcessTxOutsFunc()
	}
}
