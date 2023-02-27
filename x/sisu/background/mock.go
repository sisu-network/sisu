package background

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

///// ManagerContainer

func MockManagerContainer(args ...interface{}) ManagerContainer {
	mc := &DefaultManagerContainer{}

	for _, arg := range args {
		switch t := arg.(type) {
		case components.PostedMessageManager:
			mc.pmm = t
		case components.GlobalData:
			mc.globalData = t
		case external.DeyesClient:
			mc.deyesClient = t
		case external.DheartClient:
			mc.dheartClient = t
		case config.Config:
			mc.config = t
		case components.TxSubmit:
			mc.txSubmit = t
		case components.AppKeys:
			mc.appKeys = t
		case components.TxTracker:
			mc.txTracker = t
		case keeper.Keeper:
			mc.keeper = t
		case sdk.Context:
			mc.readOnlyContext.Store(t)
		case chains.TxOutputProducer:
			mc.txOutProducer = t
		case components.ValidatorManager:
			mc.valsManager = t
		case Background:
			mc.background = t
		case chains.BridgeManager:
			mc.bridgeManager = t
		case keeper.PrivateDb:
			mc.privateDb = t
		}
	}

	return mc
}

////// TxOutputProducer

type MockTxOutputProducer struct {
	GetTxOutsFunc                     func(ctx sdk.Context, chain string, transfers []*types.TransferDetails) ([]*types.TxOut, error)
	PauseContractFunc                 func(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)
	ResumeContractFunc                func(ctx sdk.Context, chain string, hash string) (*types.TxOutMsg, error)
	ContractChangeOwnershipFunc       func(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutMsg, error)
	ContractSetLiquidPoolAddressFunc  func(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutMsg, error)
	ContractEmergencyWithdrawFundFunc func(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutMsg, error)
}

func (m *MockTxOutputProducer) GetTxOuts(ctx sdk.Context, chain string, transfers []*types.TransferDetails) ([]*types.TxOut, error) {
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

////// Background

type MockBackground struct {
	StartFunc         func()
	UpdateFunc        func(ctx sdk.Context)
	AddVoteTxOutFunc  func(height int64, msg *types.TxOutMsg)
	AddRetryTxOutFunc func(height int64, txOut *types.TxOut)
}

func (m *MockBackground) Start() {
	if m.StartFunc != nil {
		m.StartFunc()
	}
}

func (m *MockBackground) Update(ctx sdk.Context) {
	if m.UpdateFunc != nil {
		m.UpdateFunc(ctx)
	}
}

func (m *MockBackground) AddVoteTxOut(height int64, msg *types.TxOutMsg) {
	if m.AddVoteTxOutFunc != nil {
		m.AddVoteTxOutFunc(height, msg)
	}
}

func (m *MockBackground) AddRetryTxOut(height int64, txOut *types.TxOut) {
	if m.AddRetryTxOutFunc != nil {
		m.AddRetryTxOutFunc(height, txOut)
	}
}
