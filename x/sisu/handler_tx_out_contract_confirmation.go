package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutContractConfirmation struct {
	pmm         PostedMessageManager
	keeper      keeper.Keeper
	deyesClient tssclients.DeyesClient
}

func NewHandlerTxOutContractConfirmation(mc ManagerContainer) *HandlerTxOutContractConfirmation {
	return &HandlerTxOutContractConfirmation{
		keeper:      mc.Keeper(),
		pmm:         mc.PostedMessageManager(),
		deyesClient: mc.DeyesClient(),
	}
}

func (h *HandlerTxOutContractConfirmation) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutContractConfirmWithSigner) (*sdk.Result, error) {
	process, hash, err := h.pmm.PreProcessingMsg(ctx, signerMsg)
	if err != nil {
		return &sdk.Result{}, err
	}

	if !process {
		return &sdk.Result{}, nil
	}

	data, err := h.doTxOutContractConfirm(ctx, signerMsg)
	h.keeper.ProcessTxRecord(ctx, hash)

	return &sdk.Result{Data: data}, err
}

func (h *HandlerTxOutContractConfirmation) doTxOutContractConfirm(ctx sdk.Context, msgWithSigner *types.TxOutContractConfirmWithSigner) ([]byte, error) {
	msg := msgWithSigner.Data
	if h.keeper.IsTxOutConfirmExisted(ctx, msg.OutChain, msg.OutHash) {
		// The message has been processed
		return nil, nil
	}

	log.Info("Delivering TxOutContractConfirm")

	// Save this to keeper and private db
	h.keeper.SaveTxOutConfirm(ctx, msg)

	txOut := h.keeper.GetTxOut(ctx, msg.OutChain, msg.OutHash)
	if txOut == nil {
		log.Critical("cannot find txout from txOutConfirm message, chain & hash = ",
			msg.OutChain, msg.OutHash)
		return nil, nil
	}

	log.Info("txOut.ContractHash = ", txOut.ContractHash)

	// Update the address for the contract.
	contract := h.keeper.GetContract(ctx, txOut.OutChain, txOut.ContractHash, false)
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
	h.keeper.SaveContract(ctx, contract, false)

	// Create a new entry with contract & address as key for easy txOut look up.
	h.keeper.CreateContractAddress(ctx, txOut.OutChain, txOut.OutHash, msg.ContractAddress)

	// Add the address to deyes to watch
	h.deyesClient.AddWatchAddresses(msg.OutChain, []string{msg.ContractAddress})

	return nil, nil
}
