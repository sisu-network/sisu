package keeper

import (
	adstore "github.com/cosmos/cosmos-sdk/store/dbadapter"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cosmostypes "github.com/cosmos/cosmos-sdk/store/types"
	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	dbm "github.com/tendermint/tm-db"
)

// go:generate mockgen -source x/sisu/keeper/storage.go -destination=tests/mock/tss/storage.go -package=mock
type Storage interface {
	// Debug
	PrintStore(name string)
	PrintStoreKeys(name string)

	// TxRecord
	SaveTxRecord(hash []byte, signer string) int

	// TxRecordProcessed
	ProcessTxRecord(hash []byte)
	IsTxRecordProcessed(hash []byte) bool

	// Keygen
	SaveKeygen(msg *types.Keygen)
	IsKeygenExisted(keyType string, index int) bool
	IsKeygenAddress(keyType string, address string) bool
	GetKeygenPubkey(keyType string) []byte
	GetAllKeygenPubkeys() map[string][]byte

	// Keygen Result
	SaveKeygenResult(signerMsg *types.KeygenResultWithSigner)
	GetAllKeygenResult(keygenType string, index int32) []*types.KeygenResultWithSigner

	// Contract
	SaveContract(msg *types.Contract, saveByteCode bool)
	SaveContracts(msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(msg *types.Contract) bool
	GetContract(chain string, hash string, includeByteCode bool) *types.Contract
	GetPendingContracts(chain string) []*types.Contract
	UpdateContractAddress(chain string, hash string, address string)
	UpdateContractsStatus(chain string, contractHash string, status string)
	GetLatestContractAddressByName(chain, name string) string

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

	// TxOutSig
	// TODO: Add unconfirmed tx store
	SaveTxOutSig(msg *types.TxOutSig)
	GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig

	// TxOutConfirm
	SaveTxOutConfirm(msg *types.TxOutContractConfirm)
	IsTxOutConfirmExisted(outChain, hash string) bool

	// Gas Price Record
	SetGasPrice(msg *types.GasPriceMsg)
	GetGasPriceRecord(chain string, height int64) *types.GasPriceRecord

	// Chain
	SaveChain(chain *types.Chain)
	GetChain(chain string) *types.Chain
	GetAllChains() map[string]*types.Chain

	// Token Price
	SetTokenPrices(blockHeight uint64, msg *types.UpdateTokenPrice)
	GetAllTokenPricesRecord() map[string]*types.TokenPriceRecord

	// Set Tokens
	SetTokens(map[string]*types.Token)
	GetTokens([]string) map[string]*types.Token
	GetAllTokens() map[string]*types.Token

	// Nodes
	SaveNode(node *types.Node)
	LoadValidators() []*types.Node

	// Liquidities
	SetLiquidities(map[string]*types.Liquidity)
	GetLiquidity(chain string) *types.Liquidity
	GetAllLiquidities() map[string]*types.Liquidity

	// Params
	SaveParams(*types.Params)
	GetParams() *types.Params
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

	// prefixTxRecord
	prefixes[string(prefixTxRecord)] = prefix.NewStore(parent, prefixTxRecord)
	// prefixTxRecordProcessed
	prefixes[string(prefixTxRecordProcessed)] = prefix.NewStore(parent, prefixTxRecordProcessed)
	// prefixKeygen
	prefixes[string(prefixKeygen)] = prefix.NewStore(parent, prefixKeygen)
	// prefixKeygenResult
	prefixes[string(prefixKeygenResultWithSigner)] = prefix.NewStore(parent, prefixKeygenResultWithSigner)
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
	// prefixTxOutSig
	prefixes[string(prefixTxOutSig)] = prefix.NewStore(parent, prefixTxOutSig)
	// prefixTxOutContractConfirm
	prefixes[string(prefixTxOutContractConfirm)] = prefix.NewStore(parent, prefixTxOutContractConfirm)
	// prefixContractName
	prefixes[string(prefixContractName)] = prefix.NewStore(parent, prefixContractName)
	// prefixGasPrice
	prefixes[string(prefixGasPrice)] = prefix.NewStore(parent, prefixGasPrice)
	// prefixChain
	prefixes[string(prefixChain)] = prefix.NewStore(parent, prefixChain)
	// prefixTokenPrices
	prefixes[string(prefixTokenPrices)] = prefix.NewStore(parent, prefixTokenPrices)
	// prefixToken
	prefixes[string(prefixToken)] = prefix.NewStore(parent, prefixToken)
	// prefixNode
	prefixes[string(prefixNode)] = prefix.NewStore(parent, prefixNode)
	// prefixLiquidity
	prefixes[string(prefixLiquidity)] = prefix.NewStore(parent, prefixLiquidity)
	// prefixParams
	prefixes[string(prefixParams)] = prefix.NewStore(parent, prefixParams)

	return prefixes
}

///// TxRecord
func (db *defaultStorage) SaveTxRecord(hash []byte, signer string) int {
	store := db.prefixes[string(prefixTxRecord)]
	return saveTxRecord(store, hash, signer)
}

///// TxRecordProcessed
func (db *defaultStorage) ProcessTxRecord(hash []byte) {
	store := db.prefixes[string(prefixTxRecordProcessed)]
	processTxRecord(store, hash)
}

func (db *defaultStorage) IsTxRecordProcessed(hash []byte) bool {
	store := db.prefixes[string(prefixTxRecordProcessed)]
	return isTxRecordProcessed(store, hash)
}

///// Keygen

func (db *defaultStorage) SaveKeygen(msg *types.Keygen) {
	store := db.prefixes[string(prefixKeygen)]
	saveKeygen(store, msg)
}

func (db *defaultStorage) IsKeygenExisted(keyType string, index int) bool {
	store := db.prefixes[string(prefixKeygen)]
	return isKeygenExisted(store, keyType, index)
}

func (db *defaultStorage) IsKeygenAddress(keyType string, address string) bool {
	store := db.prefixes[string(prefixKeygen)]
	return isKeygenAddress(store, keyType, address)
}

func (db *defaultStorage) GetKeygenPubkey(keyType string) []byte {
	store := db.prefixes[string(prefixKeygen)]
	return getKeygenPubkey(store, keyType)
}

func (db *defaultStorage) GetAllKeygenPubkeys() map[string][]byte {
	store := db.prefixes[string(prefixKeygen)]
	return getAllKeygenPubkeys(store)
}

///// Keygen Result

func (db *defaultStorage) SaveKeygenResult(signerMsg *types.KeygenResultWithSigner) {
	store := db.prefixes[string(prefixKeygenResultWithSigner)]
	saveKeygenResult(store, signerMsg)
}

func (db *defaultStorage) GetAllKeygenResult(keygenType string, index int32) []*types.KeygenResultWithSigner {
	store := db.prefixes[string(prefixKeygenResultWithSigner)]
	return getAllKeygenResult(store, keygenType, index)
}

///// Contract
func (db *defaultStorage) SaveContract(msg *types.Contract, saveByteCode bool) {
	db.SaveContracts([]*types.Contract{msg}, saveByteCode)
}

func (db *defaultStorage) SaveContracts(msgs []*types.Contract, saveByteCode bool) {
	contractStore := db.prefixes[string(prefixContract)]
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = db.prefixes[string(prefixContractByteCode)]
	}

	saveContracts(contractStore, byteCodeStore, msgs)

	// After saving contracts, also save contract address for each contract type
	contractNameStore := db.prefixes[string(prefixContractName)]
	for _, msg := range msgs {
		saveContractAddressForName(contractNameStore, msg)
	}
}

