package keeper

import (
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

// go:generate mockgen -source x/tss/keeper/keeper.go -destination=tests/mock/tss/keeper.go -package=mock
type Keeper interface {
	// Keygen
	SaveKeygen(ctx sdk.Context, msg *types.Keygen)
	IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool
	GetAllPubKeys(ctx sdk.Context) map[string][]byte

	// Keygen Result
	SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner)
	IsKeygenResultSuccess(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) bool

	// Contracts
	SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(ctx sdk.Context, msg *types.Contract) bool

	GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract
	UpdateContractsStatus(ctx sdk.Context, msgs []*types.Contract, status string)

	// TxIn
	SaveTxIn(ctx sdk.Context, msg *types.TxIn)
	IsTxInExisted(ctx sdk.Context, msg *types.TxIn) bool

	// TxOut
	SaveTxOut(ctx sdk.Context, msg *types.TxOut)
	IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool
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

func (k *DefaultKeeper) getTxInKey(chain string, height int64, hash string) []byte {
	// chain, height, hash
	return []byte(fmt.Sprintf("%s__%d__%s", chain, height, hash))
}

func (k *DefaultKeeper) getTxOutKey(inChain string, outChain string, outHash string) []byte {
	// inChain, outChain, height, hash
	return []byte(fmt.Sprintf("%s__%s__%s", inChain, outChain, outHash))
}

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.Keygen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	saveKeygen(store, msg)
}

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return isKeygenExisted(store, keyType, index)
}

func (k *DefaultKeeper) SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)
	saveKeygenResult(store, signerMsg)
}

// Keygen is considered successful if at least there is at least 1 successful KeygenReslut in the
// KVStore.
func (k *DefaultKeeper) IsKeygenResultSuccess(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)
	return isKeygenResultSuccess(store, signerMsg)
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return getAllPubKeys(store)
}

func (k *DefaultKeeper) SaveTxIn(ctx sdk.Context, msg *types.TxIn) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	key := k.getTxInKey(msg.Chain, msg.BlockHeight, msg.TxHash)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal TxIn")
		return
	}

	store.Set(key, bz)
}

///// Contracts

func (k *DefaultKeeper) SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	saveContracts(contractStore, byteCodeStore, msgs, saveByteCode)
}

func (k *DefaultKeeper) IsContractExisted(ctx sdk.Context, msg *types.Contract) bool {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	return isContractExisted(contractStore, msg)
}

func (k *DefaultKeeper) GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (k *DefaultKeeper) UpdateContractsStatus(ctx sdk.Context, msgs []*types.Contract, status string) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	updateContractsStatus(contractStore, msgs, status)
}

///// TxIn

func (k *DefaultKeeper) IsTxInExisted(ctx sdk.Context, tx *types.TxIn) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	key := k.getTxInKey(tx.GetChain(), tx.GetBlockHeight(), tx.GetTxHash())

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveTxOut(ctx sdk.Context, msg *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	key := k.getTxOutKey(msg.InChain, msg.OutChain, msg.GetHash())

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal tx out")
		return
	}

	store.Set(key, bz)
}

func (k *DefaultKeeper) IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	key := k.getTxOutKey(msg.InChain, msg.OutChain, msg.GetHash())

	return store.Get(key) != nil
}
