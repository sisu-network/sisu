package keeper

import (
	"testing"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveAndGetTxOut(t *testing.T) {
	keeper, ctx := GetTestKeeperAndContext()

	chain := "bitcoin"
	hash := utils.RandomHeximalString(32)

	original := &types.TxOut{
		Content: &types.TxOutContent{
			OutChain: chain,
			OutHash:  hash,
			OutBytes: []byte("Hash"),
		},
	}

	keeper.SaveTxOut(ctx, original)
	txOut := keeper.GetTxOut(ctx, chain, hash)
	require.Equal(t, original, txOut)

	// Any chain in OutChain, BlockHeight, OutBytes would not retrieve the txOut.
	txOut = keeper.GetTxOut(ctx, "eth", hash)
	require.Nil(t, txOut)

	txOut = keeper.GetTxOut(ctx, chain, utils.RandomHeximalString(48))
	require.Nil(t, txOut)
}

func TestVault(t *testing.T) {
	keeper, ctx := GetTestKeeperAndContext()

	ethVault := &types.Vault{
		Id:      "eth0",
		Chain:   "eth",
		Address: "0x-eth0",
		Token:   "Token0",
	}

	solVault0 := &types.Vault{
		Id:      "solana0",
		Chain:   "solana-devnet",
		Address: "0x-solana0",
		Token:   "Token0",
	}

	solVault1 := &types.Vault{
		Id:      "solana1",
		Chain:   "solana-devnet",
		Address: "0x-solana1",
		Token:   "Token1",
	}

	vaults := []*types.Vault{
		ethVault, solVault0, solVault1,
	}

	keeper.SetVaults(ctx, vaults)

	vault := keeper.GetVault(ctx, "solana-devnet", "Token0")
	assert.Equal(t, solVault0, vault)
	vault = keeper.GetVault(ctx, "solana-devnet", "Token1")
	assert.Equal(t, solVault1, vault)

	solVaults := keeper.GetAllVaultsForChain(ctx, "solana-devnet")
	assert.Equal(t, []*types.Vault{solVault0, solVault1}, solVaults)
}
