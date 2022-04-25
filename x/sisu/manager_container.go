package sisu

import (
	"fmt"
	"sync/atomic"

	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/world"
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
	WorldState() world.WorldState
	TxTracker() TxTracker
	Keeper() keeper.Keeper
	BankKeeper() bankkeeper.Keeper

	SetReadOnlyContext(ctx sdk.Context)
	GetReadOnlyContext() (sdk.Context, error)
}

type DefaultManagerContainer struct {
	pmm           PostedMessageManager
	partyManager  PartyManager
	dheartClient  tssclients.DheartClient
	deyesClient   tssclients.DeyesClient
	globalData    common.GlobalData
	txSubmit      common.TxSubmit
	config        config.TssConfig
	appKeys       common.AppKeys
	txOutProducer TxOutputProducer
	worldState    world.WorldState
	txTracker     TxTracker
	keeper        keeper.Keeper
	bankKeeper    bankkeeper.Keeper

	readOnlyContext atomic.Value
}

func NewManagerContainer(pmm PostedMessageManager, partyManager PartyManager,
	dheartClient tssclients.DheartClient, deyesClient tssclients.DeyesClient,
	globalData common.GlobalData, txSubmit common.TxSubmit, cfg config.TssConfig,
	appKeys common.AppKeys, txOutProducer TxOutputProducer, worldState world.WorldState, txTracker TxTracker, keeper keeper.Keeper, bankKeeper bankkeeper.Keeper) ManagerContainer {
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
		worldState:    worldState,
		txTracker:     txTracker,
		keeper:        keeper,
		bankKeeper:    bankKeeper,
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

func (mc *DefaultManagerContainer) WorldState() world.WorldState {
	return mc.worldState
}

func (mc *DefaultManagerContainer) TxTracker() TxTracker {
	return mc.txTracker
}

func (mc *DefaultManagerContainer) Keeper() keeper.Keeper {
	return mc.keeper
}

func (mc *DefaultManagerContainer) BankKeeper() bankkeeper.Keeper {
	return mc.bankKeeper
}

func (mc *DefaultManagerContainer) SetReadOnlyContext(ctx sdk.Context) {
	mc.readOnlyContext.Store(ctx)
}

func (mc *DefaultManagerContainer) GetReadOnlyContext() (sdk.Context, error) {
	val := mc.readOnlyContext.Load()
	if val == nil {
		return sdk.Context{}, fmt.Errorf("Read only context is not set")
	}

	return val.(sdk.Context), nil
}