func (db *defaultStorage) IsContractExisted(msg *types.Contract) bool {
	contractStore := db.prefixes[string(prefixContract)]
	return isContractExisted(contractStore, msg)
}

func (db *defaultStorage) GetContract(chain string, hash string, includeByteCode bool) *types.Contract {
	contractStore := db.prefixes[string(prefixContract)]
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if includeByteCode {
		byteCodeStore = db.prefixes[string(prefixContractByteCode)]
	}

	return getContract(contractStore, byteCodeStore, chain, hash)
}

func (db *defaultStorage) GetPendingContracts(chain string) []*types.Contract {
	contractStore := db.prefixes[string(prefixContract)]
	byteCodeStore := db.prefixes[string(prefixContractByteCode)]

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (db *defaultStorage) UpdateContractAddress(chain string, hash string, address string) {
	contractStore := db.prefixes[string(prefixContract)]
	updateContractAddress(contractStore, chain, hash, address)
}

func (db *defaultStorage) UpdateContractsStatus(chain string, contractHash string, status string) {
	contractStore := db.prefixes[string(prefixContract)]
	updateContractsStatus(contractStore, chain, contractHash, status)
}

func (db *defaultStorage) GetLatestContractAddressByName(chain, name string) string {
	contractNameStore := db.prefixes[string(prefixContractName)]
	return getContractAddressByName(contractNameStore, chain, name)
}

///// Contract Address
func (db *defaultStorage) CreateContractAddress(chain string, txOutHash string, address string) {
	caStore := db.prefixes[string(prefixContractAddress)]
	txOutStore := db.prefixes[string(prefixTxOut)]

	createContractAddress(caStore, txOutStore, chain, txOutHash, address)
}

func (db *defaultStorage) IsContractExistedAtAddress(chain string, address string) bool {
	caStore := db.prefixes[string(prefixContractAddress)]

	return isContractExistedAtAddress(caStore, chain, address)
}

///// TxIn
func (db *defaultStorage) SaveTxIn(msg *types.TxIn) {
	store := db.prefixes[string(prefixTxIn)]
	saveTxIn(store, msg)
}

func (db *defaultStorage) IsTxInExisted(msg *types.TxIn) bool {
	store := db.prefixes[string(prefixTxIn)]
	return isTxInExisted(store, msg)
}

///// TxOut
func (db *defaultStorage) SaveTxOut(msg *types.TxOut) {
	store := db.prefixes[string(prefixTxOut)]
	saveTxOut(store, msg)
}

func (db *defaultStorage) IsTxOutExisted(msg *types.TxOut) bool {
	store := db.prefixes[string(prefixTxOut)]
	return isTxOutExisted(store, msg)
}

func (db *defaultStorage) GetTxOut(outChain, hash string) *types.TxOut {
	store := db.prefixes[string(prefixTxOut)]
	return getTxOut(store, outChain, hash)
}

func (db *defaultStorage) GetTxOutSig(outChain, hashWithSig string) *types.TxOutSig {
	withSigStore := db.prefixes[string(prefixTxOutSig)]
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	return txOutSig
}

///// TxOutSig
func (db *defaultStorage) SaveTxOutSig(msg *types.TxOutSig) {
	store := db.prefixes[string(prefixTxOutSig)]
	saveTxOutSig(store, msg)
}

///// TxOutConfirm
func (db *defaultStorage) SaveTxOutConfirm(msg *types.TxOutContractConfirm) {
	store := db.prefixes[string(prefixTxOutContractConfirm)]
	saveTxOutConfirm(store, msg)
}

func (db *defaultStorage) IsTxOutConfirmExisted(chain, hash string) bool {
	store := db.prefixes[string(prefixTxOutContractConfirm)]
	return isTxOutConfirmExisted(store, chain, hash)
}

///// GasPrice
func (db *defaultStorage) SetGasPrice(msg *types.GasPriceMsg) {
	store := db.prefixes[string(prefixGasPrice)]
	saveGasPrice(store, msg)
}

func (db *defaultStorage) GetGasPriceRecord(chain string, height int64) *types.GasPriceRecord {
	store := db.prefixes[string(prefixGasPrice)]
	return getGasPriceRecord(store, chain, height)
}

///// Network gas price

func (db *defaultStorage) SaveChain(chain *types.Chain) {
	store := db.prefixes[string(prefixChain)]

	saveChain(store, chain)
}

func (db *defaultStorage) GetChain(chain string) *types.Chain {
	store := db.prefixes[string(prefixChain)]
	return getChain(store, chain)
}

func (db *defaultStorage) GetAllChains() map[string]*types.Chain {
	store := db.prefixes[string(prefixChain)]
	return getAllChains(store)
}

///// Token Prices

func (db *defaultStorage) SetTokenPrices(blockHeight uint64, msg *types.UpdateTokenPrice) {
	store := db.prefixes[string(prefixTokenPrices)]
	setTokenPrices(store, blockHeight, msg)
}

func (db *defaultStorage) GetAllTokenPricesRecord() map[string]*types.TokenPriceRecord {
	store := db.prefixes[string(prefixTokenPrices)]
	return getAllTokenPrices(store)
}

///// Calculated token prices

func (db *defaultStorage) SetTokens(prices map[string]*types.Token) {
	store := db.prefixes[string(prefixToken)]
	setTokens(store, prices)
}

func (db *defaultStorage) GetTokens(tokenIds []string) map[string]*types.Token {
	store := db.prefixes[string(prefixToken)]
	return getTokens(store, tokenIds)
}

func (db *defaultStorage) GetAllTokens() map[string]*types.Token {
	store := db.prefixes[string(prefixToken)]
	return getAllTokens(store)
}

///// Nodes

func (db *defaultStorage) SaveNode(node *types.Node) {
	store := db.prefixes[string(prefixNode)]
	saveNode(store, node)
}

func (db *defaultStorage) LoadValidators() []*types.Node {
	store := db.prefixes[string(prefixNode)]
	return loadValidators(store)
}

///// Liquidities

func (db *defaultStorage) SetLiquidities(liquids map[string]*types.Liquidity) {
	store := db.prefixes[string(prefixLiquidity)]
	setLiquidities(store, liquids)
}

func (db *defaultStorage) GetLiquidity(chain string) *types.Liquidity {
	store := db.prefixes[string(prefixLiquidity)]
	return getLiquidity(store, chain)
}

func (db *defaultStorage) GetAllLiquidities() map[string]*types.Liquidity {
	store := db.prefixes[string(prefixLiquidity)]
	return getAllLiquidities(store)
}

///// Params
func (db *defaultStorage) SaveParams(params *types.Params) {
	store := db.prefixes[string(prefixParams)]
	saveParams(store, params)
}

func (db *defaultStorage) GetParams() *types.Params {
	store := db.prefixes[string(prefixParams)]
	return getParams(store)
}

///// Debug

func (db *defaultStorage) getStoreFromName(name string) cstypes.KVStore {
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
func (db *defaultStorage) PrintStore(name string) {
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

func (db *defaultStorage) PrintStoreKeys(name string) {
	log.Info("======== DEBUGGING PrintStoreKeys")
	log.Info("Printing ALL values in store ", name)
	store := db.getStoreFromName(name)
	if store != nil {
		printStoreKeys(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}
