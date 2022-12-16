package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func GetTestDefaultContext(key sdk.StoreKey, tkey sdk.StoreKey) sdk.Context {
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

func GetTestKeeperAndContext() (*DefaultKeeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey("store_key")
	transientKey := sdk.NewTransientStoreKey("transient_key")
	ctx := GetTestDefaultContext(storeKey, transientKey)
	keeper := NewKeeper(storeKey).(*DefaultKeeper)

	return keeper, ctx
}

func GetTestStorage() Storage {
	return NewStorageDb(".", dbm.MemDBBackend)
}
