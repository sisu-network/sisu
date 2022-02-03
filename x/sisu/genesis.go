package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.DefaultKeeper, publicDb keeper.Storage, genState types.GenesisState) []abci.ValidatorUpdate {
	validators := make([]abci.ValidatorUpdate, len(genState.Nodes))

	for i, node := range genState.Nodes {
		validators[i] = abci.Ed25519ValidatorUpdate(node.ConsensusKey.Bytes, 100)
	}

	fmt.Println("End of genesis, validators size = ", len(validators))

	return validators
}

// ExportGenesis returns the capability module's exported genesis
func ExportGenesis(ctx sdk.Context, k keeper.DefaultKeeper) *types.GenesisState {
	genesis := types.DefaultGenesis()

	return genesis
}
