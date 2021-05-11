package ethchain

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"

	ethLog "github.com/ethereum/go-ethereum/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/rlp"
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
	config "github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
)

type ChainState int

const (
	TX_MAX_SIZE    = 128 * 1024
	COMMIT_TIMEOUT = time.Second * 5
)

var (
	BlackholeAddr = common.Address{
		1, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	}
)

var (
	lastAcceptedKey = []byte("lastAccepted")
)

type ETHChain struct {
	chainConfig    *config.ETHConfig
	backend        *eth.Ethereum
	cb             *dummy.ConsensusCallbacks
	mcb            *miner.MinerCallbacks
	backendCb      *types.BackendAPICallback
	chainState     ChainState
	chainMode      string
	signer         types.EIP155Signer
	gasLowestLimit *big.Int
	lastBlockState *state.StateDB
	genBlockDoneCh chan bool

	// Soft state
	softState *SoftState
	ssLock    *sync.RWMutex

	mu        *sync.RWMutex
	lastBlock *types.Block
	chainDb   ethdb.Database
}

// TODO: Remove initGenesis in both this repo & coreth. The coreth already had support for removing
// this variable. Keep it for now to check potential bugs.
func NewETHChain(
	chainConfig *config.ETHConfig,
	settings eth.Settings,
	onTxSubmitted func(*types.Transaction),
) *ETHChain {
	node, err := node.New(chainConfig.Node)
	if err != nil {
		panic(err)
	}

	chainDb, err := getChainDb(chainConfig)
	if err != nil {
		panic(err)
	}

	cb := new(dummy.ConsensusCallbacks)
	mcb := new(miner.MinerCallbacks)
	backendCb := new(types.BackendAPICallback)
	backendCb.OnTxSubmitted = onTxSubmitted
	backend, err := eth.New(node, chainConfig.Eth, cb, mcb, backendCb, chainDb, settings, true)
	if err != nil {
		panic(fmt.Sprintf("failed to create new eth backend due to %s", err))
	}
	backend.SetEtherbase(BlackholeAddr)

	chain := &ETHChain{
		chainConfig:    chainConfig,
		backend:        backend,
		cb:             cb,
		mcb:            mcb,
		chainDb:        chainDb,
		mu:             &sync.RWMutex{},
		ssLock:         &sync.RWMutex{},
		genBlockDoneCh: make(chan bool),
		signer:         types.NewEIP155Signer(chainConfig.Eth.Genesis.Config.ChainID),

		gasLowestLimit: new(big.Int).SetUint64(chainConfig.Eth.TxPool.PriceLimit),
	}
	chain.mcb.OnSealFinish = chain.OnSealFinish

	return chain
}

func getChainDb(chainConfig *config.ETHConfig) (ethdb.Database, error) {
	var db ethdb.Database
	var err error

	if chainConfig.UseInMemDb {
		utils.LogInfo("Use In memory for ETH")
		db = rawdb.NewMemoryDatabase()
	} else {
		utils.LogInfo("Use real DB for ETH")
		// Use level DB.
		// TODO: Create new configs.
		db, err = rawdb.NewLevelDBDatabase(chainConfig.DbPath, 1024, 500, "metrics_")
	}

	return db, err
}

