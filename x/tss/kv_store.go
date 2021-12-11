package tss

import (
	"encoding/json"
	"strings"

	"github.com/sisu-network/cosmos-sdk/store/prefix"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
)

const KVSeparator = "__"

var (
	PrefixTxOut = []byte("txout")
)

var _ KVDatabase = (*KVStore)(nil)

type KVDatabase interface {
	InsertTxOuts(ctx sdk.Context, txs []*tsstypes.TxOutEntity)
}

type KVStore struct {
	storeKey sdk.StoreKey
}

func NewDefaultKVStore(storeKey sdk.StoreKey) *KVStore {
	return &KVStore{storeKey: storeKey}
}

func (s *KVStore) InsertTxOuts(ctx sdk.Context, txs []*tsstypes.TxOutEntity) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	for _, tx := range txs {
		key := getTxOutKey(tx.InChain, tx.HashWithoutSig)
		_ = saveRecord(store, key, tx)
	}
}

func getTxOutKey(chain, hashWithoutSig string) []byte {
	return []byte(strings.Join([]string{chain, hashWithoutSig}, KVSeparator))
}

func saveRecord(store prefix.Store, key []byte, entity interface{}) error {
	bz, err := json.Marshal(entity)
	if err != nil {
		log.Error(err)
		return err
	}

	store.Set(key, bz)
	return nil
}
