package keeper

import (
	"encoding/base64"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"

	"github.com/cosmos/cosmos-sdk/store"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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
	keeper := NewKeeper(storeKey).(*DefaultKeeper)

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

func TestDefaultKeeper_IncAndDecSlashToken(t *testing.T) {
	t.Parallel()

	keeper, ctx := getTestKeeperAndContext()

	addr := []byte("0x1")
	require.NoError(t, keeper.IncSlashToken(ctx, addr, 1))
	curSlash, err := keeper.GetSlashToken(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(1), curSlash)

	require.NoError(t, keeper.DecSlashToken(ctx, addr, 1))
	curSlash, err = keeper.GetSlashToken(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), curSlash)

	require.NoError(t, keeper.DecSlashToken(ctx, addr, 1))
	curSlash, err = keeper.GetSlashToken(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), curSlash)
}

func TestDefaultKeeper_IncAndDecNodeBalance(t *testing.T) {
	t.Parallel()

	keeper, ctx := getTestKeeperAndContext()

	addr := []byte("0x1")
	require.NoError(t, keeper.IncBondBalance(ctx, addr, 1))
	curBalance, err := keeper.GetBondBalance(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(1), curBalance)

	require.NoError(t, keeper.DecBondBalance(ctx, addr, 1))
	curBalance, err = keeper.GetBondBalance(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), curBalance)

	require.NoError(t, keeper.DecBondBalance(ctx, addr, 1))
	curBalance, err = keeper.GetBondBalance(ctx, addr)
	require.NoError(t, err)
	require.Equal(t, int64(0), curBalance)
}

func TestDefaultKeeper_GetTopBalance(t *testing.T) {
	t.Parallel()

	keeper, ctx := getTestKeeperAndContext()

	addr1 := []byte("0x1")
	require.NoError(t, keeper.IncBondBalance(ctx, addr1, 1))
	addr2 := []byte("0x2")
	require.NoError(t, keeper.IncBondBalance(ctx, addr2, 3))
	addr3 := []byte("0x3")
	require.NoError(t, keeper.IncBondBalance(ctx, addr3, 2))

	top1Balance := keeper.GetTopBondBalance(ctx, 1)
	require.Len(t, top1Balance, 1)
	require.Equal(t, addr2, top1Balance[0].Bytes())

	top2Balances := keeper.GetTopBondBalance(ctx, 2)
	require.Len(t, top2Balances, 2)
	require.Equal(t, addr3, top2Balances[0].Bytes())
	require.Equal(t, addr2, top2Balances[1].Bytes())

	// addr1 = 1, addr2 = 0, addr3 = 2
	require.NoError(t, keeper.DecBondBalance(ctx, addr2, 3))
	top2Balances = keeper.GetTopBondBalance(ctx, 2)
	require.Len(t, top2Balances, 2)
	require.Equal(t, addr1, top2Balances[0].Bytes())
	require.Equal(t, addr3, top2Balances[1].Bytes())
}

func TestDefaultKeeper_SaveAndGetValidatorUpdates(t *testing.T) {
	t.Parallel()

	keeper, ctx := getTestKeeperAndContext()

	bz1, err := base64.StdEncoding.DecodeString("sk9Ab7wGydi2YEXJzOWAHlKBc/3un2i78hXEl0Mlohs=")
	require.NoError(t, err)
	v1 := abci.Ed25519ValidatorUpdate(bz1, 100)

	bz2, err := base64.StdEncoding.DecodeString("FAmIvjG2OsZ3wAdoqoXOwJPud9LCofVmF7CZfIXCO2k=")
	require.NoError(t, err)
	v2 := abci.Ed25519ValidatorUpdate(bz2, 0)

	valUpdates := []abci.ValidatorUpdate{v1, v2}
	require.NoError(t, keeper.SaveIncomingValidatorUpdates(ctx, valUpdates))

	afterSaved := keeper.GetIncomingValidatorUpdates(ctx)
	require.Len(t, afterSaved, 2)
}
