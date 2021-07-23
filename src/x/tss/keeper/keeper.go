package keeper

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	KEY_RECORDED_CHAIN            = "recored_chain"
	KEY_OBSERVED_TX_VALIDATOR_SET = "observed_tx_%s_%d_%s" // chain - block height - tx hash
)

type Keeper struct {
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) *Keeper {
	return &Keeper{
		storeKey: storeKey,
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
		return err
	}

	store.Set([]byte(KEY_RECORDED_CHAIN), bz)
	return nil
}

// This updates the set of validators that attest to observe a specific tx (identified by its hash)
// on a specific chain.
func (k *Keeper) UpdateObservedTxCount(ctx sdk.Context, msg *types.ObservedTx, signer string) {
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
			return
		}
	}

	if !validators[signer] {
		validators[signer] = true
		bz, err := json.Marshal(validators)
		if err != nil {
			utils.LogError("Cannot marshal validator set")
			return
		}

		store.Set(key, bz)
	}
}
