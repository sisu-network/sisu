package keeper

import (
	"encoding/json"
	"fmt"
	"strings"

	tssTypes "github.com/sisu-network/sisu/x/tss/types"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	OBSERVED_TX_CACHE_SIZE = 2500
)

// TODO: clean up this list
var (
	PREFIX_RECORDED_CHAIN            = []byte{0x01}
	PREFIX_OBSERVED_TX               = []byte{0x02}
	PREFIX_OBSERVED_TX_VALIDATOR_SET = []byte{0x03}
	PREFIX_PENDING_OBSERVED_TX       = []byte{0x04}
	PREFIX_PROCESSED_OBSERVED_TX     = []byte{0x05}
	PREFIX_PUBLIC_KEY_BYTES          = []byte{0x06}
	PREFIX_PENDING_KEYGEN_TX         = []byte{0x07}
	PREFIX_ETH_KEY_ADDRESS           = []byte{0x08}
	PREFIX_TX_OUT                    = []byte{0x09}

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

// Get a list of chains that this node supported and have generated private key through TSS.
func (k *Keeper) GetRecordedChainsOnSisu(ctx sdk.Context) (*types.ChainsInfo, error) {
	// store := ctx.KVStore(k.storeKey)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_RECORDED_CHAIN)
	bz := store.Get([]byte(PREFIX_RECORDED_CHAIN))

	chainsInfo := &types.ChainsInfo{}
	err := chainsInfo.Unmarshal(bz)

	return chainsInfo, err
}

// Saves a list of chains that this node supports.
func (k *Keeper) SetChainsInfo(ctx sdk.Context, chainsInfo *types.ChainsInfo) error {
	utils.LogInfo("Keeper: Saving chain info into KV store", chainsInfo.Chains)

	store := ctx.KVStore(k.storeKey)
	bz, err := chainsInfo.Marshal()
	if err != nil {
		utils.LogError("Cannot set chains info. Err = ", err)
		return err
	}

	store.Set([]byte(PREFIX_RECORDED_CHAIN), bz)
	return nil
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

func (k *Keeper) IsContractDeployingOrDeployed(ctx sdk.Context, chain string, hash string) bool {
	hash = strings.ToLower(hash)

	deployingKey := fmt.Sprintf(KEY_DEPLOYING_CONTRACT, chain, hash)

	if _, ok := k.contractQueue[deployingKey]; ok {
		return true
	}

	return false
}

// Save a contract with a specific hash into a queue for later deployment.
func (k *Keeper) EnqueueContract(ctx sdk.Context, chain string, hash string, content string) {
	key := fmt.Sprintf(KEY_CONTRACT_QUEUE, chain, hash)
	k.contractQueue[key] = content
}

// Get all contract hashes in a pending queue for a particular chain.
func (k *Keeper) GetContractQueueHashes(ctx sdk.Context, chain string) []string {
	hashes := make([]string, 0)
	for key, value := range k.contractQueue {
		if len(key) <= len("contract_queue_") {
			utils.LogError("Invalid key:", key)
			continue
		}

		suffix := key[len("contract_queue_"):]
		index := strings.Index(suffix, "_")
		if index <= 0 {
			utils.LogError("Invalid suffix:", suffix)
			continue
		}

		c := suffix[0:index]
		if c != chain {
			continue
		}

		hashes = append(hashes, value)
	}

	return hashes
}

// Delete all the contracts in the queue
func (k *Keeper) ClearContractQueue(ctx sdk.Context) {
	k.contractQueue = make(map[string]string)
}

func (k *Keeper) DequeueContract(ctx sdk.Context, chain string, hash string) {
	key := fmt.Sprintf(KEY_CONTRACT_QUEUE, chain, hash)
	delete(k.contractQueue, key)
}

// Adds a list of hashes
func (k *Keeper) AddDeployingContract(ctx sdk.Context, chain string, hash string, txBytes []byte, createdBlock int64) {
	key := fmt.Sprintf(KEY_DEPLOYING_CONTRACT, chain, hash)
	k.deployingContracts[key] = &deployContractWrapper{
		data:         txBytes,
		createdBlock: createdBlock,
	}
}

func (k *Keeper) SaveEthKeyAddrs(ctx sdk.Context, chain string, keyAddrs map[string]bool) {
	key := fmt.Sprintf(KEY_ETH_KEY_ADDRESS, chain)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_ETH_KEY_ADDRESS)

	bz, err := json.Marshal(keyAddrs)
	if err != nil {
		utils.LogError("cannot marshal key addrs, err =", err)
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
			utils.LogError("cannot unmarshal value with key", iter.Key())
			continue
		}
		m[string(iter.Key())] = m2
	}

	return m
}

func (k *Keeper) SaveTxOut(ctx sdk.Context, txOut *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_TX_OUT)
	bz, err := txOut.Marshal()
	if err != nil {
		utils.LogError("cannot marshal tx out, err =", err)
		return
	}
	store.Set([]byte(txOut.GetHash()), bz)
}

func (k *Keeper) GetTxOut(ctx sdk.Context, hash string) (*types.TxOut, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_TX_OUT)
	bz := store.Get([]byte(hash))
	if bz == nil {
		return nil, fmt.Errorf("cannot find tx with hash %s", hash)
	}

	tx := &types.TxOut{}
	err := tx.Unmarshal(bz)
	return tx, err
}
