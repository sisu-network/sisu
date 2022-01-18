package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

func (p *Processor) checkTxOutConfirm(ctx sdk.Context, msgWithSigner *types.TxOutConfirmWithSigner) error {
	msg := msgWithSigner.Data
	if !p.privateDb.IsTxOutConfirmExisted(msg.OutChain, msg.OutHash) {
		fmt.Println("Cannot find the tx in our db")
		return ErrCannotFindMessage
	}

	if p.keeper.IsTxOutConfirmExisted(ctx, msg.OutChain, msg.OutHash) {
		fmt.Println("Tx has been included into the keeper")
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

	log.Info("Delivering TxOutConfirm, msg.TxType = ", msg.TxType)

	// Save this to keeper and private db
	p.keeper.SaveTxOutConfirm(ctx, msg)
	p.privateDb.SaveTxOutConfirm(msg)

	// If this is a contract deployment, update the address for the contract.
	if msg.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		txOut := p.keeper.GetTxOut(ctx, msg.OutChain, msg.OutHash)
		if txOut == nil {
			log.Critical("cannot find txout from txOutConfirm message, chain & hash = ",
				msg.OutChain, msg.OutHash)
			return nil, nil
		}

		log.Info("txOut.ContractHash = ", txOut.ContractHash)

		// Update the address for the contract.
		contract := p.keeper.GetContract(ctx, txOut.OutChain, txOut.ContractHash, false)
		if contract == nil {
			log.Critical("cannot find contract hash with hash ", txOut.ContractHash, " on chain ", txOut.OutChain)
			return nil, nil
		}

		if len(msg.ContractAddress) == 0 {
			log.Critical("contract address is nil")
			return nil, nil
		}

		contract.Address = msg.ContractAddress
		log.Infof("Contract address for chain %s = %s ", contract.Chain, msg.ContractAddress)

		// Save the contract (with address)
		p.keeper.SaveContract(ctx, contract, false)
		p.privateDb.SaveContract(contract, false)

		// Create a new entry with contract & address as key for easy txOut look up.
		p.keeper.CreateContractAddress(ctx, txOut.OutChain, txOut.OutHash, msg.ContractAddress)
		p.privateDb.CreateContractAddress(txOut.OutChain, txOut.OutHash, msg.ContractAddress)

		// Add the address to deyes to watch
		p.AddWatchAddresses(msg.OutChain, msg.ContractAddress)
	}

	return nil, nil
}

func (p *Processor) AddWatchAddresses(chain, address string) {
}
