package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var _ Keeper = (*DefaultKeeper)(nil)

type Keeper interface {
	// Debug
	PrintStore(ctx sdk.Context, name string)
	PrintStoreKeys(ctx sdk.Context, name string)

	// TxRecord
	SaveTxRecord(ctx sdk.Context, hash []byte, signer string) int

	// TxRecordProcessed
	ProcessTxRecord(ctx sdk.Context, hash []byte)
	IsTxRecordProcessed(ctx sdk.Context, hash []byte) bool

	// Keygen
	SaveKeygen(ctx sdk.Context, msg *types.Keygen)
	IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool
	IsKeygenAddress(ctx sdk.Context, keyType string, address string) bool
	GetKeygenPubkey(ctx sdk.Context, keyType string) []byte
	GetAllKeygenPubkeys(ctx sdk.Context) map[string][]byte

	// Keygen Result
	SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner)
	GetAllKeygenResult(ctx sdk.Context, keygenType string, index int32) []*types.KeygenResultWithSigner

	// TxOut
	SaveTxOut(ctx sdk.Context, msg *types.TxOut)
	GetTxOut(ctx sdk.Context, outChain, hash string) *types.TxOut

	// TxOutSig
	SaveTxOutSig(ctx sdk.Context, msg *types.TxOutSig)
	GetTxOutSig(ctx sdk.Context, outChain, hashWithSig string) *types.TxOutSig

	// Gas Price Record
	SetGasPrice(ctx sdk.Context, msg *types.GasPriceMsg)
	GetGasPriceRecord(ctx sdk.Context, height int64) *types.GasPriceRecord

	// Chain
	SaveChain(ctx sdk.Context, chain *types.Chain)
	GetChain(ctx sdk.Context, chain string) *types.Chain
	GetAllChains(ctx sdk.Context) map[string]*types.Chain

	// Token Price
	SetTokenPrices(ctx sdk.Context, blockHeight uint64, msg *types.UpdateTokenPrice)
	GetAllTokenPricesRecord(ctx sdk.Context) map[string]*types.TokenPriceRecords

	// Token
	SetTokens(ctx sdk.Context, tokens map[string]*types.Token)
	GetTokens(ctx sdk.Context, tokens []string) map[string]*types.Token
	GetAllTokens(ctx sdk.Context) map[string]*types.Token

	// Nodes
	SaveNode(ctx sdk.Context, node *types.Node)
	LoadValidators(ctx sdk.Context) []*types.Node

	// Vaults
	SetVaults(ctx sdk.Context, vaults []*types.Vault)
	GetVault(ctx sdk.Context, chain string, token string) *types.Vault
	GetAllVaultsForChain(ctx sdk.Context, chain string) []*types.Vault

	// MPC Address
	SetMpcAddress(ctx sdk.Context, chain string, address string)
	GetMpcAddress(ctx sdk.Context, chain string) string

	// Params
	SaveParams(ctx sdk.Context, params *types.Params)
	GetParams(ctx sdk.Context) *types.Params

	// Gateway checkpoint
	AddGatewayCheckPoint(ctx sdk.Context, checkPoint *types.GatewayCheckPoint)
	GetGatewayCheckPoint(ctx sdk.Context, chain string) *types.GatewayCheckPoint
	GetAllGatewayCheckPoints(ctx sdk.Context) map[string]*types.GatewayCheckPoint

	// Command Queue
	SetCommandQueue(ctx sdk.Context, chain string, commands []*types.Command)
	GetCommandQueue(ctx sdk.Context, chain string) []*types.Command

	// Transfer
	AddTransfers(ctx sdk.Context, transfers []*types.Transfer)
	GetTransfer(ctx sdk.Context, id string) *types.Transfer
	GetTransfers(ctx sdk.Context, ids []string) []*types.Transfer

	// Transfer Queue
	SetTransferQueue(ctx sdk.Context, chain string, transfers []*types.Transfer)
	GetTransferQueue(ctx sdk.Context, chain string) []*types.Transfer

	// TxOutQueue
	SetTxOutQueue(ctx sdk.Context, chain string, txOuts []*types.TxOut)
	GetTxOutQueue(ctx sdk.Context, chain string) []*types.TxOut

	// PendingTxOut
	SetPendingTxOutInfo(ctx sdk.Context, chain string, txOut *types.PendingTxOutInfo)
	GetPendingTxOutInfo(ctx sdk.Context, chain string) *types.PendingTxOutInfo
}

type DefaultKeeper struct {
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) Keeper {
	keeper := &DefaultKeeper{
		storeKey: storeKey,
	}

	return keeper
}

