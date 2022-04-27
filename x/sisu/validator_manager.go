package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const SlashPointThreshold = 100

type ValidatorManager interface {
	AddNode(ctx sdk.Context, node *types.Node)
	UpdateNodeStatus(ctx sdk.Context, accAddress string, consKey []byte, status types.NodeStatus)
	IsValidator(ctx sdk.Context, signer string) bool
	SetValidators(ctx sdk.Context, nodes []*types.Node) error
	GetVals(ctx sdk.Context) map[string]*types.Node
	GetExceedSlashThresholdValidators(ctx sdk.Context) ([]*types.Node, error)
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

func (m *DefaultValidatorManager) GetVals(ctx sdk.Context) map[string]*types.Node {
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

func (m *DefaultValidatorManager) AddNode(ctx sdk.Context, node *types.Node) {
	m.keeper.SaveNode(ctx, node)

	vals := m.GetVals(ctx)

	m.valLock.Lock()
	vals[node.AccAddress] = node
	m.vals = vals
	m.valLock.Unlock()
}

func (m *DefaultValidatorManager) UpdateNodeStatus(ctx sdk.Context, accAddress string, consKey []byte, status types.NodeStatus) {
	m.keeper.UpdateNodeStatus(ctx, consKey, status)
	vals := m.GetVals(ctx)

	m.valLock.RLock()
	node := vals[accAddress]
	m.valLock.RUnlock()

	if node == nil {
		return
	}

	m.valLock.Lock()
	defer m.valLock.Unlock()
	node.Status = status
	if status == types.NodeStatus_Validator {
		node.IsValidator = true
		return
	}
	node.IsValidator = false
}

func (m *DefaultValidatorManager) IsValidator(ctx sdk.Context, signer string) bool {
	vals := m.GetVals(ctx)

	m.valLock.RLock()
	defer m.valLock.RUnlock()

	return vals[signer] != nil
}

func (m *DefaultValidatorManager) SetValidators(ctx sdk.Context, nodes []*types.Node) error {
	validVals, err := m.keeper.SetValidators(ctx, nodes)
	if err != nil {
		return err
	}

	newVals := make(map[string]*types.Node)
	for _, val := range validVals {
		newVals[val.AccAddress] = val
	}

	m.valLock.Lock()
	defer m.valLock.Unlock()

	m.vals = newVals
	return nil
}

// GetExceedSlashThresholdValidators return validators who has too much slash points (exceed threshold)
func (m *DefaultValidatorManager) GetExceedSlashThresholdValidators(ctx sdk.Context) ([]*types.Node, error) {
	slashValidators := make([]*types.Node, 0)
	validators := m.GetVals(ctx)

	for _, validator := range validators {
		slashPoint, err := m.keeper.GetSlashToken(ctx, validator.ConsensusKey.GetBytes())
		if err != nil {
			return nil, err
		}

		if slashPoint <= SlashPointThreshold {
			continue
		}

		slashValidators = append(slashValidators, validator)
	}

	return slashValidators, nil
}
