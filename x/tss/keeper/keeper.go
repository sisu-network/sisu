package keeper

import (
	"encoding/json"
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	OBSERVED_TX_CACHE_SIZE = 2500
)

// TODO: clean up this list
var (
	prefixObservedTx      = []byte{0x01}
	prefixPublicKeyBytes  = []byte{0x02}
	prefixEthKeyAddresses = []byte{0x03}
	prefixTxOut           = []byte{0x04}
	prefixKeygenProposal  = []byte{0x05}

	// Deprecated
	KEY_ETH_KEY_ADDRESS = "eth_key_address_%s" // chain
)

// go:generate mockgen -source x/tss/keeper/keeper.go -destination=tests/mock/tss/keeper.go -package=mock
type Keeper interface {
	SaveKeygenProposal(ctx sdk.Context, msg *types.KeygenProposal)

	// Observed Tx
	SaveObservedTx(ctx sdk.Context, msg *types.ObservedTx)
	IsObservedTxExisted(ctx sdk.Context, msg *types.ObservedTx) bool

	// TxOut
	SaveTxOut(ctx sdk.Context, msg *types.TxOut)
	IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool

	SavePubKey(ctx sdk.Context, chain string, keyBytes []byte)
	GetAllPubKeys(ctx sdk.Context) map[string][]byte
	SaveEthKeyAddrs(ctx sdk.Context, chain string, keyAddrs map[string]bool) error
	GetAllEthKeyAddrs(ctx sdk.Context) (map[string]map[string]bool, error)
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

func (k *DefaultKeeper) getObservedTxKey(chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	return []byte(fmt.Sprintf("%s__%d__%s", chain, height, hash))
}

func (k *DefaultKeeper) getTxOutKey(inChain string, outChain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	return []byte(fmt.Sprintf("%s__%s__%d__%s", inChain, outChain, height, hash))
}

func (k *DefaultKeeper) getKeygenProposalKey(keyType string, id string, createdBlock int64) []byte {
	return []byte(fmt.Sprintf("%s__%s__%d", keyType, id, createdBlock))
}

func (k *DefaultKeeper) SaveKeygenProposal(ctx sdk.Context, msg *types.KeygenProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenProposalKey(msg.KeyType, msg.Id, msg.CreatedBlock)

	// TODO: Fix this
	store.Set(key, []byte(""))
}

func (k *DefaultKeeper) SaveObservedTx(ctx sdk.Context, msg *types.ObservedTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixObservedTx)
	key := k.getObservedTxKey(msg.Chain, msg.BlockHeight, msg.TxHash)

	bz := msg.SerializeWithoutSigner()
	store.Set(key, bz)
}

// func (k *DefaultKeeper) GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte {
func (k *DefaultKeeper) IsObservedTxExisted(ctx sdk.Context, tx *types.ObservedTx) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixObservedTx)
	key := k.getObservedTxKey(tx.GetChain(), tx.GetBlockHeight(), tx.GetTxHash())

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveTxOut(ctx sdk.Context, msg *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	outHash := utils.KeccakHash32(string(msg.OutBytes))
	key := k.getTxOutKey(msg.InChain, msg.OutChain, msg.InBlockHeight, outHash)

	bz := msg.SerializeWithoutSigner()
	store.Set(key, bz)
}

func (k *DefaultKeeper) IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	outHash := utils.KeccakHash32(string(msg.OutBytes))
	key := k.getTxOutKey(msg.InChain, msg.OutChain, msg.InBlockHeight, outHash)

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SavePubKey(ctx sdk.Context, keyType string, keyBytes []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixPublicKeyBytes)
	store.Set([]byte(keyType), keyBytes)
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixPublicKeyBytes)
	iter := store.Iterator(nil, nil)
	ret := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		ret[string(iter.Key())] = iter.Value()
	}

	return ret
}

func (k *DefaultKeeper) SaveEthKeyAddrs(ctx sdk.Context, chain string, keyAddrs map[string]bool) error {
	key := fmt.Sprintf(KEY_ETH_KEY_ADDRESS, chain)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixEthKeyAddresses)

	bz, err := json.Marshal(keyAddrs)
	if err != nil {
		log.Error("cannot marshal key addrs, err =", err)
		return err
	}

	store.Set([]byte(key), bz)
	return nil
}

func (k *DefaultKeeper) GetAllEthKeyAddrs(ctx sdk.Context) (map[string]map[string]bool, error) {
	m := make(map[string]map[string]bool)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixEthKeyAddresses)

	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		m2 := make(map[string]bool)
		err := json.Unmarshal(iter.Value(), &m2)
		if err != nil {
			log.Error("cannot unmarshal value with key", iter.Key())
			return nil, err
		}
		m[string(iter.Key())] = m2
	}

	return m, nil
}
