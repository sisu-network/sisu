package keeper

import (
	"fmt"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

var (
	prefixKeygenProposal   = []byte{0x01}
	prefixKeygen           = []byte{0x02}
	prefixKeygenResult     = []byte{0x03}
	prefixContract         = []byte{0x04}
	prefixContractByteCode = []byte{0x05}
	prefixTxIn             = []byte{0x06}
	prefixTxOut            = []byte{0x07}
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

func (k *DefaultKeeper) getKeygenKey(keyType string, index int) []byte {
	// keyType + id
	return []byte(fmt.Sprintf("%s__%d", keyType, index))
}

func (k *DefaultKeeper) getKeygenResultKey(keyType string, index int, from string) []byte {
	// keyType
	return []byte(fmt.Sprintf("%s__%d__%s", keyType, index, from))
}

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.Keygen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenKey(msg.KeyType, int(msg.Index))

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("SaveKeygenProposal: cannot marshal keygen proposal, err = ", err)
	}
	store.Set(key, bz)
}

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenKey(keyType, index)

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)
	key := k.getKeygenResultKey(signerMsg.Keygen.KeyType, int(signerMsg.Keygen.Index), signerMsg.Data.From)

	bz, err := signerMsg.Data.Marshal()
	if err != nil {
		log.Error("SaveKeygenResult: Cannot marshal KeygenResult message, err = ", err)
		return
	}

	store.Set(key, bz)
}

// Keygen is considered successful if at least there is at least 1 successful KeygenReslut in the
// KVStore.
func (k *DefaultKeeper) IsKeygenResultSuccess(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) bool {
	msg := signerMsg.Keygen
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResult)

	begin := []byte(fmt.Sprintf("%s__%d__", msg.KeyType, int(msg.Index)))
	end := []byte(fmt.Sprintf("%s__%d__~", msg.KeyType, int(msg.Index)))

	iter := store.Iterator(begin, end)
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		msg := &types.KeygenResult{}
		err := msg.Unmarshal(bz)
		if err != nil {
			log.Error("Cannot unmarshal keygen result")
			continue
		}

		if msg.Result == types.KeygenResult_SUCCESS {
			return true
		}
	}

	return false
}

func (k *DefaultKeeper) GetAllPubKeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)

	iter := store.Iterator(nil, nil)
	ret := make(map[string][]byte)
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		msg := &types.Keygen{}
		err := msg.Unmarshal(bz)
		if err != nil {
			log.Error("cannot unmarshal KeygenResult message, err = ", err)
			continue
		}

		ret[string(iter.Key())] = msg.PubKeyBytes
	}

	return ret
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
