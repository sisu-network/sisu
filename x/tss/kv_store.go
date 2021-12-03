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

// KVDatabase has the same interface with db.Database with 1 more parameter: sdk.Context
type KVDatabase interface {
	Init() error
	Close() error

	// Keygen
	CreateKeygen(ctx sdk.Context, chain string) error
	UpdateKeygenAddress(ctx sdk.Context, chain, address string, pubKey []byte)

	IsKeyExisted(ctx sdk.Context, chain string) bool
	IsChainKeyAddress(ctx sdk.Context, chain, address string) bool
	GetPubKey(ctx sdk.Context, chain string) []byte
	UpdateKeygenStatus(ctx sdk.Context, chain, status string)
	GetKeygenStatus(ctx sdk.Context, chain string) (string, error)

	// Contracts
	InsertContracts(ctx sdk.Context, contracts []*tsstypes.ContractEntity)
	GetPendingDeployContracts(ctx sdk.Context, chain string) []*tsstypes.ContractEntity
	GetContractFromAddress(ctx sdk.Context, chain, address string) *tsstypes.ContractEntity
	GetContractFromHash(ctx sdk.Context, chain, hash string) *tsstypes.ContractEntity
	UpdateContractsStatus(ctx sdk.Context, contracts []*tsstypes.ContractEntity, status string) error
	UpdateContractAddress(ctx sdk.Context, chain, hash, address string)

	// Txout
	InsertTxOuts(ctx sdk.Context, txs []*tsstypes.TxOutEntity)
	GetTxOutWithHash(ctx sdk.Context, chain string, hash string, isHashWithSig bool) *tsstypes.TxOutEntity
	IsContractDeployTx(ctx sdk.Context, chain string, hashWithoutSig string) (bool, error)
	UpdateTxOutSig(ctx sdk.Context, chain, hashWithoutSign, hashWithSig string, sig []byte) error
	UpdateTxOutStatus(ctx sdk.Context, chain, hash string, status tsstypes.TxOutStatus, isHashWithSig bool) error

	// Mempool tx
	InsertMempoolTxHash(ctx sdk.Context, hash string, blockHeight int64)
	MempoolTxExisted(ctx sdk.Context, hash string) bool
	MempoolTxExistedRange(ctx sdk.Context, hash string, minBlock int64, maxBlock int64) bool
}

type KVStore struct {
	storeKey sdk.StoreKey
}

func NewDefaultKVStore(store sdk.StoreKey) *KVStore {
	return &KVStore{storeKey: store}
}

func (s *KVStore) Init() error {
	return nil
}

func (s *KVStore) Close() error {
	return nil
}

func (s *KVStore) CreateKeygen(ctx sdk.Context, chain string) error {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixKeygen)

	kgEntity := &tsstypes.KeygenEntity{
		Chain: chain,
	}

}

func (s *KVStore) UpdateKeygenAddress(ctx sdk.Context, chain, address string, pubKey []byte) {
	panic("implement me")
}

func (s *KVStore) IsKeyExisted(ctx sdk.Context, chain string) bool {
	panic("implement me")
}

func (s *KVStore) IsChainKeyAddress(ctx sdk.Context, chain, address string) bool {
	panic("implement me")
}

func (s *KVStore) GetPubKey(ctx sdk.Context, chain string) []byte {
	panic("implement me")
}

func (s *KVStore) UpdateKeygenStatus(ctx sdk.Context, chain, status string) {
	panic("implement me")
}

func (s *KVStore) GetKeygenStatus(ctx sdk.Context, chain string) (string, error) {
	panic("implement me")
}

func (s *KVStore) InsertContracts(ctx sdk.Context, contracts []*tsstypes.ContractEntity) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	for _, contract := range contracts {
		bz, err := json.Marshal(contract)
		if err != nil {
			log.Error(err)
			continue
		}
		key := getContractKey(contract.Chain, contract.Hash)
		store.Set(key, bz)
	}
}

func (s *KVStore) GetPendingDeployContracts(ctx sdk.Context, chain string) []*tsstypes.ContractEntity {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	pendingContracts := make([]*tsstypes.ContractEntity, 0)
	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		var contract *tsstypes.ContractEntity
		if err := json.Unmarshal(iter.Value(), contract); err != nil {
			log.Error(err)
			continue
		}

		if contract != nil && strings.EqualFold(contract.Chain, chain) && len(contract.Status) == 0 {
			pendingContracts = append(pendingContracts, contract)
		}
	}

	return pendingContracts
}

func (s *KVStore) GetContractFromAddress(ctx sdk.Context, chain, address string) *tsstypes.ContractEntity {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		var contract *tsstypes.ContractEntity
		if err := json.Unmarshal(iter.Value(), contract); err != nil {
			log.Error(err)
			continue
		}

		if contract != nil && strings.EqualFold(contract.Chain, chain) && strings.EqualFold(contract.Address, address) {
			return contract
		}
	}

	return nil
}

