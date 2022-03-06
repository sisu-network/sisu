package sisu

import (
	"fmt"
	"strings"
	"sync"
	"time"

	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/config"

	"github.com/sisu-network/lib/log"
)

type TxTrackerType int64

const (
	TxTrackerTxIn TxTrackerType = iota
	TxTrackerTxOut
)

type TxStatus int64

const (
	TxStatusUnknown TxStatus = iota
	TxStatusCreated
	TxStatusDelivered
	TxStatusSigned
	TxStatusSignFailed
	TxStatusDepoyed // transaction has been sent to blockchain but not confirmed yet.
	TxStatusConfirmed
)

const (
	ExpireDuration = time.Minute * 5 // 5 minutes
)

type txObject struct {
	txType  TxTrackerType
	chain   string
	hash    string // hash without signature
	status  TxStatus
	content []byte

	addedTime time.Time
}

func newTxObject(txType TxTrackerType, chain string, hash string, content []byte) *txObject {
	return &txObject{
		txType:    txType,
		chain:     chain,
		hash:      hash,
		content:   content,
		status:    TxStatusCreated,
		addedTime: time.Now(),
	}
}

// TxTracker is used to track failed transaction. This includes both TxIn and TxOut. The tracked txs
// are in-memory only.
type TxTracker interface {
	AddTransaction(txType TxTrackerType, chain string, hash string, content []byte, extra interface{})
	UpdateStatus(txType TxTrackerType, chain string, hash string, status TxStatus)
	RemoveTransaction(txType TxTrackerType, chain string, hash string)

	OnTxFailed(txType TxTrackerType, chain string, hash string, status TxStatus)
}

type DefaultTxTracker struct {
	txs         *sync.Map
	emailConfig config.EmailAlertConfig
}

func NewTxTracker(emailConfig config.EmailAlertConfig) TxTracker {
	return &DefaultTxTracker{
		txs:         &sync.Map{},
		emailConfig: emailConfig,
	}
}

func (t *DefaultTxTracker) getTxoKey(chain string, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}

func (t *DefaultTxTracker) AddTransaction(txType TxTrackerType, chain string, hash string, content []byte, extra interface{}) {
	key := t.getTxoKey(chain, hash)

	if _, ok := t.txs.Load(key); ok {
		log.Warnf("Tx has been added to the tracker, chain and hash = ", chain, hash)
		return
	}

	log.Verbosef("Adding tx to tracking %s %s", chain, hash)

	t.txs.Store(key, newTxObject(txType, chain, hash, content))
}

func (t *DefaultTxTracker) UpdateStatus(txType TxTrackerType, chain string, hash string, status TxStatus) {
	value, ok := t.txs.Load(t.getTxoKey(chain, hash))
	if ok {
		tx := value.(*txObject)
		tx.status = status
	}
}

func (t *DefaultTxTracker) CheckExpiredTransaction() {
	go t.checkExpiredTransaction()
}

func (t *DefaultTxTracker) RemoveTransaction(txType TxTrackerType, chain string, hash string) {
	key := t.getTxoKey(chain, hash)
	t.txs.Delete(key)

	// Print size
	count := 0
	t.txs.Range(func(key, value interface{}) bool {
		count += 1
		return true
	})

	log.Verbosef("TxTracker: Removing tx from tracking %s %s", chain, hash)
	log.Verbosef("TxTracker: Remaining count = %d", count)
}

func (t *DefaultTxTracker) OnTxFailed(txType TxTrackerType, chain string, hash string, status TxStatus) {
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
	log.Warnf("Processing failed transaction: %s %s %v", txo.chain, txo.hash, txo.status)

	key := t.getTxoKey(txo.chain, txo.hash)
	t.txs.Delete(key)
}
