package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/sisu-network/lib/log"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
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
	KEY_CONTRACT_QUEUE     = "contract_queue_%s_%s"     // chain
	KEY_DEPLOYING_CONTRACT = "deploying_contract_%s_%s" // chain - contract hash
	KEY_ETH_KEY_ADDRESS    = "eth_key_address_%s"       // chain
)

type deployContractWrapper struct {
	data         []byte
	createdBlock int64 // Sisu block when the contract is created.
	// id of the designated validator that is supposed to post the tx out to the Sisu chain.
	designatedValidator string
}

type Keeper struct {
	storeKey sdk.StoreKey

	// List of contracts that waits to be deployed.
	contractQueue map[string]string

	// List of contracts that are being deployed.
	deployingContracts map[string]*deployContractWrapper

	// List of contracts that have been deployed. (This should be saved in a KV store?)
	deployedContracts map[string]*deployContractWrapper
}

func NewKeeper(storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey:           storeKey,
		contractQueue:      make(map[string]string),
		deployingContracts: make(map[string]*deployContractWrapper),
		deployedContracts:  make(map[string]*deployContractWrapper),
	}
}

func (k *Keeper) getKey(chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	chain = strings.Replace(chain, "_", "*", -1)
	return []byte(fmt.Sprintf("%s_%d_%s", chain, height, hash))
}

func (k *Keeper) SaveObservedTx(ctx sdk.Context, tx *tssTypes.ObservedTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(tx.Chain, tx.BlockHeight, tx.TxHash)

	store.Set(key, tx.Serialized)
}

func (k *Keeper) GetObservedTx(ctx sdk.Context, chain string, blockHeight int64, hash string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX)
	key := k.getKey(chain, blockHeight, hash)

	return store.Get(key)
}

func (k *Keeper) SavePubKey(ctx sdk.Context, chain string, keyBytes []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	store.Set([]byte(chain), keyBytes)
}

func (k *Keeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLIC_KEY_BYTES)
	iter := store.Iterator(nil, nil)
	ret := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		ret[string(iter.Key())] = iter.Value()
	}

	return ret
}

func (k *Keeper) SaveEthKeyAddrs(ctx sdk.Context, chain string, keyAddrs map[string]bool) {
	key := fmt.Sprintf(KEY_ETH_KEY_ADDRESS, chain)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_ETH_KEY_ADDRESS)

	bz, err := json.Marshal(keyAddrs)
	if err != nil {
		log.Error("cannot marshal key addrs, err =", err)
		return
	}

	store.Set([]byte(key), bz)
}

func (k *Keeper) GetAllEthKeyAddrs(ctx sdk.Context) map[string]map[string]bool {
	m := make(map[string]map[string]bool)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_ETH_KEY_ADDRESS)

	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		m2 := make(map[string]bool)
		err := json.Unmarshal(iter.Value(), &m2)
		if err != nil {
			log.Error("cannot unmarshal value with key", iter.Key())
			continue
		}
		m[string(iter.Key())] = m2
	}

	return m
}
