package keeper

import (
	adstore "github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cosmostypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	dbm "github.com/tendermint/tm-db"
)

type Storage interface {
	// TxOutSig
	SaveTxOutSig(msg *types.TxOutSig)
	GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig
}

type defaultStorage struct {
	store    cosmostypes.KVStore
	prefixes map[string]prefix.Store
}

func NewStorageDb(dbDir string) Storage {
	log.Info("Private db dir = ", dbDir)
	db, err := dbm.NewDB("storage", dbm.GoLevelDBBackend, dbDir)
	if err != nil {
		panic(err)
	}

	store := &adstore.Store{
		DB: db,
	}

	return &defaultStorage{
		store:    store,
		prefixes: initPrefixes(store),
	}
}

func initPrefixes(parent cosmostypes.KVStore) map[string]prefix.Store {
	prefixes := make(map[string]prefix.Store)

	// prefixTxOutSig
	prefixes[string(prefixTxOutSig)] = prefix.NewStore(parent, prefixTxOutSig)

	return prefixes
}

///// TxOutSig
func (db *defaultStorage) GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig {
	withSigStore := db.prefixes[string(prefixTxOutSig)]
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	return txOutSig
}

func (db *defaultStorage) SaveTxOutSig(msg *types.TxOutSig) {
	store := db.prefixes[string(prefixTxOutSig)]
	saveTxOutSig(store, msg)
}
