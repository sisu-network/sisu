package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

const (
	OBSERVED_TX_CACHE_SIZE = 2500
)

// TODO: clean up this list
var (
	PREFIX_OBSERVED_TX      = []byte{0x01}
	PREFIX_PUBLIC_KEY_BYTES = []byte{0x02}
	PREFIX_ETH_KEY_ADDRESS  = []byte{0x03}

	// List of on memory keys. These data are not persisted into kvstore.
	// List of contracts that need to be deployed to a chain.
	KEY_ETH_KEY_ADDRESS    = "eth_key_address_%s"       // chain
	KEY_CONTRACT_QUEUE     = "contract_queue_%s_%s"     // chain
	KEY_DEPLOYING_CONTRACT = "deploying_contract_%s_%s" // chain - contract hash
)

type Keeper interface {
	SaveObservedTx(ctx sdk.Context, tx *tssTypes.ObservedTx)
	GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte
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

func (k *DefaultKeeper) getKey(chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	chain = strings.Replace(chain, "_", "*", -1)
	return []byte(fmt.Sprintf("%s_%d_%s", chain, height, hash))
}

func (k *DefaultKeeper) SaveObservedTx(ctx sdk.Context, tx *tssTypes.ObservedTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(tx.Chain, tx.BlockHeight, tx.TxHash)

	store.Set(key, tx.Serialized)
}

func (k *DefaultKeeper) GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(chain, blockHeight, hash)

	return store.Get(key)
}

func (k *DefaultKeeper) SavePubKey(ctx sdk.Context, keyType string, keyBytes []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	store.Set([]byte(keyType), keyBytes)
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	ret := make(map[string][]byte)
	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		ret[string(iter.Key())] = iter.Value()
	}

	log.Info("Length of all pubkeys: ", len(ret))
	return ret
}

func (k *DefaultKeeper) SaveEthKeyAddrs(ctx sdk.Context, chain string, keyAddrs map[string]bool) error {
	key := fmt.Sprintf(KEY_ETH_KEY_ADDRESS, chain)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_ETH_KEY_ADDRESS)

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
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_ETH_KEY_ADDRESS)

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
