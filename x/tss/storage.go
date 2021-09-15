package tss

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

// const (
// 	KEY_OBSERVE_TX = "observed_tx_%s_%d_%s"
// )

// Data structure that wraps around pending tx outs.
type PendingTxOutWrapper struct {
	InBlockHeight int64
	InChain       string
	OutChain      string
	InHash        string
	OutBytes      []byte
}

// Local storage for data of this specific node. This is not a db for application state.
type TssStorage struct {
	db *leveldb.Database

	pendingTx map[string]*tssTypes.ObservedTx

	// A map that remembers what transaction is assigned to which validators.
	pendingTxout map[int64][]*PendingTxOutWrapper // blockHeight -> txInHash (as string) -> validator address
}

func NewTssStorage(file string) (*TssStorage, error) {
	utils.LogInfo("Initializing TSS storage...")
	db, err := leveldb.New(file, 1024, 500, "metrics_", false)
	if err != nil {
		return nil, err
	}
	return &TssStorage{
		db:           db,
		pendingTx:    make(map[string]*tssTypes.ObservedTx),
		pendingTxout: make(map[int64][]*PendingTxOutWrapper),
	}, nil
}

func (s *TssStorage) getKey(format string, chain string, height int64, hash string) []byte {
	// Replace all the _ in the chain.
	chain = strings.Replace(chain, "_", "*", -1)
	return []byte(fmt.Sprintf(format, chain, height, hash))
}

// func (s *TssStorage) SaveObservedTx(chain string, blockHeight int64, hash string, txBytes []byte) {
// 	key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, chain, blockHeight, hash))
// 	s.db.Put(key, txBytes)
// }

// func (s *TssStorage) GetObservedTx(chain string, blockHeight int64, hash string) []byte {
// 	key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, chain, blockHeight, hash))
// 	bz, err := s.db.Get(key)
// 	if err != nil {
// 		return nil
// 	}

// 	return bz
// }

// Adds pending in txs to be processed at the end of the block.
func (s *TssStorage) AddPendingTx(msg *tssTypes.ObservedTx) {
	key := s.getKey("%s__%d__%s", msg.Chain, msg.BlockHeight, msg.TxHash)
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

// Saves an assignment of a validator for a particular out tx.
func (s *TssStorage) AddPendingTxOut(blockHeight int64, inChain string, inHash string, outChain string, outBytes []byte) {
	m := s.pendingTxout[blockHeight]
	if m == nil {
		m = make([]*PendingTxOutWrapper, 0)
	}
	newTxWrapper := &PendingTxOutWrapper{
		InBlockHeight: blockHeight,
		InChain:       inChain,
		OutChain:      outChain,
		InHash:        inHash,
		OutBytes:      outBytes,
	}
	m = append(m, newTxWrapper)

	s.pendingTxout[blockHeight] = m
}

// Returns a list of txs(both in and out) assigned to a specific validator at a particular block
// height.
func (s *TssStorage) GetPendingTxOutForValidator(blockHeight int64) []*PendingTxOutWrapper {
	return s.pendingTxout[blockHeight]
}

func (s *TssStorage) GetPendingTxOUt(blockHeight int64, inHash string) *PendingTxOutWrapper {
	m := s.pendingTxout[blockHeight]
	for _, tx := range m {
		if tx.InHash == inHash {
			return tx
		}
	}

	return nil
}

func (s *TssStorage) RemovePendingTxOut(blockHeight int64, inHash string) {
	arr := s.pendingTxout[blockHeight]
	for i, tx := range arr {
		if tx.InHash == inHash {
			s.pendingTxout[blockHeight] = append(arr[:i], arr[i+1:]...)
			break
		}
	}
}
