package tss

import (
	"bytes"
	"encoding/hex"
	"fmt"

	tTypes "github.com/sisu-network/dheart/types"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/contracts/eth/dummy"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

// Produces response for an observed tx. This has to be deterministic based on all the data that
// the processor has.
func (p *Processor) CreateTxOuts(ctx sdk.Context, tx *types.ObservedTx) {
	var outMsgs []*tssTypes.TxOut
	var err error

	switch tx.Chain {
	case "eth":
		outMsgs, err = p.getEthResponse(ctx, tx)

		if err != nil {
			utils.LogError("Cannot get response for an eth tx")
		}
	}

	validators := p.globalData.GetValidatorSet()
	for _, msg := range outMsgs {
		// Find a validator that the network expects to post out tx in the next 1 or 2 blocks. If
		// the assigned validator does not post, everyone in the network will broadcast the out tx
		// and the assigned validator has minor slashing.

		// TODO: Use online/active validators instead the whole validator sets.
		valAddr := p.getAssignedValidator(p.currentHeight, string(msg.OutBytes), validators)

		fmt.Println("p.currentHeight = ", p.currentHeight)

		// Save this valAddr for later block check.
		p.storage.AddPendingTxOut(
			p.currentHeight,
			tx.Chain,
			tx.TxHash,
			msg.OutChain,
			msg.OutBytes,
			valAddr,
		)
	}
}

// Get ETH out from an observed tx. Only do this if this is a validator node.
func (p *Processor) getEthResponse(ctx sdk.Context, tx *types.ObservedTx) ([]*tssTypes.TxOut, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(tx.Serialized)
	if err != nil {
		utils.LogError("Failed to unmarshall eth tx. err =", err)
		return nil, err
	}

	outMsgs := make([]*tssTypes.TxOut, 0)
	// Process different kind of eth transaction.
	// 1. Check if the To address of our public key. This is likely a tx to provide ETH for our
	// account to deploy contracts. Check if we have some pending contracts and deploy if needed.
	if ethTx.To().String() == p.keyAddress {
		// TODO: Check balance required to deploy all these contracts.
		// Get all contract in the pending queue.
		contracts := p.keeper.GetContractQueueHashes(ctx, tx.Chain)

		if len(contracts) > 0 {
			// Get the list of deploy transactions. Those txs need to posted and verified (by validators)
			// to the Sisu chain
			outEthTxs := p.checkEthDeployContract(ctx, tx.Chain, ethTx, contracts)

			for _, outTx := range outEthTxs {
				bz, err := outTx.MarshalBinary()
				if err != nil {
					utils.LogError("Cannot marshall binary")
					continue
				}

				fmt.Println("Adding to txBytes")

				outMsgs = append(outMsgs, tssTypes.NewMsgTxOut(
					p.appKeys.GetSignerAddress().String(),
					tx.BlockHeight,
					tx.Chain,
					tx.TxHash,
					tx.Chain,
					bz,
				))
			}
		}
	}

	return outMsgs, nil
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

// Processes all txs at the end of a block that are added in the current block.
func (p *Processor) processPendingTxs(ctx sdk.Context) {
	observedTxList := p.storage.GetAllPendingTxs()
	// 1. Creates all tx out for all observed txs in this blocks.
	for _, tx := range observedTxList {
		p.CreateTxOuts(ctx, tx)
	}
	p.storage.ClearPendingTxs()

	// 2. Broadcast all txs out that have been assigned to this node.
	p.broadcastAssignedTxOuts()

	// 3. Check if there is txs out that is supposed to be included in this block or the previous
	// block. If there is, broadcast that tx out.
}

// Broadcasts all txouts that have been assigned to this validator.
func (p *Processor) broadcastAssignedTxOuts() {
	myValidatorAddr := p.globalData.GetMyTendermintValidatorAddr()

	txWrappers := p.storage.GetPendingTxOutForValidator(p.currentHeight, myValidatorAddr)

	for _, tx := range txWrappers {
		go func(tx *PendingTxOutWrapper) {
			p.txSubmit.SubmitMessage(
				tssTypes.NewMsgTxOut(
					p.appKeys.GetSignerAddress().String(),
					tx.InBlockHeight,
					tx.InChain,
					tx.InHash,
					tx.OutChain,
					tx.OutBytes,
				),
			)
		}(tx)
	}
}

func (p *Processor) CheckTxOut(ctx sdk.Context, msg *types.TxOut) error {
	fmt.Println("Checking Txout...")

	txWrapper := p.storage.GetPendingTxOUt(msg.InBlockHeight, msg.InHash)
	if txWrapper == nil {
		utils.LogError("Cannot find txWrapper", msg.InBlockHeight, msg.InHash)
		return fmt.Errorf("Transaction not found")
	}

	if bytes.Compare(txWrapper.OutBytes, msg.OutBytes) != 0 {
		utils.LogError("Txouts do not match.")
		return fmt.Errorf("OutBytes do not match")
	}

	fmt.Println("Txout is good")

	return nil
}

func (p *Processor) DeliverTxOut(ctx sdk.Context, msg *types.TxOut) ([]byte, error) {
	utils.LogVerbose("Delivering TXOUT")

	outHash, err := utils.GetTxHash(msg.OutChain, msg.OutBytes)
	if err != nil {
		utils.LogCritical("Cannot get tx hash for tx with serialized data: ", hex.EncodeToString(msg.OutBytes), "err = ", err)
		return nil, err
	}

	// TODO: bring this logic back after debugging.
	// if p.keeper.IsPendingKeygenTxExisted(ctx, msg.OutChain, p.currentHeight, outHash) {
	if false {
		// This transaction has been processed and keysigned. No need to do keysign again.
		return nil, nil
	}

	// 1. Remove the tx out from the storage.
	p.storage.RemovePendingTxOut(msg.InBlockHeight, msg.InHash)

	// 2. Mark the tx as processed and save it to KVStore.
	p.keeper.AddProcessedTx(ctx, msg)

	// 3. Add it to a queue to do keygen.
	p.keeper.AddPendingKeygenTx(ctx, msg.OutChain, p.currentHeight, outHash)

	// 4. Broadcast it to Dheart for processing.
	err = p.dheartClient.KeySign(&tTypes.KeysignRequest{
		OutChain:       msg.OutChain,
		OutBlockHeight: p.currentHeight,
		OutHash:        outHash,
		OutBytes:       msg.OutBytes,
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}
