package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"sort"
	"strings"

	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/abci/types"
)

// TODO: Move txout's byte into separate store.
var (
	prefixTxRecord               = []byte{0x01} // Vote for a tx by different nodes
	prefixTxRecordProcessed      = []byte{0x02}
	prefixKeygen                 = []byte{0x03}
	prefixKeygenResultWithSigner = []byte{0x04}
	prefixContract               = []byte{0x05}
	prefixContractByteCode       = []byte{0x06}
	prefixContractAddress        = []byte{0x07}
	prefixTxIn                   = []byte{0x08}
	prefixTxOut                  = []byte{0x09}
	prefixTxOutSig               = []byte{0x0A}
	prefixTxOutContractConfirm   = []byte{0x0B}
	prefixContractName           = []byte{0x0C}
	prefixGasPrice               = []byte{0x0D}
	prefixChain                  = []byte{0x0E}
	prefixToken                  = []byte{0x0F}
	prefixTokenPrices            = []byte{0x10}
	prefixNode                   = []byte{0x11}
	prefixLiquidity              = []byte{0x12}
	prefixParams                 = []byte{0x13}
	prefixSlash                  = []byte{0x14}
	prefixNodeBalance            = []byte{0x15}
	prefixValidatorUpdate        = []byte{0x16}
	prefixDheartIPAddress        = []byte{0x17}
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

func getGasPriceKey(chain string, height int64) []byte {
	// chain, height
	return []byte(fmt.Sprintf("%s__%d", chain, height))
}

func getValidatorUpdatesKey() []byte {
	return []byte(fmt.Sprintf("validator_update"))
}

///// TxRecord

func saveTxRecord(store cstypes.KVStore, hash []byte, validator string) int {
	vals := make([]string, 0)
	bz := store.Get(hash)
	if bz != nil {
		vals = strings.Split(string(bz), ",")
	}

	if strings.Index(validator, ",") >= 0 {
		return len(vals)
	}

	found := false
	for _, val := range vals {
		if val == validator {
			found = true
			break
		}
	}

	// Only save the result when the validator has not posted the tx record yet.
	if !found {
		vals = append(vals, validator)
		bz = []byte(strings.Join(vals, ","))
		store.Set(hash, bz)
	}

	return len(vals)
}

func getVotersString(store cstypes.KVStore, hash []byte) []string {
	bz := store.Get(hash)
	if bz == nil {
		return []string{}
	}

	return strings.Split(string(bz), ",")
}

func processTxRecord(store cstypes.KVStore, hash []byte) {
	store.Set(hash, []byte{})
}

func isTxRecordProcessed(store cstypes.KVStore, hash []byte) bool {
	return store.Has(hash)
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

	return store.Has(key)
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
		} else {
			log.Warn("msg.PubKeyBytes is empty")
		}
	}

	return result
}

///// Keygen Result With Signer

func saveKeygenResult(store cstypes.KVStore, signerMsg *types.KeygenResultWithSigner) {
	key := getKeygenResultKey(signerMsg.Keygen.KeyType, int(signerMsg.Keygen.Index), signerMsg.Data.From)

	bz, err := signerMsg.Marshal()
	if err != nil {
		log.Error("SaveKeygenResult: Cannot marshal KeygenResult message, err = ", err)
		return
	}

	store.Set(key, bz)
}

func getAllKeygenResult(store cstypes.KVStore, keygenType string, index int32) []*types.KeygenResultWithSigner {
	begin := []byte(fmt.Sprintf("%s__%06d__", keygenType, index))
	end := []byte(fmt.Sprintf("%s__%06d__~", keygenType, index))

	results := make([]*types.KeygenResultWithSigner, 0)

	iter := store.Iterator(begin, end)
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		msg := &types.KeygenResultWithSigner{}
		err := msg.Unmarshal(bz)
		if err != nil {
			log.Error("isKeygenResultSuccess: cannot unmarshal keygen result")
			continue
		}

		results = append(results, msg)
	}

	return results
}

///// Contract

