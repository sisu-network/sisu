package tss

import (
	"encoding/json"
	"testing"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestKvStore(t *testing.T) {
	t.Parallel()

	storeKey := sdk.NewKVStoreKey("store_key")
	transientKey := sdk.NewTransientStoreKey("transient_key")
	ctx := defaultContext(storeKey, transientKey)
	store := NewDefaultKVStore(storeKey)
	txOutEntities := []*types.TxOutEntity{{
		OutChain:        utils.RandomHeximalString(10),
		HashWithoutSig:  utils.RandomHeximalString(10),
		HashWithSig:     utils.RandomHeximalString(10),
		InChain:         utils.RandomHeximalString(10),
		InHash:          utils.RandomHeximalString(10),
		BytesWithoutSig: nil,
		Status:          utils.RandomHeximalString(10),
		Signature:       utils.RandomHeximalString(10),
		ContractHash:    utils.RandomHeximalString(10),
	}}
	require.NotPanics(t, func() {
		store.InsertTxOuts(ctx, txOutEntities)
	})

	expected := txOutEntities[0]
	key := getTxOutKey(expected.InChain, expected.HashWithoutSig)
	bz := getKVStoreValueByKey(ctx, storeKey, key, PrefixTxOut)
	saved := &types.TxOutEntity{}
	err := json.Unmarshal(bz, saved)
	require.NoError(t, err)

	assertEqualTxOutEntity(t, expected, saved)
}

func assertEqualTxOutEntity(t *testing.T, left, right *types.TxOutEntity) {
	require.Equal(t, left.OutChain, right.OutChain)
	require.Equal(t, left.HashWithoutSig, right.HashWithoutSig)
	require.Equal(t, left.HashWithSig, right.HashWithSig)
	require.Equal(t, left.InChain, right.InChain)
	require.Equal(t, left.InHash, right.InHash)
	require.Equal(t, left.Status, right.Status)
	require.Equal(t, left.Signature, right.Signature)
	require.Equal(t, left.ContractHash, right.ContractHash)
}

func getKVStoreValueByKey(ctx sdk.Context, storeKey sdk.StoreKey, key, prefixKey []byte) []byte {
	store := prefix.NewStore(ctx.KVStore(storeKey), prefixKey)

	return store.Get(key)
}
