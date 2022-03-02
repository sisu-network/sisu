package sisu

import (
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/world"
)

type ManagerContainer interface {
	PostedMessageManager() PostedMessageManager
	PublicDb() keeper.Storage
	PartyManager() PartyManager
	DheartClient() tssclients.DheartClient
	DeyesClient() tssclients.DeyesClient
	GlobalData() common.GlobalData
	TxSubmit() common.TxSubmit
	Config() config.TssConfig
	AppKeys() common.AppKeys
	TxOutProducer() TxOutputProducer
	WorldState() world.WorldState
	TransactionTracker() TransactionTracker
}

type DefaultManagerContainer struct {
	pmm               PostedMessageManager
	publicDb          keeper.Storage
	majorityThreshold int
	partyManager      PartyManager
	dheartClient      tssclients.DheartClient
	deyesClient       tssclients.DeyesClient
	globalData        common.GlobalData
	txSubmit          common.TxSubmit
	config            config.TssConfig
	appKeys           common.AppKeys
	txOutProducer     TxOutputProducer
	worldState        world.WorldState
	txTracker         TransactionTracker
}

func NewManagerContainer(pmm PostedMessageManager, publicDb keeper.Storage, majorityThreshold int, partyManager PartyManager,
	dheartClient tssclients.DheartClient, deyesClient tssclients.DeyesClient,
	globalData common.GlobalData, txSubmit common.TxSubmit, cfg config.TssConfig,
	appKeys common.AppKeys, txOutProducer TxOutputProducer, worldState world.WorldState, txTracker TransactionTracker) ManagerContainer {
	return &DefaultManagerContainer{
		pmm:               pmm,
		publicDb:          publicDb,
		majorityThreshold: majorityThreshold,
		partyManager:      partyManager,
		dheartClient:      dheartClient,
		deyesClient:       deyesClient,
		globalData:        globalData,
		txSubmit:          txSubmit,
		config:            cfg,
		appKeys:           appKeys,
		txOutProducer:     txOutProducer,
		worldState:        worldState,
		txTracker:         txTracker,
	}
}

func (mc *DefaultManagerContainer) PostedMessageManager() PostedMessageManager {
	return mc.pmm
}

func (mc *DefaultManagerContainer) PublicDb() keeper.Storage {
	return mc.publicDb
}

func (mc *DefaultManagerContainer) PartyManager() PartyManager {
	return mc.partyManager
}

func (mc *DefaultManagerContainer) DheartClient() tssclients.DheartClient {
	return mc.dheartClient
}

func (mc *DefaultManagerContainer) GlobalData() common.GlobalData {
	return mc.globalData
}

func (mc *DefaultManagerContainer) TxSubmit() common.TxSubmit {
	return mc.txSubmit
}

func (mc *DefaultManagerContainer) Config() config.TssConfig {
	return mc.config
}

func (mc *DefaultManagerContainer) AppKeys() common.AppKeys {
	return mc.appKeys
}

func (mc *DefaultManagerContainer) TxOutProducer() TxOutputProducer {
	return mc.txOutProducer
}

func (mc *DefaultManagerContainer) DeyesClient() tssclients.DeyesClient {
	return mc.deyesClient
}

func (mc *DefaultManagerContainer) WorldState() world.WorldState {
	return mc.worldState
}

func (mc *DefaultManagerContainer) TransactionTracker() TransactionTracker {
	return mc.txTracker
}
