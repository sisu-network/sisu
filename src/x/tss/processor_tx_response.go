package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/contracts/eth/dummy"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) GetObserveTxResponse(ctx sdk.Context, tx *types.ObservedTx) {
	var txBytes [][]byte
	var err error

	switch tx.Chain {
	case "eth":
		txBytes, err = p.getEthResponse(ctx, tx)

		if err != nil {
			utils.LogError("Cannot get response for an eth tx")
		}
	}

	validators := p.globalData.GetValidatorSet()
	for _, bz := range txBytes {
		// Find a validator that the network expects to post out tx in the next 1 or 2 blocks. If
		// the assigned validator does not post, everyone in the network will broadcast the out tx
		// and the assigned validator has minor slashing.

		// TODO: Use online/active validators instead the whole validator sets.
		valAddr := p.getAssignedValidator(p.currentHeight, string(bz), validators)

		// Save this valAddr for later block check.
		p.keeper.AddAssignedValForOutTx(p.currentHeight, bz, valAddr)
	}
}

// Get ETH response from an observed tx. Only do this if this is a validator node.
func (p *Processor) getEthResponse(ctx sdk.Context, tx *types.ObservedTx) ([][]byte, error) {
	bz := p.storage.GetObservedTx(tx.Chain, tx.BlockHeight, tx.TxHash)
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(bz)
	if err != nil {
		utils.LogError("Failed to unmarshall eth tx. err =", err)
		return nil, err
	}

	txBytes := make([][]byte, 0)

	// Process different kind of eth transaction.
	// 1. Check if the To address of our public key. This is likely a tx to provide ETH for our
	// account to deploy contracts. Check if we have some pending contracts and deploy if needed.
	if ethTx.To().String() == p.keyAddress {
		// Get all contract in the pending queue.
		contractHashes := p.keeper.GetContractQueueHashes(ctx, tx.Chain)
		if len(contractHashes) > 0 {
			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain
			outEthTxs := p.checkEthDeployContract(ctx, tx.Chain, ethTx, contractHashes)

			for _, tx := range outEthTxs {
				bz, err := tx.MarshalBinary()
				if err != nil {
					utils.LogError("Cannot marshall binary")
					continue
				}

				txBytes = append(txBytes, bz)
			}
		}
	}

	return txBytes, nil
}

// Get one validator from the validator list based on blockHeight and a hash. This is one way to
// get "random" validator in the deterministic world of crypto.
func (p *Processor) getAssignedValidator(blockHeight int64, hash string, validators []*common.Validator) string {
	index := utils.GetRandomIndex(blockHeight, hash, len(validators))
	return validators[index].Address
}

// Check if we can deploy contract after seeing some ETH being sent to our ethereum address.
func (p *Processor) checkEthDeployContract(ctx sdk.Context, chain string, ethTx *ethTypes.Transaction,
	hashes []string) []*ethTypes.Transaction {
	txs := make([]*ethTypes.Transaction, 0)

	nonce := int64(0)
	for _, hash := range hashes {
		switch hash {
		case dummy.DummyABI:
			rawTx := p.logic.PrepareEthContractDeployment(chain, nonce)
			txs = append(txs, rawTx)
			nonce++

			// Save it to the deploying list.
			bz, err := rawTx.MarshalBinary()
			if err == nil {
				// Delete all the contracts in the pending queue and move them to deploying set.
				p.keeper.DequeueContract(ctx, chain, hash)
				p.keeper.AddDeployingContract(ctx, chain, hash, bz, p.currentHeight)
			}
		}
	}

	return txs
}
