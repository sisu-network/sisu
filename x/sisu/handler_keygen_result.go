package sisu

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeygenResult struct {
	mc ManagerContainer

	publicDb    keeper.Storage
	pmm         PostedMessageManager
	globalData  common.GlobalData
	deyesClient tssclients.DeyesClient
	config      config.TssConfig
	txSubmit    common.TxSubmit
	appKeys     common.AppKeys
}

func NewHandlerKeygenResult(mc ManagerContainer) *HandlerKeygenResult {
	return &HandlerKeygenResult{
		mc:         mc,
		publicDb:   mc.PublicDb(),
		pmm:        mc.PostedMessageManager(),
		globalData: mc.GlobalData(),
		config:     mc.Config(),
		txSubmit:   mc.TxSubmit(),
		appKeys:    mc.AppKeys(),
	}
}

func (h *HandlerKeygenResult) DeliverMsg(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) (*sdk.Result, error) {
	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		h.doKeygenResult(ctx, signerMsg)
		h.publicDb.ProcessTxRecord(hash)
	}

	return nil, nil

}

func (h *HandlerKeygenResult) doKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) ([]byte, error) {
	msg := signerMsg.Data

	log.Info("Delivering keygen result, result = ", msg.Result)

	result := h.getKeygenResult(ctx, signerMsg)

	// TODO: Get majority of the votes here.
	if result == types.KeygenResult_SUCCESS {
		log.Info("Keygen succeeded")

		// Save result to KVStore & private db
		h.publicDb.SaveKeygenResult(signerMsg)

		if !h.globalData.IsCatchingUp() {
			h.createContracts(ctx, signerMsg.Keygen)
		}
	} else {
		// TODO: handle failure case
	}

	return nil, nil
}

func (h *HandlerKeygenResult) getKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) types.KeygenResult_Result {
	results := h.publicDb.GetAllKeygenResult(signerMsg.Keygen.KeyType, signerMsg.Keygen.Index)

	// Check the majority of the results
	successCount := 0
	for _, result := range results {
		if result.Data.Result == types.KeygenResult_SUCCESS {
			successCount += 1
		}
	}

	if successCount >= (len(results)+1)/2 {
		// TODO: Choose the address with most vote.
		// Save keygen Address
		log.Info("Saving keygen...")
		h.publicDb.SaveKeygen(signerMsg.Keygen)

		return types.KeygenResult_SUCCESS
	}

	return types.KeygenResult_FAILURE
}

// createContracts creates and broadcast pending contracts. All nodes need to agree what
// contracts to deploy on what chains.
func (h *HandlerKeygenResult) createContracts(ctx sdk.Context, msg *types.Keygen) {
	log.Info("Create and broadcast contracts...")

	// We want the final contracts array to be deterministic. We need to sort the list of chains
	// and list of contracts alphabetically.
	// Sort all chains alphabetically.
	chains := make([]string, len(h.config.SupportedChains))
	i := 0
	for chain := range h.config.SupportedChains {
		chains[i] = chain
		i += 1
	}
	sort.Strings(chains)

	// Sort all contracts name alphabetically
	names := make([]string, len(SupportedContracts))
	i = 0
	for contract := range SupportedContracts {
		names[i] = contract
		i += 1
	}
	sort.Strings(names)

	// Create contracts
	contracts := make([]*types.Contract, 0)
	for _, chain := range chains {
		if libchain.GetKeyTypeForChain(chain) == msg.KeyType {
			log.Info("Saving contracts for chain ", chain)

			for _, name := range names {
				c := SupportedContracts[name]
				contract := &types.Contract{
					Chain:     chain,
					Hash:      c.AbiHash,
					Name:      name,
					ByteCodes: []byte(c.Bin),
				}

				contracts = append(contracts, contract)
			}
		}
	}

	h.txSubmit.SubmitMessageAsync(types.NewContractsWithSigner(
		h.appKeys.GetSignerAddress().String(),
		contracts,
	))
}