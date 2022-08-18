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
	log.Infof("Delivering keygen result of type %s, from %s, result = %v",
		signerMsg.Keygen.KeyType, signerMsg.Data.From, signerMsg.Data.Result)

	// Save result to KVStore & private db
	h.keeper.SaveKeygenResult(ctx, signerMsg)

	// TODO: Handler keygen failure and check result expiration time.
	results := h.keeper.GetAllKeygenResult(ctx, signerMsg.Keygen.KeyType, signerMsg.Keygen.Index)
	if len(results) == h.valsMgr.GetValidatorLength(ctx) {
		h.doKeygenResult(ctx, signerMsg.Keygen, results)
	}

	return &sdk.Result{}, nil
}

func (h *HandlerKeygenResult) doKeygenResult(ctx sdk.Context, keygen *types.Keygen,
	results []*types.KeygenResultWithSigner) ([]byte, error) {
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
		h.keeper.SaveKeygen(ctx, keygen)

		log.Infof("Keygen %s succeeded", keygen.KeyType)

		if !h.globalData.IsCatchingUp() {
			switch keygen.KeyType {
			case libchain.KEY_TYPE_ECDSA:
				h.createContracts(ctx, keygen)
			}
		}
	} else {
		// TODO: handle failure case
		return nil, nil
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
