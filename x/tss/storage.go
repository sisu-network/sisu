package tss

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/sisu-network/sisu/utils"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
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
func (s *TssStorage) AddTxOut(txOut *tssTypes.TxOut) {
	hash := txOut.GetHash()
	bz, err := txOut.Marshal()
	if err != nil {
		utils.LogError("cannot marshal txout, err =", err)
		return
	}

	key := fmt.Sprintf("tx_out_%s", hash)
	s.db.Put([]byte(key), bz)
}

func (s *TssStorage) GetTxOut(hash string) *tssTypes.TxOut {
	key := fmt.Sprintf("tx_out_%s", hash)
	bz, err := s.db.Get([]byte(key))
	if err != nil {
		utils.LogError("txout not found, hash =", hash)
		return nil
	}

	txOut := &tssTypes.TxOut{}
	if err := txOut.Unmarshal(bz); err != nil {
		utils.LogError("cannot unmashal txout, err =", err)
		return nil
	}

	return txOut
}

func (s *TssStorage) SavePubKey(chain string, pubKey []byte) {
	key := fmt.Sprintf("pubkey_%s", chain)
	s.db.Put([]byte(key), pubKey)
}

func (s *TssStorage) GetPubKey(chain string) []byte {
	key := fmt.Sprintf("pubkey_%s", chain)

	pubKey, err := s.db.Get([]byte(key))
	if err != nil {
		return nil
	}

	return pubKey
}
