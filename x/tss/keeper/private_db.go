package keeper

import (
	adstore "github.com/sisu-network/cosmos-sdk/store/dbadapter"
	"github.com/sisu-network/cosmos-sdk/store/prefix"
	cosmostypes "github.com/sisu-network/cosmos-sdk/store/types"
	cstypes "github.com/sisu-network/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
	dbm "github.com/tendermint/tm-db"
)

// go:generate mockgen -source x/tss/keeper/private_db.go -destination=tests/mock/tss/private_db.go -package=mock
type PrivateDb interface {
	// Debug
	PrintStore(name string)

	// Keygen
	SaveKeygen(msg *types.Keygen)
	IsKeygenExisted(keyType string, index int) bool
	IsKeygenAddress(keyType string, address string) bool
	GetKeygenPubkey(keyType string) []byte

	// Keygen Result
	SaveKeygenResult(signerMsg *types.KeygenResultWithSigner)
	IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner) bool

	// Contracts
	SaveContracts(msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(msg *types.Contract) bool

	GetPendingContracts(chain string) []*types.Contract
	UpdateContractsStatus(msgs []*types.Contract, status string)

	// TxIn
	SaveTxIn(msg *types.TxIn)
	IsTxInExisted(msg *types.TxIn) bool

	// TxOut
	SaveTxOut(msg *types.TxOut)
	IsTxOutExisted(msg *types.TxOut) bool
	GetTxOut(inChain string, outChain, hash string) *types.TxOut
}

type defaultPrivateDb struct {
	store    cosmostypes.KVStore
	prefixes map[string]prefix.Store
}

func NewPrivateDb(dbDir string) PrivateDb {
	log.Info("Private db dir = ", dbDir)
	db, err := dbm.NewDB("private", dbm.GoLevelDBBackend, dbDir)
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

	// prefixKeygen
	prefixes[string(prefixKeygen)] = prefix.NewStore(parent, prefixKeygen)
	// prefixKeygenResult
	prefixes[string(prefixKeygenResult)] = prefix.NewStore(parent, prefixKeygenResult)
	// prefixContract
	prefixes[string(prefixContract)] = prefix.NewStore(parent, prefixContract)
	// prefixContractByteCode
	prefixes[string(prefixContractByteCode)] = prefix.NewStore(parent, prefixContractByteCode)
	// prefixTxIn
	prefixes[string(prefixTxIn)] = prefix.NewStore(parent, prefixTxIn)
	// prefixTxOut
	prefixes[string(prefixTxOut)] = prefix.NewStore(parent, prefixTxOut)

	return prefixes
}

///// Keygen

func (db *defaultPrivateDb) SaveKeygen(msg *types.Keygen) {
	store := db.prefixes[string(prefixKeygen)]
	saveKeygen(store, msg)
}

func (db *defaultPrivateDb) IsKeygenExisted(keyType string, index int) bool {
	store := db.prefixes[string(prefixKeygen)]
	return isKeygenExisted(store, keyType, index)
}

func (db *defaultPrivateDb) IsKeygenAddress(keyType string, address string) bool {
	store := db.prefixes[string(prefixKeygen)]
	return isKeygenAddress(store, keyType, address)
}

func (db *defaultPrivateDb) GetKeygenPubkey(keyType string) []byte {
	store := db.prefixes[string(prefixKeygen)]
	return getKeygenPubkey(store, keyType)
}

///// Keygen Result

func (db *defaultPrivateDb) SaveKeygenResult(signerMsg *types.KeygenResultWithSigner) {
	store := db.prefixes[string(prefixKeygenResult)]
	saveKeygenResult(store, signerMsg)
}

func (db *defaultPrivateDb) IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner) bool {
	store := db.prefixes[string(prefixKeygenResult)]
	return isKeygenResultSuccess(store, signerMsg)
}

///// Contract

func (db *defaultPrivateDb) SaveContracts(msgs []*types.Contract, saveByteCode bool) {
	contractStore := db.prefixes[string(prefixContract)]
	byteCodeStore := db.prefixes[string(prefixContractByteCode)]

	saveContracts(contractStore, byteCodeStore, msgs, saveByteCode)
}

func (db *defaultPrivateDb) IsContractExisted(msg *types.Contract) bool {
	contractStore := db.prefixes[string(prefixContract)]
	return isContractExisted(contractStore, msg)
}

func (db *defaultPrivateDb) GetPendingContracts(chain string) []*types.Contract {
	contractStore := db.prefixes[string(prefixContract)]
	byteCodeStore := db.prefixes[string(prefixContractByteCode)]

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (db *defaultPrivateDb) UpdateContractsStatus(msgs []*types.Contract, status string) {
	contractStore := db.prefixes[string(prefixContract)]
	updateContractsStatus(contractStore, msgs, status)
}

///// TxIn
func (db *defaultPrivateDb) SaveTxIn(msg *types.TxIn) {
	store := db.prefixes[string(prefixTxIn)]
	saveTxIn(store, msg)
}

func (db *defaultPrivateDb) IsTxInExisted(msg *types.TxIn) bool {
	store := db.prefixes[string(prefixTxIn)]
	return isTxInExisted(store, msg)
}

///// TxOut
func (db *defaultPrivateDb) SaveTxOut(msg *types.TxOut) {
	store := db.prefixes[string(prefixTxOut)]
	saveTxOut(store, msg)
}

func (db *defaultPrivateDb) IsTxOutExisted(msg *types.TxOut) bool {
	store := db.prefixes[string(prefixTxOut)]
	return isTxOutExisted(store, msg)
}

func (db *defaultPrivateDb) GetTxOut(inChain string, outChain, hash string) *types.TxOut {
	store := db.prefixes[string(prefixTxOut)]
	return getTxOut(store, inChain, outChain, hash)
}

///// Debug

// PrintStore is a debug function
func (db *defaultPrivateDb) PrintStore(name string) {
	log.Info("======== DEBUGGING")
	log.Info("Printing ALL values in store ", name)
	var store cstypes.KVStore
	switch name {
	case "keygen":
		store = db.prefixes[string(prefixKeygen)]
	case "contract":
		store = db.prefixes[string(prefixContract)]
	}

	printStore(store)
	log.Info("======== END OF DEBUGGING")
}
