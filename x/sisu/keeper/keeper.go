package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	cstypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

var _ Keeper = (*DefaultKeeper)(nil)

// go:generate mockgen -source x/sisu/keeper/keeper.go -destination=tests/mock/x/sisu/keeper/keeper.go -package=mock
type Keeper interface {
	// Debug
	PrintStore(ctx sdk.Context, name string)
	PrintStoreKeys(ctx sdk.Context, name string)

	// TxRecord
	SaveTxRecord(ctx sdk.Context, hash []byte, signer string) int
	GetVoters(ctx sdk.Context, hash []byte) []string
	GetVotersInAccAddress(ctx sdk.Context, hash []byte) []sdk.AccAddress

	// TxRecordProcessed
	ProcessTxRecord(ctx sdk.Context, hash []byte)
	IsTxRecordProcessed(ctx sdk.Context, hash []byte) bool

	// Keygen
	SaveKeygen(ctx sdk.Context, msg *types.Keygen)
	IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool
	IsKeygenAddress(ctx sdk.Context, keyType string, address string) bool
	GetKeygenPubkey(ctx sdk.Context, keyType string) []byte
	GetAllKeygenPubkeys(ctx sdk.Context) map[string][]byte

	// Keygen Result
	SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner)
	GetAllKeygenResult(ctx sdk.Context, keygenType string, index int32) []*types.KeygenResultWithSigner

	// Contract
	SaveContract(ctx sdk.Context, msg *types.Contract, saveByteCode bool)
	SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool)
	IsContractExisted(ctx sdk.Context, msg *types.Contract) bool
	GetContract(ctx sdk.Context, chain string, hash string, includeByteCode bool) *types.Contract
	GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract
	UpdateContractAddress(ctx sdk.Context, chain string, hash string, address string)
	UpdateContractsStatus(ctx sdk.Context, chain string, contractHash string, status string)
	GetLatestContractAddressByName(ctx sdk.Context, chain, name string) string

	// Contract Address
	CreateContractAddress(ctx sdk.Context, chain string, txOutHash string, address string)
	IsContractExistedAtAddress(ctx sdk.Context, chain string, address string) bool

	// TxIn
	SaveTxIn(ctx sdk.Context, msg *types.TxIn)
	IsTxInExisted(ctx sdk.Context, msg *types.TxIn) bool

	// TxOut
	SaveTxOut(ctx sdk.Context, msg *types.TxOut)
	IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool
	GetTxOut(ctx sdk.Context, outChain, hash string) *types.TxOut

	// TxOutSig
	// TODO: Add unconfirmed tx store
	SaveTxOutSig(ctx sdk.Context, msg *types.TxOutSig)
	GetTxOutSig(ctx sdk.Context, outChain, hashWithSig string) *types.TxOutSig

	// TxOutConfirm
	SaveTxOutConfirm(ctx sdk.Context, msg *types.TxOutContractConfirm)
	IsTxOutConfirmExisted(ctx sdk.Context, outChain, hash string) bool

	// Gas Price Record
	SetGasPrice(ctx sdk.Context, msg *types.GasPriceMsg)
	GetGasPriceRecord(ctx sdk.Context, height int64) *types.GasPriceRecord

	// Chain
	SaveChain(ctx sdk.Context, chain *types.Chain)
	GetChain(ctx sdk.Context, chain string) *types.Chain
	GetAllChains(ctx sdk.Context) map[string]*types.Chain

	// Token Price
	SetTokenPrices(ctx sdk.Context, blockHeight uint64, msg *types.UpdateTokenPrice)
	GetAllTokenPricesRecord(ctx sdk.Context) map[string]*types.TokenPriceRecords

	// Set Tokens
	SetTokens(ctx sdk.Context, tokens map[string]*types.Token)
	GetTokens(ctx sdk.Context, tokens []string) map[string]*types.Token
	GetAllTokens(ctx sdk.Context) map[string]*types.Token

	// Nodes
	SaveNode(ctx sdk.Context, node *types.Node)
	LoadNodesByStatus(ctx sdk.Context, status types.NodeStatus) []*types.Node
	SetValidators(ctx sdk.Context, nodes []*types.Node) ([]*types.Node, error)
	UpdateNodeStatus(ctx sdk.Context, consKey []byte, status types.NodeStatus)

	// Liquidities
	SetLiquidities(ctx sdk.Context, liquidities map[string]*types.Liquidity)
	GetLiquidity(ctx sdk.Context, chain string) *types.Liquidity
	GetAllLiquidities(ctx sdk.Context) map[string]*types.Liquidity

	// Params
	SaveParams(ctx sdk.Context, params *types.Params)
	GetParams(ctx sdk.Context) *types.Params

	// Slash
	IncSlashToken(ctx sdk.Context, amount int64, addresses ...sdk.AccAddress) error
	DecSlashToken(ctx sdk.Context, amount int64, addresses ...sdk.AccAddress) error
	GetSlashToken(ctx sdk.Context, address sdk.AccAddress) (int64, error)

	// Nodes balance
	IncBondBalance(ctx sdk.Context, address sdk.AccAddress, amount int64) error
	DecBondBalance(ctx sdk.Context, address sdk.AccAddress, amount int64) error
	GetBondBalance(ctx sdk.Context, address sdk.AccAddress) (int64, error)
	GetTopBondBalance(ctx sdk.Context, n int) []sdk.AccAddress

	// Validator update
	SaveIncomingValidatorUpdates(ctx sdk.Context, validatorUpdates abci.ValidatorUpdates) error
	GetIncomingValidatorUpdates(ctx sdk.Context) abci.ValidatorUpdates
	GetValidatorUpdateIndex(ctx sdk.Context) int
	ClearValidatorUpdates(ctx sdk.Context)

	// Dheart IP address
	SaveDheartIPAddress(ctx sdk.Context, address sdk.AccAddress, ip string) error
	GetAllDheartIPAddresses(ctx sdk.Context) []AccAddressDheartIP
}

