package sisu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"

	"github.com/sisu-network/lib/log"
)

const (
	ExpireDuration = time.Minute * 5 // 5 minutes
)

type txObject struct {
	txOut  *types.TxOut
	status types.TxStatus
	txIn   *types.TxIn

	addedTime time.Time
}

func newTxObject(txOut *types.TxOut, txIn *types.TxIn) *txObject {
	return &txObject{
		txOut:     txOut,
		status:    types.TxStatusCreated,
		txIn:      txIn,
		addedTime: time.Now(),
	}
}

//go:generate mockgen -source=./x/sisu/transaction_tracker.go -destination=./tests/mock/x/sisu/transaction_tracker.go -package=mock
// TxTracker is used to track failed transaction. This includes both TxIn and TxOut. The tracked txs
// are in-memory only.
type TxTracker interface {
	AddTransaction(txOut *types.TxOut, txIn *types.TxIn)
	UpdateStatus(chain string, hash string, status types.TxStatus)
	RemoveTransaction(chain string, hash string)
	OnTxFailed(chain string, hash string, status types.TxStatus)

	CheckExpiredTransaction()
}

type DefaultTxTracker struct {
	txs         map[string]*txObject
	txsLock     *sync.RWMutex
	emailConfig config.EmailAlertConfig
	worldState  world.WorldState
}

func NewTxTracker(emailConfig config.EmailAlertConfig, worldState world.WorldState) TxTracker {
	return &DefaultTxTracker{
		txs:         make(map[string]*txObject),
		txsLock:     &sync.RWMutex{},
		emailConfig: emailConfig,
		worldState:  worldState,
	}
}

func (t *DefaultTxTracker) getTxoKey(chain string, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}

func (t *DefaultTxTracker) AddTransaction(txOut *types.TxOut, txIn *types.TxIn) {
	chain := txOut.OutChain
	hash := txOut.OutHash
	key := t.getTxoKey(chain, hash)

	t.txsLock.Lock()
	defer t.txsLock.Unlock()

	if t.txs[key] != nil {
		log.Warnf("Tx has been added to the tracker, chain and hash = ", chain, hash)
	} else {
		log.Verbosef("Adding tx to tracking %s %s", chain, hash)

		t.txs[key] = newTxObject(txOut, txIn)
	}
}

func (t *DefaultTxTracker) UpdateStatus(chain string, hash string, status types.TxStatus) {
	t.txsLock.Lock()
	defer t.txsLock.Unlock()

	txo := t.txs[t.getTxoKey(chain, hash)]
	if txo == nil {
		return
	}

	txo.status = status
}

func (t *DefaultTxTracker) CheckExpiredTransaction() {
	go t.checkExpiredTransaction()
}

func (t *DefaultTxTracker) RemoveTransaction(chain string, hash string) {
	key := t.getTxoKey(chain, hash)

	t.txsLock.Lock()
	delete(t.txs, key)
	t.txsLock.Unlock()
}

func (t *DefaultTxTracker) OnTxFailed(chain string, hash string, status types.TxStatus) {
	key := t.getTxoKey(chain, hash)
	t.txsLock.RLock()
	txo := t.txs[key]
	t.txsLock.RUnlock()

	if txo == nil {
		log.Warnf("cannot find transaction in tracker with chain %s and hash %s", chain, hash)
		return
	}

	go t.processFailedTx(txo)
}

func (t *DefaultTxTracker) checkExpiredTransaction() {
	toRemove := make(map[string]*txObject)
	now := time.Now()

	t.txsLock.RLock()
	for key, txo := range t.txs {
		expire := txo.addedTime.Add(ExpireDuration)

		if expire.Before(now) {
			// This transcation has expired.
			toRemove[key] = txo
		}
	}
	t.txsLock.RUnlock()

	if len(toRemove) > 0 {
		log.Warnf("There are some transactions that are not observed on remote blockchain, size = ", len(toRemove))
	}

	// Broadcast the failure
	t.txsLock.Lock()
	for key, txo := range toRemove {
		delete(t.txs, key)
		go t.processFailedTx(txo)
	}
	t.txsLock.Unlock()
}

func (t *DefaultTxTracker) processFailedTx(txo *txObject) {
	inChain := ""
	inHash := ""
	if txo.txIn != nil {
		inChain = txo.txIn.Chain
		inHash = txo.txIn.TxHash
	}
	log.Warnf("Processing failed transaction. outChain: %s, outHash :%s, status: %v, inChain: %s, inHash: %s",
		txo.txOut.OutChain, txo.txOut.OutHash, txo.status, inChain, inHash)

	key := t.getTxoKey(txo.txOut.OutChain, txo.txOut.OutHash)
	t.txsLock.Lock()
	delete(t.txs, key)
	t.txsLock.Unlock()

	// Send email if needed
	if len(t.emailConfig.Url) > 0 && len(t.emailConfig.Secret) > 0 {
		t.sendAlertEmail(t.emailConfig.Url, t.emailConfig.Secret, t.emailConfig.Email, txo)
	}
}

