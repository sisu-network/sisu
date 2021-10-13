package ethchain

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"sync"
	"time"

	ethLog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
	lru "github.com/hashicorp/golang-lru"
	"github.com/sisu-network/dcore/accounts"
	"github.com/sisu-network/dcore/accounts/keystore"
	"github.com/sisu-network/dcore/consensus/dummy"
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/core/rawdb"
	"github.com/sisu-network/dcore/core/state"
	"github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/extra"
	"github.com/sisu-network/dcore/miner"
	"github.com/sisu-network/dcore/node"
	"github.com/sisu-network/dcore/rpc"
	sisuCommon "github.com/sisu-network/sisu/common"
	config "github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type ChainState int

const (
	TX_MAX_SIZE       = 128 * 1024
	COMMIT_TIMEOUT    = time.Second * 5
	TX_CACHE_SIZE     = 4096
	SUBMIT_TX_TIMEOUT = time.Second * 6
)

var (
	BlackholeAddr = common.Address{
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}

	ERR_SHUTTING_DOWN = errors.New("Chain is shutting down")
	ERR_NOT_FOUND     = errors.New("not found")
)

var (
	lastAcceptedKey = []byte("lastAccepted")
	emptyRootHash   = common.Hash{}
)

type ETHChain struct {
	chainConfig     *config.ETHConfig
	backend         *eth.Ethereum
	cb              *dummy.ConsensusCallbacks
	mcb             *miner.MinerCallbacks
	backendCb       *types.BackendAPICallback
	chainState      ChainState
	chainMode       string
	signer          types.Signer
	gasLowestLimit  *big.Int
	lastBlockState  *state.StateDB
	lastBlockLock   *sync.RWMutex
	genBlockDoneCh  chan bool
	stopping        bool
	txSubmit        sisuCommon.TxSubmit
	acceptedTxCache *lru.Cache

	chainDb ethdb.Database
}

// TODO: Remove initGenesis in both this repo & coreth. The coreth already had support for removing
// this variable. Keep it for now to check potential bugs.
func NewETHChain(
	chainConfig *config.ETHConfig,
	settings eth.Settings,
	txSubmit sisuCommon.TxSubmit,
) *ETHChain {
	node, err := node.New(chainConfig.Node)
	if err != nil {
		panic(err)
	}

	chainDb, err := getChainDb(chainConfig)
	if err != nil {
		panic(err)
	}

	// TODO: Handle corrupted database here.
	_, lastAcceptedErr := chainDb.Get(lastAcceptedKey)
	initGenesis := lastAcceptedErr == ERR_NOT_FOUND
	utils.LogInfo("initGenesis = ", initGenesis)

	cb := new(dummy.ConsensusCallbacks)
	mcb := new(miner.MinerCallbacks)
	backendCb := new(types.BackendAPICallback)
	backend, err := eth.New(node, chainConfig.Eth, cb, mcb, backendCb, chainDb, settings, initGenesis)
	if err != nil {
		panic(fmt.Sprintf("failed to create new eth backend due to %s", err))
	}
	backend.SetEtherbase(BlackholeAddr)

	txCache, err := lru.New(TX_CACHE_SIZE)
	if err != nil {
		panic(err)
	}

	chain := &ETHChain{
		chainConfig:     chainConfig,
		backend:         backend,
		cb:              cb,
		mcb:             mcb,
		chainDb:         chainDb,
		txSubmit:        txSubmit,
		lastBlockLock:   &sync.RWMutex{},
		acceptedTxCache: txCache,
		genBlockDoneCh:  make(chan bool),
		signer:          types.NewEIP2930Signer(chainConfig.Eth.Genesis.Config.ChainID),

		gasLowestLimit: new(big.Int).SetUint64(chainConfig.Eth.TxPool.PriceLimit),
	}
	chain.mcb.OnSealFinish = chain.OnSealFinish
	backendCb.OnTxSubmitted = chain.onEthTxSubmitted

	return chain
}

func getChainDb(chainConfig *config.ETHConfig) (ethdb.Database, error) {
	var db ethdb.Database
	var err error

	dbPath := filepath.Join(chainConfig.Dir, "leveldb")
	db, err = rawdb.NewLevelDBDatabase(dbPath, 1024, 500, "metrics_", false)

	return db, err
}

func (self *ETHChain) Initialize() error {
	// Setting log level
	ethLog.Root().SetHandler(ethLog.LvlFilterHandler(
		ethLog.LvlCrit, ethLog.StreamHandler(os.Stderr, ethLog.TerminalFormat(false))))

	// TODO: handle corrupted DB
	lastAcceptedBytes, lastAcceptedErr := self.chainDb.Get(lastAcceptedKey)
	var lastAccepted *types.Block
	utils.LogInfo("lastAcceptedErr = ", lastAcceptedErr)

	if lastAcceptedErr == nil {
		var hash common.Hash
		if err := rlp.DecodeBytes(lastAcceptedBytes, &hash); err == nil {
			if block := self.GetBlockByHash(hash); block == nil {
				utils.LogInfo("lastAccepted block not found in chaindb")
			} else {
				lastAccepted = block
			}
		}
	}

	switch {
	case lastAccepted == nil:
		lastAccepted = self.GetGenesisBlock()
	}

	if err := self.Accept(lastAccepted); err != nil {
		return fmt.Errorf("could not initialize VM with last accepted hash %s: %w", lastAccepted.Hash(), err)
	}

	self.lastBlockState, _ = self.backend.BlockChain().State()
	utils.LogInfo("lastBlock = ", lastAccepted.Number())

	return nil
}

