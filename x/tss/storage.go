package tss

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/sisu-network/sisu/utils"
)

// TODO: Remove this class. Use db instead
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