// Sends alert email using SendGrid. This could be moved into an interface to use any mail service.
func (t *DefaultTxTracker) sendAlertEmail(url, secret, email string, txo *txObject) {
	type Email struct {
		Email string `json:"email"`
	}

	type ContentElement struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type To struct {
		To []Email `json:"to"`
	}

	type SendGrid struct {
		Personalizations []To             `json:"personalizations"`
		From             Email            `json:"from"`
		Subject          string           `json:"subject"`
		Content          []ContentElement `json:"content"`
	}

	body, err := t.getEmailBodyString(txo)
	if err != nil {
		log.Error("sendAlertEmail: cannot pretty print body, err = ", err)
		return
	}

	value := SendGrid{
		Personalizations: []To{
			{
				To: []Email{
					{
						Email: email,
					},
				},
			},
		},
		From: Email{
			Email: email,
		},
		Subject: "Not observed transaction",
		Content: []ContentElement{
			{
				Type:  "text/plain",
				Value: body,
			},
		},
	}

	json_data, err := json.Marshal(value)
	if err != nil {
		log.Error("sendAlertEmail: failed to marshal value, err = ", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	if err != nil {
		log.Error("sendAlertEmail: failed to create new http post request, err = ", err)
		return
	}

	var client http.Client
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secret))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Error("sendAlertEmail: failed to send request, err = ", err)
		return
	}
	defer resp.Body.Close() // ignore response
}

func (t *DefaultTxTracker) getEmailBodyString(txo *txObject) (string, error) {
	type TxInData struct {
		Chain        string `json:"chain"`
		TokenAddress string `json:"token_address"`
		Recipient    string `json:"recipient"`
		Amount       string `json:"amount"`
	}

	type TxOutData struct {
		Type            string `json:"type"`
		Chain           string `json:"chain"`
		ContractAddress string `json:"contract_address"`
		TokenAddress    string `json:"token_address"`
		Recipient       string `json:"recipient"`
		Amount          string `json:"amount"`
	}

	type Body struct {
		TxType     string    `json:"type"`
		Chain      string    `json:"chain"`
		Hash       string    `json:"hash"`
		LastStatus string    `json:"last_status"`
		TxOutData  TxOutData `json:"tx_out_data"`
		TxInData   TxInData  `json:"tx_in_data"`
	}

	body := Body{}
	body.Chain = txo.txOut.OutChain
	body.Hash = txo.txOut.OutHash
	body.LastStatus = types.StatusStrings[txo.status]

	txIn := txo.txIn
	if txIn != nil {
		if txo.txOut.TxType == types.TxOutType_TRANSFER_OUT {
			// Trying its best to Deserialize the txIn
			data, err := t.getEThTransferIn(txIn.Serialized)
			if err == nil {
				txInData := TxInData{
					Chain:        txIn.Chain,
					TokenAddress: data.tokenAddr.String(),
					Recipient:    data.recipient.String(),
					Amount:       data.amount.String(),
				}
				body.TxInData = txInData
			} else {
				log.Error("cannot parse eth transfer in, err = ", err)
			}
		}
	}

	switch txo.txOut.TxType {
	case types.TxOutType_TRANSFER_OUT:
		data, err := t.getEthTransferIn(txo.txOut.OutBytes)
		if err != nil {
			log.Error("Cannot parse transfer in data, err = ", err)
			return "", err
		}

		body.TxOutData = TxOutData{
			Type:         "TRANSFER_OUT",
			Chain:        txo.txOut.OutChain,
			TokenAddress: data.token.String(),
			Recipient:    data.recipient.String(),
			Amount:       data.amount.String(),
		}
	case types.TxOutType_CONTRACT_DEPLOYMENT:
		return fmt.Sprintf("contract deployment failed, hash = %s", txo.txOut.OutHash), nil
	}

	return utils.PrettyStruct(body)
}

func (t *DefaultTxTracker) getEThTransferIn(bz []byte) (*transferOutData, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(bz)
	if err != nil {
		return nil, err
	}

	return parseEthTransferOut(ethTx, t.worldState)
}

func (t *DefaultTxTracker) getEthTransferIn(bz []byte) (*transferInData, error) {
	ethTx := &ethTypes.Transaction{}

	err := ethTx.UnmarshalBinary(bz)
	if err != nil {
		return nil, err
	}

	return parseTransferInData(ethTx, t.worldState)
}