///// TxRecord
func (k *DefaultKeeper) SaveTxRecord(ctx sdk.Context, hash []byte, signer string) int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecord)
	return saveTxRecord(store, hash, signer)
}

///// TxRecordProcessed
func (k *DefaultKeeper) ProcessTxRecord(ctx sdk.Context, hash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecordProcessed)
	processTxRecord(store, hash)
}

func (k *DefaultKeeper) IsTxRecordProcessed(ctx sdk.Context, hash []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecordProcessed)
	return isTxRecordProcessed(store, hash)
}

///// Keygen

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.Keygen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	saveKeygen(store, msg)
}

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return isKeygenExisted(store, keyType, index)
}

func (k *DefaultKeeper) IsKeygenAddress(ctx sdk.Context, keyType string, address string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return isKeygenAddress(store, keyType, address)
}

func (k *DefaultKeeper) GetKeygenPubkey(ctx sdk.Context, keyType string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return getKeygenPubkey(store, keyType)
}

func (k *DefaultKeeper) GetAllKeygenPubkeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return getAllKeygenPubkeys(store)
}

///// Keygen Result

func (k *DefaultKeeper) SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResultWithSigner)
	saveKeygenResult(store, signerMsg)
}

func (k *DefaultKeeper) GetAllKeygenResult(ctx sdk.Context, keygenType string, index int32) []*types.KeygenResultWithSigner {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResultWithSigner)
	return getAllKeygenResult(store, keygenType, index)
}

///// TxOut
func (k *DefaultKeeper) SaveTxOut(ctx sdk.Context, msg *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	saveTxOut(store, msg)
}

func (k *DefaultKeeper) GetTxOut(ctx sdk.Context, outChain, hash string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	return getTxOut(store, outChain, hash)
}

func (k *DefaultKeeper) GetTxOutSig(ctx sdk.Context, outChain, hashWithSig string) *types.TxOutSig {
	withSigStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutSig)
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	return txOutSig
}

///// TxOutSig
func (k *DefaultKeeper) SaveTxOutSig(ctx sdk.Context, msg *types.TxOutSig) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutSig)
	saveTxOutSig(store, msg)
}

///// GasPrice
func (k *DefaultKeeper) SetGasPrice(ctx sdk.Context, msg *types.GasPriceMsg) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGasPrice)
	saveGasPrice(store, msg)
}

func (k *DefaultKeeper) GetGasPriceRecord(ctx sdk.Context, height int64) *types.GasPriceRecord {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGasPrice)
	return getGasPriceRecord(store, height)
}

///// Network gas price

func (k *DefaultKeeper) SaveChain(ctx sdk.Context, chain *types.Chain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)

	saveChain(store, chain)
}

func (k *DefaultKeeper) GetChain(ctx sdk.Context, chain string) *types.Chain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)
	return getChain(store, chain)
}

func (k *DefaultKeeper) GetAllChains(ctx sdk.Context) map[string]*types.Chain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)
	return getAllChains(store)
}

///// Token Prices

func (k *DefaultKeeper) SetTokenPrices(ctx sdk.Context, blockHeight uint64, msg *types.UpdateTokenPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTokenPrices)
	setTokenPrices(store, blockHeight, msg)
}

func (k *DefaultKeeper) GetAllTokenPricesRecord(ctx sdk.Context) map[string]*types.TokenPriceRecords {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTokenPrices)
	return getAllTokenPrices(store)
}

func (k *DefaultKeeper) SetTokens(ctx sdk.Context, tokens map[string]*types.Token) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	setTokens(store, tokens)
}

func (k *DefaultKeeper) GetTokens(ctx sdk.Context, tokenIds []string) map[string]*types.Token {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	return getTokens(store, tokenIds)
}

func (k *DefaultKeeper) GetAllTokens(ctx sdk.Context) map[string]*types.Token {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	return getAllTokens(store)
}

///// Nodes

func (k *DefaultKeeper) SaveNode(ctx sdk.Context, node *types.Node) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	saveNode(store, node)
}

func (k *DefaultKeeper) LoadValidators(ctx sdk.Context) []*types.Node {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	return loadValidators(store)
}

///// Vaults

func (k *DefaultKeeper) SetVaults(ctx sdk.Context, vaults []*types.Vault) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixVault)
	setVaults(store, vaults)
}

func (k *DefaultKeeper) GetVault(ctx sdk.Context, chain string, token string) *types.Vault {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixVault)
	return getVault(store, chain, token)
}

