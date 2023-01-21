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
	prefixPrivateTxOutSig       = []byte{0x01}
	prefixPrivateTransferState  = []byte{0x04}
	prefixPrivateTxOutState     = []byte{0x05}
	prefixPrivateHoldProcessing = []byte{0x06}
)

type PrivateDb interface {
	// TxOutSig
	SaveTxOutSig(msg *types.TxOutSig)
	GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig

	// Transfer State
	SetTransferState(id string, state types.TransferState)
	GetTransferState(id string) types.TransferState

	// TxOut State
	SetTxOutState(id string, state types.TxOutState)
	GetTxOutState(id string) types.TxOutState

	// Hold Processing job on a chain (e.g. Transfer or TxOut Queue)
	SetHoldProcessing(jobType, chain string, hold bool)
	GetHoldProcessing(jobType, chain string) bool
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
	prefixes[string(prefixPrivateTransferState)] = prefix.NewStore(parent, prefixPrivateTransferState)
	prefixes[string(prefixPrivateTxOutState)] = prefix.NewStore(parent, prefixPrivateTxOutState)
	prefixes[string(prefixPrivateHoldProcessing)] = prefix.NewStore(parent, prefixPrivateHoldProcessing)

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

///// Transfer State
func (db *defaultPrivateDb) SetTransferState(id string, state types.TransferState) {
	store := db.prefixes[string(prefixPrivateTransferState)]
	setState(store, id, int(state))
}

func (db *defaultPrivateDb) GetTransferState(id string) types.TransferState {
	store := db.prefixes[string(prefixPrivateTransferState)]
	return types.TransferState(getState(store, id))
}

///// TxOut State
func (db *defaultPrivateDb) SetTxOutState(id string, state types.TxOutState) {
	store := db.prefixes[string(prefixPrivateTxOutState)]
	setState(store, id, int(state))
}

func (db *defaultPrivateDb) GetTxOutState(id string) types.TxOutState {
	store := db.prefixes[string(prefixPrivateTxOutState)]
	return types.TxOutState(getState(store, id))
}

///// Hold Processing job on a chain (e.g. Transfer or TxOut Queue)
func (db *defaultPrivateDb) SetHoldProcessing(jobType, chain string, hold bool) {
	store := db.prefixes[string(prefixPrivateHoldProcessing)]
	setHoldProcessing(store, jobType, chain, hold)
}

func (db *defaultPrivateDb) GetHoldProcessing(jobType, chain string) bool {
	store := db.prefixes[string(prefixPrivateHoldProcessing)]
	return getHoldProcessing(store, jobType, chain)
}
