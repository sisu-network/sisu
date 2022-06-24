package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type ValidatorManager interface {
	AddValidator(ctx sdk.Context, node *types.Node)
	IsValidator(ctx sdk.Context, signer string) bool
	GetValidatorLength(ctx sdk.Context) int
}

type DefaultValidatorManager struct {
	keeper  keeper.Keeper
	vals    map[string]*types.Node
	valLock *sync.RWMutex
}

func NewValidatorManager(keeper keeper.Keeper) ValidatorManager {
	return &DefaultValidatorManager{
		keeper:  keeper,
		vals:    make(map[string]*types.Node),
		valLock: &sync.RWMutex{},
	}
}

func (m *DefaultValidatorManager) getVals(ctx sdk.Context) map[string]*types.Node {
	var vals map[string]*types.Node
	m.valLock.RLock()
	vals = m.vals
	m.valLock.RUnlock()

	if vals != nil && len(vals) > 0 {
		return vals
	}

	valsArr := m.keeper.LoadValidators(ctx)
	vals = make(map[string]*types.Node)
	for _, val := range valsArr {
		vals[val.AccAddress] = val
	}

	m.valLock.Lock()
	m.vals = vals
	m.valLock.Unlock()

	return vals
}

func (m *DefaultValidatorManager) GetValidatorLength(ctx sdk.Context) int {
	vals := m.getVals(ctx)
	m.valLock.RLock()
	defer m.valLock.RUnlock()

	return len(vals)
}

func (m *DefaultValidatorManager) AddValidator(ctx sdk.Context, node *types.Node) {
	m.keeper.SaveNode(ctx, node)

	vals := m.getVals(ctx)

	m.valLock.Lock()
	vals[node.AccAddress] = node
	m.vals = vals
	m.valLock.Unlock()
}

func (m *DefaultValidatorManager) IsValidator(ctx sdk.Context, signer string) bool {
	vals := m.getVals(ctx)

	m.valLock.RLock()
	defer m.valLock.RUnlock()

	return vals[signer] != nil
}