func (self *ETHChain) Initialize() error {
	// Setting log level
	ethLog.Root().SetHandler(ethLog.LvlFilterHandler(
		ethLog.LvlDebug, ethLog.StreamHandler(os.Stderr, ethLog.TerminalFormat(false))))

	lastAcceptedBytes, lastAcceptedErr := self.chainDb.Get(lastAcceptedKey)
	var lastAccepted *types.Block
	// TODO: handle corrupted DB
	if lastAcceptedErr == nil {
		var hash common.Hash
		if err := rlp.DecodeBytes(lastAcceptedBytes, &hash); err == nil {
			if block := self.GetBlockByHash(hash); block == nil {
				utils.LogDebug("lastAccepted block not found in chaindb")
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

	self.BlockChain().UnlockIndexing()
	self.lastBlockState, _ = self.backend.BlockChain().State()

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

	handler := self.NewRPCHandler(time.Second * 10)
	handler.RegisterName("web3", &extra.Web3API{})
	handler.RegisterName("net", &extra.NetAPI{NetworkId: "1"})
	handler.RegisterName("evm", &extra.EvmApi{})

	self.AttachEthService(handler, []string{"eth", "personal", "txpool", "debug"})

	s.Initialize(self.chainConfig.Host, uint16(self.chainConfig.Port), []string{}, handler)

	go s.Dispatch()
}

func (self *ETHChain) Stop() {
	self.backend.Stop()
}

func (self *ETHChain) UnlockIndexing() {
	self.backend.BlockChain().UnlockIndexing()
}

func (self *ETHChain) NewRPCHandler(maximumDuration time.Duration) *rpc.Server {
	return rpc.NewServer(maximumDuration)
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

// BeginBlock is called when at the beginning of tendermint block. It prepares
func (self *ETHChain) BeginBlock() error {
	// Create a new softstate for execution in this block.
	self.createNewSoftState()

	return nil
}

func (self *ETHChain) CheckTx(tx *types.Transaction) error {
	self.ssLock.Lock()
	defer self.ssLock.Unlock()

	return self.valdiateTx(tx)
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
		return core.ErrInvalidSender
	}

	// Make sure gas price is higher than lowest limit.
	if tx.GasPriceIntCmp(self.gasLowestLimit) < 0 {
		return core.ErrUnderpriced
	}

	// Ensure the transaction adheres to nonce ordering
	if self.lastBlockState.GetNonce(from) > tx.Nonce() {
		return core.ErrNonceTooLow
	}

	// Transactor should have enough funds to cover the costs
	// cost == V + GP * GL
	if balance := self.lastBlockState.GetBalance(from); balance.Cmp(tx.Cost()) < 0 {
		return fmt.Errorf("insufficient funds for gas * price + value, balance: %d, cost: %d", balance, tx.Cost())
	}

	// Ensure the transaction has more gas than the basic tx fee.
	intrGas, err := core.IntrinsicGas(tx.Data(), tx.To() == nil, true, true)
	if err != nil {
		return err
	}
	if tx.Gas() < intrGas {
		return core.ErrIntrinsicGas
	}

	return nil
}

// DeliverTx adds a tx to the ETH tx pool. It does not do actual execution. The TX execution and
// db state change is done in the Commit function.
func (self *ETHChain) DeliverTx(tx *types.Transaction) (uint64, error) {
	utils.LogDebug("Delivering tx.....")

	var gasUsed uint64
	receipt, err := self.softState.ApplyTx(tx)
	if err == nil {
		gasUsed = receipt.GasUsed
	} else {
		gasUsed = 0
	}

	errs := self.backend.TxPool().AddRemotesSync([]*types.Transaction{tx})
	return gasUsed, errs[0]
}

// EndBlock tries to generate an ETH block
func (self *ETHChain) EndBlock() error {
	return nil
}

// Commit executes txs in the ETH mempool.
func (self *ETHChain) Commit() {
	utils.LogDebug("Start gen block")

	self.backend.Miner().GenBlock()

	// Block until we receive onSealFinish
	select {
	case <-self.genBlockDoneCh:
	case <-time.After(COMMIT_TIMEOUT):
	}
}

func (self *ETHChain) transition(newState ChainState) {
	self.chainState = newState
}

func (self *ETHChain) OnSealFinish(block *types.Block) error {
	utils.LogDebug("Done one seal")

	self.mu.Lock()
	self.lastBlock = block
	self.mu.Unlock()

	if err := self.Accept(block); err != nil {
		utils.LogError(err)
		return err
	}

	self.backend.BlockChain().InsertChain([]*types.Block{block})

	lastState, err := self.backend.BlockChain().State()
	if err == nil {
		self.mu.Lock()
		self.lastBlockState = lastState
		self.mu.Unlock()
	} else {
		utils.LogError("Cannot get last block state.")
	}

	self.genBlockDoneCh <- true

	return nil
}

func (self *ETHChain) createNewSoftState() {
	self.ssLock.Lock()
	defer self.ssLock.Unlock()

	dbState, err := self.backend.BlockChain().State()
	if err != nil {
		return
	}

	self.softState = NewSoftState(
		dbState,
		self.backend.BlockChain().LastAcceptedBlock(),
		self.chainConfig.Eth.Genesis.Config,
		*self.backend.BlockChain().GetVMConfig(),
		self.backend.BlockChain(),
	)
}

func (self *ETHChain) GetLastBlockDetails() ([]byte, *big.Int) {
	self.mu.RLock()
	defer self.mu.RUnlock()

	return self.lastBlock.Hash().Bytes(), self.lastBlock.Number()
}

func (self *ETHChain) GetBlockByHash(hash common.Hash) *types.Block {
	return self.backend.BlockChain().GetBlockByHash(hash)
}

func (self *ETHChain) Accept(block *types.Block) error {
	if err := self.BlockChain().Accept(block); err != nil {
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