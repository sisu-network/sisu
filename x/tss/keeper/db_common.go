package keeper

import (
	"fmt"
	"strings"

	cstypes "github.com/sisu-network/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

var (
	prefixKeygen           = []byte{0x01}
	prefixKeygenResult     = []byte{0x02}
	prefixContract         = []byte{0x03}
	prefixContractByteCode = []byte{0x04}
	prefixTxIn             = []byte{0x05}
	prefixTxOut            = []byte{0x06}
)

func getKeygenKey(keyType string, index int) []byte {
	// keyType + id
	return []byte(fmt.Sprintf("%s__%06d", keyType, index))
}

func getKeygenResultKey(keyType string, index int, from string) []byte {
	// keyType
	return []byte(fmt.Sprintf("%s__%06d__%s", keyType, index, from))
}

func getContractKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

func getContractByteCodeKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

func getTxInKey(chain string, height int64, hash string) []byte {
	// chain, height, hash
	return []byte(fmt.Sprintf("%s__%d__%s", chain, height, hash))
}

func getTxOutKey(inChain string, outChain string, outHash string) []byte {
	// inChain, outChain, height, hash
	return []byte(fmt.Sprintf("%s__%s__%s", inChain, outChain, outHash))
}

///// Keygen

func saveKeygen(store cstypes.KVStore, msg *types.Keygen) {
	key := getKeygenKey(msg.KeyType, int(msg.Index))

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("SaveKeygenProposal: cannot marshal keygen proposal, err = ", err)
	}
	store.Set(key, bz)
}

func isKeygenExisted(store cstypes.KVStore, keyType string, index int) bool {
	key := getKeygenKey(keyType, index)

	return store.Get(key) != nil
}

func saveKeygenResult(store cstypes.KVStore, signerMsg *types.KeygenResultWithSigner) {
	key := getKeygenResultKey(signerMsg.Keygen.KeyType, int(signerMsg.Keygen.Index), signerMsg.Data.From)

	bz, err := signerMsg.Data.Marshal()
	if err != nil {
		log.Error("SaveKeygenResult: Cannot marshal KeygenResult message, err = ", err)
		return
	}

	store.Set(key, bz)
}

func isKeygenAddress(store cstypes.KVStore, keyType string, address string) bool {
	begin := []byte(fmt.Sprintf("%s__", keyType))

	iter := store.ReverseIterator(begin, nil)
	for ; iter.Valid(); iter.Next() {
		msg := &types.Keygen{}
		if err := msg.Unmarshal(iter.Value()); err != nil {
			log.Error("IsKeygenAddress: cannot unmarshal keygen")
			continue
		}

		if msg.Address == address {
			return true
		}
	}

	return false
}

func getKeygenPubkey(store cstypes.KVStore, keyType string) []byte {
	begin := []byte(fmt.Sprintf("%s__", keyType))

	iter := store.ReverseIterator(begin, nil)
	for ; iter.Valid(); iter.Next() {
		msg := &types.Keygen{}
		if err := msg.Unmarshal(iter.Value()); err != nil {
			log.Error("IsKeygenAddress: cannot unmarshal keygen")
			continue
		}

		if msg.PubKeyBytes != nil && len(msg.PubKeyBytes) > 0 {
			return msg.PubKeyBytes
		}
	}

	return nil
}

func getAllKeygenPubkeys(store cstypes.KVStore) map[string][]byte {
	iter := store.Iterator(nil, nil)

	result := make(map[string][]byte)

	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key())
		index := strings.Index(key, "__")
		if index < 0 {
			continue
		}

		keyType := key[0:index]
		msg := &types.Keygen{}
		if err := msg.Unmarshal(iter.Value()); err != nil {
			log.Error("getAllKeygenPubkeys: cannot unmarshal keygen")
			continue
		}
		if len(msg.PubKeyBytes) > 0 {
			result[keyType] = msg.PubKeyBytes
		}
	}

	return result
}

///// Keygen Result

