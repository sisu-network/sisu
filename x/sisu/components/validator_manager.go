package components

import (
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type ValidatorManager interface {
	AddValidator(ctx sdk.Context, node *types.Node)
	IsValidator(ctx sdk.Context, signer string) bool
	GetValidatorLength(ctx sdk.Context) int
	GetValidators(ctx sdk.Context) []*types.Node
	GetAssignedValidator(ctx sdk.Context, uniqId string) (*types.Node, error)
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
	vals = m.vals

	if vals != nil && len(vals) > 0 {
		return vals
	}

	valsArr := m.keeper.LoadValidators(ctx)
	vals = make(map[string]*types.Node)
	for _, val := range valsArr {
		vals[val.AccAddress] = val
	}

	m.vals = vals

	return vals
}

func (m *DefaultValidatorManager) GetValidatorLength(ctx sdk.Context) int {
	m.valLock.RLock()
	defer m.valLock.RUnlock()

	vals := m.getVals(ctx)

	return len(vals)
}

func (m *DefaultValidatorManager) AddValidator(ctx sdk.Context, node *types.Node) {
	m.valLock.Lock()
	defer m.valLock.Unlock()

	m.keeper.SaveNode(ctx, node)

	vals := m.getVals(ctx)
	vals[node.AccAddress] = node
	m.vals = vals
}

func (m *DefaultValidatorManager) IsValidator(ctx sdk.Context, signer string) bool {
	vals := m.getVals(ctx)

	m.valLock.RLock()
	defer m.valLock.RUnlock()

	return vals[signer] != nil
}

// GetValAccAddrs implements ValidatorManager interface. It returns the list of signer account
// addresses.
func (m *DefaultValidatorManager) GetValidators(ctx sdk.Context) []*types.Node {
	m.valLock.RLock()
	defer m.valLock.RUnlock()

	vals := m.getVals(ctx)
	// Convert map to array
	arr := make([]*types.Node, 0, len(vals))
	for _, value := range vals {
		arr = append(arr, value)
	}

	return arr
}

func (m *DefaultValidatorManager) GetAssignedValidator(ctx sdk.Context, uniqId string) (*types.Node, error) {
	transferId, retryNum := types.GetIdFromUniqId(uniqId)
	threshold := int(m.keeper.GetParams(ctx).GetMaxRejectedTransferRetry())
	if retryNum%threshold == 0 && retryNum > 0 {
		return nil, fmt.Errorf("Exceed the maximum number of retry rejected transfer, id = %s", transferId)
	}

	sorted := utils.GetSortedValidators(uniqId, m.GetValidators(ctx))
	return sorted[0], nil
}
