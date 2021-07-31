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

var (
	PREFIX_RECORDED_CHAIN            = []byte{0x01}
	PREFIX_OBSERVED_TX               = []byte{0x02}
	PREFIX_OBSERVED_TX_VALIDATOR_SET = []byte{0x03}
	PREFIX_PENDING_OBSERVED_TX       = []byte{0x04}
	PREFIX_PROCESSED_OBSERVED_TX     = []byte{0x05}
	PREFIX_PUBLICK_KEY_BYTES         = []byte{0x06}

	// List of transactions that have enough observation and pending for output.
	// KEY_PENDING_OBSERVED_TX = "pending_observed_tx_%s_%d_%s" // chain - block height - tx hash

	// List of transactions that have been processed.
	// KEY_PROCESSED_OBSERVED_TX = "processed_observed_tx_%s_%d_%s" // chain - block height - tx hash

	// KEY_PUBLICK_KEY_BYTES = "public_key_bytes_%s"

	// List of on memory keys. These data are not persisted into kvstore.
	// List of contracts that need to be deployed to a chain.
	KEY_CONTRACT_QUEUE     = "contract_queue_%s_%s"     // chain
	KEY_DEPLOYING_CONTRACT = "deploying_contract_%s_%s" // chain - contract hash
)

type deployContractWrapper struct {
	data         []byte
	createdBlock int64 // Sisu block when the contract is created.
	// id of the designated validator that is supposed to post the tx out to the Sisu chain.
	designatedValidator string
}

// Data structure that wraps around pending tx outs.
type pendingTxOutWrapper struct {
	valAddr string
	txInt   []byte
	txOut   []byte
}

type Keeper struct {
	storeKey sdk.StoreKey

	// List of contracts that waits to be deployed.
	contractQueue map[string]string

	// List of contracts that are being deployed.
	deployingContracts map[string]*deployContractWrapper

	// List of contracts that have been deployed. (This should be saved in a KV store?)
	deployedContracts map[string]*deployContractWrapper

	// A map that remembers what transaction is assigned to which validators.
	assignedValidators map[int64]map[string]*pendingTxOutWrapper // blockHeight -> tx bytes (as string) -> validator address
}

func NewKeeper(storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey:           storeKey,
		contractQueue:      make(map[string]string),
		deployingContracts: make(map[string]*deployContractWrapper),
		deployedContracts:  make(map[string]*deployContractWrapper),
		assignedValidators: make(map[int64]map[string]*pendingTxOutWrapper),
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

// @Deprecated. TODO: Remove
// This updates the set of validators that attest to observe a specific tx (identified by its hash)
// on a specific chain.
func (k *Keeper) UpdateObservedTxCount(ctx sdk.Context, msg *types.ObservedTx, signer string) (int, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_OBSERVED_TX_VALIDATOR_SET)

	key := k.getKey(msg.Chain, msg.BlockHeight, msg.TxHash)
	bz := store.Get(key)

	var validators map[string]bool
	if bz == nil || len(bz) == 0 {
		validators = make(map[string]bool)
	} else {
		err := json.Unmarshal(bz, &validators)
		if err != nil {
			utils.LogError("Cannot unmarshall validator sets")
			return 0, err
		}
	}

	if !validators[signer] {
		validators[signer] = true
		bz, err := json.Marshal(validators)
		if err != nil {
			utils.LogError("Cannot marshal validator set")
			return 0, err
		}

		store.Set(key, bz)
	}

	return len(validators), nil
}

// Checks if a tx has been processed or in the pending list. Returns true if either of the condition
// meets.
func (k *Keeper) IsObservedTxPendingOrProcessed(ctx sdk.Context, msg *types.ObservedTx) bool {
	// Check Pending list.
	if k.IsObservedTxPending(ctx, msg) {
		return true
	}

	// Check processed list.
	if k.IsObservedTxProcessed(ctx, msg) {
		return true
	}

	return false
}

func (k *Keeper) IsObservedTxPending(ctx sdk.Context, msg *types.ObservedTx) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PENDING_OBSERVED_TX)
	key := k.getKey(msg.Chain, msg.BlockHeight, msg.TxHash)
	bz := store.Get(key)
	if bz != nil {
		return true
	}

	return false
}

// Returns true if an observed tx has been processed.
func (k *Keeper) IsObservedTxProcessed(ctx sdk.Context, msg *types.ObservedTx) bool {
	// Check processed list.
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PROCESSED_OBSERVED_TX)
	key := k.getKey(msg.Chain, msg.BlockHeight, msg.TxHash)
	bz := store.Get(key)
	if bz != nil {
		return true
	}

	return false
}

func (k *Keeper) AddObservedTxToPending(ctx sdk.Context, msg *types.ObservedTx) {
	if k.IsObservedTxProcessed(ctx, msg) {
		// Transaction has been processed, there is no need to add it to pending.
		utils.LogVerbose("Transaction has been processed, there is no need to add it to pending.")
		return
	}

	store := ctx.KVStore(k.storeKey)
	key := k.getKey(msg.Chain, msg.BlockHeight, msg.TxHash)

	bz, err := msg.Marshal()
	if err != nil {
		utils.LogError("Cannot marshal observed tx, err = ", err)
		return
	}

	store.Set(key, bz)
}

func (k *Keeper) GetAndClearObservedTxPendingList(ctx sdk.Context) []*types.ObservedTx {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PENDING_OBSERVED_TX)
	itr := store.Iterator(nil, nil)
	keys := make([][]byte, 0)

	txs := make([]*types.ObservedTx, 0)
	for ; itr.Valid(); itr.Next() {
		bz := itr.Value()
		msg := &types.ObservedTx{}
		err := msg.Unmarshal(bz)
		if err != nil {
			utils.LogError("Cannot unmarshall message with key ", string(itr.Key()))
			continue
		}
		txs = append(txs, msg)
		keys = append(keys, itr.Key())
	}
	itr.Close()

	// Delete the list.
	for _, key := range keys {
		store.Delete(key)
	}

	return txs
}

func (k *Keeper) SavePubKey(ctx sdk.Context, chain string, keyBytes []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), PREFIX_PUBLICK_KEY_BYTES)
	store.Set([]byte(chain), keyBytes)
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

// Saves an assignment of a validator for a particular out tx.
func (k *Keeper) AddAssignedValForOutTx(blockHeight int64, txBytes []byte, valAddr string) {
	m := k.assignedValidators[blockHeight]
	if m == nil {
		m = make(map[string]*pendingTxOutWrapper)
	}
	m[string(txBytes)] = &pendingTxOutWrapper{
		valAddr: valAddr,
		txOut:   txBytes,
	}

	k.assignedValidators[blockHeight] = m
}
