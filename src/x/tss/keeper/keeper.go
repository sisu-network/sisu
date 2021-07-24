package keeper

import (
	"encoding/json"
	"fmt"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	KEY_RECORDED_CHAIN = "recored_chain"

	// Set of validators that attest this transaction.
	KEY_OBSERVED_TX_VALIDATOR_SET = "observed_tx_%s_%d_%s" // chain - block height - tx hash

	// List of transactions that have enough observation and pending for output.
	KEY_PENDING_OBSERVED_TX = "pending_observed_tx_%s_%d_%s" // chain - block height - tx hash

	// List of transactions that have been processed.
	KEY_PROCESSED_OBSERVED_TX = "processed_observed_tx_%s_%d_%s" // chain - block height - tx hash

	KEY_PUBLICK_KEY_BYTES = "public_key_bytes_%s"
)

type Keeper struct {
	storeKey sdk.StoreKey

	// TODO: Use on memory cache to speed up read operation for both pending & processed tx list.
	pendingObservedTxLock   *sync.RWMutex
	processedObservedTxLock *sync.RWMutex
}

func NewKeeper(storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey:                storeKey,
		pendingObservedTxLock:   &sync.RWMutex{},
		processedObservedTxLock: &sync.RWMutex{},
	}
}

// Get a list of chains that this node supported and have generated private key through TSS.
func (k *Keeper) GetRecordedChainsOnSisu(ctx sdk.Context) (*types.ChainsInfo, error) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get([]byte(KEY_RECORDED_CHAIN))

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

	store.Set([]byte(KEY_RECORDED_CHAIN), bz)
	return nil
}

// This updates the set of validators that attest to observe a specific tx (identified by its hash)
// on a specific chain.
func (k *Keeper) UpdateObservedTxCount(ctx sdk.Context, msg *types.ObservedTx, signer string) (int, error) {
	store := ctx.KVStore(k.storeKey)

	key := []byte(fmt.Sprintf(KEY_OBSERVED_TX_VALIDATOR_SET, msg.Chain, msg.BlockHeight, msg.TxHash))
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
	k.pendingObservedTxLock.RLock()
	defer k.pendingObservedTxLock.RUnlock()

	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf(KEY_PENDING_OBSERVED_TX, msg.Chain, msg.BlockHeight, msg.TxHash))
	bz := store.Get(key)
	if bz != nil {
		return true
	}

	return false
}

// Returns true if an observed tx has been processed.
func (k *Keeper) IsObservedTxProcessed(ctx sdk.Context, msg *types.ObservedTx) bool {
	k.processedObservedTxLock.RLock()
	defer k.processedObservedTxLock.RUnlock()

	// Check processed list.
	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf(KEY_PROCESSED_OBSERVED_TX, msg.Chain, msg.BlockHeight, msg.TxHash))
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

	k.pendingObservedTxLock.Lock()
	defer k.pendingObservedTxLock.Unlock()

	store := ctx.KVStore(k.storeKey)
	key := []byte(fmt.Sprintf(KEY_PENDING_OBSERVED_TX, msg.Chain, msg.BlockHeight, msg.TxHash))

	bz, err := msg.Marshal()
	if err != nil {
		utils.LogError("Cannot marshal observed tx, err = ", err)
		return
	}
	store.Set(key, bz)
}

func (k *Keeper) GetObservedTxPendingList(ctx sdk.Context) {
	k.pendingObservedTxLock.RLock()
	defer k.pendingObservedTxLock.RUnlock()

	store := ctx.KVStore(k.storeKey)
	itr := store.Iterator(nil, nil)

	for ; itr.Valid(); itr.Next() {
		// TODO: Complete this.
	}
}

func (k *Keeper) SavePubKey(ctx sdk.Context, chain string, keyBytes []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set([]byte(fmt.Sprintf(KEY_PUBLICK_KEY_BYTES, chain)), keyBytes)
}
