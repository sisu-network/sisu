package components

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	// TODO: put these values into config file.
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 100_000_000
	UnInitializedSeq     = 18446744073709551615 // Max of uint64. This means it's not initialized
)

var (
	QueueTime = time.Second / 2
	ErrNone   = errors.New("This is not an error")
)

type TxSubmit interface {
	SubmitMessageAsync(msg sdk.Msg) error
	SubmitMessageSync(msg sdk.Msg) error
}

type TxSubmitter struct {
	sisuHome string
	appKeys  *DefaultAppKeys
	cfg      config.Config

	// internal
	clientCtx   client.Context
	factory     tx.Factory
	fromAccount sdk.AccAddress
	httpClient  *retryablehttp.Client

	// Sequence
	sequenceLock *sync.RWMutex

	// Current sequence is the current sequence that will be used for transaction. It's possible that
	// multiple transactions could submitted within a block and account's sequence could be out
	// of sync with account keeper.
	curSequence   uint64
	accountNumber uint64

	// Tx queue
	queue           []sdk.Msg
	queueLock       *sync.RWMutex
	submitRequestCh chan bool
}

var (
	nodeAddress = "http://0.0.0.0:26657"
)

func NewTxSubmitter(cfg config.Config, appKeys *DefaultAppKeys) *TxSubmitter {
	httpClient := retryablehttp.NewClient()
	httpClient.Logger = nil

	t := &TxSubmitter{
		appKeys:         appKeys,
		cfg:             cfg,
		httpClient:      httpClient,
		queueLock:       &sync.RWMutex{},
		queue:           make([]sdk.Msg, 0),
		submitRequestCh: make(chan bool, 20),
		sequenceLock:    &sync.RWMutex{},
		curSequence:     UnInitializedSeq,
	}

	var err error
	t.fromAccount = appKeys.GetSignerInfo().GetAddress()
	t.clientCtx, err = t.buildClientCtx(appKeys.GetSignerInfo().GetName())
	t.factory = newFactory(t.clientCtx)

	if err != nil {
		panic(err)
	}

	return t
}

func (t *TxSubmitter) SubmitMessageAsync(msg sdk.Msg) error {
	go t.submitMessage(msg)
	return nil
}

func (t *TxSubmitter) SubmitMessageSync(msg sdk.Msg) error {
	return t.submitMessage(msg)
}

func (t *TxSubmitter) submitMessage(msg sdk.Msg) error {
	log.Verbose("Submitting tx ....")
	t.addMessage(msg)

	// Delay a short period to accumulate more transactions before sending.
	t.schedule()

	return nil
}

func (t *TxSubmitter) addMessage(msg sdk.Msg) {
	t.addMessages([]sdk.Msg{msg})
}

func (t *TxSubmitter) addMessages(msgs []sdk.Msg) {
	t.queueLock.Lock()
	defer t.queueLock.Unlock()

	for _, msg := range msgs {
		t.queue = append(t.queue, msg)
	}
}

func (t *TxSubmitter) schedule() {
	t.submitRequestCh <- true
}

func (t *TxSubmitter) Start() {
	for {
		select {
		case <-t.submitRequestCh:
			// 1. Gets all pending messages in the queue.
			// Use read lock since it's cheaper
			t.queueLock.Lock()
			if len(t.queue) == 0 {
				t.queueLock.Unlock()
				continue
			}
			copy := t.queue
			t.queue = make([]sdk.Msg, 0) // Clear the queue
			t.queueLock.Unlock()

			if len(copy) == 0 {
				continue
			}

			log.Info("Queue size = ", len(copy))

			// 2. Get account sequence
			seq := t.getSequence()
			log.Info("Sequence = ", seq)
			t.factory = t.factory.WithSequence(seq)

			// 3. Send all messages
			res, err := t.trySubmitTx(copy)
			if err != nil || (res != nil && res.Code != 0) {
				code := -1
				if res != nil {
					code = int(res.Code)
				} else {
					log.Error("responsed from sisu is nil")
				}
				log.Errorf("Cannot broadcast transaction, code = %d and err = %v", code, err)

				// Do retry if the error we get is incorrect sequence number
				// List of error code here: https://github.com/cosmos/cosmos-sdk/blob/v0.42.1/types/errors/errors.go
				if res != nil && res.Code == 32 { // incorrect sequence number
					accNumber, newSequence, err := t.getLatestSequence()
					if err == nil {
						log.Info("New sequence = ", newSequence)
						// Update account number and the sequence.
						t.updateAccNumberAndSquence(accNumber, newSequence)

						// Retry the second time
						res2, err2 := t.trySubmitTx(copy)
						if err2 != nil || (res2 != nil && res2.Code != 0) {
							log.Errorf("Retry failed, code = %d and err = %v", res2.Code, err2)
						} else {
							log.Verbose("Retry succeeded")
						}
					} else {
						log.Error("cannot get sequence, err = ", err)
					}
				} else {
					log.Error("We cannot sequence number. We will readded all transactions in the queue again.")
					t.addMessages(copy)
				}
			}
		}
	}
}

