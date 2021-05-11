package ethchain

import (
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/core/state"
	"github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/core/vm"
	"github.com/sisu-network/dcore/params"
	"github.com/sisu-network/sisu/utils"
)

// SoftState is a temporary state db that is not connected to the chain. It is only used for
// validating transactions in a block and will be discarded after block seals.
type SoftState struct {
	db           *state.StateDB
	block        *types.Block
	chainConfig  *params.ChainConfig
	vmConfig     vm.Config
	chainContext core.ChainContext
	gasPool      *core.GasPool
	index        int
}

func NewSoftState(db *state.StateDB, block *types.Block, chainConfig *params.ChainConfig,
	cfg vm.Config, chainContext core.ChainContext) *SoftState {
	return &SoftState{
		db:           db,
		block:        block,
		chainConfig:  chainConfig,
		vmConfig:     cfg,
		chainContext: chainContext,
		// TODO: fix this
		// gasPool:  new(core.GasPool).AddGas(block.GasLimit()),
		gasPool: new(core.GasPool).AddGas(10000000000000),
		index:   0,
	}
}

func (s *SoftState) ApplyTx(tx *types.Transaction) (*types.Receipt, error) {
	s.db.Prepare(tx.Hash(), s.block.Hash(), s.index)

	usedGas := new(uint64)

	utils.LogDebug("Applying ETH TX...")

	receipt, err := core.ApplyTransaction(s.chainConfig, s.chainContext, nil, s.gasPool, s.db,
		s.block.Header(), tx, usedGas, s.vmConfig)

	if err != nil {
		utils.LogDebug(err)
	}

	utils.LogDebug("Done Applying ETH TX....")

	s.index++

	return receipt, err
}
