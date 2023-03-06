package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	prefixTxRecord               = []byte{0x01} // Vote for a tx by different nodes
	prefixTxRecordProcessed      = []byte{0x02}
	prefixKeygen                 = []byte{0x03}
	prefixKeygenResultWithSigner = []byte{0x04}
	prefixTxOut                  = []byte{0x05}
	prefixGasPrice               = []byte{0x06}
	prefixChain                  = []byte{0x07}
	prefixToken                  = []byte{0x08}
	prefixNode                   = []byte{0x0A}
	prefixVault                  = []byte{0x0B}
	prefixParams                 = []byte{0x0C}
	prefixMpcAddress             = []byte{0x0D}
	prefixTxOutQueue             = []byte{0x10}
	prefixCommandQueue           = []byte{0x11}
	prefixTransfer               = []byte{0x12}
	prefixTransferQueue          = []byte{0x13}
	prefixSignerNonce            = []byte{0x15}
	prefixBlockHeight            = []byte{0x16}
	prefixVoteResult             = []byte{0x19}
	prefixProposedTxOut          = []byte{0x1A}
	prefixMpcPublicKey           = []byte{0x1B}
	prefixExpirationBlock        = []byte{0x1C}
	prefixFinalizedTxOut         = []byte{0x1D}
	prefixKeySignRetryCount      = []byte{0x1E}
	prefixTransferCounter        = []byte{0x1F}
	prefixFailedTransfer         = []byte{0x20}
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

	// Chain
	SaveChain(ctx sdk.Context, chain *types.Chain)
	GetChain(ctx sdk.Context, chain string) *types.Chain
	GetAllChains(ctx sdk.Context) map[string]*types.Chain

	// Token
	SetTokens(ctx sdk.Context, tokens map[string]*types.Token)
	GetToken(ctx sdk.Context, token string) *types.Token
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

	// MPC Public Key
	SetMpcPublicKey(ctx sdk.Context, chain string, pubkey []byte)
	GetMpcPublicKey(ctx sdk.Context, chain string) []byte

	// Params
	SaveParams(ctx sdk.Context, params *types.Params)
	GetParams(ctx sdk.Context) *types.Params

	// Reported Mpc Nonce by each signer.
	SetSignerNonce(ctx sdk.Context, chain string, signer string, nonce uint64)
	GetAllSignerNonces(ctx sdk.Context, chain string) []uint64

	// Command Queue
	SetCommandQueue(ctx sdk.Context, chain string, commands []*types.Command)
	GetCommandQueue(ctx sdk.Context, chain string) []*types.Command

	// Transfer
	AddTransfers(ctx sdk.Context, transfers []*types.TransferDetails)
	IncTransferRetryNum(ctx sdk.Context, transferId string)
	GetTransfer(ctx sdk.Context, id string) *types.TransferDetails
	GetTransfers(ctx sdk.Context, ids []string) []*types.TransferDetails

	// Transfer Queue
	SetTransferQueue(ctx sdk.Context, chain string, transfers []*types.TransferDetails)
	GetTransferQueue(ctx sdk.Context, chain string) []*types.TransferDetails

	// Failed Transfer
	AddFailedTransfer(ctx sdk.Context, transferId string)
	RemoveFailedTransfer(ctx sdk.Context, transferId string)
	HasFailedTransfer(ctx sdk.Context, transferId string) bool

	// TxOutQueue
	SetTxOutQueue(ctx sdk.Context, chain string, txOuts []*types.TxOut)
	GetTxOutQueue(ctx sdk.Context, chain string) []*types.TxOut

	// Max Block height that all nodes observed (Not all chains need this property)
	SetBlockHeight(ctx sdk.Context, chain string, height int64, hash string)
	GetBlockHeight(ctx sdk.Context, chain string) *types.BlockHeight

	// Vote Result
	AddVoteResult(ctx sdk.Context, txOutId string, signer string, result types.VoteResult)
	GetVoteResults(ctx sdk.Context, txOutId string) map[string]types.VoteResult

	// Proposed TxOut
	AddProposedTxOut(ctx sdk.Context, msg *types.TxOut)
	GetProposedTxOut(ctx sdk.Context, id string) *types.TxOut

	// Finalized TxOut
	SetFinalizedTxOut(ctx sdk.Context, txOut *types.TxOut)
	GetFinalizedTxOut(ctx sdk.Context, id string) *types.TxOut
	GetFinalizedTxOutFromHash(ctx sdk.Context, chain, hash string) *types.TxOut

	// Expiration block
	SetExpirationBlock(ctx sdk.Context, opType string, key string, height int64)
	GetExpirationBlock(ctx sdk.Context, opType string, key string) int64
	RemoveExpirationBlock(ctx sdk.Context, opType string, id string)

	// Keysign retry count.
	SetKeySignRetryCount(ctx sdk.Context, txOutId string, count int)
	GetKeySignRetryCount(ctx sdk.Context, txOutId string) int
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

