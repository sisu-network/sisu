package tss

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	deTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/utils"
)

const (
	KEY_OBSERVE_TX = "observed_tx_%s_%d_%s"
)

// Local storage for data of this specific node. This is not a db for application state.
type TssStorage struct {
	db *leveldb.Database
}

func NewTssStorage(file string) (*TssStorage, error) {
	utils.LogInfo("Initializing TSS storage...")
	db, err := leveldb.New(file, 1024, 500, "metrics_", false)
	if err != nil {
		return nil, err
	}
	return &TssStorage{
		db: db,
	}, nil
}

func (s *TssStorage) SaveTxs(txs *deTypes.Txs) {
	for _, tx := range txs.Arr {
		key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, txs.Chain, txs.Block, tx.Hash))
		utils.LogVerbose("Saving items with key", string(key))
		s.db.Put(key, tx.Serialized)
	}
}

func (s *TssStorage) GetObservedTx(chain string, blockHeight int64, hash string) []byte {
	key := []byte(fmt.Sprintf(KEY_OBSERVE_TX, chain, blockHeight, hash))
	bz, err := s.db.Get(key)
	if err != nil {
		utils.LogError("Cannot find item for key", string(key), ". Error = ", err)
		return nil
	}

	return bz
}
