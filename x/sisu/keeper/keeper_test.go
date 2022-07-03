package keeper

import (
	"sort"
	"strings"
	"testing"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/stretchr/testify/require"
)

func TestKeeper_SaveAndGetObservedTx(t *testing.T) {
	// t.Parallel()
	// keeper, ctx := GetTestKeeperAndContext()

	// observedTx := &types.TxIn{
	// 	Chain:       "eth",
	// 	BlockHeight: 1,
	// 	TxHash:      "Hash",
	// 	Serialized:  []byte("Serialized"),
	// }

	// // Save observed Tx
	// keeper.SaveTxIn(ctx, observedTx)

	// // Check Observed Tx
	// require.Equal(t, true, keeper.IsTxInExisted(ctx, observedTx))

	// // Different signer would not change the observedTx retrieval
	// other := *observedTx
	// other.Chain = "signer2"
	// require.Equal(t, true, keeper.IsTxInExisted(ctx, observedTx))

	// // Any change in the chain, block height or tx hash would not retrieve the observed tx.
	// other = *observedTx
	// other.Chain = "bitcoin"
	// require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))

	// other = *observedTx
	// other.BlockHeight = 2
	// require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))

	// other = *observedTx
	// other.TxHash = "Hash2"
	// require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))
}

func TestKeeper_SaveAndGetTxOut(t *testing.T) {
	t.Parallel()
	keeper, ctx := GetTestKeeperAndContext()

	txOutWithSigner := &types.TxOutWithSigner{
		Signer: "signer",
		Data: &types.TxOut{
			InChain:       "eth",
			OutChain:      "bitcoin",
			OutHash:       utils.RandomHeximalString(32),
			InBlockHeight: 1,
			OutBytes:      []byte("Hash"),
		},
	}

	keeper.SaveTxOut(ctx, txOutWithSigner.Data)
	require.Equal(t, true, keeper.IsTxOutExisted(ctx, txOutWithSigner.Data))

	// Different signer would not change the observedTx retrieval
	other := *txOutWithSigner.Data
	require.Equal(t, true, keeper.IsTxOutExisted(ctx, txOutWithSigner.Data))

	// Any chain in OutChain, BlockHeight, OutBytes would not retrieve the txOut.
	other = *txOutWithSigner.Data
	other.OutChain = "sisu"
	require.Equal(t, false, keeper.IsTxOutExisted(ctx, &other))

	other = *txOutWithSigner.Data
	other.OutHash = utils.RandomHeximalString(48)
	require.Equal(t, false, keeper.IsTxOutExisted(ctx, &other))
}

func TestKeeper_BlockHeights(t *testing.T) {
	keeper, ctx := GetTestKeeperAndContext()
	keeper.SaveBlockHeights(ctx, "signer1", &types.BlockHeightRecord{
		BlockHeights: []*types.BlockHeight{
			{
				Chain: "ganache1",
			},
			{
				Chain: "ganache2",
			},
		},
	})

	keeper.SaveBlockHeights(ctx, "signer2", &types.BlockHeightRecord{
		BlockHeights: []*types.BlockHeight{
			{
				Chain: "ganache1",
			},
		},
	})
	keeper.SaveBlockHeights(ctx, "signer3", &types.BlockHeightRecord{
		BlockHeights: []*types.BlockHeight{
			{
				Chain: "ganache1",
			},
		},
	})

	blockHeightRecord := keeper.GetBlockHeightRecord(ctx, "signer1")
	require.Equal(t, []*types.BlockHeight{
		{
			Chain: "ganache1",
		},
		{
			Chain: "ganache2",
		},
	}, blockHeightRecord.BlockHeights)

	blockHeightsMap := keeper.GetBlockHeightsForChain(ctx, "ganache1", []string{"ganache1", "ganache2"})
	_, blockHeights := types.ConvertBlockHeightsMapToArray(blockHeightsMap)
	sort.Slice(blockHeights, func(i, j int) bool {
		return strings.Compare(blockHeights[i].Chain, blockHeights[j].Chain) < 0
	})
	require.Equal(t, []*types.BlockHeight{
		{
			Chain: "ganache1",
		},
		{
			Chain: "ganache2",
		},
	}, blockHeightRecord.BlockHeights)
}