///// Token

func (k *DefaultKeeper) SetTokens(ctx sdk.Context, tokens map[string]*types.Token) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	setTokens(store, tokens)
}

func (k *DefaultKeeper) GetToken(ctx sdk.Context, tokenId string) *types.Token {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	return getToken(store, tokenId)
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

///// MPC Address
func (k *DefaultKeeper) SetMpcAddress(ctx sdk.Context, chain string, address string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcAddress)
	setKeyValue(store, []byte(chain), []byte(address))
}

func (k *DefaultKeeper) GetMpcAddress(ctx sdk.Context, chain string) string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcAddress)
	return string(getKeyValue(store, []byte(chain)))
}

///// MPC Public Key
func (k *DefaultKeeper) SetMpcPublicKey(ctx sdk.Context, chain string, pubkey []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcPublicKey)
	setKeyValue(store, []byte(chain), pubkey)
}

func (k *DefaultKeeper) GetMpcPublicKey(ctx sdk.Context, chain string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixMpcPublicKey)
	return getKeyValue(store, []byte(chain))
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

///// Signer nonce

func (k *DefaultKeeper) SetSignerNonce(ctx sdk.Context, chain string, signer string, nonce uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixSignerNonce)
	setSignerNonce(store, chain, signer, nonce)
}

func (k *DefaultKeeper) GetAllSignerNonces(ctx sdk.Context, chain string) []uint64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixSignerNonce)
	return getSignerNonces(store, chain)
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
func (k *DefaultKeeper) AddTransfers(ctx sdk.Context, transfers []*types.TransferDetails) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	addTransfers(store, transfers)
}

func (k *DefaultKeeper) IncTransferRetryNum(ctx sdk.Context, transferId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	transfers := getTransfers(store, []string{transferId})
	if len(transfers) == 0 {
		return
	}

	transfers[0].RetryNum++
	addTransfers(store, transfers)
}

func (k *DefaultKeeper) GetTransfer(ctx sdk.Context, id string) *types.TransferDetails {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	transfers := getTransfers(store, []string{id})
	if len(transfers) == 0 {
		return nil
	}

	return transfers[0]
}

func (k *DefaultKeeper) GetTransfers(ctx sdk.Context, ids []string) []*types.TransferDetails {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	return getTransfers(store, ids)
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

///// Block Height
func (k *DefaultKeeper) SetBlockHeight(ctx sdk.Context, chain string, height int64, hash string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixBlockHeight)
	setBlockHeight(store, chain, &types.BlockHeight{
		Height: height,
		Hash:   hash,
	})
}

func (k *DefaultKeeper) GetBlockHeight(ctx sdk.Context, chain string) *types.BlockHeight {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixBlockHeight)
	return getBlockHeight(store, chain)
}

///// Transfer Queue
func (k *DefaultKeeper) SetTransferQueue(ctx sdk.Context, chain string, transfers []*types.TransferDetails) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransferQueue)
	setTransferQueue(store, chain, transfers)
}

func (k *DefaultKeeper) GetTransferQueue(ctx sdk.Context, chain string) []*types.TransferDetails {
	transferStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	queueStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransferQueue)

	return getTransferQueue(queueStore, transferStore, chain)
}

