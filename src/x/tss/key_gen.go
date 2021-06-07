package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	PROPOSE_BLOCK_INTERVAL = 1000
)

type KeyGen struct {
	keeper keeper.Keeper
	config config.TssConfig

	lastProposeBlockHeight int
}

func NewKeyGen(keeper keeper.Keeper, config config.TssConfig) *KeyGen {
	return &KeyGen{
		keeper: keeper,
		config: config,
	}
}

func (kg *KeyGen) CheckTssKeygen(ctx sdk.Context, blockHeight int) {
	chainsInfo, err := kg.keeper.GetRecordedChainsOnSisu(ctx)
	if err != nil {
		return
	}

	recordedChains := make(map[string]*types.ChainInfo, len(chainsInfo.Chains))

	// Compare what we have in chains info and what we have in the config
	for _, chain := range chainsInfo.Chains {
		recordedChains[chain.Symbol] = chain
	}

	unavailableChains := make([]string, 0)
	for _, chainConfig := range kg.config.SupportedChains {
		if recordedChains[chainConfig.Symbol] == nil {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	if kg.lastProposeBlockHeight == 0 || blockHeight-kg.lastProposeBlockHeight > PROPOSE_BLOCK_INTERVAL {
		// Broadcast a message.
		utils.LogInfo("Broadcasting TSS Keygen Proposal message.")
	}
}
