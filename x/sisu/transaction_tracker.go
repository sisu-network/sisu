package sisu

import (
	"fmt"
	"strings"
	"sync"
	"time"

	libchain "github.com/sisu-network/lib/chain"

	"github.com/sisu-network/lib/log"
)

type TxStatus int64

const (
	TxStatusUnknown TxStatus = iota
	TxStatusCreated
	TxStatusSigned
	TxStatusSignFailed
	TxStatusConfirmed
)

const (
	ExpireDuration = time.Minute * 5 // 5 minutes
)

type txObject struct {
	status TxStatus
	chain  string
	hash   string // hash without signature

	addedTime time.Time
}

func newTxObject(chain string, hash string) *txObject {
	return &txObject{
		status:    TxStatusCreated,
		chain:     chain,
		hash:      hash,
		addedTime: time.Now(),
	}
}

type TxTracker interface {
	AddTransaction(chain string, hash string)
	UpdateStatus(chain string, hash string, status TxStatus)
	RemoveTransaction(chain string, hash string)

	OnTxFailed(chain string, hash string, status TxStatus)
}

type DefaultTxTracker struct {
	txs *sync.Map
}

func NewTxTracker() TxTracker {
	return &DefaultTxTracker{
		txs: &sync.Map{},
	}
}

func (t *DefaultTxTracker) getTxoKey(chain string, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}

func (t *DefaultTxTracker) AddTransaction(chain string, hash string) {
	key := t.getTxoKey(chain, hash)

	if _, ok := t.txs.Load(key); ok {
		log.Warnf("Tx has been added to the tracker, chain and hash = ", chain, hash)
		return
	}

	t.txs.Store(key, newTxObject(chain, hash))
}

func (t *DefaultTxTracker) UpdateStatus(chain string, hash string, status TxStatus) {
	value, ok := t.txs.Load(t.getTxoKey(chain, hash))
	if ok {
		tx := value.(*txObject)
		tx.status = status
	}
}

func (t *DefaultTxTracker) CheckExpiredTransaction() {
	go t.checkExpiredTransaction()
}

func (t *DefaultTxTracker) RemoveTransaction(chain string, hash string) {
	key := t.getTxoKey(chain, hash)
	t.txs.Delete(key)
}

func (t *DefaultTxTracker) OnTxFailed(chain string, hash string, status TxStatus) {
	key := t.getTxoKey(chain, hash)
	if val, ok := t.txs.Load(key); ok {
		t.processFailedTx(val.(*txObject))
	} else {
		log.Warnf("cannot find transaction in tracker with chain %s and hash %s", chain, hash)
	}
}

func (t *DefaultTxTracker) checkExpiredTransaction() {
	toRemove := make(map[string]*txObject)

	now := time.Now()

	t.txs.Range(func(key, value interface{}) bool {
		txo := value.(*txObject)
		expire := txo.addedTime.Add(ExpireDuration)

		if expire.Before(now) {
			// This transcation has expired.
			toRemove[key.(string)] = txo
		}

		return true
	})

	// Broadcast the failure
	for key, txo := range toRemove {
		t.txs.Delete(key)

		index := strings.Index(key, "__")
		if index < 0 {
			continue
		}

		chain := key[:index]
		if libchain.IsETHBasedChain(chain) {
			t.processFailedTx(txo)
		}
	}
}

func (t *DefaultTxTracker) processFailedTx(txo *txObject) {
	key := t.getTxoKey(txo.chain, txo.hash)

	switch txo.status {
	case TxStatusCreated, TxStatusSignFailed:
	case TxStatusSigned:
	}

	t.txs.Delete(key)
}
