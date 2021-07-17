package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

const (
	KEY_RECORDED_CHAIN = "recored_chain"
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
