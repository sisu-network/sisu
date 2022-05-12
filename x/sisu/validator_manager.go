package sisu

import (
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

const SlashPointThreshold = 100

type ValidatorManager interface {
	AddNode(ctx sdk.Context, node *types.Node)
	UpdateNodeStatus(ctx sdk.Context, consensusKey []byte, status types.NodeStatus)
	SetValidators(ctx sdk.Context, nodes []*types.Node) error
	GetNodesByStatus(status types.NodeStatus) map[string]*types.Node
	GetExceedSlashThresholdValidators(ctx sdk.Context) ([]*types.Node, error)
	GetPotentialCandidates(ctx sdk.Context, n int) []*types.Node
}

type DefaultValidatorManager struct {
	keeper keeper.Keeper
	// key: consensus public key
	nodes    map[string]*types.Node
	nodeLock *sync.RWMutex
}

func NewValidatorManager(keeper keeper.Keeper) ValidatorManager {
	return &DefaultValidatorManager{
		keeper:   keeper,
		nodes:    make(map[string]*types.Node),
		nodeLock: &sync.RWMutex{},
	}
}

// GetNodesByStatus returns all nodes whose status is either unkonwn or equal to a specific status.
func (m *DefaultValidatorManager) GetNodesByStatus(status types.NodeStatus) map[string]*types.Node {
	filteredNodes := make(map[string]*types.Node)
	m.nodeLock.RLock()
	for key, node := range m.nodes {
		if status != types.NodeStatus_Unknown && node.Status != status {
			continue
		}

		filteredNodes[key] = node
	}
	m.nodeLock.RUnlock()

	return filteredNodes
}

func (m *DefaultValidatorManager) AddNode(ctx sdk.Context, node *types.Node) {
	m.keeper.SaveNode(ctx, node)

	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()

	m.nodes[getNodeKey(node)] = node
}

func (m *DefaultValidatorManager) UpdateNodeStatus(ctx sdk.Context, consensusKey []byte, status types.NodeStatus) {
	vals := m.GetNodesByStatus(types.NodeStatus_Unknown)

	m.nodeLock.RLock()
	node := vals[string(consensusKey)]
	m.nodeLock.RUnlock()
	if node == nil {
		return
	}

	// Update keeper
	m.keeper.UpdateNodeStatus(ctx, node.ConsensusKey.GetBytes(), status)

	// Update cache
	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()
	node.Status = status
	if status == types.NodeStatus_Validator {
		node.IsValidator = true
	} else {
		node.IsValidator = false
	}
	vals[string(consensusKey)] = node
}

func (m *DefaultValidatorManager) SetValidators(ctx sdk.Context, nodes []*types.Node) error {
	validVals, err := m.keeper.SetValidators(ctx, nodes)
	if err != nil {
		return err
	}

	newVals := make(map[string]*types.Node)
	for _, val := range validVals {
		newVals[string(val.ConsensusKey.GetBytes())] = val
	}

	m.nodeLock.Lock()
	defer m.nodeLock.Unlock()

	m.nodes = newVals
	return nil
}

// GetExceedSlashThresholdValidators return validators who has too much slash points (exceed threshold)
func (m *DefaultValidatorManager) GetExceedSlashThresholdValidators(ctx sdk.Context) ([]*types.Node, error) {
	slashValidators := make([]*types.Node, 0)
	validators := m.keeper.LoadNodesByStatus(ctx, types.NodeStatus_Validator)

	for _, validator := range validators {
		addr, err := sdk.AccAddressFromBech32(validator.AccAddress)
		if err != nil {
			log.Error("error when parsing addr. error = ", err)
			return nil, err
		}

		slashPoint, err := m.keeper.GetSlashToken(ctx, addr)
		if err != nil {
			return nil, err
		}

		log.Debugf("slash point of address %s is %d", addr, slashPoint)

		if slashPoint <= SlashPointThreshold {
			continue
		}

		slashValidators = append(slashValidators, validator)
	}

	return slashValidators, nil
}

func (m *DefaultValidatorManager) GetPotentialCandidates(ctx sdk.Context, n int) []*types.Node {
	topBalances := m.keeper.GetTopBondBalance(ctx, -1)
	if len(topBalances) == 0 {
		return []*types.Node{}
	}

	newNodes := make([]*types.Node, 0)

	for _, accAddr := range topBalances {
		// enough node already
		if len(newNodes) == n {
			return newNodes
		}

		if nodeFromAddr := m.getNodeByAccAddr(accAddr, types.NodeStatus_Candidate); nodeFromAddr != nil {
			newNodes = append(newNodes, nodeFromAddr)
		}
	}

	return newNodes
}

func (m *DefaultValidatorManager) getNodeByAccAddr(addr sdk.AccAddress, status types.NodeStatus) *types.Node {
	nodes := m.GetNodesByStatus(status)
	for _, node := range nodes {
		if node.AccAddress != addr.String() {
			continue
		}

		return node
	}

	return nil
}

func getNodeKey(node *types.Node) string {
	return string(node.ConsensusKey.GetBytes())
}