func (t *TxSubmitter) trySubmitTx(msgs []sdk.Msg) (*sdk.TxResponse, error) {
	if res, err := t.submitMsgs(msgs); err != nil || (res != nil && res.Code != 0) {
		return res, err
	} else {
		t.incSequence()

		return res, err
	}
}

func (t *TxSubmitter) submitMsgs(msgs []sdk.Msg) (*sdk.TxResponse, error) {
	builder, err := tx.BuildUnsignedTx(t.factory, msgs...)
	if err != nil {
		return nil, err
	}

	err = tx.Sign(t.factory, t.clientCtx.GetFromName(), builder, true)
	if err != nil {
		return nil, err
	}

	txBytes, err := t.clientCtx.TxConfig.TxEncoder()(builder.GetTx())
	if err != nil {
		return nil, err
	}

	res, err := t.clientCtx.BroadcastTx(txBytes)
	return res, err
}

// getLatestSequence makes a request to tendermint and get the correct sequence for the current account.
func (t *TxSubmitter) getLatestSequence() (uint64, uint64, error) {
	url := fmt.Sprintf("http://127.0.0.1:1317/auth/accounts/%s", t.fromAccount)

	type AccountResp struct {
		Height string `json:"height"`
		Result struct {
			Value struct {
				AccountNumber uint64 `json:"account_number,string"`
				Sequence      uint64 `json:"sequence,string"`
			} `json:"value"`
		} `json:"result"`
	}
	resp := &AccountResp{}
	body, _, err := utils.HttpGet(t.httpClient, url)
	if err != nil {
		return 0, 0, err
	}

	if err := json.Unmarshal(body, &resp); err != nil {
		return 0, 0, err
	}

	return resp.Result.Value.AccountNumber, resp.Result.Value.Sequence, nil
}

func (t *TxSubmitter) SyncBlockSequence(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	t.sequenceLock.RLock()
	seq := t.curSequence
	t.sequenceLock.RUnlock()

	if seq != UnInitializedSeq {
		return
	}

	if t.fromAccount == nil {
		log.Error("fromAccount is not set yet")
		return
	}

	// We create a new context with a new gas meter since ak.GetAccount consumes different gas amount
	// for different length of the t.fromAccount.
	copyCtx := ctx.WithGasMeter(sdk.NewInfiniteGasMeter())
	account := ak.GetAccount(copyCtx, t.fromAccount)

	if account == nil {
		log.Error("cannot find account in the keeper, account =", t.fromAccount)
		return
	}

	t.updateAccNumberAndSquence(account.GetAccountNumber(), account.GetSequence())
}

func (t *TxSubmitter) getSequence() uint64 {
	t.sequenceLock.RLock()
	defer t.sequenceLock.RUnlock()

	return t.curSequence
}

func (t *TxSubmitter) incSequence() {
	t.sequenceLock.Lock()
	defer t.sequenceLock.Unlock()

	t.curSequence++
}

func (t *TxSubmitter) updateAccNumberAndSquence(newAccountNumber, newSeq uint64) {
	t.sequenceLock.Lock()
	defer t.sequenceLock.Unlock()

	t.curSequence = newSeq
	t.accountNumber = newAccountNumber
	t.factory = t.factory.WithAccountNumber(newAccountNumber)
	t.factory = t.factory.WithSequence(newSeq)
}

func (t *TxSubmitter) buildClientCtx(accountName string) (client.Context, error) {
	info := t.appKeys.GetSignerInfo()

	client, err := rpchttp.New(nodeAddress, "/websocket")
	if err != nil {
		panic(err)
	}

	chainId := t.cfg.Sisu.ChainId
	clientCtx := NewClientCtx(t.appKeys.GetKeyring(), client, &bytes.Buffer{}, t.sisuHome, chainId)

	return clientCtx.
		WithFromName(accountName).
		WithFromAddress(info.GetAddress()), nil
}

func NewClientCtx(kr keyring.Keyring, c *rpchttp.HTTP, out io.Writer, home, chainId string) client.Context {
	encodingConfig := params.MakeEncodingConfig()
	authtypes.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	sdk.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	staking.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cryptocodec.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	return client.Context{}.
		WithChainID(chainId).
		WithKeyring(kr).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithJSONMarshaler(encodingConfig.Marshaler).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithOutput(out).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(home).
		WithClient(c).
		WithSkipConfirmation(true)
}

func newFactory(clientCtx client.Context) tx.Factory {
	return tx.Factory{}.
		WithChainID(clientCtx.ChainID).
		WithKeybase(clientCtx.Keyring).
		WithGas(defaultGasLimit).
		WithGasAdjustment(defaultGasAdjustment).
		WithSignMode(signing.SignMode_SIGN_MODE_UNSPECIFIED).
		WithAccountRetriever(clientCtx.AccountRetriever).
		WithTxConfig(clientCtx.TxConfig).
		WithSequence(0)
}
