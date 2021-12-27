package keeper

import (
	adstore "github.com/sisu-network/cosmos-sdk/store/dbadapter"
	"github.com/sisu-network/cosmos-sdk/store/prefix"
	cosmostypes "github.com/sisu-network/cosmos-sdk/store/types"
	"github.com/sisu-network/sisu/x/tss/types"
	dbm "github.com/tendermint/tm-db"
)

type PrivateDb interface {
	// Keygen
	SaveKeygen(msg *types.Keygen)
	IsKeygenExisted(keyType string, index int) bool
	GetAllPubKeys() map[string][]byte

	// Keygen Result
	SaveKeygenResult(signerMsg *types.KeygenResultWithSigner)
	IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner) bool

	// Contracts
	SaveContracts(msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(msg *types.Contract) bool

	GetPendingContracts(chain string) []*types.Contract
	UpdateContractsStatus(msgs []*types.Contract, status string)
}

type defaultPrivateDb struct {
	store    cosmostypes.KVStore
	prefixes map[string]prefix.Store
}

func NewPrivateDb(dbDir string) PrivateDb {
	db, err := dbm.NewDB("private_db", dbm.GoLevelDBBackend, dbDir)
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

func (db *defaultPrivateDb) SaveKeygen(msg *types.Keygen) {
	store := db.prefixes[string(prefixKeygen)]
	saveKeygen(store, msg)
}

func (db *defaultPrivateDb) IsKeygenExisted(keyType string, index int) bool {
	store := db.prefixes[string(prefixKeygen)]
	return isKeygenExisted(store, keyType, index)
}

func (db *defaultPrivateDb) GetAllPubKeys() map[string][]byte {
	store := db.prefixes[string(prefixKeygen)]
	return getAllPubKeys(store)
}

func (db *defaultPrivateDb) SaveKeygenResult(signerMsg *types.KeygenResultWithSigner) {
	store := db.prefixes[string(prefixKeygenResult)]
	saveKeygenResult(store, signerMsg)
}

func (db *defaultPrivateDb) IsKeygenResultSuccess(signerMsg *types.KeygenResultWithSigner) bool {
	store := db.prefixes[string(prefixKeygenResult)]
	return isKeygenResultSuccess(store, signerMsg)
}

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
