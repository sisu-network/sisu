package sisu

import (
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
)

type ManagerContainer interface {
	PostedMessageManager() PostedMessageManager
	PublicDb() keeper.Storage
	PartyManager() PartyManager
	DheartClient() tssclients.DheartClient
	GlobalData() common.GlobalData
	TxSubmit() common.TxSubmit
	Config() config.TssConfig
	AppKeys() common.AppKeys
}

type DefaultManagerContainer struct {
	publicDb          keeper.Storage
	majorityThreshold int
	partyManager      PartyManager
	dheartClient      tssclients.DheartClient
	globalData        common.GlobalData
	txSubmit          common.TxSubmit
	config            config.TssConfig
	appKeys           common.AppKeys
}

func NewManagerContainer(publicDb keeper.Storage, majorityThreshold int, partyManager PartyManager,
	dheartClient tssclients.DheartClient, globalData common.GlobalData, txSubmit common.TxSubmit,
	cfg config.TssConfig, appKeys common.AppKeys) ManagerContainer {
	return &DefaultManagerContainer{
		publicDb:          publicDb,
		majorityThreshold: majorityThreshold,
		partyManager:      partyManager,
		dheartClient:      dheartClient,
		globalData:        globalData,
		txSubmit:          txSubmit,
		config:            cfg,
		appKeys:           appKeys,
	}
}

func (mc *DefaultManagerContainer) PostedMessageManager() PostedMessageManager {
	return NewPostedMessageManager(mc.publicDb, mc.majorityThreshold)
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