type DefaultKeeper struct {
	storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) Keeper {
	keeper := &DefaultKeeper{
		storeKey: storeKey,
	}

	return keeper
}

///// TxRecord
func (k *DefaultKeeper) SaveTxRecord(ctx sdk.Context, hash []byte, signer string) int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecord)
	return saveTxRecord(store, hash, signer)
}

func (k *DefaultKeeper) GetVoters(ctx sdk.Context, hash []byte) []string {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecord)
	return getVotersString(store, hash)
}

func (k *DefaultKeeper) GetVotersInAccAddress(ctx sdk.Context, hash []byte) []sdk.AccAddress {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecord)
	voters := getVotersString(store, hash)

	addrs := make([]sdk.AccAddress, 0)
	for _, v := range voters {
		a, err := sdk.AccAddressFromBech32(v)
		if err != nil {
			log.Error("invalid address: ", v)
			continue
		}

		addrs = append(addrs, a)
	}

	return addrs
}

///// TxRecordProcessed
func (k *DefaultKeeper) ProcessTxRecord(ctx sdk.Context, hash []byte) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecordProcessed)
	processTxRecord(store, hash)
}

func (k *DefaultKeeper) IsTxRecordProcessed(ctx sdk.Context, hash []byte) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxRecordProcessed)
	return isTxRecordProcessed(store, hash)
}

///// Keygen

func (k *DefaultKeeper) SaveKeygen(ctx sdk.Context, msg *types.Keygen) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	saveKeygen(store, msg)
}

func (k *DefaultKeeper) IsKeygenExisted(ctx sdk.Context, keyType string, index int) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return isKeygenExisted(store, keyType, index)
}

func (k *DefaultKeeper) IsKeygenAddress(ctx sdk.Context, keyType string, address string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return isKeygenAddress(store, keyType, address)
}

func (k *DefaultKeeper) GetKeygenPubkey(ctx sdk.Context, keyType string) []byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return getKeygenPubkey(store, keyType)
}