func (self *ETHChain) Start() {
	self.backend.StartMining()
	self.backend.Start()
	self.startApiServer()

	self.BlockChain().Accept(self.GetGenesisBlock())
}

func (self *ETHChain) startApiServer() {
	s := &Server{}

	handler := self.NewRPCHandler()
	handler.RegisterName("web3", &extra.Web3API{})
	handler.RegisterName("net", &extra.NetAPI{NetworkId: "1"})
	handler.RegisterName("evm", &extra.EvmApi{})

	self.AttachEthService(handler, []string{"eth", "personal", "txpool", "debug"})

	s.Initialize(self.chainConfig.Host, uint16(self.chainConfig.Port), []string{}, handler)

	go s.Dispatch()
}

func (self *ETHChain) Stop() {
	self.stopping = true
	utils.LogInfo("Stopping ETH backend....")
	self.backend.Stop()
	utils.LogInfo("ETH backend stopped")
}

func (self *ETHChain) NewRPCHandler() *rpc.Server {
	return rpc.NewServer()
}

func (self *ETHChain) BlockChain() *core.BlockChain {
	return self.backend.BlockChain()
}

func (self *ETHChain) AttachEthService(handler *rpc.Server, namespaces []string) {
	nsmap := make(map[string]bool)
	for _, ns := range namespaces {
		nsmap[ns] = true
	}
	for _, api := range self.backend.APIs() {
		if nsmap[api.Namespace] {
			handler.RegisterName(api.Namespace, api.Service)
		}
	}
}

func (self *ETHChain) GetGenesisBlock() *types.Block {
	return self.backend.BlockChain().Genesis()
}

func (self *ETHChain) BeginBlock(timestamp time.Time) error {
	self.backend.Miner().PrepareNewBlock(timestamp)
	return nil
}

// Validates a transaction. Many part of this function is borrowed from tx_pool.validateTx().
func (self *ETHChain) valdiateTx(tx *types.Transaction) error {
	// Reject transactions over defined size to prevent DOS attacks
	if uint64(tx.Size()) > TX_MAX_SIZE {
		return core.ErrOversizedData
	}

	// Transactions can't be negative. This may never happen using RLP decoded
	// transactions but may occur if you create a transaction using the RPC.
	if tx.Value().Sign() < 0 {
		return core.ErrNegativeValue
	}

	// Ensure the transaction doesn't exceed the current block limit gas.
	if self.chainConfig.Eth.Genesis.GasLimit < tx.Gas() {
		return core.ErrGasLimit
	}

	// Make sure the transaction is signed properly
	from, err := types.Sender(self.signer, tx)
	if err != nil {
		utils.LogError("Validation failed: ", err)
		return core.ErrInvalidSender
	}

	// Make sure gas price is higher than lowest limit.
	if tx.GasPriceIntCmp(self.gasLowestLimit) < 0 {
		return core.ErrUnderpriced
	}

	if err := self.checkNonceAndBalance(tx, from); err != nil {
		return err
	}

	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := core.IntrinsicGas(tx.Data(), tx.AccessList(), tx.To() == nil, true, true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return core.ErrIntrinsicGas
	}

	return nil
}

