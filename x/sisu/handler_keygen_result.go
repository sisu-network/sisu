package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	liskcrypto "github.com/sisu-network/deyes/chains/lisk/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/background"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

type HandlerKeygenResult struct {
	keeper      keeper.Keeper
	pmm         components.PostedMessageManager
	globalData  components.GlobalData
	deyesClient external.DeyesClient
	config      config.TssConfig
	txSubmit    components.TxSubmit
	appKeys     components.AppKeys
	valsMgr     components.ValidatorManager
}

func NewHandlerKeygenResult(mc background.ManagerContainer) *HandlerKeygenResult {
	return &HandlerKeygenResult{
		keeper:      mc.Keeper(),
		pmm:         mc.PostedMessageManager(),
		globalData:  mc.GlobalData(),
		config:      mc.Config().Tss,
		txSubmit:    mc.TxSubmit(),
		appKeys:     mc.AppKeys(),
		valsMgr:     mc.ValidatorManager(),
		deyesClient: mc.DeyesClient(),
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
		log.Verbose("Keygen result: from: ", result.Data.From, " type = ", result.Keygen.KeyType,
			" success = ", result.Data.Result)
		if result.Data.Result == types.KeygenResult_SUCCESS {
			successCount += 1
		}
	}

	if successCount == h.valsMgr.GetValidatorLength(ctx) {
		// TODO: Make sure that everyone has the same address and pubkey.
		// Save keygen Address
		h.keeper.SaveKeygen(ctx, keygen)

		// Setting vaults & mpc address
		params := h.keeper.GetParams(ctx)
		for _, chain := range params.SupportedChains {
			if keygen.KeyType == utils.GetKeyTypeForChain(chain) {
				h.setVault(ctx, chain, keygen)
				h.setMpcAddress(ctx, chain, keygen)
			}
		}

		log.Infof("Keygen %s succeeded", keygen.KeyType)
	} else {
		// TODO: handle failure case
		return nil, nil
	}
	return nil, nil
}

func (h *HandlerKeygenResult) setMpcAddress(ctx sdk.Context, chain string, keygen *types.Keygen) {
	address := ""

	// Calculate the MPC address
	switch {
	case libchain.IsETHBasedChain(chain):
		// Calculate the ETH address
		address = keygen.Address
	case libchain.IsCardanoChain(chain):
		address = utils.GetCardanoAddressFromPubkey(keygen.PubKeyBytes).String()
	case libchain.IsSolanaChain(chain):
		address = utils.GetSolanaAddressFromPubkey(keygen.PubKeyBytes)
	case libchain.IsLiskChain(chain):
		address = liskcrypto.GetLisk32AddressFromPublickey(keygen.PubKeyBytes)
	default:
		log.Errorf("Unknown chain type %s", chain)
	}

	if address != "" {
		log.Verbosef("Setting mpc address for chain %s, addr = %s", chain, address)
		h.keeper.SetMpcAddress(ctx, chain, address)
		h.keeper.SetMpcPublicKey(ctx, chain, keygen.PubKeyBytes)
	}
}

func (h *HandlerKeygenResult) setVault(ctx sdk.Context, chain string, keygen *types.Keygen) {
	if libchain.IsSolanaChain(chain) {
		// In solana, each token has its own vault address. We have to loop through all token vaults
		vaults := h.keeper.GetAllVaultsForChain(ctx, chain)
		for _, vault := range vaults {
			log.Verbosef("Setting vault for %s %s %s", chain, vault.Address, vault.Token)
			h.deyesClient.SetVaultAddress(chain, vault.Address, vault.Token)
		}

	} else if libchain.IsCardanoChain(chain) {
		address := utils.GetCardanoAddressFromPubkey(keygen.PubKeyBytes)
		h.deyesClient.SetVaultAddress(chain, address.String(), "")

	} else if libchain.IsETHBasedChain(chain) {
		vault := h.keeper.GetVault(ctx, chain, "")
		log.Verbosef("Setting vault for %s %s", chain, vault.Address)
		h.deyesClient.SetVaultAddress(chain, vault.Address, "")

	} else if libchain.IsLiskChain(chain) {
		address := liskcrypto.GetLisk32AddressFromPublickey(keygen.PubKeyBytes)
		h.deyesClient.SetVaultAddress(chain, address, "")

	} else {
		// Unknown chains
		log.Error("Unknown chain: ", chain)
	}
}
