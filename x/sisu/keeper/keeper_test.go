package keeper

import (
	"testing"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
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

func TestKeeper_BlockHeights(t *testing.T) {
	// keeper, ctx := GetTestKeeperAndContext()
	// keeper.SaveBlockHeights(ctx, "signer1", &types.BlockHeightRecord{
	// 	BlockHeights: []*types.BlockHeight{
	// 		{
	// 			Chain: "ganache1",
	// 		},
	// 		{
	// 			Chain: "ganache2",
	// 		},
	// 	},
	// })

	// keeper.SaveBlockHeights(ctx, "signer2", &types.BlockHeightRecord{
	// 	BlockHeights: []*types.BlockHeight{
	// 		{
	// 			Chain: "ganache1",
	// 		},
	// 	},
	// })
	// keeper.SaveBlockHeights(ctx, "signer3", &types.BlockHeightRecord{
	// 	BlockHeights: []*types.BlockHeight{
	// 		{
	// 			Chain: "ganache1",
	// 		},
	// 	},
	// })

	// blockHeightRecord := keeper.GetBlockHeightRecord(ctx, "signer1")
	// require.Equal(t, []*types.BlockHeight{
	// 	{
	// 		Chain: "ganache1",
	// 	},
	// 	{
	// 		Chain: "ganache2",
	// 	},
	// }, blockHeightRecord.BlockHeights)

	// blockHeightsMap := keeper.GetBlockHeightsForChain(ctx, "ganache1", []string{"ganache1", "ganache2"})
	// _, blockHeights := types.ConvertBlockHeightsMapToArray(blockHeightsMap)

	// sort.Slice(blockHeights, func(i, j int) bool {
	// 	return strings.Compare(blockHeights[i].Chain, blockHeights[j].Chain) < 0
	// })
	// require.Equal(t, []*types.BlockHeight{
	// 	{
	// 		Chain: "ganache1",
	// 	},
	// 	{
	// 		Chain: "ganache2",
	// 	},
	// }, blockHeightRecord.BlockHeights)
}
