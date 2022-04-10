package utils

import (
	"testing"

	tlog "github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
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
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, tlog.NewNopLogger())
	return ctx
}

func TestCloneSdkContext(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("store_key")
	transientKey := sdk.NewTransientStoreKey("transient_key")
	ctx := defaultContext(storeKey, transientKey)

	key := []byte("key")
	value := []byte("value")
	value2 := []byte("value two")

	store := prefix.NewStore(ctx.KVStore(storeKey), []byte{0x01})
	store.Set(key, value)

	clone := CloneSdkContext(ctx)
	store2 := prefix.NewStore(clone.KVStore(storeKey), []byte{0x01})
	store2.Set(key, value2)

	require.Equal(t, value, store.Get(key))   // value
	require.Equal(t, value2, store2.Get(key)) // value two
}
