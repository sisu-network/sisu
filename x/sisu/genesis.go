package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.DefaultKeeper, genState types.GenesisState) []abci.ValidatorUpdate {
	fmt.Println("len(genState) = ", len(genState.Nodes))

	validators := make([]abci.ValidatorUpdate, len(genState.Nodes))

	for i, node := range genState.Nodes {
		fmt.Println(node.Key.Type, node.Key.Bytes)

		pk, err := utils.GetCosmosPubKey(node.Key.Type, node.Key.Bytes)
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