// Keygen is considered successful if at least there is at least 1 successful KeygenReslut in the
// KVStore.
func isKeygenResultSuccess(store cstypes.KVStore, signerMsg *types.KeygenResultWithSigner) bool {
	msg := signerMsg.Keygen
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

func getAllPubKeys(store cstypes.KVStore) map[string][]byte {
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

func saveContracts(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, msgs []*types.Contract, saveByteCode bool) {
	log.Info("Saving contracts, contracts length = ", len(msgs))

	for _, msg := range msgs {
		log.Infof("Saving contract on chain %s with hash = %s", msg.Chain, msg.Hash)

		bz, err := msg.Marshal()
		if err != nil {
			log.Error("Cannot marshal contract message, err = ", err)
			continue
		}

		// Save byte code into separate store since it's rarely read.
		copy := &types.Contract{}
		if msg.ByteCodes == nil {
			// ByteCode is nil, the copy is the same object reference as message
			copy = msg
		} else {
			// ByteCode is not nil, we need to remove the bytecode from the copy.
			err = copy.Unmarshal(bz)
			if err != nil {
				log.Error("Cannot unmarshal contract message, err = ", err)
				continue
			}

			// Set bytecode to nil
			copy.ByteCodes = nil
		}

		// Get the serialized bytes of copy
		bz, err = copy.Marshal()
		if err != nil {
			log.Error("Cannot marshal contract copy, err = ", err)
			continue
		}

		contractKey := getContractKey(msg.Chain, msg.Hash)
		contractStore.Set(contractKey, bz)

		// Save byte code
		if saveByteCode && msg.ByteCodes != nil {
			byteCodeKey := getContractByteCodeKey(msg.Chain, msg.Hash)
			byteCodeStore.Set(byteCodeKey, msg.ByteCodes)
		}
	}
}

func isContractExisted(contractStore cstypes.KVStore, msg *types.Contract) bool {
	contractKey := getContractKey(msg.Chain, msg.Hash)
	return contractStore.Has(contractKey)
}

func getPendingContracts(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, chain string) []*types.Contract {
	contracts := make([]*types.Contract, 0)

	iter := contractStore.Iterator([]byte(fmt.Sprintf("%s__", chain)), []byte(fmt.Sprintf("%s__~", chain)))

	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		bz := iter.Value()

		contract := &types.Contract{}
		err := contract.Unmarshal(bz)
		if err != nil {
			log.Error("GetPendingContracts: Cannot unmarshal contract bytes")
			continue
		}

		if contract.Status != "" {
			continue
		}

		bz = byteCodeStore.Get(key)
		contract.ByteCodes = bz

		contracts = append(contracts, contract)
	}

	return contracts
}

func updateContractsStatus(contractStore cstypes.KVStore, msgs []*types.Contract, status string) {
	for _, msg := range msgs {
		key := getContractByteCodeKey(msg.Chain, msg.GetHash())

		bz := contractStore.Get(key)
		contract := &types.Contract{}
		err := contract.Unmarshal(bz)

		if err != nil {
			log.Error("UpdateContractsStatus: Cannot unmarshal contract bytes")
			return
		}

		contract.Status = status
		bz, err = contract.Marshal()
		if err != nil {
			log.Error("UpdateContractsStatus: Cannot marshal contract bytes")
			return
		}

		contractStore.Set(key, bz)
	}
}

///// TxIn
func saveTxIn(store cstypes.KVStore, msg *types.TxIn) {
	key := getTxInKey(msg.Chain, msg.BlockHeight, msg.TxHash)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal TxIn")
		return
	}

	store.Set(key, bz)
}

func isTxInExisted(store cstypes.KVStore, msg *types.TxIn) bool {
	key := getTxInKey(msg.GetChain(), msg.GetBlockHeight(), msg.GetTxHash())
	return store.Has(key)
}

///// TxOut
func saveTxOut(store cstypes.KVStore, msg *types.TxOut) {
	key := getTxOutKey(msg.InChain, msg.OutChain, msg.GetHash())
	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal tx out")
		return
	}

	store.Set(key, bz)
}

func isTxOutExisted(store cstypes.KVStore, msg *types.TxOut) bool {
	key := getTxOutKey(msg.InChain, msg.OutChain, msg.GetHash())
	return store.Has(key)
}

func getTxOut(store cstypes.KVStore, inChain string, outChain, hash string) *types.TxOut {
	key := getTxOutKey(inChain, outChain, hash)
	bz := store.Get(key)

	if bz == nil {
		return nil
	}

	txOut := &types.TxOut{}
	err := txOut.Unmarshal(bz)
	if err != nil {
		log.Error("getTxOUt: Cannot unmasharl txout")
		return nil
	}

	return txOut
}

/// Debug functions
func printStore(store cstypes.KVStore) {
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		log.Info("key = ", string(iter.Key()))
		log.Info("value = ", string(iter.Value()))
	}
}
