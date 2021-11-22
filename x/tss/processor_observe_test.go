package tss

import (
	"sync/atomic"
	"testing"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/tests/mock"
	"github.com/sisu-network/sisu/utils"
)

func TestProcessor_OnObservedTxs(t *testing.T) {
	t.Parallel()

	t.Run("empty_txs", func(t *testing.T) {
		t.Parallel()

		processor := initProcessor()
		processor.OnObservedTxs(&eyesTypes.Txs{})
	})

	t.Run("success_new_observe_txs", func(t *testing.T) {
		t.Parallel()

		mockDb := &mock.Database{}
		mockDb.IsChainKeyAddressFunc = func(chain, address string) bool {
			return true
		}

		txs := &eyesTypes.Txs{
			Chain: utils.RandomString(10, utils.AlphaNumericCharacters),
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: nil,
				To:         utils.RandomHeximalString(64),
				From:       utils.RandomHeximalString(64),
			}},
		}
		processor := initProcessor()
		processor.db = mockDb

		processor.OnObservedTxs(txs)
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
