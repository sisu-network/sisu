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
	keeper      keeper.Keeper
	pmm         PostedMessageManager
	globalData  common.GlobalData
	deyesClient tssclients.DeyesClient
	config      config.TssConfig
	txSubmit    common.TxSubmit
	appKeys     common.AppKeys
	valsMgr     ValidatorManager
}

func NewHandlerKeygenResult(mc ManagerContainer) *HandlerKeygenResult {
	return &HandlerKeygenResult{
		keeper:     mc.Keeper(),
		pmm:        mc.PostedMessageManager(),
		globalData: mc.GlobalData(),
		config:     mc.Config(),
		txSubmit:   mc.TxSubmit(),
		appKeys:    mc.AppKeys(),
		valsMgr:    mc.ValidatorManager(),
	}
}

func (h *HandlerKeygenResult) DeliverMsg(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) (*sdk.Result, error) {
	// Save result to KVStore & private db
	h.keeper.SaveKeygenResult(ctx, signerMsg)

	if process, hash := h.pmm.ShouldProcessMsg(ctx, signerMsg); process {
		data, err := h.doKeygenResult(ctx, signerMsg)
		h.keeper.ProcessTxRecord(ctx, hash)

		return &sdk.Result{Data: data}, err
	}

	return &sdk.Result{}, nil
}

func (h *HandlerKeygenResult) doKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) ([]byte, error) {
	msg := signerMsg.Data

	log.Infof("Delivering keygen result of type %s, result = %v", signerMsg.Keygen.KeyType, msg.Result)

	// Get the majority of of votes
	// TODO: Only count result from validator node. Otherwise, anyone could post fake result.
	results := h.keeper.GetAllKeygenResult(ctx, signerMsg.Keygen.KeyType, signerMsg.Keygen.Index)
	if len(results) == h.valsMgr.GetValidatorLength(ctx) {
		// Check the majority of the results
		successCount := 0
		for _, result := range results {
			if result.Data.Result == types.KeygenResult_SUCCESS {
				successCount += 1
			}
		}

		if successCount == h.valsMgr.GetValidatorLength(ctx) {
			// TODO: Make sure that everyone has the same address and pubkey.
			// Save keygen Address
			log.Info("Saving keygen...")
			h.keeper.SaveKeygen(ctx, signerMsg.Keygen)

			log.Infof("Keygen %s succeeded", signerMsg.Keygen.KeyType)

			if !h.globalData.IsCatchingUp() {
				switch signerMsg.Keygen.KeyType {
				case libchain.KEY_TYPE_ECDSA:
					h.createContracts(ctx, signerMsg.Keygen)
				}
			}
		} else {
			// TODO: handle failure case
			return nil, nil
		}
	}

	return nil, nil
}

// createContracts creates and broadcast pending contracts. All nodes need to agree what
// contracts to deploy on what chains.
func (h *HandlerKeygenResult) createContracts(ctx sdk.Context, msg *types.Keygen) {
	log.Info("Create and broadcast contracts...")

	// Sort all contracts name alphabetically
	names := make([]string, len(SupportedContracts))
	i := 0
	for contract := range SupportedContracts {
		names[i] = contract
		i += 1
	}
	sort.Strings(names)

	params := h.keeper.GetParams(ctx)

	// Create contracts
	for _, chain := range params.SupportedChains {
		if libchain.GetKeyTypeForChain(chain) == msg.KeyType {
			log.Info("Saving contracts for chain ", chain)

			for _, name := range names {
				c := SupportedContracts[name]
				if !c.IsDeployBySisu {
					continue
				}

				contract := &types.Contract{
					Chain:     chain,
					Hash:      c.AbiHash,
					Name:      name,
					ByteCodes: []byte(c.Bin),
				}

				h.txSubmit.SubmitMessageAsync(types.NewContractsMsg(
					h.appKeys.GetSignerAddress().String(),
					[]*types.Contract{contract},
				))
			}
		}
	}
}