///// Failed Transfer Queue
func (k *DefaultKeeper) AddFailedTransfer(ctx sdk.Context, transferId string) {
	transferStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTransfer)
	failedStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFailedTransfer)
	addFailedTransfer(failedStore, transferStore, transferId)
}

func (k *DefaultKeeper) RemoveFailedTransfer(ctx sdk.Context, transferId string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFailedTransfer)
	removeFailedTransfer(store, transferId)
}

func (k *DefaultKeeper) HasFailedTransfer(ctx sdk.Context, transferId string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFailedTransfer)
	return hasFailedTransfer(store, transferId)
}

///// Vote Result
func (k *DefaultKeeper) AddVoteResult(ctx sdk.Context, txOutId string, signer string,
	result types.VoteResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixVoteResult)
	addVoteResult(store, txOutId, signer, result)
}

func (k *DefaultKeeper) GetVoteResults(ctx sdk.Context, txOutId string) map[string]types.VoteResult {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixVoteResult)
	return getVoteResults(store, txOutId)
}

///// Proposed TxOut

func (k *DefaultKeeper) AddProposedTxOut(ctx sdk.Context, txOut *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixProposedTxOut)
	addProposedTxOut(store, txOut)
}

func (k *DefaultKeeper) GetProposedTxOut(ctx sdk.Context, id string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixProposedTxOut)
	return getProposedTxOut(store, id)
}

///// Finalized TxOut
func (k *DefaultKeeper) SetFinalizedTxOut(ctx sdk.Context, txOut *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFinalizedTxOut)
	currentTxOut := k.GetFinalizedTxOutFromHash(ctx, txOut.Content.OutChain, txOut.Content.OutHash)
	if currentTxOut != nil {
		store.Delete([]byte(currentTxOut.GetId()))
	}

	bz, err := txOut.Marshal()
	if err != nil {
		log.Errorf("Failed to marshal txout, err = %s", err)
		return
	}

	setKeyValue(store, []byte(txOut.GetId()), bz)
}

func (k *DefaultKeeper) GetFinalizedTxOut(ctx sdk.Context, id string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFinalizedTxOut)

	bz := store.Get([]byte(id))
	if bz == nil {
		return nil
	}

	txOut := &types.TxOut{}
	err := txOut.Unmarshal(bz)
	if err != nil {
		log.Errorf("Failed to unmarshal txout, err = %s", err)
		return nil
	}

	return txOut
}

func (k *DefaultKeeper) GetFinalizedTxOutFromHash(ctx sdk.Context, chain, hash string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixFinalizedTxOut)

	begin := []byte(fmt.Sprintf("%s__%s__", chain, hash))
	end := []byte(fmt.Sprintf("%s__%s__~", chain, hash))

	counter := 0
	var id []byte
	for iter := store.Iterator(begin, end); iter.Valid(); iter.Next() {
		counter++
		id = iter.Key()
	}

	if counter == 0 {
		return nil
	}

	if counter > 1 {
		log.Errorf("Have %d txouts with outchain = %s, outhash = %s", counter, chain, hash)
		return nil
	}

	return k.GetFinalizedTxOut(ctx, string(id))
}

///// Expiration block
func (k *DefaultKeeper) SetExpirationBlock(ctx sdk.Context, objectType string, id string,
	height int64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixExpirationBlock)
	setExpirationBlock(store, objectType, id, height)
}

func (k *DefaultKeeper) GetExpirationBlock(ctx sdk.Context, objectType string, id string) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixExpirationBlock)
	return getExpirationBlock(store, objectType, id)
}

func (k *DefaultKeeper) RemoveExpirationBlock(ctx sdk.Context, objectType string, id string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixExpirationBlock)
	removeExpirationBlock(store, objectType, id)
}

///// Keysign retry count.
func (k *DefaultKeeper) SetKeySignRetryCount(ctx sdk.Context, txOutId string, count int) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeySignRetryCount)
	store.Set([]byte(txOutId), []byte{byte(count)})
}

func (k *DefaultKeeper) GetKeySignRetryCount(ctx sdk.Context, txOutId string) int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeySignRetryCount)
	bz := store.Get([]byte(txOutId))
	if bz == nil || len(bz) != 1 {
		return 0
	}

	return int(bz[0])
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
