package tss

import (
	"sync/atomic"
	"testing"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/tests/mock"
)

func TestProcessor_OnObservedTxs(t *testing.T) {
	t.Parallel()

	processor := initProcessor()
	t.Run("empty_txs", func(t *testing.T) {
		t.Parallel()

		processor.OnObservedTxs(&eyesTypes.Txs{})
	})

	t.Run("success_new_observe_txs", func(t *testing.T) {
		t.Parallel()
	})
}

func initProcessor() *Processor {
	return &Processor{
		keeper:                 nil,
		config:                 config.TssConfig{},
		tendermintPrivKey:      nil,
		txSubmit:               &mock.TxSubmitter{},
		lastProposeBlockHeight: 0,
		appKeys:                nil,
		globalData:             nil,
		currentHeight:          0,
		partyManager:           nil,
		txOutputProducer:       nil,
		lastContext:            atomic.Value{},
		checkDuplicatedTxFunc:  nil,
		txDecoder:              nil,
		keyAddress:             "",
		dheartClient:           nil,
		deyesClients:           nil,
		worldState:             nil,
		keygenVoteResult:       nil,
		keygenBlockPairs:       nil,
		db:                     &mock.Database{},
	}
}
