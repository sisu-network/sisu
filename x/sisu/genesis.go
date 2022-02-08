package sisu

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, publicDb keeper.Storage, valsMgr ValidatorManager, genState types.GenesisState) []abci.ValidatorUpdate {
	// Saves initial token configs from genesis file.
	tokenIds := make([]string, 0)
	m := make(map[string]*types.Token)
	for _, token := range genState.Tokens {
		m[token.Id] = token
		tokenIds = append(tokenIds, token.Id)
	}
	publicDb.SetTokens(m)

	log.Info("Tokens in the genesis file: ", strings.Join(tokenIds, ", "))

	// Create validator nodes
	validators := make([]abci.ValidatorUpdate, len(genState.Nodes))
	for i, node := range genState.Nodes {
		if !node.IsValidator {
			continue
		}

		pk, err := utils.GetCosmosPubKey(node.ConsensusKey.Type, node.ConsensusKey.Bytes)
		if err != nil {
			panic(err)
		}

		validators[i] = abci.Ed25519ValidatorUpdate(pk.Bytes(), 100)
		valsMgr.AddValidator(node)
	}

	return validators
}

// ExportGenesis returns the capability module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.DefaultKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