func (k *DefaultKeeper) GetAllKeygenPubkeys(ctx sdk.Context) map[string][]byte {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	return getAllKeygenPubkeys(store)
}

///// Keygen Result

func (k *DefaultKeeper) SaveKeygenResult(ctx sdk.Context, signerMsg *types.KeygenResultWithSigner) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResultWithSigner)
	saveKeygenResult(store, signerMsg)
}

func (k *DefaultKeeper) GetAllKeygenResult(ctx sdk.Context, keygenType string, index int32) []*types.KeygenResultWithSigner {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygenResultWithSigner)
	return getAllKeygenResult(store, keygenType, index)
}

///// Contract
func (k *DefaultKeeper) SaveContract(ctx sdk.Context, msg *types.Contract, saveByteCode bool) {
	k.SaveContracts(ctx, []*types.Contract{msg}, saveByteCode)
}

func (k *DefaultKeeper) SaveContracts(ctx sdk.Context, msgs []*types.Contract, saveByteCode bool) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if saveByteCode {
		byteCodeStore = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)
	}

	saveContracts(contractStore, byteCodeStore, msgs)

	// After saving contracts, also save contract address for each contract type
	contractNameStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractName)
	for _, msg := range msgs {
		saveContractAddressForName(contractNameStore, msg)
	}
}

func (k *DefaultKeeper) IsContractExisted(ctx sdk.Context, msg *types.Contract) bool {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	return isContractExisted(contractStore, msg)
}

func (k *DefaultKeeper) GetContract(ctx sdk.Context, chain string, hash string, includeByteCode bool) *types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	var byteCodeStore cstypes.KVStore
	byteCodeStore = nil
	if includeByteCode {
		byteCodeStore = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)
	}

	return getContract(contractStore, byteCodeStore, chain, hash)
}

func (k *DefaultKeeper) GetPendingContracts(ctx sdk.Context, chain string) []*types.Contract {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	byteCodeStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractByteCode)

	return getPendingContracts(contractStore, byteCodeStore, chain)
}

func (k *DefaultKeeper) UpdateContractAddress(ctx sdk.Context, chain string, hash string, address string) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	updateContractAddress(contractStore, chain, hash, address)
}

func (k *DefaultKeeper) UpdateContractsStatus(ctx sdk.Context, chain string, contractHash string, status string) {
	contractStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	updateContractsStatus(contractStore, chain, contractHash, status)
}

func (k *DefaultKeeper) GetLatestContractAddressByName(ctx sdk.Context, chain, name string) string {
	contractNameStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractName)
	return getContractAddressByName(contractNameStore, chain, name)
}

///// Contract Address
func (k *DefaultKeeper) CreateContractAddress(ctx sdk.Context, chain string, txOutHash string, address string) {
	caStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractAddress)
	txOutStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)

	createContractAddress(caStore, txOutStore, chain, txOutHash, address)
}

func (k *DefaultKeeper) IsContractExistedAtAddress(ctx sdk.Context, chain string, address string) bool {
	caStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixContractAddress)

	return isContractExistedAtAddress(caStore, chain, address)
}

///// TxIn
func (k *DefaultKeeper) SaveTxIn(ctx sdk.Context, msg *types.TxIn) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	saveTxIn(store, msg)
}

func (k *DefaultKeeper) IsTxInExisted(ctx sdk.Context, msg *types.TxIn) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxIn)
	return isTxInExisted(store, msg)
}

///// TxOut
func (k *DefaultKeeper) SaveTxOut(ctx sdk.Context, msg *types.TxOut) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	saveTxOut(store, msg)
}

func (k *DefaultKeeper) IsTxOutExisted(ctx sdk.Context, msg *types.TxOut) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	return isTxOutExisted(store, msg)
}

func (k *DefaultKeeper) GetTxOut(ctx sdk.Context, outChain, hash string) *types.TxOut {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	return getTxOut(store, outChain, hash)
}

func (k *DefaultKeeper) GetTxOutSig(ctx sdk.Context, outChain, hashWithSig string) *types.TxOutSig {
	withSigStore := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutSig)
	txOutSig := getTxOutSig(withSigStore, outChain, hashWithSig)

	return txOutSig
}

