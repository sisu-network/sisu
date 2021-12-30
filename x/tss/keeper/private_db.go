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
	PrintStoreKeys(name string)

	// Keygen
	SaveKeygen(msg *types.Keygen)
	IsKeygenExisted(keyType string, index int) bool
	IsKeygenAddress(keyType string, address string) bool
	GetKeygenPubkey(keyType string) []byte

	// Keygen Result
	SaveKeygenResult(signerMsg *types.KeygenResultWithSigner)
	IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner, self string) bool

	// Contract
	SaveContract(msg *types.Contract, saveByteCode bool)
	SaveContracts(msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(msg *types.Contract) bool
	GetContract(chain string, hash string, includeByteCode bool) *types.Contract

	GetPendingContracts(chain string) []*types.Contract
	UpdateContractAddress(chain string, hash string, address string)

	UpdateContractsStatus(msgs []*types.Contract, status string)

	// Contract Address
	CreateContractAddress(chain string, txOutHash string, address string)
	IsContractExistedAtAddress(chain string, address string) bool

	// TxIn
	SaveTxIn(msg *types.TxIn)
	IsTxInExisted(msg *types.TxIn) bool

	// TxOut
	SaveTxOut(msg *types.TxOut)
	IsTxOutExisted(msg *types.TxOut) bool
	GetTxOut(outChain, hash string) *types.TxOut
	GetTxOutFromSigHash(outChain, hashWithSig string) *types.TxOut

	// TODO: Add unconfirmed tx store
	// TxOutSig
	SaveTxOutSig(msg *types.TxOutSig)

	// TxOutConfirm
	SaveTxOutConfirm(msg *types.TxOutConfirm)
	IsTxOutConfirmExisted(outChain, hash string) bool
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
	// prefixContractAddress
	prefixes[string(prefixContractAddress)] = prefix.NewStore(parent, prefixContractAddress)
	// prefixTxIn
	prefixes[string(prefixTxIn)] = prefix.NewStore(parent, prefixTxIn)
	// prefixTxOut
	prefixes[string(prefixTxOut)] = prefix.NewStore(parent, prefixTxOut)
	// prefixTxOutConfirm
	prefixes[string(prefixTxOutConfirm)] = prefix.NewStore(parent, prefixTxOutConfirm)

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

func (db *defaultPrivateDb) IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner, self string) bool {
	store := db.prefixes[string(prefixKeygenResult)]
	return isKeygenResultSuccess(store, signerMsg, self)
}

///// Contract
func (db *defaultPrivateDb) SaveContract(msg *types.Contract, saveByteCode bool) {
	contractStore := db.prefixes[string(prefixContract)]
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = db.prefixes[string(prefixContractByteCode)]
	}

	saveContract(contractStore, byteCodeStore, msg)
}

func (db *defaultPrivateDb) SaveContracts(msgs []*types.Contract, saveByteCode bool) {
	contractStore := db.prefixes[string(prefixContract)]
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = db.prefixes[string(prefixContractByteCode)]
	}

	saveContracts(contractStore, byteCodeStore, msgs)
}

func (db *defaultPrivateDb) IsContractExisted(msg *types.Contract) bool {
	contractStore := db.prefixes[string(prefixContract)]
	return isContractExisted(contractStore, msg)
}

func (db *defaultPrivateDb) GetContract(chain string, hash string, includeByteCode bool) *types.Contract {
	contractStore := db.prefixes[string(prefixContract)]
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if includeByteCode {
		byteCodeStore = db.prefixes[string(prefixContractByteCode)]
	}

	return getContract(contractStore, byteCodeStore, chain, hash)
}

func (db *defaultPrivateDb) GetPendingContracts(chain string) []*types.Contract {
	contractStore := db.prefixes[string(prefixContract)]
	byteCodeStore := db.prefixes[string(prefixContractByteCode)]

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (db *defaultPrivateDb) UpdateContractAddress(chain string, hash string, address string) {
	contractStore := db.prefixes[string(prefixContract)]
	updateContractAddress(contractStore, chain, hash, address)
}

func (db *defaultPrivateDb) UpdateContractsStatus(msgs []*types.Contract, status string) {
	contractStore := db.prefixes[string(prefixContract)]
	updateContractsStatus(contractStore, msgs, status)
}

///// Contract Address
func (db *defaultPrivateDb) CreateContractAddress(chain string, txOutHash string, address string) {
	caStore := db.prefixes[string(prefixContractAddress)]
	txOutStore := db.prefixes[string(prefixTxOut)]

	createContractAddress(caStore, txOutStore, chain, txOutHash, address)
}

func (db *defaultPrivateDb) IsContractExistedAtAddress(chain string, address string) bool {
	caStore := db.prefixes[string(prefixContractAddress)]

	return isContractExistedAtAddress(caStore, chain, address)
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

func (db *defaultPrivateDb) GetTxOut(outChain, hash string) *types.TxOut {
	store := db.prefixes[string(prefixTxOut)]
	return getTxOut(store, outChain, hash)
}

func (db *defaultPrivateDb) GetTxOutFromSigHash(outChain, hashWithSig string) *types.TxOut {
	withSigStore := db.prefixes[string(prefixTxOutSig)]
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	noSigStore := db.prefixes[string(prefixTxOut)]
	return getTxOut(noSigStore, outChain, txOutSig.HashNoSig)
}

///// TxOutSig
func (db *defaultPrivateDb) SaveTxOutSig(msg *types.TxOutSig) {
	store := db.prefixes[string(prefixTxOutSig)]
	saveTxOutSig(store, msg)
}

///// TxOutConfirm
func (db *defaultPrivateDb) SaveTxOutConfirm(msg *types.TxOutConfirm) {
	store := db.prefixes[string(prefixTxOutConfirm)]
	saveTxOutConfirm(store, msg)
}

func (db *defaultPrivateDb) IsTxOutConfirmExisted(chain, hash string) bool {
	store := db.prefixes[string(prefixTxOutConfirm)]
	return isTxOutConfirmExisted(store, chain, hash)
}

///// Debug

func (db *defaultPrivateDb) getStoreFromName(name string) cstypes.KVStore {
	var store cstypes.KVStore
	switch name {
	case "keygen":
		store = db.prefixes[string(prefixKeygen)]
	case "contract":
		store = db.prefixes[string(prefixContract)]
	case "txOut":
		store = db.prefixes[string(prefixTxOut)]
	}

	return store
}

// PrintStore is a debug function
func (db *defaultPrivateDb) PrintStore(name string) {
	log.Info("======== DEBUGGING PrintStore")
	log.Info("Printing ALL values in store ", name)

	store := db.getStoreFromName(name)
	if store != nil {
		printStore(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}

func (db *defaultPrivateDb) PrintStoreKeys(name string) {
	log.Info("======== DEBUGGING PrintStoreKeys")
	store := db.getStoreFromName(name)
	if store != nil {
		printStoreKeys(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}
