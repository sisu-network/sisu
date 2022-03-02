package sisu

import (
	"fmt"
	"sync"

	"github.com/sisu-network/lib/log"
)

type TxStatus int64

const (
	TxStatusUnknown TxStatus = iota
	TxStatusCreated
	TxStatusSigned
	TxStatusConfirmed
)

type txObject struct {
	status TxStatus
	chain  string
	hash   string // hash without signature
}

func newTxObject(chain string, hash string) *txObject {
	return &txObject{
		status: TxStatusCreated,
		chain:  chain,
		hash:   hash,
	}
}

type TransactionTracker interface {
	AddTransaction(chain string, hash string)
	UpdateStatus(chain string, hash string, status TxStatus)
}

type DefaultTransactionTracker struct {
	txs *sync.Map
}

func NewTransactionTracker() TransactionTracker {
	return &DefaultTransactionTracker{
		txs: &sync.Map{},
	}
}

func (t *DefaultTransactionTracker) getTxoKey(chain string, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}

func (t *DefaultTransactionTracker) AddTransaction(chain string, hash string) {
	key := t.getTxoKey(chain, hash)

	if _, ok := t.txs.Load(key); ok {
		log.Warnf("Tx has been added to the tracker, chain and hash = ", chain, hash)
		return
	}

	t.txs.Store(key, newTxObject(chain, hash))
}

func (t *DefaultTransactionTracker) UpdateStatus(chain string, hash string, status TxStatus) {
	value, ok := t.txs.Load(t.getTxoKey(chain, hash))
	if ok {
		tx := value.(*txObject)
		tx.status = status
	}
}

func (t *DefaultTransactionTracker) CheckExpiredTransaction() {
	go t.checkExpiredTransaction()
}

func (t *DefaultTransactionTracker) checkExpiredTransaction() {
}