///// TxOutSig
func (k *DefaultKeeper) SaveTxOutSig(ctx sdk.Context, msg *types.TxOutSig) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutSig)
	saveTxOutSig(store, msg)
}

///// TxOutConfirm
func (k *DefaultKeeper) SaveTxOutConfirm(ctx sdk.Context, msg *types.TxOutContractConfirm) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutContractConfirm)
	saveTxOutConfirm(store, msg)
}

func (k *DefaultKeeper) IsTxOutConfirmExisted(ctx sdk.Context, chain, hash string) bool {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOutContractConfirm)
	return isTxOutConfirmExisted(store, chain, hash)
}

///// GasPrice
func (k *DefaultKeeper) SetGasPrice(ctx sdk.Context, msg *types.GasPriceMsg) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGasPrice)
	saveGasPrice(store, msg)
}

func (k *DefaultKeeper) GetGasPriceRecord(ctx sdk.Context, height int64) *types.GasPriceRecord {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixGasPrice)
	return getGasPriceRecord(store, height)
}

///// Network gas price

func (k *DefaultKeeper) SaveChain(ctx sdk.Context, chain *types.Chain) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)

	saveChain(store, chain)
}

func (k *DefaultKeeper) GetChain(ctx sdk.Context, chain string) *types.Chain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)
	return getChain(store, chain)
}

func (k *DefaultKeeper) GetAllChains(ctx sdk.Context) map[string]*types.Chain {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixChain)
	return getAllChains(store)
}

///// Token Prices

func (k *DefaultKeeper) SetTokenPrices(ctx sdk.Context, blockHeight uint64, msg *types.UpdateTokenPrice) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTokenPrices)
	setTokenPrices(store, blockHeight, msg)
}

func (k *DefaultKeeper) GetAllTokenPricesRecord(ctx sdk.Context) map[string]*types.TokenPriceRecords {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixTokenPrices)
	return getAllTokenPrices(store)
}

///// Calculated token prices

func (k *DefaultKeeper) SetTokens(ctx sdk.Context, prices map[string]*types.Token) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	setTokens(store, prices)
}

func (k *DefaultKeeper) GetTokens(ctx sdk.Context, tokenIds []string) map[string]*types.Token {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	return getTokens(store, tokenIds)
}

func (k *DefaultKeeper) GetAllTokens(ctx sdk.Context) map[string]*types.Token {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixToken)
	return getAllTokens(store)
}

///// Nodes

func (k *DefaultKeeper) SaveNode(ctx sdk.Context, node *types.Node) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	saveNode(store, node)
}

func (k *DefaultKeeper) LoadNodesByStatus(ctx sdk.Context, status types.NodeStatus) []*types.Node {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	return loadNodesByStatus(store, status)
}

func (k *DefaultKeeper) SetValidators(ctx sdk.Context, vals []*types.Node) ([]*types.Node, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	return setValidators(store, vals)
}

func (k *DefaultKeeper) UpdateNodeStatus(ctx sdk.Context, consKey []byte, status types.NodeStatus) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNode)
	updateNodeStatus(store, consKey, status)
}

///// Liquidities

func (k *DefaultKeeper) SetLiquidities(ctx sdk.Context, liquids map[string]*types.Liquidity) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixLiquidity)
	setLiquidities(store, liquids)
}

func (k *DefaultKeeper) GetLiquidity(ctx sdk.Context, chain string) *types.Liquidity {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixLiquidity)
	return getLiquidity(store, chain)
}

func (k *DefaultKeeper) GetAllLiquidities(ctx sdk.Context) map[string]*types.Liquidity {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixLiquidity)
	return getAllLiquidities(store)
}

///// Params
func (k *DefaultKeeper) SaveParams(ctx sdk.Context, params *types.Params) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixParams)
	saveParams(store, params)
}

func (k *DefaultKeeper) GetParams(ctx sdk.Context) *types.Params {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixParams)
	return getParams(store)
}

