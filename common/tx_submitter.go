package common

import (
	"bytes"
	"errors"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	staking "github.com/cosmos/cosmos-sdk/x/staking/types"
	rpchttp "github.com/sisu-network/tendermint/rpc/client/http"
)

//go:generate mockgen -source=tx_submitter.go -destination=../tests/mock/tx_submitter.go -package=mock

const (
	// TODO: put these values into config file.
	defaultGasAdjustment = 1.0
	defaultGasLimit      = 3_000_000
)

var (
	QueueTime = time.Second / 2
	ErrNone   = errors.New("This is not an error")
)

type QElementPair struct {
	msg   sdk.Msg
	index int64
}

// mockgen -source common/tx_submitter.go -destination=tests/mock/common/tx_submitter.go -package=mock
type TxSubmit interface {
	SubmitMessage(msg sdk.Msg) error
}

type TxSubmitter struct {
	sisuHome string
	appKeys  *DefaultAppKeys
	cfg      config.Config

	// internal
	clientCtx   client.Context
	factory     tx.Factory
	fromAccount sdk.AccAddress

	// Tx queue
	queue           []*QElementPair
	queueLock       *sync.RWMutex
	msgIndex        int64
	msgStatuses     map[int64]error
	submitRequestCh chan bool
}

var (
	nodeAddress = "http://0.0.0.0:26657"
)

func NewTxSubmitter(cfg config.Config, appKeys *DefaultAppKeys) *TxSubmitter {
	t := &TxSubmitter{
		appKeys:         appKeys,
		cfg:             cfg,
		queueLock:       &sync.RWMutex{},
		queue:           make([]*QElementPair, 0),
		submitRequestCh: make(chan bool),
		msgStatuses:     make(map[int64]error),
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

func (t *TxSubmitter) SubmitMessage(msg sdk.Msg) error {
	log.Debug("Submitting tx ....")

	index := t.addMessage(msg)
	var err error

	// Delay a short period to accumulate more transactions before sending.
	t.schedule()

	for {
		time.Sleep(QueueTime)

		t.queueLock.RLock()
		err = t.msgStatuses[index]
		t.queueLock.RUnlock()
		if err != nil {
			break
		}
	}
	defer t.removeMessage(index)

	if err != ErrNone {
		return err
	}

	return nil
}

func (t *TxSubmitter) addMessage(msg sdk.Msg) int64 {
	t.queueLock.Lock()
	defer t.queueLock.Unlock()

	t.msgIndex++
	t.queue = append(t.queue, &QElementPair{
		msg:   msg,
		index: t.msgIndex,
	})

	return t.msgIndex
}

func (t *TxSubmitter) removeMessage(msgIndex int64) {
	t.queueLock.Lock()
	defer t.queueLock.Unlock()

	delete(t.msgStatuses, msgIndex)
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
			t.queueLock.RLock()
			if len(t.queue) == 0 {
				t.queueLock.RUnlock()
				continue
			}
			copy := t.queue
			t.queue = make([]*QElementPair, 0) // Clear the queue
			t.queueLock.RUnlock()

			if len(copy) == 0 {
				continue
			}

			log.Info("Queue size = ", len(copy))

			// 3. Send all messages
			msgs := convert(copy)
			if err := tx.BroadcastTx(t.clientCtx, t.factory, msgs...); err != nil {
				log.Error("Cannot broadcast transaction, err = ", err)
				t.updateStatus(copy, err)
			} else {
				log.Debug("Tx submitted successfully")
				t.updateStatus(copy, ErrNone)
			}
		}
	}
}

func (t *TxSubmitter) updateStatus(list []*QElementPair, err error) {
	t.queueLock.Lock()
	defer t.queueLock.Unlock()

	for _, pair := range list {
		t.msgStatuses[pair.index] = err
	}
}

func convert(list []*QElementPair) []sdk.Msg {
	msgs := make([]sdk.Msg, len(list))
	for i, pair := range list {
		msgs[i] = pair.msg
	}
	return msgs
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
