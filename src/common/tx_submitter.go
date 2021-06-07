package common

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/types"

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
	defaultGasLimit      = 300000
	UN_INITIALIZED_SEQ   = 18446744073709551615 // Max of uint64. This means it's not initialized
)

var (
	QUEUE_TIME = time.Second / 2
	ERR_NONE   = errors.New("This is not an error")
)

type QElementPair struct {
	msg   sdk.Msg
	index int64
}

type TxSubmit interface {
	SubmitEThTx(data []byte) error
	SubmitTssTx(data []byte) error
}

type TxSubmitter struct {
	sisuHome string
	appKeys  *AppKeys
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

	// Sequence
	sequenceLock *sync.RWMutex
	// Current sequence is the current sequence that will be used for transaction. It's possible that
	// multiple transactions could submitted within a block and account's sequence could be out
	// of sync with account keeper.
	curSequence uint64
	// blockSequence is the sequence for the account in the last block.
	blockSequence uint64
}

var (
	nodeAddress = "http://0.0.0.0:26657"
)

func NewTxSubmitter(cfg config.Config, appKeys *AppKeys) *TxSubmitter {
	t := &TxSubmitter{
		appKeys:         appKeys,
		cfg:             cfg,
		sequenceLock:    &sync.RWMutex{},
		queueLock:       &sync.RWMutex{},
		queue:           make([]*QElementPair, 0),
		submitRequestCh: make(chan bool),
		msgStatuses:     make(map[int64]error),
		curSequence:     UN_INITIALIZED_SEQ,
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

func (t *TxSubmitter) submitMessage(msg sdk.Msg) error {
	seq := t.getSequence()
	if seq == UN_INITIALIZED_SEQ {
		return fmt.Errorf("Server is not ready")
	}

	index := t.addMessage(msg)
	var err error

	// Delay a short period to accumulate more transactions before sending.
	t.schedule()

	for {
		time.Sleep(QUEUE_TIME)

		t.queueLock.RLock()
		err = t.msgStatuses[index]
		t.queueLock.RUnlock()
		if err != nil {
			break
		}
	}
	defer t.removeMessage(index)

	if err != ERR_NONE {
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
			t.queueLock.RUnlock()

			t.queueLock.Lock()
			t.queue = make([]*QElementPair, 0) // Clear the queue
			t.queueLock.Unlock()

			if len(copy) == 0 {
				continue
			}

			utils.LogInfo("Queue size = ", len(copy))

			// 2. Get account sequence
			seq := t.getSequence()
			utils.LogInfo("Sequence = ", seq)
			t.factory = t.factory.WithSequence(seq)

			// 3. Send all messages
			msgs := convert(copy)
			if err := tx.BroadcastTx(t.clientCtx, t.factory, msgs...); err != nil {
				utils.LogError("Cannot broadcast transaction", err)
				t.updateStatus(copy, err)

				// Use block sequence for the sequence.
				t.sequenceLock.Lock()
				t.curSequence = t.blockSequence
				t.sequenceLock.Unlock()
			} else {
				utils.LogDebug("Tx submitted successfully")
				t.updateStatus(copy, ERR_NONE)
				t.incSequence()
			}
		}
	}
}

func (t *TxSubmitter) SyncBlockSequence(ctx sdk.Context, ak authkeeper.AccountKeeper) {
	t.sequenceLock.Lock()
	defer t.sequenceLock.Unlock()

	t.blockSequence = ak.GetAccount(ctx, t.fromAccount).GetSequence()
	if t.curSequence == UN_INITIALIZED_SEQ {
		t.curSequence = t.blockSequence
	}
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

func (t *TxSubmitter) updateStatus(list []*QElementPair, err error) {
	t.queueLock.Lock()
	defer t.queueLock.Unlock()

	for _, pair := range list {
		t.msgStatuses[pair.index] = err
	}
}

func (t *TxSubmitter) SubmitEThTx(data []byte) error {
	msg := types.NewMsgEthTx(t.clientCtx.GetFromAddress().String(), data)
	return t.submitMessage(msg)
}

func (t *TxSubmitter) SubmitTssTx(data []byte) error {
	return nil
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

	chainId := t.cfg.GetSisuConfig().ChainId
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
		WithTxConfig(clientCtx.TxConfig)
}
