package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerContract struct {
	pmm      PostedMessageManager
	publicDb keeper.Storage
}

func NewHandlerContract(mc ManagerContainer) *HandlerContract {
	return &HandlerContract{
		publicDb: mc.PublicDb(),
		pmm:      mc.PostedMessageManager(),
	}
}

func (h *HandlerContract) DeliverMsg(ctx sdk.Context, signerMsg *types.ContractsWithSigner) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		h.doContracts(ctx, signerMsg)
		h.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil
}

func (h *HandlerContract) doContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	// TODO: Don't do duplicated delivery
	log.Info("Deliver pending contracts")

	for _, contract := range wrappedMsg.Data.Contracts {
		if h.publicDb.IsContractExisted(contract) {
			log.Infof("Contract %s has been processed", contract.Name)
			return nil, nil
		}
	}

	log.Info("Saving contracts, contracts length = ", len(wrappedMsg.Data.Contracts))

	// Save into KVStore & private db
	h.publicDb.SaveContracts(wrappedMsg.Data.Contracts, true)

	return nil, nil
}
