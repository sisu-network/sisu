package tss

import (
	"github.com/sisu-network/cosmos-sdk/store"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/tendermint/libs/log"
	tmproto "github.com/sisu-network/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
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