func saveContracts(contractStore cstypes.KVStore, byteCodeStore cstypes.KVStore, msgs []*types.Contract) {
	for _, msg := range msgs {
		saveContract(contractStore, byteCodeStore, msg)
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

		// Pending contracts are contracts that do not have address or status
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
		log.Errorf("getContract: serialized contract is nil, chain = %s, contract hash = %s", chain, contractHash)
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

///// Gas Price
func saveGasPrice(store cstypes.KVStore, msg *types.GasPriceMsg) {
	var (
		record      *types.GasPriceRecord
		savedRecord []byte
		err         error
	)
	key := getGasPriceKey(msg.Chain, msg.BlockHeight)
	currentRecord := getGasPriceRecord(store, msg.Chain, msg.BlockHeight)
	if currentRecord == nil {
		record = &types.GasPriceRecord{
			Messages: []*types.GasPriceMsg{msg},
		}
		savedRecord, err = record.Marshal()
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		record = &types.GasPriceRecord{
			Messages: append(currentRecord.Messages, msg),
		}
		savedRecord, err = record.Marshal()
		if err != nil {
			log.Error(err)
			return
		}
	}

	if savedRecord == nil {
		return
	}

	store.Set(key, savedRecord)
}

func getGasPriceRecord(store cstypes.KVStore, chain string, height int64) *types.GasPriceRecord {
	key := getGasPriceKey(chain, height)
	bz := store.Get(key)
	record := &types.GasPriceRecord{}
	if err := record.Unmarshal(bz); err != nil {
		log.Error(err)
		return nil
	}

	return record
}

///// Chain
func saveChain(store cstypes.KVStore, chain *types.Chain) {
	bz, err := chain.Marshal()
	if err != nil {
		log.Error("saveChain: failed to save chain, chain = ", chain.Id)
		return
	}

	store.Set([]byte(chain.Id), bz)
}

func getChain(store cstypes.KVStore, chainId string) *types.Chain {
	chain := &types.Chain{}
	bz := store.Get([]byte(chainId))
	if bz == nil {
		return nil
	}

	if err := chain.Unmarshal(bz); err != nil {
		log.Error("getChain: failed to unmarshal bytes for chain ", chainId)
		return nil
	}

	return chain
}

func getAllChains(store cstypes.KVStore) map[string]*types.Chain {
	m := make(map[string]*types.Chain)

	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		chain := &types.Chain{}

		if err := chain.Unmarshal(iter.Value()); err != nil {
			log.Error("getAllChains: failed to unmarshal bytes for chain ", string(iter.Key()))
			return nil
		}

		m[string(iter.Key())] = chain
	}

	return m
}

///// TxOutConfirm
func saveTxOutConfirm(store cstypes.KVStore, msg *types.TxOutContractConfirm) {
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

///// Token Prices
func setTokenPrices(store cstypes.KVStore, blockHeight uint64, msg *types.UpdateTokenPrice) {
	key := []byte(msg.Signer)
	value := store.Get(key)

	var record *types.TokenPriceRecords
	if value == nil {
		record = new(types.TokenPriceRecords)
		record.Records = make([]*types.TokenPriceRecord, 0)
	} else {
		err := record.Unmarshal(value)
		if err != nil {
			log.Error("cannot unmarshal record for signer ", msg.Signer)
			return
		}
	}

	indexes := make(map[string]int)
	for i, record := range record.Records {
		indexes[record.Token] = i
	}

	for _, tokenPrice := range msg.TokenPrices {
		if index, ok := indexes[tokenPrice.Id]; ok {
			record.Records[index].BlockHeight = blockHeight
			record.Records[index].Price = tokenPrice.Price
		} else {
			record.Records = append(record.Records, &types.TokenPriceRecord{
				Token:       tokenPrice.Id,
				BlockHeight: blockHeight,
				Price:       tokenPrice.Price,
			})
		}
	}

	bz, err := record.Marshal()
	if err != nil {
		log.Error("cannot unmarshal token price record for signer ", msg.Signer)
		return
	}

	store.Set(key, bz)
}

// getAllTokenPrices gets all the token prices all of all signers.
func getAllTokenPrices(store cstypes.KVStore) map[string]*types.TokenPriceRecords {
	result := make(map[string]*types.TokenPriceRecords)

	for iter := store.Iterator(nil, nil); iter.Valid(); iter.Next() {
		// Key is signer.
		signer := string(iter.Key())
		bz := iter.Value()
		record := new(types.TokenPriceRecords)
		err := record.Unmarshal(bz)
		if err != nil {
			log.Error("cannot unmarshal token price record for signer ", signer, " err = ", err)
			continue
		}

		result[signer] = record
	}

	return result
}

///// Tokens

func setTokens(store cstypes.KVStore, tokens map[string]*types.Token) {
	for id, token := range tokens {
		bz, err := token.Marshal()
		if err != nil {
			log.Error("cannot marshal token ", id)
			continue
		}

		store.Set([]byte(id), bz)
	}
}

func getTokens(store cstypes.KVStore, tokenIds []string) map[string]*types.Token {
	tokens := make(map[string]*types.Token)
	for _, id := range tokenIds {
		bz := store.Get([]byte(id))

		token := &types.Token{}
		err := token.Unmarshal(bz)
		if err != nil {
			log.Error("getTokens: cannot unmarshal token ", id)
			continue
		}

		tokens[id] = token
	}

	return tokens
}

func getAllTokens(store cstypes.KVStore) map[string]*types.Token {
	tokens := make(map[string]*types.Token)

	iter := store.Iterator(nil, nil)

	for ; iter.Valid(); iter.Next() {
		token := &types.Token{}
		err := token.Unmarshal(iter.Value())
		if err != nil {
			log.Error("cannot unmarshal token ", string(iter.Key()))
			continue
		}

		tokens[string(iter.Key())] = token
	}

	return tokens
}

///// Node
func saveNode(store cstypes.KVStore, node *types.Node) {
	bz, err := node.Marshal()
	if err != nil {
		log.Error("cannot marshal node, err = ", err)
		return
	}

	store.Set(node.ConsensusKey.Bytes, bz)
}

func updateNodeStatus(store cstypes.KVStore, consKey []byte, nodeStatus types.NodeStatus) {
	bz := store.Get(consKey)
	if bz == nil {
		return
	}

	node := &types.Node{}
	if err := node.Unmarshal(bz); err != nil {
		log.Error("can not unmarshal node, err = ", err)
		return
	}

	if nodeStatus == types.NodeStatus_Validator {
		node.IsValidator = true
	} else {
		node.IsValidator = false
	}

	node.Status = nodeStatus
	saveNode(store, node)
}

func loadNodesByStatus(store cstypes.KVStore, status types.NodeStatus) []*types.Node {
	vals := make([]*types.Node, 0)

	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		node := &types.Node{}
		err := node.Unmarshal(iter.Value())
		if err != nil {
			log.Error("cannot unmarshal node, err = ", err)
			continue
		}

		if status == types.NodeStatus_Unknown || node.Status == status {
			vals = append(vals, node)
		}
	}

	return vals
}

func setValidators(store cstypes.KVStore, vals []*types.Node) ([]*types.Node, error) {
	if err := clearOldValidators(store); err != nil {
		log.Error("error when clear old validators set. error = ", err)
		return nil, err
	}

	res := make([]*types.Node, 0)
	for _, val := range vals {
		if !val.IsValidator {
			continue
		}

		saveNode(store, val)
		res = append(res, val)
	}

	return res, nil
}

func clearOldValidators(store cstypes.KVStore) error {
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		node := &types.Node{}
		err := node.Unmarshal(iter.Value())
		if err != nil {
			log.Error("cannot unmarshal node, err = ", err)
			continue
		}

		if !node.IsValidator {
			continue
		}

		store.Delete(iter.Key())
	}

	return iter.Close()
}