///// Slash
func (k *DefaultKeeper) IncSlashToken(ctx sdk.Context, amount int64, addresses ...sdk.AccAddress) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixSlash)

	for _, a := range addresses {
		if err := incOrDecSlashToken(store, amount, a); err != nil {
			return err
		}
	}

	return nil
}

func (k *DefaultKeeper) DecSlashToken(ctx sdk.Context, amount int64, addresses ...sdk.AccAddress) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixSlash)

	for _, a := range addresses {
		if err := incOrDecSlashToken(store, -amount, a); err != nil {
			return err
		}
	}

	return nil
}

func (k *DefaultKeeper) GetSlashToken(ctx sdk.Context, address sdk.AccAddress) (int64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixSlash)
	return getCurSlashToken(store, address)
}

///// Bond balances
func (k *DefaultKeeper) IncBondBalance(ctx sdk.Context, address sdk.AccAddress, amount int64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNodeBalance)
	return incOrDecBondBalance(store, address, amount)
}

func (k *DefaultKeeper) DecBondBalance(ctx sdk.Context, address sdk.AccAddress, amount int64) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNodeBalance)
	return incOrDecBondBalance(store, address, -amount)
}

func (k *DefaultKeeper) GetBondBalance(ctx sdk.Context, address sdk.AccAddress) (int64, error) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNodeBalance)
	return getBondBalance(store, address)
}

// GetTopBalance if n is -1, get all balances
func (k *DefaultKeeper) GetTopBondBalance(ctx sdk.Context, n int) []sdk.AccAddress {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixNodeBalance)
	return getTopBalances(store, n)
}

///// Validator updates
func (k *DefaultKeeper) SaveIncomingValidatorUpdates(ctx sdk.Context, validatorUpdates abci.ValidatorUpdates) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixValidatorUpdate)
	return saveValidatorUpdates(store, validatorUpdates)
}

func (k *DefaultKeeper) ClearValidatorUpdates(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixValidatorUpdate)
	resetValidatorUpdate(store)
}

func (k *DefaultKeeper) GetIncomingValidatorUpdates(ctx sdk.Context) abci.ValidatorUpdates {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixValidatorUpdate)
	return getValidatorUpdates(store)
}

func (k *DefaultKeeper) GetValidatorUpdateIndex(ctx sdk.Context) int {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixValidatorUpdateIndex)
	return increaseValidatorUpdateIndex(store)
}

///// Set Dheart IP address
func (k *DefaultKeeper) SaveDheartIPAddress(ctx sdk.Context, accAddr sdk.AccAddress, ip string) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixDheartIPAddress)
	return saveDheartIPAddress(store, accAddr, ip)
}

func (k *DefaultKeeper) GetAllDheartIPAddresses(ctx sdk.Context) []AccAddressDheartIP {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), prefixDheartIPAddress)
	return getAllDheartIPAddresses(store)
}

///// Debug

func (k *DefaultKeeper) getStoreFromName(ctx sdk.Context, name string) cstypes.KVStore {
	var store cstypes.KVStore
	switch name {
	case "keygen":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixKeygen)
	case "contract":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixContract)
	case "txOut":
		store = prefix.NewStore(ctx.KVStore(k.storeKey), prefixTxOut)
	}

	return store
}

// PrintStore is a debug function
func (k *DefaultKeeper) PrintStore(ctx sdk.Context, name string) {
	log.Info("======== DEBUGGING PrintStore")
	log.Info("Printing ALL values in store ", name)

	store := k.getStoreFromName(ctx, name)
	if store != nil {
		printStore(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}

func (k *DefaultKeeper) PrintStoreKeys(ctx sdk.Context, name string) {
	log.Info("======== DEBUGGING PrintStoreKeys")
	log.Info("Printing ALL values in store ", name)
	store := k.getStoreFromName(ctx, name)
	if store != nil {
		printStoreKeys(store)
	} else {
		log.Info("Invalid name")
	}

	log.Info("======== END OF DEBUGGING")
}
