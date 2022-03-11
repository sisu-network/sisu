package sisu

import (
	"fmt"

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
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		h.doContracts(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)
	}

	return nil, nil
}

func (h *HandlerContract) doContracts(ctx sdk.Context, wrappedMsg *types.ContractsWithSigner) ([]byte, error) {
	// TODO: Don't do duplicated delivery
	log.Info("Deliver pending contracts")

	for _, contract := range wrappedMsg.Data.Contracts {
		if h.keeper.IsContractExisted(ctx, contract) {
			log.Infof("Contract %s has been processed", contract.Name)
			return nil, nil
		}
	}

	log.Info("Saving contracts, contracts length = ", len(wrappedMsg.Data.Contracts))
	for _, contract := range wrappedMsg.Data.Contracts {
		fmt.Println("AAA contract hash = ", contract.Hash)
	}

	// Save into KVStore & private db
	h.keeper.SaveContracts(ctx, wrappedMsg.Data.Contracts, true)

	for _, contract := range wrappedMsg.Data.Contracts {
		c := h.keeper.GetContract(ctx, contract.Chain, contract.Hash, false)
		if c != nil {
			fmt.Println("Contract is found", contract.Chain, contract.Hash)
		} else {
			fmt.Println("Contract is NOT found")
		}
	}

	return nil, nil
}
