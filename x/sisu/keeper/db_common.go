package keeper

import (
	"fmt"
	"strings"

	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
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

func getLiquidityKey(chain string) []byte {
	// chain
	return []byte(chain)
}

func getVaultKey(chain string, token string) []byte {
	return []byte(fmt.Sprintf("%s__%s", chain, token))
}

func getChainMetadataKey(chain, signer string) []byte {
	return []byte(fmt.Sprintf("%s__%s", chain, signer))
}

func getGasPriceKey(chain, signer string) []byte {
	return []byte(fmt.Sprintf("%s__%s", chain, signer))
}

func getVoteResultKey(hash, signer string) []byte {
	return []byte(fmt.Sprintf("%s__%s", hash, signer))
}

///// TxREcord

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
	begin := append([]byte(keyType), byte(255))
	end := []byte(keyType)

	iter := store.ReverseIterator(end, begin)
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

///// TxOut
func saveTxOut(store cstypes.KVStore, msg *types.TxOut) {
	key := getTxOutKey(msg.Content.OutChain, msg.Content.OutHash)
	bz, err := msg.Marshal()
	if err != nil {
		log.Error("Cannot marshal tx out")
		return
	}

	store.Set(key, bz)
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
func saveGasPrice(store cstypes.KVStore, chain, signer string, record *types.GasPriceRecord) {
	key := getGasPriceKey(chain, signer)
	bz, err := record.Marshal()
	if err != nil {
		log.Errorf("saveGasPrice: Failed to marshal gas price")
		return
	}

	store.Set(key, bz)
}

func getGasPrices(store cstypes.KVStore, chain string) map[string]*types.GasPriceRecord {
	ret := make(map[string]*types.GasPriceRecord)
	begin := []byte(fmt.Sprintf("%s__", chain))
	end := []byte(fmt.Sprintf("%s__~", chain))

	iter := store.Iterator(begin, end)
	for ; iter.Valid(); iter.Next() {
		key := string(iter.Key())
		index := strings.Index(key, "__")
		if index < 0 || index >= len(key)-2 {
			log.Errorf("getGasPrices: Invalid key %s", key)
			continue
		}
		signer := key[index+2:]

		bz := iter.Value()
		record := new(types.GasPriceRecord)
		if err := record.Unmarshal(bz); err != nil {
			log.Errorf("Failed to unmarshal gas price record, err = %s", err)
			continue
		}

		ret[signer] = record
	}

	return ret
}

///// Chain
func saveChain(store cstypes.KVStore, chain *types.Chain) {
	bz, err := chain.Marshal()
	if err != nil {
		log.Error("saveChain: failed to save chain, chaain = ", chain.Id)
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

///// Token Prices
func setTokenPrices(store cstypes.KVStore, blockHeight uint64, msg *types.UpdateTokenPrice) {
	key := []byte(msg.Signer)
	value := store.Get(key)

	record := &types.TokenPriceRecords{Records: make([]*types.TokenPriceRecord, 0)}
	if len(value) > 0 {
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
	}

	store.Set(node.ValPubkey.Bytes, bz)
}

func loadValidators(store cstypes.KVStore) []*types.Node {
	vals := make([]*types.Node, 0)

	iter := store.Iterator(nil, nil)
	for ; iter.Valid(); iter.Next() {
		node := &types.Node{}
		err := node.Unmarshal(iter.Value())
		if err != nil {
			log.Error("cannot unmarshal node, err = ", err)
			continue
		}

		if node.IsValidator {
			vals = append(vals, node)
		}
	}

	return vals
}

///// Vault
func setVaults(store cstypes.KVStore, vaults []*types.Vault) {
	for i, vault := range vaults {
		bz, err := vaults[i].Marshal()
		if err != nil {
			log.Error("cannot marshal vault on chain ", vault.Chain)
			continue
		}

		key := getVaultKey(vault.Chain, vault.Token)

		store.Set(key, bz)
	}
}

func getVault(store cstypes.KVStore, chain string, token string) *types.Vault {
	key := getVaultKey(chain, token)
	bz := store.Get(key)
	if bz == nil {
		return nil
	}

	vault := &types.Vault{}
	if err := vault.Unmarshal(bz); err != nil {
		log.Errorf("getVault: error when unmarshal liquid for chain: %s", chain)
		return nil
	}

	return vault
}

func getAllVaultsForChain(store cstypes.KVStore, chain string) []*types.Vault {
	iter := store.Iterator([]byte(fmt.Sprintf("%s__", chain)), nil)
	vaults := make([]*types.Vault, 0)

	// Go through all vaults
	for ; iter.Valid(); iter.Next() {
		bz := iter.Value()
		vault := new(types.Vault)
		if err := vault.Unmarshal(bz); err != nil {
			log.Errorf("getAllVaultsForChain: error when unmarshal liquid for chain: %s", chain)
			return nil
		}

		vaults = append(vaults, vault)
	}

	return vaults
}

///// MPC Address

func setMpcAddress(store cstypes.KVStore, chain string, address string) {
	store.Set([]byte(chain), []byte(address))
}

func getMpcAddress(store cstypes.KVStore, chain string) string {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return ""
	}

	return string(bz)
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

///// Gateway

func setSisuAccount(store cstypes.KVStore, chain string, address string) {
	key := fmt.Sprintf("%s__%s", chain, "erc20") // we only have erc20 gateway at the moment
	store.Set([]byte(key), []byte(address))
}

func getSisuAccount(store cstypes.KVStore, chain string) string {
	key := fmt.Sprintf("%s__%s", chain, "erc20")
	bz := store.Get([]byte(key))
	if bz == nil {
		return ""
	}

	return string(bz)
}

///// Signer nonce
func setSignerNonce(store cstypes.KVStore, chain string, signer string, nonce uint64) {
	key := fmt.Sprintf("%s__%s", chain, signer)
	store.Set([]byte(key), utils.Uint64ToBytes(nonce))
}

func getSignerNonces(store cstypes.KVStore, chain string) []uint64 {
	begin := []byte(fmt.Sprintf("%s__", chain))
	end := []byte(fmt.Sprintf("%s__~", chain))

	nonces := make([]uint64, 0)

	iter := store.Iterator(begin, end)
	for ; iter.Valid(); iter.Next() {
		nonce := utils.BytesToUint64(iter.Value())
		nonces = append(nonces, nonce)
	}

	return nonces
}

///// Mpc Nonce
func setMpcNonce(store cstypes.KVStore, mpcNonce *types.MpcNonce) {
	if mpcNonce.Chain == "" {
		log.Errorf("setMpcNonce: chain is not defined")
		return
	}

	bz, err := mpcNonce.Marshal()
	if err != nil {
		log.Error("cannot marshal mpcNonce")
		return
	}

	store.Set([]byte(mpcNonce.Chain), bz)
}

func getMpcNonce(store cstypes.KVStore, chain string) *types.MpcNonce {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return nil
	}

	mpcNonce := &types.MpcNonce{}
	err := mpcNonce.Unmarshal(bz)
	if err != nil {
		log.Error("Failed to unmarshal mpcNonce, err = ", err)
		return nil
	}

	return mpcNonce
}

///// Command Queue
func setCommandQueue(store cstypes.KVStore, chain string, commands []*types.Command) {
	cmds := &types.Commands{
		List: commands,
	}

	bz, err := cmds.Marshal()
	if err != nil {
		log.Error("saveTranferQueue: faield to marshal transfer batch")
		return
	}

	store.Set([]byte(chain), bz)
}

func getCommandQueue(store cstypes.KVStore, chain string) []*types.Command {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return nil
	}

	cmds := &types.Commands{}
	err := cmds.Unmarshal(bz)
	if err != nil {
		log.Error("getCommandQueue: failed to unmarshal command")
		return nil
	}

	return cmds.List
}

///// Transfer
func addTransfers(store cstypes.KVStore, transfers []*types.TransferDetails) {
	for _, transfer := range transfers {
		bz, err := transfer.Marshal()
		if err != nil {
			log.Error("addTransfer: failed to marshal transfer, err = ", err)
			continue
		}

		store.Set([]byte(transfer.Id), bz)
	}
}

func getTransfers(store cstypes.KVStore, ids []string) []*types.TransferDetails {
	transfers := make([]*types.TransferDetails, len(ids))

	for i, id := range ids {
		bz := store.Get([]byte(id))
		if bz == nil {
			transfers[i] = nil
			log.Error("Transfer is nil for id ", id)
			continue
		}

		transfer := &types.TransferDetails{}
		err := transfer.Unmarshal(bz)
		if err != nil {
			log.Error("getTransfer: Failed to unmarshal transfer out, err = ", err)
			transfers[i] = nil
			continue
		}

		transfers[i] = transfer
	}

	return transfers
}

///// Transfer Queue
func setTranferQueue(store cstypes.KVStore, chain string, transfers []*types.TransferDetails) {
	if len(transfers) == 0 {
		store.Delete([]byte(chain))
		return
	}

	ids := make([]string, len(transfers))
	for i, transfer := range transfers {
		ids[i] = transfer.Id
		fmt.Println("setTranferQueue transfer.Id = ", transfer.Id)
	}

	s := strings.Join(ids, ",")
	store.Set([]byte(chain), []byte(s))
}

func getTransferQueue(queueStore, transferStore cstypes.KVStore, chain string) []*types.TransferDetails {
	bz := queueStore.Get([]byte(chain))
	if bz == nil {
		return []*types.TransferDetails{}
	}

	s := string(bz)
	ids := strings.Split(s, ",")

	return getTransfers(transferStore, ids)
}

///// TxOutQueue
func setTxOutQueue(store cstypes.KVStore, chain string, txOuts []*types.TxOut) {
	queue := &types.TxOutQueue{
		TxOuts: txOuts,
	}
	bz, err := queue.Marshal()
	if err != nil {
		log.Error("setTxOutQueue: failed to marshal queue")
		return
	}

	store.Set([]byte(chain), bz)
}

func getTxOutQueue(store cstypes.KVStore, chain string) []*types.TxOut {
	bz := store.Get([]byte(chain))
	queue := &types.TxOutQueue{}
	err := queue.Unmarshal(bz)
	if err != nil {
		log.Error("getTxOutQueue: failed to unmarshal TxOutQueue")
		return nil
	}

	if queue.TxOuts == nil {
		queue.TxOuts = make([]*types.TxOut, 0)
	}

	return queue.TxOuts
}

///// Pending TxOut
func setPendingTxOut(store cstypes.KVStore, chain string, txOutInfo *types.PendingTxOutInfo) {
	if txOutInfo == nil {
		store.Delete([]byte(chain))
		return
	}

	bz, err := txOutInfo.Marshal()
	if err != nil {
		log.Error("setPendingTxOut: failed to marshal txOut")
		return
	}

	store.Set([]byte(txOutInfo.TxOut.Content.OutChain), bz)
}

func getPendingTxOutInfo(store cstypes.KVStore, chain string) *types.PendingTxOutInfo {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return nil
	}

	txOutInfo := &types.PendingTxOutInfo{}
	err := txOutInfo.Unmarshal(bz)
	if err != nil {
		log.Error("getPendingTxOut: failed to unmarshal txout")
		return nil
	}

	return txOutInfo
}

///// Chain Metadata
func setSolanaConfirmedBlock(store cstypes.KVStore, chain, signer, hash string, height int64) {
	key := getChainMetadataKey(chain, signer)

	meta := &types.ChainMetadata{}
	bz := store.Get(key)
	if bz == nil {
		meta.Chain = chain
		meta.Signer = signer
	} else {
		err := meta.Unmarshal(bz)
		if err != nil {
			log.Error("setSolanaConfirmedBlock: failed to unmarshal into chain metadata")
			return
		}

	}

	meta.SolanaRecentBlockHash = hash
	meta.SolanaRecentBlockHeight = height
	bz, err := meta.Marshal()
	if err != nil {
		log.Error("setSolanaConfirmedBlock: cannot marshal meta")
		return
	}

	store.Set(key, bz)
}

func getAllSolanaConfirmedBlock(store cstypes.KVStore, chain string) map[string]*types.ChainMetadata {
	ret := make(map[string]*types.ChainMetadata)

	begin := []byte(fmt.Sprintf("%s__", chain))
	end := []byte(fmt.Sprintf("%s__~", chain))

	iter := store.Iterator(begin, end)
	for ; iter.Valid(); iter.Next() {
		meta := &types.ChainMetadata{}
		err := meta.Unmarshal(iter.Value())
		if err == nil {
			ret[meta.Signer] = meta
		} else {
			log.Error("getAllSolanaConfirmedBlock: cannot unmarshal bz")
		}
	}

	return ret
}

///// Block Height
func setBlockHeight(store cstypes.KVStore, chain string, block *types.BlockHeight) {
	bz, err := block.Marshal()
	if err != nil {
		log.Errorf("Failed to save block height for chain %s", chain)
		return
	}

	store.Set([]byte(chain), bz)
}

func getBlockHeight(store cstypes.KVStore, chain string) *types.BlockHeight {
	bz := store.Get([]byte(chain))
	if bz == nil {
		return nil
	}

	block := &types.BlockHeight{}
	err := block.Unmarshal(bz)
	if err != nil {
		log.Errorf("Failed to unmarshal block height for chain %s", chain)
		return nil
	}

	return block
}

///// Tx Hash Index
func setTxHashIndex(store cstypes.KVStore, key string, value uint32) {
	store.Set([]byte(key), utils.Uint32ToBytes(value))
}

func getTxHashIndex(store cstypes.KVStore, key string) uint32 {
	bz := store.Get([]byte(key))
	if bz == nil {
		return 0
	}

	return utils.BytesToUint32(bz)
}

///// TxInDetails
func setTxInDetails(store cstypes.KVStore, txInId string, txIn *types.TxInDetails) {
	fmt.Println("msg.Data.TxIn.Id = ", txInId)

	bz, err := txIn.Marshal()
	if err != nil {
		log.Errorf("setTxInDetails: failed to marshal msg, err = %s", err)
		return
	}

	store.Set([]byte(txInId), bz)
}

func getTxInDetails(store cstypes.KVStore, txInId string) *types.TxInDetailsMsg {
	bz := store.Get([]byte(txInId))
	if bz == nil {
		return nil
	}

	msg := new(types.TxInDetailsMsg)
	err := msg.Unmarshal(bz)
	if err != nil {
		log.Errorf("getTxInDetails: failed to unmarshal TxInDetails with id %s, err = %s", txInId, err)
		return nil
	}

	return msg
}

///// Confirmed TxIn
func setConfirmedTxIn(store cstypes.KVStore, tx *types.ConfirmedTxIn) {
	bz, err := tx.Marshal()
	if err != nil {
		log.Errorf("setConfirmedTxIn: failed to marhsal confirmed TxIn")
		return
	}

	store.Set([]byte(tx.TxInId), bz)
}

func getConfirmedTxIn(store cstypes.KVStore, id string) *types.ConfirmedTxIn {
	bz := store.Get([]byte(id))
	if bz == nil {
		return nil
	}

	tx := new(types.ConfirmedTxIn)
	if err := tx.Unmarshal(bz); err != nil {
		log.Errorf("getConfirmedTxIn: failed to unmarshal confirmed TxIn")
		return nil
	}

	return tx
}

///// Vote Result
func addVoteResult(store cstypes.KVStore, hash string, signer string, result *types.VoteResult) {
	bz := utils.ToByte(result)
	if bz == nil {
		log.Errorf("addVoteResult: failed to convert result to byte")
		return
	}

	key := getVoteResultKey(hash, signer)
	store.Set(key, bz)
}

func getVoteResults(store cstypes.KVStore, hash string) map[string]types.VoteResult {
	begin := []byte(fmt.Sprintf("%s__", hash))
	end := []byte(fmt.Sprintf("%s__~", hash))

	ret := make(map[string]types.VoteResult)
	for iter := store.Iterator(begin, end); iter.Valid(); iter.Next() {
		key := string(iter.Key())
		if len(key) <= len(hash)+2 {
			continue
		}

		signer := key[len(hash)+2:]
		bz := utils.FromByteToInt(iter.Value())
		result := types.VoteResult(bz)
		ret[signer] = result
	}

	return ret
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
