package tss

import (
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) checkTxOutConfirm(ctx sdk.Context, msgWithSigner *types.TxOutConfirmWithSigner) error {
	msg := msgWithSigner.Data
	if !p.privateDb.IsTxOutConfirmExisted(msg.OutChain, msg.OutHash) {
		return ErrCannotFindMessage
	}

	if p.keeper.IsTxOutConfirmExisted(ctx, msg.OutChain, msg.OutHash) {
		return ErrMessageHasBeenProcessed
	}

	return nil
}

func (p *Processor) deliverTxOutConfirm(ctx sdk.Context, msgWithSigner *types.TxOutConfirmWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data
	if p.keeper.IsTxOutConfirmExisted(ctx, msg.OutChain, msg.OutHash) {
		// The message has been processed
		return nil, nil
	}

	// Save this to keeper and private db
	p.keeper.SaveTxOutConfirm(ctx, msg)
	p.privateDb.SaveTxOutConfirm(msg)

	// If this is a contract deployment, update the address for the contract.
	if msg.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		txOut := p.keeper.GetTxOut(ctx, msg.OutChain, msg.OutHash)

		// Update the address for the contract.
		contract := p.keeper.GetContract(ctx, txOut.OutChain, txOut.ContractHash, false)
		contract.Address = msg.ContractAddress

		p.keeper.SaveContract(ctx, contract, false)
		p.privateDb.SaveContract(contract, false)

		// Create a new entry with contract & address as key for easy look up.
		p.keeper.CreateContractAddress(ctx, txOut.OutChain, txOut.GetHash(), msg.ContractAddress)
		p.privateDb.CreateContractAddress(txOut.OutChain, txOut.GetHash(), msg.ContractAddress)

		// Add the address to deyes to watch
		p.AddWatchAddresses(msg.OutChain, msg.ContractAddress)
	}

	return nil, nil
}

func (p *Processor) AddWatchAddresses(chain, address string) {
	deyeClient := p.deyesClients[chain]
	deyeClient.AddWatchAddresses(chain, []string{address})
}
