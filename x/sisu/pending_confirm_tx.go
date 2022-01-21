package sisu

import (
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	eyesTypes "github.com/sisu-network/deyes/types"
)

// TODO: Complete implementation of this file.

// A struct that holds all the tx that we see (from deyes) but we cannot confirm yet since keysign
// result has not come back to us. It's possible that the transaction has been deployed by some
// other nodes.
//
// This struct uses on-memory storage only.
type PendingConfirmTx struct {
	lock  *sync.RWMutex
	cache *lru.Cache
}

type txData struct {
	tx    *eyesTypes.Tx
	chain string
	block int64
}

func NewPendingConfirmTx() *PendingConfirmTx {
	cache, err := lru.New(300)
	if err != nil {
		panic(err)
	}

	return &PendingConfirmTx{
		lock:  &sync.RWMutex{},
		cache: cache,
	}
}

func (p *PendingConfirmTx) getKey(chain, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}

func (p *PendingConfirmTx) AddTx(tx *eyesTypes.Tx, chain string, block int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	key := p.getKey(chain, tx.Hash)
	p.cache.Add(key, &txData{
		tx:    tx,
		chain: chain,
		block: block,
	})
}
