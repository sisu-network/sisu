package keeper

import (
	"fmt"
	"strings"

	cstypes "github.com/sisu-network/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/tss/types"
)

// TODO: Move txout's byte into separate store.
var (
	prefixKeygen           = []byte{0x01}
	prefixKeygenResult     = []byte{0x02}
	prefixContract         = []byte{0x03}
	prefixContractByteCode = []byte{0x04}
	prefixContractAddress  = []byte{0x05}
	prefixTxIn             = []byte{0x06}
	prefixTxOut            = []byte{0x07}
	prefixTxOutSig         = []byte{0x07}
	prefixTxOutConfirm     = []byte{0x08}
	prefixMempoolTx        = []byte{0x09}
	prefixContractName     = []byte{0x10}
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

func getContractNameKey(chain, name string) []byte {
	// chain + contract name
	return []byte(fmt.Sprintf("%s__%s", chain, name))
}

func getTxInKey(chain string, height int64, hash string) []byte {
	// chain, height, hash
	return []byte(fmt.Sprintf("%s__%d__%s", chain, height, hash))
}

func getTxOutKey(outChain string, outHash string) []byte {
	// outChain, hash
	return []byte(fmt.Sprintf("%s__%s", outChain, outHash))
}

func getTxOutSigKey(outChain string, hashWithSig string) []byte {
	// outChain, hash with sig
	return []byte(fmt.Sprintf("%s__%s", outChain, hashWithSig))
}

func getTxOutConfirmKey(outChain string, outHash string) []byte {
	// outChain, hash
	return []byte(fmt.Sprintf("%s__%s", outChain, outHash))
}

func getContractAddressKey(chain string, address string) []byte {
	// chain, address
	return []byte(fmt.Sprintf("%s__%s", chain, address))
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

	log.Debug(result)
	return result
}

///// Keygen Result

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
func isKeygenResultSuccess(store cstypes.KVStore, signerMsg *types.KeygenResultWithSigner, self string) bool {
	msg := signerMsg.Keygen
	begin := []byte(fmt.Sprintf("%s__%06d__", msg.KeyType, int(msg.Index)))
	end := []byte(fmt.Sprintf("%s__%06d__~", msg.KeyType, int(msg.Index)))

	iter := store.Iterator(begin, end)
	count := 0
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		msg := &types.KeygenResult{}
		err := msg.Unmarshal(bz)
		if err != nil {
			log.Error("isKeygenResultSuccess: cannot unmarshal keygen result")
			continue
		}
		count += 1

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

///// Contract

func saveContracts(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, msgs []*types.Contract) {
	for _, msg := range msgs {
		saveContract(contractStore, byteCodeStore, msg)
	}
}

func saveContractAddressesForName(contractHashStore cstypes.KVStore, msgs []*types.Contract) {
	for _, msg := range msgs {
		saveContractAddressForName(contractHashStore, msg)
	}
}

func saveContract(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, msg *types.Contract) {
	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal contract message, err = ", err)
		return
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
			return
		}

		// Set bytecode to nil
		copy.ByteCodes = nil
	}

	// Get the serialized bytes of copy
	bz, err = copy.Marshal()
	if err != nil {
		log.Error("Cannot marshal contract copy, err = ", err)
		return
	}

	contractKey := getContractKey(msg.Chain, msg.Hash)
	contractStore.Set(contractKey, bz)

	// Save byte code
	if byteCodeStore != nil && len(msg.ByteCodes) > 0 {
		byteCodeKey := getContractKey(msg.Chain, msg.Hash)
		byteCodeStore.Set(byteCodeKey, msg.ByteCodes)
	}
}

func saveContractAddressForName(contractStore cstypes.KVStore, msg *types.Contract) {
	key := getContractNameKey(msg.Chain, msg.Name)
	contractStore.Set(key, []byte(msg.Address))
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

		if contract.Address != "" || contract.Status != "" {
			continue
		}

		bz = byteCodeStore.Get(key)
		contract.ByteCodes = bz

		contracts = append(contracts, contract)
	}

	return contracts
}