///// Liquidity
func setLiquidities(store cstypes.KVStore, liquidities map[string]*types.Liquidity) {
	for id, liquid := range liquidities {
		bz, err := liquid.Marshal()
		if err != nil {
			log.Error("cannot marshal liquidity ", id)
			continue
		}

		store.Set([]byte(id), bz)
	}
}

func getLiquidity(store cstypes.KVStore, chain string) *types.Liquidity {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return nil
	}

	liquid := &types.Liquidity{}
	if err := liquid.Unmarshal(bz); err != nil {
		log.Errorf("getLiquidity: error when unmarshal liquid for chain: %s", chain)
		return nil
	}

	return liquid
}

func getAllLiquidities(store cstypes.KVStore) map[string]*types.Liquidity {
	liquids := make(map[string]*types.Liquidity)

	iter := store.Iterator(nil, nil)

	for ; iter.Valid(); iter.Next() {
		liq := &types.Liquidity{}
		err := liq.Unmarshal(iter.Value())
		if err != nil {
			log.Error("cannot unmarshal liquidity ", string(iter.Key()))
			continue
		}

		liquids[string(iter.Key())] = liq
	}

	_ = iter.Close()

	return liquids
}

///// Params
func saveParams(store cstypes.KVStore, params *types.Params) {
	bz, err := params.Marshal()
	if err != nil {
		log.Error("cannot marshal params ")
	}

	store.Set([]byte("params"), bz)
}

func getParams(store cstypes.KVStore) *types.Params {
	bz := store.Get([]byte("params"))
	if bz == nil {
		return nil
	}

	params := &types.Params{}
	if err := params.Unmarshal(bz); err != nil {
		log.Errorf("getParams: error when unmarshal params")
		return nil
	}

	return params
}

///// Slash
func incOrDecSlashToken(store cstypes.KVStore, amount int64, address sdk.AccAddress) error {
	oldAmt, err := getCurSlashToken(store, address)
	if err != nil {
		return err
	}

	newAmt := oldAmt + amount
	if newAmt < 0 {
		newAmt = 0
	}

	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, uint64(newAmt))
	store.Set(address.Bytes(), bz)
	return nil
}

func getCurSlashToken(store cstypes.KVStore, address sdk.AccAddress) (int64, error) {
	addrKey := address.Bytes()
	bz := store.Get(addrKey)
	if bz == nil {
		return 0, nil
	}

	cur := binary.LittleEndian.Uint64(bz)
	return int64(cur), nil
}

