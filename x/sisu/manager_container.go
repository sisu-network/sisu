package sisu

import (
	"sync/atomic"

	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/service"
)

type ManagerContainer interface {
	PostedMessageManager() PostedMessageManager
	PartyManager() components.PartyManager
	DheartClient() external.DheartClient
	DeyesClient() external.DeyesClient
	GlobalData() components.GlobalData
	TxSubmit() components.TxSubmit
	Config() config.Config
	AppKeys() components.AppKeys
	TxOutProducer() chains.TxOutputProducer
	TxTracker() components.TxTracker
	Keeper() keeper.Keeper
	ValidatorManager() components.ValidatorManager
	TransferQueue() TransferQueue
	BridgeManager() chains.BridgeManager
	ChainPolling() service.ChainPolling
	PrivateDb() keeper.PrivateDb
}

type DefaultManagerContainer struct {
	readOnlyContext atomic.Value

	pmm              PostedMessageManager
	partyManager     components.PartyManager
	dheartClient     external.DheartClient
	deyesClient      external.DeyesClient
	globalData       components.GlobalData
	txSubmit         components.TxSubmit
	config           config.Config
	appKeys          components.AppKeys
	txOutProducer    chains.TxOutputProducer
	txTracker        components.TxTracker
	keeper           keeper.Keeper
	valsManager      components.ValidatorManager
	transferOutQueue TransferQueue
	bridgeManager    chains.BridgeManager
	chainPolling     service.ChainPolling
	privateDb        keeper.PrivateDb
}

func NewManagerContainer(pmm PostedMessageManager, partyManager components.PartyManager,
	dheartClient external.DheartClient, deyesClient external.DeyesClient,
	globalData components.GlobalData, txSubmit components.TxSubmit, cfg config.Config,
	appKeys components.AppKeys, txOutProducer chains.TxOutputProducer, txTracker components.TxTracker,
	keeper keeper.Keeper, valsManager components.ValidatorManager, txInQueue TransferQueue,
	bridgeManager chains.BridgeManager, chainPolling service.ChainPolling,
	privateDb keeper.PrivateDb) ManagerContainer {
	return &DefaultManagerContainer{
		pmm:              pmm,
		partyManager:     partyManager,
		dheartClient:     dheartClient,
		deyesClient:      deyesClient,
		globalData:       globalData,
		txSubmit:         txSubmit,
		config:           cfg,
		appKeys:          appKeys,
		txOutProducer:    txOutProducer,
		txTracker:        txTracker,
		keeper:           keeper,
		valsManager:      valsManager,
		transferOutQueue: txInQueue,
		bridgeManager:    bridgeManager,
		chainPolling:     chainPolling,
		privateDb:        privateDb,
	}
}

func (mc *DefaultManagerContainer) PostedMessageManager() PostedMessageManager {
	return mc.pmm
}

func (mc *DefaultManagerContainer) PartyManager() components.PartyManager {
	return mc.partyManager
}

func (mc *DefaultManagerContainer) DheartClient() external.DheartClient {
	return mc.dheartClient
}

func (mc *DefaultManagerContainer) GlobalData() components.GlobalData {
	return mc.globalData
}

func (mc *DefaultManagerContainer) TxSubmit() components.TxSubmit {
	return mc.txSubmit
}

func (mc *DefaultManagerContainer) Config() config.Config {
	return mc.config
}

func (mc *DefaultManagerContainer) AppKeys() components.AppKeys {
	return mc.appKeys
}

func (mc *DefaultManagerContainer) TxOutProducer() chains.TxOutputProducer {
	return mc.txOutProducer
}

func (mc *DefaultManagerContainer) DeyesClient() external.DeyesClient {
	return mc.deyesClient
}

func (mc *DefaultManagerContainer) TxTracker() components.TxTracker {
	return mc.txTracker
}

func (mc *DefaultManagerContainer) Keeper() keeper.Keeper {
	return mc.keeper
}

func (mc *DefaultManagerContainer) ValidatorManager() components.ValidatorManager {
	return mc.valsManager
}

func (mc *DefaultManagerContainer) TransferQueue() TransferQueue {
	return mc.transferOutQueue
}

func (mc *DefaultManagerContainer) BridgeManager() chains.BridgeManager {
	return mc.bridgeManager
}

func (mc *DefaultManagerContainer) ChainPolling() service.ChainPolling {
	return mc.chainPolling
}

func (mc *DefaultManagerContainer) PrivateDb() keeper.PrivateDb {
	return mc.privateDb
}
