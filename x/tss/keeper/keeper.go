package keeper

import (
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
)

var (
	prefixKeygenProposal = []byte{0x01}
	prefixKeygen         = []byte{0x02}
	prefixObservedTx     = []byte{0x03}
	prefixTxOut          = []byte{0x04}
)

// go:generate mockgen -source x/tss/keeper/keeper.go -destination=tests/mock/tss/keeper.go -package=mock
type Keeper interface {
	// KeygenProposal
	SaveKeygenProposal(ctx sdk.Context, msg *types.KeygenProposal)
	IsKeygenProposalExisted(ctx sdk.Context, msg *types.KeygenProposal) bool

	// Keygen
	SaveKeygen(ctx sdk.Context, msg *types.KeygenResult)
	GetAllPubKeys(ctx sdk.Context) map[string][]byte
	IsKeygenExisted(ctx sdk.Context, msg *types.KeygenResult) bool

	// Observed Tx
	SaveObservedTx(ctx sdk.Context, msg *types.ObservedTx)
	IsObservedTxExisted(ctx sdk.Context, msg *types.ObservedTx) bool

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

func (k *DefaultKeeper) getObservedTxKey(chain string, height int64, hash string) []byte {
	// chain, height, hash
	return []byte(fmt.Sprintf("%s__%d__%s", chain, height, hash))
}

func (k *DefaultKeeper) getTxOutKey(inChain string, outChain string, height int64, hash string) []byte {
	// inChain, outChain, height, hash
	return []byte(fmt.Sprintf("%s__%s__%d__%s", inChain, outChain, height, hash))
}

func (k *DefaultKeeper) getKeygenProposalKey(keyType string, id string) []byte {
	// keyType + id
	return []byte(fmt.Sprintf("%s__%s", keyType, id))
}

func (k *DefaultKeeper) getKeygenResultKey(keyType string) []byte {
	// keyType
	return []byte(fmt.Sprintf("%s", keyType))
}

func (k *DefaultKeeper) SaveKeygenProposal(ctx sdk.Context, msg *types.KeygenProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenProposalKey(msg.KeyType, msg.Id)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("SaveKeygenProposal: cannot marshal keygen proposal, err = ", err)
	}
	store.Set(key, bz)
}

func (k *DefaultKeeper) IsKeygenProposalExisted(ctx sdk.Context, msg *types.KeygenProposal) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenProposalKey(msg.KeyType, msg.Id)

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.KeygenResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	key := k.getKeygenResultKey(msg.Keygen.KeyType)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal KeygenResult message, err = ", err)
		return
	}

	store.Set(key, bz)
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)

	iter := store.Iterator(nil, nil)
	ret := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		msg := &types.KeygenResult{}
		err := msg.Unmarshal(bz)
		if err != nil {
			log.Error("cannot unmarshal KeygenResult message, err = ", err)
			continue
		}

		ret[string(iter.Key())] = msg.Keygen.PubKeyBytes
	}

	return ret
}

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, msg *types.KeygenResult) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	key := k.getKeygenResultKey(msg.Keygen.KeyType)

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveObservedTx(ctx sdk.Context, msg *types.ObservedTx) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixObservedTx)
	key := k.getObservedTxKey(msg.Chain, msg.BlockHeight, msg.TxHash)

	bz := msg.SerializeWithoutSigner()
	store.Set(key, bz)
}

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
