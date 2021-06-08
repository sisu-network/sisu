package tss

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	PROPOSE_BLOCK_INTERVAL = 1000
)

type KeyGen struct {
	keeper                 keeper.Keeper
	config                 config.TssConfig
	txSubmit               common.TxSubmit
	lastProposeBlockHeight int64
	appKeys                *common.AppKeys
}

func NewKeyGen(keeper keeper.Keeper,
	config config.TssConfig,
	appKeys *common.AppKeys,
	txSubmit common.TxSubmit,
) *KeyGen {
	return &KeyGen{
		keeper:   keeper,
		appKeys:  appKeys,
		config:   config,
		txSubmit: txSubmit,
	}
}

func (kg *KeyGen) CheckTssKeygen(ctx sdk.Context, blockHeight int64) {
	chainsInfo, err := kg.keeper.GetRecordedChainsOnSisu(ctx)
	if err != nil {
		return
	}

	recordedChains := make(map[string]*types.ChainInfo, len(chainsInfo.Chains))

	// Compare what we have in chains info and what we have in the config
	for _, chain := range chainsInfo.Chains {
		recordedChains[chain.Symbol] = chain
	}

	utils.LogInfo("recordedChains = ", recordedChains)

	unavailableChains := make([]string, 0)
	for _, chainConfig := range kg.config.SupportedChains {
		if recordedChains[chainConfig.Symbol] == nil {
			unavailableChains = append(unavailableChains, chainConfig.Symbol)
		}
	}

	if kg.lastProposeBlockHeight == 0 || blockHeight-kg.lastProposeBlockHeight > PROPOSE_BLOCK_INTERVAL {
		// Broadcast a message.
		utils.LogInfo("Broadcasting TSS Keygen Proposal message. len(unavailableChains) = ", len(unavailableChains))
		signer := kg.appKeys.GetSignerAddress()

		for _, chain := range unavailableChains {
			proposal := types.NewMsgKeygenProposal(signer.String(), chain)
			utils.LogDebug("Submitting proposal message for chain", chain)
			kg.txSubmit.SubmitMessage(proposal)
		}

		kg.lastProposeBlockHeight = blockHeight
	}
}
