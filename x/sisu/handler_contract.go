package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerContract struct {
	pmm    PostedMessageManager
	keeper keeper.Keeper
}

func NewHandlerContract(mc ManagerContainer) *HandlerContract {
	return &HandlerContract{
		keeper: mc.Keeper(),
		pmm:    mc.PostedMessageManager(),
	}
}

func (h *HandlerContract) DeliverMsg(ctx sdk.Context, signerMsg *types.ContractsWithSigner) (*sdk.Result, error) {
	if process, hash, err := h.pmm.PreProcessingMsg(ctx, signerMsg); process {
		if err != nil {
			return &sdk.Result{}, err
		}

		data, err := h.doContracts(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerContract) doContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	log.Info("Saving contracts, contracts length = ", len(wrappedMsg.Data.Contracts))

	for _, contract := range wrappedMsg.Data.Contracts {
		if h.keeper.IsContractExisted(ctx, contract) {
			log.Infof("Contract %s has been processed", contract.Name)
			return nil, nil
		}

		log.Info("Saving contarct ", contract.Name, " on chain ", contract.Chain)
	}

	// Save into KVStore & private db
	h.keeper.SaveContracts(ctx, wrappedMsg.Data.Contracts, true)

	return nil, nil
}