func (k *DefaultKeeper) GetAllVaultsForChain(ctx sdk.Context, chain string) []*types.Vault {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixVault)
	return getAllVaultsForChain(store, chain)
}

///// Vaults
func (k *DefaultKeeper) SetMpcAddress(ctx sdk.Context, chain string, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcAddress)
	setMpcAddress(store, chain, address)
}

func (k *DefaultKeeper) GetMpcAddress(ctx sdk.Context, chain string) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcAddress)
	return getMpcAddress(store, chain)
}

///// Params
func (k *DefaultKeeper) SaveParams(ctx sdk.Context, params *types.Params) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixParams)
	saveParams(store, params)
}

func (k *DefaultKeeper) GetParams(ctx sdk.Context) *types.Params {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixParams)
	return getParams(store)
}

///// Gateway Checkpoint
func (k *DefaultKeeper) AddGatewayCheckPoint(ctx sdk.Context, checkPoint *types.GatewayCheckPoint) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGatewayCheckPoint)
	addCheckPoint(store, checkPoint)
}

func (k *DefaultKeeper) GetGatewayCheckPoint(ctx sdk.Context, chain string) *types.GatewayCheckPoint {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGatewayCheckPoint)
	return getCheckPoint(store, chain)
}

func (k *DefaultKeeper) GetAllGatewayCheckPoints(ctx sdk.Context) map[string]*types.GatewayCheckPoint {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGatewayCheckPoint)
	return getAllGatewayCheckPoints(store)
}

///// Command Queue
func (k *DefaultKeeper) SetCommandQueue(ctx sdk.Context, chain string, commands []*types.Command) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixCommandQueue)
	setCommandQueue(store, chain, commands)
}

func (k *DefaultKeeper) GetCommandQueue(ctx sdk.Context, chain string) []*types.Command {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixCommandQueue)
	return getCommandQueue(store, chain)
}

///// Transfer
func (k *DefaultKeeper) AddTransfers(ctx sdk.Context, transfers []*types.Transfer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	addTransfers(store, transfers)
}

func (k *DefaultKeeper) GetTransfer(ctx sdk.Context, id string) *types.Transfer {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	transfers := getTransfers(store, []string{id})
	if len(transfers) == 0 {
		return nil
	}

	return transfers[0]
}

func (k *DefaultKeeper) GetTransfers(ctx sdk.Context, ids []string) []*types.Transfer {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	return getTransfers(store, ids)
}

///// Transfer Queue
func (k *DefaultKeeper) SetTransferQueue(ctx sdk.Context, chain string, transfers []*types.Transfer) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransferQueue)
	setTranferQueue(store, chain, transfers)
}

func (k *DefaultKeeper) GetTransferQueue(ctx sdk.Context, chain string) []*types.Transfer {
	transferStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	queueStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransferQueue)
	return getTransferQueue(queueStore, transferStore, chain)
}

///// TxOutQueue
func (k *DefaultKeeper) SetTxOutQueue(ctx sdk.Context, chain string, txOuts []*types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutQueue)
	setTxOutQueue(store, chain, txOuts)
}

func (k *DefaultKeeper) GetTxOutQueue(ctx sdk.Context, chain string) []*types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutQueue)
	return getTxOutQueue(store, chain)
}

///// PendingTxOut
func (k *DefaultKeeper) SetPendingTxOutInfo(ctx sdk.Context, chain string, txOutInfo *types.PendingTxOutInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixPendingTxOut)
	setPendingTxOut(store, chain, txOutInfo)
}

func (k *DefaultKeeper) GetPendingTxOutInfo(ctx sdk.Context, chain string) *types.PendingTxOutInfo {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixPendingTxOut)
	return getPendingTxOutInfo(store, chain)
}

///// Debug

func (k *DefaultKeeper) getStoreFromName(ctx sdk.Context, name string) cstypes.KVStore {
	var store cstypes.KVStore
	switch name {
	case "keygen":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	case "txOut":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	}

	return store
}

// PrintStore is a debug function
func (k *DefaultKeeper) PrintStore(ctx sdk.Context, name string) {
	log.Info("======== DEBUGGING PrintStore")
	log.Info("Printing ALL values in store ", name)

	store := k.getStoreFromName(ctx, name)
	if store != nil {
		printStore(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}

func (k *DefaultKeeper) PrintStoreKeys(ctx sdk.Context, name string) {
	log.Info("======== DEBUGGING PrintStoreKeys")
	log.Info("Printing ALL values in store ", name)
	store := k.getStoreFromName(ctx, name)
	if store != nil {
		printStoreKeys(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}
