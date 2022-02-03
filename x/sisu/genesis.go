package sisu

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.DefaultKeeper, publicDb keeper.Storage, genState types.GenesisState) []abci.ValidatorUpdate {
	validators := make([]abci.ValidatorUpdate, len(genState.Nodes))

	for i, node := range genState.Nodes {
		pk, err := utils.GetCosmosPubKey(node.ConsensusKey.Type, node.ConsensusKey.Bytes)
		if err != nil {
			panic(err)
		}

		validators[i] = abci.Ed25519ValidatorUpdate(pk.Bytes(), 100)
	}

	return validators
}

// ExportGenesis returns the capability module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.DefaultKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
