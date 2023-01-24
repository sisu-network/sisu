package chains

import (
	"sync"

	libchain "github.com/sisu-network/lib/chain"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	chainscardano "github.com/sisu-network/sisu/x/sisu/chains/cardano"
	chainseth "github.com/sisu-network/sisu/x/sisu/chains/eth"
	chainslisk "github.com/sisu-network/sisu/x/sisu/chains/lisk"
	chainssolana "github.com/sisu-network/sisu/x/sisu/chains/solana"

	chainstypes "github.com/sisu-network/sisu/x/sisu/chains/types"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
)

type BridgeManager interface {
	GetBridge(ctx sdk.Context, chain string) chainstypes.Bridge
}

type defaultBridgeManager struct {
	bridges map[string]chainstypes.Bridge
	lock    *sync.RWMutex
	config  config.Config

	signer      string
	keeper      keeper.Keeper
	deyesClient external.DeyesClient
}

func NewBridgeManager(signer string, k keeper.Keeper, deyesClient external.DeyesClient, cfg config.Config) BridgeManager {
	return &defaultBridgeManager{
		bridges:     make(map[string]chainstypes.Bridge),
		lock:        &sync.RWMutex{},
		signer:      signer,
		keeper:      k,
		deyesClient: deyesClient,
		config:      cfg,
	}
}

func (m *defaultBridgeManager) GetBridge(ctx sdk.Context, chain string) chainstypes.Bridge {
	m.lock.Lock()
	defer m.lock.Unlock()

	bridge := m.bridges[chain]
	if bridge != nil {
		return bridge
	}

	// Make sure that this chain is among the supported chains
	params := m.keeper.GetParams(ctx)
	found := false
	for _, paramsChain := range params.SupportedChains {
		if paramsChain == chain {
			found = true
			break
		}
	}

	if !found {
		return nil
	}

	m.bridges[chain] = m.getBridge(chain)
	return m.bridges[chain]
}

func (m *defaultBridgeManager) getBridge(chain string) chainstypes.Bridge {
	if libchain.IsETHBasedChain(chain) {
		return chainseth.NewBridge(chain, m.signer, m.keeper, m.deyesClient)
	}

	if libchain.IsCardanoChain(chain) {
		return chainscardano.NewBridge(chain, m.signer, m.keeper, m.deyesClient)
	}

	if libchain.IsSolanaChain(chain) {
		return chainssolana.NewBridge(chain, m.signer, m.keeper, m.config, m.deyesClient)
	}

	if libchain.IsLiskChain(chain) {
		return chainslisk.NewBridge(chain, m.signer, m.keeper, m.deyesClient)
	}

	return nil
}
