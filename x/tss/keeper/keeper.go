package keeper

import (
	"github.com/sisu-network/cosmos-sdk/store/prefix"
	cstypes "github.com/sisu-network/cosmos-sdk/store/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

var _ Keeper = (*DefaultKeeper)(nil)

// go:generate mockgen -source x/tss/keeper/keeper.go -destination=tests/mock/tss/keeper.go -package=mock
type Keeper interface {
	// Debug
	PrintStore(ctx sdk.Context, name string)

	// TxVote
	SaveTxVotes(ctx sdk.Context, hash []byte, val string) int

	// Keygen
	SaveKeygen(ctx sdk.Context, msg *types.Keygen)
	IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool
	IsKeygenAddress(ctx sdk.Context, keyType string, address string) bool

	// Keygen Result
	SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner)
	IsKeygenResultSuccess(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner, self string) bool

	// Contracts
	SaveContract(ctx sdk.Context, msg *types.Contract, saveByteCode bool)
	SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(ctx sdk.Context, msg *types.Contract) bool
	GetContract(ctx sdk.Context, chain string, hash string, includeByteCode bool) *types.Contract
	GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract
	UpdateContractAddress(ctx sdk.Context, chain string, hash string, address string)
	UpdateContractsStatus(ctx sdk.Context, chain string, contractHash string, status string)

	// Contract Address
	CreateContractAddress(ctx sdk.Context, chain string, txOutHash string, address string)
	IsContractExistedAtAddress(ctx sdk.Context, chain string, address string) bool
	GetLatestContractAddressByName(ctx sdk.Context, chain, name string) string

	// TxIn
	SaveTxIn(ctx sdk.Context, msg *types.TxIn)
	IsTxInExisted(ctx sdk.Context, msg *types.TxIn) bool

	// TxOut
	SaveTxOut(ctx sdk.Context, msg *types.TxOut)
	IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool
	GetTxOut(ctx sdk.Context, chain, hash string) *types.TxOut

	// TxOutConfirm
	SaveTxOutConfirm(ctx sdk.Context, msg *types.TxOutConfirm)
	IsTxOutConfirmExisted(ctx sdk.Context, outChain, hash string) bool
}

type DefaultKeeper struct {
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) *DefaultKeeper {
	keeper := &DefaultKeeper{
		storeKey: storeKey,
	}

	return keeper
}

///// TxVote
func (k *DefaultKeeper) SaveTxVotes(ctx sdk.Context, hash []byte, val string) int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxVotes)
	return saveTxVotes(store, hash, val)
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

///// Keygen Result

func (k *DefaultKeeper) SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)
	saveKeygenResult(store, signerMsg)
}

// Keygen is considered successful if at least there is at least 1 successful KeygenReslut in the
// KVStore.
func (k *DefaultKeeper) IsKeygenResultSuccess(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner, self string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)
	return isKeygenResultSuccess(store, signerMsg.Keygen.KeyType, signerMsg.Keygen.Index, self)
}

///// Contracts

func (k *DefaultKeeper) SaveContract(ctx sdk.Context, msg *types.Contract, saveByteCode bool) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)
	}

	saveContract(contractStore, byteCodeStore, msg)

	contractNameStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractName)
	saveContractAddressForName(contractNameStore, msg)
}

func (k *DefaultKeeper) SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)
	}

	saveContracts(contractStore, byteCodeStore, msgs)

	contractNameStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractName)
	for _, msg := range msgs {
		saveContractAddressForName(contractNameStore, msg)
	}
}

func (k *DefaultKeeper) IsContractExisted(ctx sdk.Context, msg *types.Contract) bool {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	return isContractExisted(contractStore, msg)
}

func (k *DefaultKeeper) GetContract(ctx sdk.Context, chain string, hash string, includeByteCode bool) *types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if includeByteCode {
		byteCodeStore = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)
	}

	return getContract(contractStore, byteCodeStore, chain, hash)
}

func (k *DefaultKeeper) GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (k *DefaultKeeper) UpdateContractAddress(ctx sdk.Context, chain string, hash string, address string) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	updateContractAddress(contractStore, chain, hash, address)
}

func (k *DefaultKeeper) UpdateContractsStatus(ctx sdk.Context, chain string, contractHash string, status string) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	updateContractsStatus(contractStore, chain, contractHash, status)
}

///// Contract Address

func (k *DefaultKeeper) CreateContractAddress(ctx sdk.Context, chain string, txOutHash string, address string) {
	caStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractAddress)
	txOutStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)

	createContractAddress(caStore, txOutStore, chain, txOutHash, address)
}

func (k *DefaultKeeper) IsContractExistedAtAddress(ctx sdk.Context, chain string, address string) bool {
	caStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractAddress)
	return isContractExistedAtAddress(caStore, chain, address)
}

func (k *DefaultKeeper) GetLatestContractAddressByName(ctx sdk.Context, chain, name string) string {
	contractNameStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractName)
	return getContractAddressByName(contractNameStore, chain, name)
}

///// TxIn
func (k *DefaultKeeper) SaveTxIn(ctx sdk.Context, msg *types.TxIn) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	saveTxIn(store, msg)
}

func (k *DefaultKeeper) IsTxInExisted(ctx sdk.Context, msg *types.TxIn) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	return isTxInExisted(store, msg)
}

///// TxOut

func (k *DefaultKeeper) SaveTxOut(ctx sdk.Context, msg *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	saveTxOut(store, msg)
}

func (k *DefaultKeeper) IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	return isTxOutExisted(store, msg)
}

func (k *DefaultKeeper) GetTxOut(ctx sdk.Context, chain, hash string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	return getTxOut(store, chain, hash)
}

///// TxOutConfirm
func (k *DefaultKeeper) SaveTxOutConfirm(ctx sdk.Context, msg *types.TxOutConfirm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutConfirm)
	saveTxOutConfirm(store, msg)
}

func (k *DefaultKeeper) IsTxOutConfirmExisted(ctx sdk.Context, outChain, hash string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutConfirm)
	return isTxOutConfirmExisted(store, outChain, hash)
}

///// Debug

// PrintStore is a debug function
func (k *DefaultKeeper) PrintStore(ctx sdk.Context, name string) {
	log.Info("======== DEBUGGING")
	log.Info("Printing ALL values in store ", name)
	var store cstypes.KVStore
	switch name {
	case "keygen":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	case "contract":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	}

	printStore(store)
	log.Info("======== END OF DEBUGGING")
}