///// Node balance
func incOrDecBondBalance(store cstypes.KVStore, address sdk.AccAddress, amount int64) error {
	oldAmt, err := getBondBalance(store, address)
	if err != nil {
		return err
	}

	newAmt := oldAmt + amount
	if newAmt < 0 {
		newAmt = 0
	}

	bz := make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, uint64(newAmt))
	store.Set(address.Bytes(), bz)
	return nil
}

func getBondBalance(store cstypes.KVStore, address sdk.AccAddress) (int64, error) {
	addrKey := address.Bytes()
	bz := store.Get(addrKey)
	if bz == nil {
		return 0, nil
	}

	cur := binary.LittleEndian.Uint64(bz)
	return int64(cur), nil
}

type NodeBalance struct {
	Addr    sdk.AccAddress
	Balance int64
}

// order from the highest balance to the lowest balance
func getTopBalances(store cstypes.KVStore, n int) []sdk.AccAddress {
	allNodeBalances := getOrderedNodeBalances(store)
	topBalanceAddrs := make([]sdk.AccAddress, 0)
	if n == -1 {
		for _, nodeBalance := range allNodeBalances {
			topBalanceAddrs = append(topBalanceAddrs, nodeBalance.Addr)
		}

		return topBalanceAddrs
	}

	if n > len(allNodeBalances) {
		n = len(allNodeBalances)
	}

	for i := 0; i < n; i++ {
		topBalanceAddrs = append(topBalanceAddrs, allNodeBalances[i].Addr)
	}
	return topBalanceAddrs
}

func getOrderedNodeBalances(store cstypes.KVStore) []*NodeBalance {
	allNodeBalances := make([]*NodeBalance, 0)
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress(iter.Key())
		balance := binary.LittleEndian.Uint64(iter.Value())
		allNodeBalances = append(allNodeBalances, &NodeBalance{
			Addr:    addr,
			Balance: int64(balance),
		})
	}
	_ = iter.Close()

	sort.SliceStable(allNodeBalances, func(i, j int) bool {
		return allNodeBalances[i].Balance > allNodeBalances[j].Balance
	})

	return allNodeBalances
}

func saveValidatorUpdates(store cstypes.KVStore, validatorUpdates abci.ValidatorUpdates) error {
	log.Debug("Save incoming validator update", validatorUpdates)
	resetValidatorUpdate(store)
	for i, v := range validatorUpdates {
		buf := bytes.NewBuffer(make([]byte, 0))
		if err := tmTypes.WriteMessage(&v, buf); err != nil {
			log.Error("error when write validator update message. error = ", err)
			return err
		}

		key := make([]byte, 4)
		binary.BigEndian.PutUint32(key, uint32(i))

		store.Set(key, buf.Bytes())
	}

	return nil
}

func getValidatorUpdates(store cstypes.KVStore) abci.ValidatorUpdates {
	validatorUpdates := make(abci.ValidatorUpdates, 0)
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		buf := bytes.NewBuffer(iter.Value())
		v := new(abci.ValidatorUpdate)
		if err := tmTypes.ReadMessage(buf, v); err != nil {
			log.Error("error when read validate update from buffer. error = ", err)
			return abci.ValidatorUpdates{}
		}

		validatorUpdates = append(validatorUpdates, *v)
	}

	return validatorUpdates
}

func resetValidatorUpdate(store cstypes.KVStore) {
	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

func increaseValidatorUpdateIndex(store cstypes.KVStore) int {
	key := []byte("msg_index")
	b := store.Get(key)
	// Set index to 0 for the first time
	if b == nil {
		index := make([]byte, 4)
		binary.LittleEndian.PutUint32(index, 0)
		store.Set(key, index)
		return 0
	}

	index := binary.LittleEndian.Uint32(b)
	index++
	indexBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(indexBytes, index)
	store.Set(key, indexBytes)
	return int(index)
}

func saveDheartIPAddress(store cstypes.KVStore, accAddr sdk.AccAddress, ip string) error {
	if len(ip) == 0 {
		return nil
	}

	bz, err := accAddr.Marshal()
	if err != nil {
		log.Error("error when marshal accAddr: ", err)
		return err
	}

	store.Set(bz, []byte(ip))
	return nil
}

type AccAddressDheartIP struct {
	Addr sdk.AccAddress
	IP   string
}

func getAllDheartIPAddresses(store cstypes.KVStore) []AccAddressDheartIP {
	r := make([]AccAddressDheartIP, 0)
	iter := store.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		addr := sdk.AccAddress{}
		if err := addr.Unmarshal(iter.Key()); err != nil {
			log.Error("error when unmarshalling acc address: ", err)
			return nil
		}

		r = append(r, AccAddressDheartIP{
			Addr: addr,
			IP:   string(iter.Value()),
		})
	}

	return r
}

///// Debug functions
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
