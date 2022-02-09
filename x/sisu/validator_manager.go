package sisu

import (
	"sync"

	"github.com/ethereum/go-ethereum/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type ValidatorManager interface {
	Init()
	AddValidator(node *types.Node)
	IsValidator(signer string) bool
}

type DefaultValidatorManager struct {
	publicDb keeper.Storage
	vals     map[string]*types.Node
	valLock  *sync.RWMutex
}

func NewValidatorManager(publicDb keeper.Storage) ValidatorManager {
	return &DefaultValidatorManager{
		publicDb: publicDb,
		vals:     make(map[string]*types.Node),
		valLock:  &sync.RWMutex{},
	}
}

func (m *DefaultValidatorManager) Init() {
	// Load all validator onto memory
	valsArr := m.publicDb.LoadValidators()
	vals := make(map[string]*types.Node)
	for _, val := range valsArr {
		vals[val.AccAddress] = val
	}

	log.Info("validators at start = ", vals)

	m.valLock.Lock()
	m.vals = vals
	m.valLock.Unlock()
}

func (m *DefaultValidatorManager) AddValidator(node *types.Node) {
	m.publicDb.SaveNode(node)

	m.valLock.Lock()
	m.vals[node.AccAddress] = node
	m.valLock.Unlock()
}

func (m *DefaultValidatorManager) IsValidator(signer string) bool {
	m.valLock.RLock()
	defer m.valLock.RUnlock()

	return m.vals[signer] != nil
}
