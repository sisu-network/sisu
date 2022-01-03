package keeper

import (
	"testing"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"

	"github.com/sisu-network/cosmos-sdk/store"
	"github.com/sisu-network/tendermint/libs/log"
	tmproto "github.com/sisu-network/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/stretchr/testify/require"
)

func defaultContext(key sdk.StoreKey, tkey sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tkey, sdk.StoreTypeTransient, db)
	err := cms.LoadLatestVersion()
	if err != nil {
		panic(err)
	}
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, log.NewNopLogger())
	return ctx
}

func getTestKeeperAndContext() (*DefaultKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("store_key")
	transientKey := sdk.NewTransientStoreKey("transient_key")
	ctx := defaultContext(storeKey, transientKey)
	keeper := NewKeeper(storeKey)

	return keeper, ctx
}

func TestKeeper_SaveAndGetObservedTx(t *testing.T) {
	t.Parallel()
	keeper, ctx := getTestKeeperAndContext()

	observedTx := &types.TxIn{
		Chain:       "eth",
		BlockHeight: 1,
		TxHash:      "Hash",
		Serialized:  []byte("Serialized"),
	}

	// Save observed Tx
	keeper.SaveTxIn(ctx, observedTx)

	// Check Observed Tx
	require.Equal(t, true, keeper.IsTxInExisted(ctx, observedTx))

	// Different signer would not change the observedTx retrieval
	other := *observedTx
	other.Chain = "signer2"
	require.Equal(t, true, keeper.IsTxInExisted(ctx, observedTx))

	// Any change in the chain, block height or tx hash would not retrieve the observed tx.
	other = *observedTx
	other.Chain = "bitcoin"
	require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))

	other = *observedTx
	other.BlockHeight = 2
	require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))

	other = *observedTx
	other.TxHash = "Hash2"
	require.Equal(t, false, keeper.IsTxInExisted(ctx, &other))
}

func TestKeeper_SaveAndGetTxOut(t *testing.T) {
	t.Parallel()
	keeper, ctx := getTestKeeperAndContext()

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
