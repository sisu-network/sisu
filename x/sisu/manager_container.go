package sisu

import (
	"sync/atomic"

	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/chains"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/service"
)

type ManagerContainer interface {
	PostedMessageManager() PostedMessageManager
	PartyManager() PartyManager
	DheartClient() external.DheartClient
	DeyesClient() external.DeyesClient
	GlobalData() common.GlobalData
	TxSubmit() common.TxSubmit
	Config() config.Config
	AppKeys() common.AppKeys
	TxOutProducer() TxOutputProducer
	TxTracker() TxTracker
	Keeper() keeper.Keeper
	ValidatorManager() ValidatorManager
	TransferQueue() TransferQueue
	BridgeManager() chains.BridgeManager
	ChainPolling() service.ChainPolling
	PrivateDb() keeper.PrivateDb
	BacgroundService() Background
}

type DefaultManagerContainer struct {
	readOnlyContext atomic.Value

	pmm               PostedMessageManager
	partyManager      PartyManager
	dheartClient      external.DheartClient
	deyesClient       external.DeyesClient
	globalData        common.GlobalData
	txSubmit          common.TxSubmit
	config            config.Config
	appKeys           common.AppKeys
	txOutProducer     TxOutputProducer
	txTracker         TxTracker
	keeper            keeper.Keeper
	valsManager       ValidatorManager
	transferOutQueue  TransferQueue
	bridgeManager     chains.BridgeManager
	chainPolling      service.ChainPolling
	privateDb         keeper.PrivateDb
	backgroundService Background
}

func NewManagerContainer(pmm PostedMessageManager, partyManager PartyManager,
	dheartClient external.DheartClient, deyesClient external.DeyesClient,
	globalData common.GlobalData, txSubmit common.TxSubmit, cfg config.Config,
	appKeys common.AppKeys, txOutProducer TxOutputProducer, txTracker TxTracker,
	keeper keeper.Keeper, valsManager ValidatorManager, txInQueue TransferQueue,
	bridgeManager chains.BridgeManager, chainPolling service.ChainPolling,
	privateDb keeper.PrivateDb, backgroundService Background) ManagerContainer {
	return &DefaultManagerContainer{
		pmm:               pmm,
		partyManager:      partyManager,
		dheartClient:      dheartClient,
		deyesClient:       deyesClient,
		globalData:        globalData,
		txSubmit:          txSubmit,
		config:            cfg,
		appKeys:           appKeys,
		txOutProducer:     txOutProducer,
		txTracker:         txTracker,
		keeper:            keeper,
		valsManager:       valsManager,
		transferOutQueue:  txInQueue,
		bridgeManager:     bridgeManager,
		chainPolling:      chainPolling,
		privateDb:         privateDb,
		backgroundService: backgroundService,
	}
}

func (mc *DefaultManagerContainer) PostedMessageManager() PostedMessageManager {
	return mc.pmm
}

func (mc *DefaultManagerContainer) PartyManager() PartyManager {
	return mc.partyManager
}

func (mc *DefaultManagerContainer) DheartClient() external.DheartClient {
	return mc.dheartClient
}

func (mc *DefaultManagerContainer) GlobalData() common.GlobalData {
	return mc.globalData
}

func (mc *DefaultManagerContainer) TxSubmit() common.TxSubmit {
	return mc.txSubmit
}

func (mc *DefaultManagerContainer) Config() config.Config {
	return mc.config
}

func (mc *DefaultManagerContainer) AppKeys() common.AppKeys {
	return mc.appKeys
}

func (mc *DefaultManagerContainer) TxOutProducer() TxOutputProducer {
	return mc.txOutProducer
}

func (mc *DefaultManagerContainer) DeyesClient() external.DeyesClient {
	return mc.deyesClient
}

func (mc *DefaultManagerContainer) TxTracker() TxTracker {
	return mc.txTracker
}

func (mc *DefaultManagerContainer) Keeper() keeper.Keeper {
	return mc.keeper
}

func (mc *DefaultManagerContainer) ValidatorManager() ValidatorManager {
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

func (mc *DefaultManagerContainer) BacgroundService() Background {
	return mc.backgroundService
}