func (self *ETHChain) checkNonceAndBalance(tx *types.Transaction, from common.Address) error {
	self.lastBlockLock.Lock()
	defer self.lastBlockLock.Unlock()

	// Ensure the transaction adheres to nonce ordering
	if self.lastBlockState.GetNonce(from) > tx.Nonce() {
		utils.LogDebug("self.lastBlockState.GetNonce(from) = ", self.lastBlockState.GetNonce(from))
		return core.ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if balance := self.lastBlockState.GetBalance(from); balance.Cmp(tx.Cost()) < 0 {
		return fmt.Errorf("insufficient funds for gas * price + value, balance: %d, cost: %d", balance, tx.Cost())
	}

	return nil
}

// EndBlock tries to generate an ETH block
func (self *ETHChain) EndBlock() error {
	if self.stopping {
		return ERR_SHUTTING_DOWN
	}

	self.backend.Miner().GenBlock()

	// Block until we receive onSealFinish
	select {
	case <-self.genBlockDoneCh:
	case <-time.After(COMMIT_TIMEOUT):
	}

	return nil
}

// Commit executes txs in the ETH mempool.
func (self *ETHChain) Commit() {
}

func (self *ETHChain) transition(newState ChainState) {
	self.chainState = newState
}

func (self *ETHChain) OnSealFinish(block *types.Block) error {
	utils.LogDebug("Block is sealed, number =", block.Number())

	if err := self.Accept(block); err != nil {
		utils.LogError(err)
		return err
	}

	self.backend.BlockChain().InsertChain([]*types.Block{block})

	lastState, err := self.backend.BlockChain().State()
	if err == nil {
		self.lastBlockLock.Lock()
		self.lastBlockState = lastState
		self.lastBlockLock.Unlock()
	} else {
		utils.LogError("Cannot get last block state.")
	}

	self.genBlockDoneCh <- true

	size, err := self.PendingSize()
	utils.LogDebug("Pending size = ", size)

	return nil
}

func (self *ETHChain) PendingSize() (int, error) {
	pending, err := self.backend.TxPool().Pending()
	count := 0
	for _, txs := range pending {
		count += len(txs)
	}
	return count, err
}

func (self *ETHChain) GetBlockByHash(hash common.Hash) *types.Block {
	return self.backend.BlockChain().GetBlockByHash(hash)
}

func (self *ETHChain) Accept(block *types.Block) error {
	if err := self.BlockChain().Accept(block); err != nil {
		utils.LogError("Failed to accept block", err)
		return err
	}

	b, err := rlp.EncodeToBytes(block.Hash())
	if err != nil {
		return err
	}

	return self.chainDb.Put(lastAcceptedKey, b)
}

func (self *ETHChain) LastAcceptedBlock() *types.Block {
	return self.BlockChain().LastAcceptedBlock()
}

// fetchKeystore retrieves the encrypted keystore from the account manager.
func fetchKeystore(am *accounts.Manager) (*keystore.KeyStore, error) {
	if ks := am.Backends(keystore.KeyStoreType); len(ks) > 0 {
		return ks[0].(*keystore.KeyStore), nil
	}
	return nil, errors.New("local keystore not used")
}

// ImportAccounts is used only in dev mode
func (self *ETHChain) ImportAccounts() {
	am := self.backend.APIBackend.AccountManager()
	ks, err := fetchKeystore(am)
	if err != nil {
		return
	}

	wallet := utils.GetLocalWallet()
	accounts := utils.GetLocalAccounts()

	if len(ks.Accounts()) <= 10 {
		utils.LogDebug("Importing accounts...")

		for _, account := range accounts {
			privateKey, err := wallet.PrivateKey(account)
			if err != nil {
				return
			}

			ks.ImportECDSA(privateKey, utils.LOCAL_KEYSTORE_PASS)
		}
	}

	// Unlocking accounts
	for _, account := range accounts {
		ks.Unlock(account, utils.LOCAL_KEYSTORE_PASS)
	}
	utils.LogDebug("Done importing. Accounts length = ", len(ks.Accounts()))
}

func (self *ETHChain) CheckTx(txs []*types.Transaction) error {
	err := fmt.Errorf("No ETH transaction is accepted")
	utils.LogDebug("Checking tx ....")

	errs := self.backend.TxPool().AddRemotesSync(txs)
	for i, tx := range txs {
		if errs[i] == nil {
			self.acceptedTxCache.Add(tx.Hash().String(), tx)
			utils.LogDebug("Tx is accepted", tx.Hash().String())
			err = nil
		} else {
			utils.LogError("Accept tx error: ", i, errs[i])
		}
	}

	return err
}

// DeliverTx adds a tx to the ETH tx pool. It does not do actual execution. The TX execution and
// db state change is done in the Commit function.
func (self *ETHChain) DeliverTx(tx *types.Transaction) (*types.Receipt, common.Hash, error) {
	if self.stopping {
		return nil, emptyRootHash, ERR_SHUTTING_DOWN
	}
	utils.LogDebug("Delivering tx.....")

	receipt, rootHash, err := self.backend.Miner().ExecuteTxSync(tx)
	if err != nil {
		utils.LogError("Failed to execute the transaction", err)
		return nil, emptyRootHash, err
	}

	return receipt, rootHash, nil
}

func (self *ETHChain) onEthTxSubmitted(tx *types.Transaction) error {
	_, ok := self.acceptedTxCache.Get(tx.Hash().String())
	if ok {
		return fmt.Errorf("The transaction is already accepted for execution.")
	}

	if err := self.valdiateTx(tx); err != nil {
		utils.LogError("Failed to validate tx", err)
		return err
	}

	js, err := tx.MarshalJSON()
	if err != nil {
		return err
	}

	if err := self.txSubmit.SubmitEThTx(js); err != nil {
		return err
	}

	deadline := time.Now().Add(SUBMIT_TX_TIMEOUT)
	for {
		_, ok = self.acceptedTxCache.Get(tx.Hash().String())
		if ok {
			break
		}

		if time.Now().After(deadline) {
			break
		}

		time.Sleep(1000)
	}

	// Check if the tx pool has the tx or not.
	_, ok = self.acceptedTxCache.Get(tx.Hash().String())
	if !ok {
		utils.LogError("Cannot find transaction in the pool.", tx.Hash().String())
		return fmt.Errorf("Failed to add transaction to the pool")
	}

	return nil
}