func (s *KVStore) GetContractFromHash(ctx sdk.Context, chain, hash string) *tsstypes.ContractEntity {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	key := getContractKey(chain, hash)
	contractBz := store.Get(key)
	var contract *tsstypes.ContractEntity
	if err := json.Unmarshal(contractBz, contract); err != nil {
		log.Error(err)
		return nil
	}

	return contract
}

func (s *KVStore) UpdateContractsStatus(ctx sdk.Context, contracts []*tsstypes.ContractEntity, status string) error {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	for _, updateContract := range contracts {
		key := getContractKey(updateContract.Chain, updateContract.Hash)
		contractBz := store.Get(key)
		if len(contractBz) == 0 {
			log.Warnf("cannot find contract for chain(%s) with hash(%s)", updateContract.Chain, updateContract.Hash)
			continue
		}

		var contract *tsstypes.ContractEntity
		if err := json.Unmarshal(contractBz, contract); err != nil {
			log.Error(err)
			return err
		}

		contract.Status = status
		if err := saveRecord(store, key, contract); err != nil {
			return err
		}
	}

	return nil
}

func (s *KVStore) UpdateContractAddress(ctx sdk.Context, chain, hash, address string) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixContract)

	key := getContractKey(chain, hash)
	contractBz := store.Get(key)
	if len(contractBz) == 0 {
		log.Warnf("cannot find contract for chain(%s), hash(%s)", chain, hash)
		return
	}

	var contract *tsstypes.ContractEntity
	if err := json.Unmarshal(contractBz, contract); err != nil {
		log.Error(err)
	}

	if contract == nil {
		return
	}

	contract.Address = address
	_ = saveRecord(store, key, contract)
}

func (s *KVStore) InsertTxOuts(ctx sdk.Context, txs []*tsstypes.TxOutEntity) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)
	for _, tx := range txs {
		bz, err := json.Marshal(tx)
		if err != nil {
			log.Error(err)
			continue
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

		var txOut *tsstypes.TxOutEntity
		if err := json.Unmarshal(txOutBz, txOut); err != nil {
			log.Error(err)
			return nil
		}

		return txOut
	}

	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		var txOut *tsstypes.TxOutEntity
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

func (s *KVStore) IsContractDeployTx(ctx sdk.Context, chain string, hashWithoutSig string) (bool, error) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	var txOut *tsstypes.TxOutEntity
	txOutBz := store.Get(getTxOutKey(chain, hashWithoutSig))
	if err := json.Unmarshal(txOutBz, txOut); err != nil {
		log.Error(err)
		return false, err
	}

	return len(txOut.ContractHash) > 0, nil
}

func (s *KVStore) UpdateTxOutSig(ctx sdk.Context, chain, hashWithoutSign, hashWithSig string, sig []byte) error {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	var txOut *tsstypes.TxOutEntity
	key := getTxOutKey(chain, hashWithoutSign)
	txOutBiz := store.Get(key)
	if err := json.Unmarshal(txOutBiz, txOut); err != nil {
		log.Error(err)
		return err
	}

	txOut.HashWithSig = hashWithSig
	txOut.Signature = string(sig)

	_ = saveRecord(store, key, txOut)
	return nil
}

func (s *KVStore) UpdateTxOutStatus(ctx sdk.Context, chain, hash string, status tsstypes.TxOutStatus, isHashWithSig bool) error {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), PrefixTxOut)

	if !isHashWithSig {
		var txOut *tsstypes.TxOutEntity
		key := getTxOutKey(chain, hash)
		txOutBiz := store.Get(key)
		if err := json.Unmarshal(txOutBiz, txOut); err != nil {
			log.Error(err)
			return err
		}

		txOut.Status = string(status)
		_ = saveRecord(store, key, txOut)
		return nil
	}

	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		var txOut *tsstypes.TxOutEntity
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

func (s *KVStore) InsertMempoolTxHash(ctx sdk.Context, hash string, blockHeight int64) {
	panic("implement me")
}

func (s *KVStore) MempoolTxExisted(ctx sdk.Context, hash string) bool {
	panic("implement me")
}

func (s *KVStore) MempoolTxExistedRange(ctx sdk.Context, hash string, minBlock int64, maxBlock int64) bool {
	panic("implement me")
}

func getTxOutKey(chain, hashWithoutSig string) []byte {
	return []byte(strings.Join([]string{chain, hashWithoutSig}, KVSeparator))
}

func getContractKey(chain, hash string) []byte {
	return []byte(strings.Join([]string{chain, hash}, KVSeparator))
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
