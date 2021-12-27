package keeper

import (
	"fmt"

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
	return []byte(fmt.Sprintf("%s__%d", keyType, index))
}

func getKeygenResultKey(keyType string, index int, from string) []byte {
	// keyType
	return []byte(fmt.Sprintf("%s__%d__%s", keyType, index, from))
}

func getContractKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

func getContractByteCodeKey(chain string, hash string) []byte {
	// chain + hash
	return []byte(fmt.Sprintf("%s__%s", chain, hash))
}

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

/// Debug functions
func printStore(store cstypes.KVStore) {
	log.Info("======== DEBUGGING")
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		log.Info("key = ", string(iter.Key()))
		log.Info("value = ", string(iter.Value()))
	}
	log.Info("======== END OF DEBUGGING")
}
