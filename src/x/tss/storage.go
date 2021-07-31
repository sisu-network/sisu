package tss

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

const (
	KEY_OBSERVE_TX = "observed_tx_%s_%d_%s"
)

// Local storage for data of this specific node. This is not a db for application state.
type TssStorage struct {
	db *leveldb.Database

	pendingTx map[string]*tssTypes.ObservedTx
}

func NewTssStorage(file string) (*TssStorage, error) {
	utils.LogInfo("Initializing TSS storage...")
	db, err := leveldb.New(file, 1024, 500, "metrics_", false)
	if err != nil {
		return nil, err
	}
	return &TssStorage{
		db:        db,
		pendingTx: make(map[string]*tssTypes.ObservedTx),
	}, nil
}

func (s *TssStorage) getKey(format string, chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	chain = strings.Replace(chain, "_", "*", -1)
	return []byte(fmt.Sprintf(format, chain, height, hash))
}

func (s *TssStorage) SaveObservedTx(chain string, blockHeight int64, hash string, txBytes []byte) {
	key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, chain, blockHeight, hash))
	s.db.Put(key, txBytes)
}

func (s *TssStorage) GetObservedTx(chain string, blockHeight int64, hash string) []byte {
	key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, chain, blockHeight, hash))
	bz, err := s.db.Get(key)
	if err != nil {
		return nil
	}

	return bz
}

// Adds pending in txs to be processed at the end of the block.
func (s *TssStorage) AddPendingTx(msg *tssTypes.ObservedTx) {
	key := s.getKey("%s_%d_%s", msg.Chain, msg.BlockHeight, msg.TxHash)
	s.pendingTx[string(key)] = msg
}

// Returns a list of txs that need to be processed at the end of a block.
func (s *TssStorage) GetAllPendingTxs() []*tssTypes.ObservedTx {
	txs := make([]*tssTypes.ObservedTx, 0)
	for _, value := range s.pendingTx {
		txs = append(txs, value)
	}

	return txs
}

// Clears all pending txs.
func (s *TssStorage) ClearPendingTxs() {
	s.pendingTx = make(map[string]*tssTypes.ObservedTx)
}
