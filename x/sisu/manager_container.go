package sisu

import (
	"sync/atomic"

	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

type ManagerContainer interface {
	PostedMessageManager() PostedMessageManager
	PartyManager() PartyManager
	DheartClient() tssclients.DheartClient
	DeyesClient() tssclients.DeyesClient
	GlobalData() common.GlobalData
	TxSubmit() common.TxSubmit
	Config() config.TssConfig
	AppKeys() common.AppKeys
	TxOutProducer() TxOutputProducer
	TxTracker() TxTracker
	Keeper() keeper.Keeper
	ValidatorManager() ValidatorManager
	TransferQueue() TransferQueue
}

type DefaultManagerContainer struct {
	readOnlyContext atomic.Value

	pmm           PostedMessageManager
	partyManager  PartyManager
	dheartClient  tssclients.DheartClient
	deyesClient   tssclients.DeyesClient
	globalData    common.GlobalData
	txSubmit      common.TxSubmit
	config        config.TssConfig
	appKeys       common.AppKeys
	txOutProducer TxOutputProducer
	txTracker     TxTracker
	keeper        keeper.Keeper
	valsManager   ValidatorManager
	txInQueue     TransferQueue
}

func NewManagerContainer(pmm PostedMessageManager, partyManager PartyManager,
	dheartClient tssclients.DheartClient, deyesClient tssclients.DeyesClient,
	globalData common.GlobalData, txSubmit common.TxSubmit, cfg config.TssConfig,
	appKeys common.AppKeys, txOutProducer TxOutputProducer, txTracker TxTracker,
	keeper keeper.Keeper, valsManager ValidatorManager, txInQueue TransferQueue) ManagerContainer {
	return &DefaultManagerContainer{
		pmm:           pmm,
		partyManager:  partyManager,
		dheartClient:  dheartClient,
		deyesClient:   deyesClient,
		globalData:    globalData,
		txSubmit:      txSubmit,
		config:        cfg,
		appKeys:       appKeys,
		txOutProducer: txOutProducer,
		txTracker:     txTracker,
		keeper:        keeper,
		valsManager:   valsManager,
		txInQueue:     txInQueue,
	}
}

func (mc *DefaultManagerContainer) PostedMessageManager() PostedMessageManager {
	return mc.pmm
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
	return mc.txInQueue
}
