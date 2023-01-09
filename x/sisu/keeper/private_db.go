package keeper

import (
	adstore "github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cosmostypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	prefixPrivateTxOutSig      = []byte{0x01}
	prefixPrivatePendingTxOut  = []byte{0x02}
	prefixPrivateTxHashIndex   = []byte{0x03}
	prefixPrivateTransferQueue = []byte{0x0F}
)

type PrivateDb interface {
	// TxOutSig
	SaveTxOutSig(msg *types.TxOutSig)
	GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig

	SetPendingTxOut(chain string, txOut *types.PendingTxOutInfo)
	GetPendingTxOut(chain string) *types.PendingTxOutInfo

	// Transaction index. This value is used to identify 2 cosmos message with the same content but
	// used at different time (similar to nonce in Ethereum).
	SetTxHashIndex(key string, value uint32)
	GetTxHashIndex(key string) uint32
}

type defaultPrivateDb struct {
	store    cosmostypes.KVStore
	prefixes map[string]prefix.Store
}

func NewPrivateDb(dbDir string, backend dbm.BackendType) PrivateDb {
	log.Info("Private db dir = ", dbDir)
	db, err := dbm.NewDB("storage", backend, dbDir)
	if err != nil {
		panic(err)
	}

	store := &adstore.Store{
		DB: db,
	}

	return &defaultPrivateDb{
		store:    store,
		prefixes: initPrefixes(store),
	}
}

func initPrefixes(parent cosmostypes.KVStore) map[string]prefix.Store {
	prefixes := make(map[string]prefix.Store)

	prefixes[string(prefixPrivateTxOutSig)] = prefix.NewStore(parent, prefixPrivateTxOutSig)
	prefixes[string(prefixPrivatePendingTxOut)] = prefix.NewStore(parent, prefixPrivatePendingTxOut)
	prefixes[string(prefixPrivateTxHashIndex)] = prefix.NewStore(parent, prefixPrivateTxHashIndex)
	prefixes[string(prefixPrivateTransferQueue)] = prefix.NewStore(parent, prefixPrivateTransferQueue)

	return prefixes
}

///// TxOutSig
func (db *defaultPrivateDb) GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig {
	withSigStore := db.prefixes[string(prefixPrivateTxOutSig)]
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	return txOutSig
}

func (db *defaultPrivateDb) SaveTxOutSig(msg *types.TxOutSig) {
	store := db.prefixes[string(prefixPrivateTxOutSig)]
	saveTxOutSig(store, msg)
}

///// Pending TxOut

func (db *defaultPrivateDb) SetPendingTxOut(chain string, txOut *types.PendingTxOutInfo) {
	store := db.prefixes[string(prefixPrivatePendingTxOut)]
	setPendingTxOut(store, chain, txOut)
}

func (db *defaultPrivateDb) GetPendingTxOut(chain string) *types.PendingTxOutInfo {
	store := db.prefixes[string(prefixPrivatePendingTxOut)]
	return getPendingTxOutInfo(store, chain)
}

///// Tx Hash Index

func (db *defaultPrivateDb) SetTxHashIndex(key string, value uint32) {
	store := db.prefixes[string(prefixPrivateTxHashIndex)]
	setTxHashIndex(store, key, value)
}

func (db *defaultPrivateDb) GetTxHashIndex(key string) uint32 {
	store := db.prefixes[string(prefixPrivateTxHashIndex)]
	return getTxHashIndex(store, key)
}
