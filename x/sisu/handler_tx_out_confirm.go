package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerTxOutConfirm struct {
	pmm           PostedMessageManager
	keeper        keeper.Keeper
	deyesClient   tssclients.DeyesClient
	transferQueue TransferQueue
}

func NewHandlerTxOutConfirm(mc ManagerContainer) *HandlerTxOutConfirm {
	return &HandlerTxOutConfirm{
		keeper:        mc.Keeper(),
		pmm:           mc.PostedMessageManager(),
		deyesClient:   mc.DeyesClient(),
		transferQueue: mc.TransferQueue(),
	}
}

func (h *HandlerTxOutConfirm) DeliverMsg(ctx sdk.Context, signerMsg *types.TxOutConfirmMsg) (*sdk.Result, error) {
	fmt.Println("There is a txout confirm")
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doTxOutConfirm(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerTxOutConfirm) doTxOutConfirm(ctx sdk.Context, msgWithSigner *types.TxOutConfirmMsg) ([]byte, error) {
	msg := msgWithSigner.Data

	log.Info("Delivering TxOutConfirm")

	txOut := h.keeper.GetTxOut(ctx, msg.OutChain, msg.OutHash)
	if txOut == nil {
		log.Critical("cannot find txout from txOutConfirm message, chain & hash = ",
			msg.OutChain, msg.OutHash)
		return nil, nil
	}

	if txOut.TxType == types.TxOutType_CONTRACT_DEPLOYMENT {
		h.confirmContractDeployment(ctx, txOut, msg.ContractAddress)
	}

	savedCheckPoint := h.keeper.GetGatewayCheckPoint(ctx, msg.OutChain)
	if savedCheckPoint == nil || savedCheckPoint.BlockHeight < msg.BlockHeight {
		// Save checkpoint
		checkPoint := &types.GatewayCheckPoint{
			Chain:       msg.OutChain,
			BlockHeight: msg.BlockHeight,
		}

		if libchain.IsETHBasedChain(msg.OutChain) {
			checkPoint.Nonce = msg.Nonce
		}

		// Update observed block height and nonce.
		h.keeper.AddGatewayCheckPoint(ctx, checkPoint)
	}

	// Clear the pending TxOut
	fmt.Println("Clearing pending out for chain", txOut.OutChain)
	h.keeper.SetPendingTxOut(ctx, txOut.OutChain, nil)

	return nil, nil
}

func (h *HandlerTxOutConfirm) confirmContractDeployment(ctx sdk.Context, txOut *types.TxOut,
	contractAddress string) ([]byte, error) {
	log.Info("txOut.ContractHash = ", txOut.ContractHash)

	// Update the address for the contract.
	contract := h.keeper.GetContract(ctx, txOut.OutChain, txOut.ContractHash, false)
	if contract == nil {
		err := fmt.Errorf("cannot find contract hash with hash %s on chain %s", txOut.ContractHash, txOut.OutChain)
		log.Critical(err)
		return nil, err
	}

	if len(contractAddress) == 0 {
		err := fmt.Errorf("contract address is nil")
		log.Critical(err)
		return nil, err
	}

	contract.Address = contractAddress
	log.Infof("Contract address for chain %s = %s ", contract.Chain, contractAddress)

	// Save the contract (with address)
	h.keeper.SaveContract(ctx, contract, false)

	// Create a new entry with contract & address as key for easy txOut look up.
	h.keeper.CreateContractAddress(ctx, txOut.OutChain, txOut.OutHash, contractAddress)

	// Add the address to deyes to watch
	h.deyesClient.SetGatewayAddress(txOut.OutChain, contractAddress)

	return nil, nil
}
