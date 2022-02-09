package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func (p *Processor) deliverTxOutConfirm(ctx sdk.Context, signerMsg *types.TxOutConfirmWithSigner) ([]byte, error) {
	if process, hash := p.shouldProcessMsg(ctx, signerMsg); process {
		p.doTxOutConfirm(ctx, signerMsg)
		p.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

func (p *Processor) doTxOutConfirm(ctx sdk.Context, msgWithSigner *types.TxOutConfirmWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data
	if p.publicDb.IsTxOutConfirmExisted(msg.OutChain, msg.OutHash) {
		// The message has been processed
		return nil, nil
	}

	log.Info("Delivering TxOutConfirm, msg.TxType = ", msg.TxType)

	// Save this to keeper and private db
	p.publicDb.SaveTxOutConfirm(msg)

	// If this is a contract deployment, update the address for the contract.
	if msg.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		txOut := p.publicDb.GetTxOut(msg.OutChain, msg.OutHash)
		if txOut == nil {
			log.Critical("cannot find txout from txOutConfirm message, chain & hash = ",
				msg.OutChain, msg.OutHash)
			return nil, nil
		}

		log.Info("txOut.ContractHash = ", txOut.ContractHash)

		// Update the address for the contract.
		contract := p.publicDb.GetContract(txOut.OutChain, txOut.ContractHash, false)
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
		p.publicDb.SaveContract(contract, false)

		// Create a new entry with contract & address as key for easy txOut look up.
		p.publicDb.CreateContractAddress(txOut.OutChain, txOut.OutHash, msg.ContractAddress)

		// Add the address to deyes to watch
		p.AddWatchAddresses(msg.OutChain, msg.ContractAddress)
	}

	return nil, nil
}

func (p *Processor) AddWatchAddresses(chain, address string) {
	p.deyesClient.AddWatchAddresses(chain, []string{address})
}
