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
	PrefixKeygen    = []byte("keygen")
	PrefixContract  = []byte("contract")
	PrefixTxOut     = []byte("txout")
	PrefixMempoolTx = []byte("mempool")
)

var _ KVDatabase = (*KVStore)(nil)

type KVDatabase interface {
	InsertTxOuts(ctx sdk.Context, txs []*tsstypes.TxOutEntity)
	GetTxOutWithHash(ctx sdk.Context, chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity
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
		bz, err := json.Marshal(tx)
		if err != nil {
			log.Error(err)
			return
		}

		key := getTxOutKey(tx.InChain, tx.HashWithoutSig)
		store.Set(key, bz)
	}
}

func (s *KVStore) GetTxOutWithHash(ctx sdk.Context, chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	if !isHashWithSig {
		txOutBz := store.Get(getTxOutKey(chain, hash))
		if len(txOutBz) == 0 {
			log.Warnf("cannot find txout for chain(%s), hash(%s), isHashWithSig(%v)", chain, hash, isHashWithSig)
			return nil
		}

		txOut := &tsstypes.TxOutEntity{}
		if err := json.Unmarshal(txOutBz, txOut); err != nil {
			log.Error(err)
			return nil
		}

		return txOut
	}

	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		txOut := &tsstypes.TxOutEntity{}
		if err := json.Unmarshal(iter.Value(), txOut); err != nil {
			log.Error(err)
			continue
		}

		if strings.EqualFold(txOut.HashWithSig, hash) {
			return txOut
		}
	}

	return nil
}

func (s *KVStore) UpdateTxOutStatus(ctx sdk.Context, chain, hash string, status tsstypes.TxOutStatus, isHashWithSig bool) error {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	if !isHashWithSig {
		txOut := &tsstypes.TxOutEntity{}
		key := getTxOutKey(chain, hash)
		txOutBz := store.Get(key)
		if len(txOutBz) == 0 {
			log.Warnf("cannot find txout for chain(%s), hash(%s), isHashWithSig(%s)", chain, hash, isHashWithSig)
			return nil
		}

		if err := json.Unmarshal(txOutBz, txOut); err != nil {
			log.Error("json unmarshal error: ", err)
			return err
		}

		txOut.Status = string(status)
		_ = saveRecord(store, key, txOut)
		return nil
	}

	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		txOut := &tsstypes.TxOutEntity{}
		if err := json.Unmarshal(iter.Value(), txOut); err != nil {
			log.Error(err)
			continue
		}

		if !strings.EqualFold(txOut.HashWithSig, hash) {
			continue
		}

		txOut.Status = string(status)
		_ = saveRecord(store, iter.Key(), txOut)
		return nil
	}

	return nil
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
