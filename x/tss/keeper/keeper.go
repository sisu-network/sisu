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
	prefixContract         = []byte{0x03}
	prefixContractByteCode = []byte{0x04}
	prefixTxIn             = []byte{0x05}
	prefixTxOut            = []byte{0x06}
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

	// Contracts
	SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool)
	GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract

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

func (k *DefaultKeeper) getKeygenProposalKey(keyType string, createdBlock int64) []byte {
	// keyType + id
	return []byte(fmt.Sprintf("%s__%d", keyType, createdBlock))
}

func (k *DefaultKeeper) getKeygenResultKey(keyType string) []byte {
	// keyType
	return []byte(fmt.Sprintf("%s", keyType))
}

func (k *DefaultKeeper) SaveKeygenProposal(ctx sdk.Context, msg *types.KeygenProposal) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenProposalKey(msg.KeyType, msg.CreatedBlock)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("SaveKeygenProposal: cannot marshal keygen proposal, err = ", err)
	}
	store.Set(key, bz)
}

func (k *DefaultKeeper) IsKeygenProposalExisted(ctx sdk.Context, msg *types.KeygenProposal) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenProposal)
	key := k.getKeygenProposalKey(msg.KeyType, msg.CreatedBlock)

	return store.Get(key) != nil
}

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.KeygenResult) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	key := k.getKeygenResultKey(msg.Keygen.KeyType)

	bz, err := msg.Keygen.Marshal()
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

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, msg *types.KeygenResult) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	key := k.getKeygenResultKey(msg.Keygen.KeyType)

	return store.Get(key) != nil
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