// TODO: Remove this. The status should be put in separate store.
func updateContractsStatus(contractStore cstypes.KVStore, chain string, contractHash string, status string) {
	key := getContractKey(chain, contractHash)

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

func getContract(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, chain string, contractHash string) *types.Contract {
	key := getContractKey(chain, contractHash)
	bz := contractStore.Get(key)

	if bz == nil {
		log.Error("getContract: serialized contract is nil")
		return nil
	}

	contract := &types.Contract{}
	err := contract.Unmarshal(bz)
	if err != nil {
		log.Error("getContract: cannot unmarshal contract, err = ", err)
		return nil
	}

	if byteCodeStore != nil {
		contract.ByteCodes = byteCodeStore.Get(key)
	}

	return contract
}

func getContractAddressByName(contractNameStore cstypes.KVStore, chain, name string) string {
	key := getContractNameKey(chain, name)
	bz := contractNameStore.Get(key)
	if bz == nil {
		log.Error("getContractAddressByName: serialized contract hash is nil")
		return ""
	}

	return string(bz)
}

func updateContractAddress(contractStore cstypes.KVStore, chain string, hash string, address string) {
	contract := getContract(contractStore, nil, chain, hash)
	if contract == nil {
		return
	}

	contract.Address = address
	saveContracts(contractStore, nil, []*types.Contract{contract})
}

///// Contract Address
func createContractAddress(caStore cstypes.KVStore, txOutStore cstypes.KVStore, chain string, txOutHash string, address string) {
	// Find the txout in the contract hash
	txOut := getTxOut(txOutStore, chain, txOutHash)
	if txOut == nil {
		log.Error("createContractAddress: cannot find txOut with hash ", txOutHash)
		return
	}

	if len(txOut.ContractHash) == 0 {
		log.Error("createContractAddress: contract hash hash length = 0")
		return
	}

	ca := &types.ContractAddress{
		Chain:        chain,
		Address:      address,
		ContractHash: txOut.ContractHash,
		TxOutHash:    txOutHash,
	}
	bz, err := ca.Marshal()
	if err != nil {
		log.Error("createContractAddress: cannot marhsal contract adress, err = ", err)
		return
	}

	key := getContractAddressKey(chain, address)
	caStore.Set(key, bz)
}

func isContractExistedAtAddress(store cstypes.KVStore, chain, address string) bool {
	key := getContractAddressKey(chain, address)

	return store.Has(key)
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
	key := getTxOutKey(msg.OutChain, msg.OutHash)
	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal tx out")
		return
	}

	store.Set(key, bz)
}

func isTxOutExisted(store cstypes.KVStore, msg *types.TxOut) bool {
	key := getTxOutKey(msg.OutChain, msg.OutHash)
	return store.Has(key)
}

func getTxOut(store cstypes.KVStore, outChain, hash string) *types.TxOut {
	key := getTxOutKey(outChain, hash)
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

///// TxOutSig
func saveTxOutSig(store cstypes.KVStore, msg *types.TxOutSig) {
	key := getTxOutSigKey(msg.Chain, msg.HashWithSig)

	bz, err := msg.Marshal()
	if err != nil {
		log.Error("saveTxOutSig: cannot marshal tx out")
		return
	}

	store.Set(key, bz)
}

func getTxOutSig(store cstypes.KVStore, chain string, hashWithSig string) *types.TxOutSig {
	key := getTxOutSigKey(chain, hashWithSig)
	bz := store.Get(key)

	if bz == nil {
		return nil
	}

	tx := &types.TxOutSig{}
	err := tx.Unmarshal(bz)
	if err != nil {
		log.Error("getTxOutSig: cannot unmarshal TxOutSig")
		return nil
	}

	return tx
}

///// TxOutConfirm
func saveTxOutConfirm(store cstypes.KVStore, msg *types.TxOutConfirm) {
	key := getTxOutConfirmKey(msg.OutChain, msg.OutHash)
	bz, err := msg.Marshal()
	if err != nil {
		log.Error("saveTxOutConfirm: Cannot marshal tx out")
		return
	}

	store.Set(key, bz)
}

func isTxOutConfirmExisted(store cstypes.KVStore, outChain, hash string) bool {
	key := getTxOutConfirmKey(outChain, hash)
	return store.Has(key)
}

/// Debug functions
func printStore(store cstypes.KVStore) {
	iter := store.Iterator(nil, nil)
	count := 0
	for ; iter.Valid(); iter.Next() {
		log.Info("key = ", string(iter.Key()))
		log.Info("value = ", string(iter.Value()))
		count += 1
	}
	log.Info("printStore: Total element count: ", count)
}

func printStoreKeys(store cstypes.KVStore) {
	iter := store.Iterator(nil, nil)
	count := 0
	for ; iter.Valid(); iter.Next() {
		log.Info("key = ", string(iter.Key()))
		count += 1
	}
	log.Info("printStoreKey: Total element count: ", count)
}
