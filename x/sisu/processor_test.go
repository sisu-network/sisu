package sisu

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForProcessorTest() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)

	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)
	txSubmit := &common.MockTxSubmit{}
	txTracker := &MockTxTracker{}

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}
	appKeys := common.NewMockAppKeys()

	mc := MockManagerContainer(k, pmm, globalData, partyManager, dheartClient, txSubmit, appKeys, ctx, txTracker)
	return ctx, mc
}

func TestProcessor_OnTxIns(t *testing.T) {
	t.Parallel()

	t.Run("empty_tx", func(t *testing.T) {
		t.Parallel()

		_, mc := mockForProcessorTest()
		processor := NewProcessor(nil, mc)

		require.NoError(t, processor.OnTxIns(&eyesTypes.Txs{}))
	})

	t.Run("success_to_our_key", func(t *testing.T) {
		t.Parallel()

		ctx, mc := mockForProcessorTest()

		k := mc.Keeper()
		k.SaveKeygen(ctx, &types.Keygen{})

		k.IsKeygenAddress(ctx, "ecdsa", "123")

		observedChain := "eth"
		toAddress := utils.RandomHeximalString(64)
		fromAddres := utils.RandomHeximalString(64)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				To:         toAddress,
				From:       fromAddres,
				Success:    true,
			}},
		}

		submitCount := 0
		txSubmit := mc.TxSubmit().(*common.MockTxSubmit)
		txSubmit.SubmitMessageAsyncFunc = func(msg sdk.Msg) error {
			submitCount = 1
			return nil
		}

		processor := NewProcessor(nil, mc)
		err := processor.OnTxIns(txs)

		require.NoError(t, err)
		require.Equal(t, 1, submitCount)
	})

	t.Run("failed_transaction", func(t *testing.T) {
		txs := &eyesTypes.Txs{
			Arr: []*eyesTypes.Tx{{
				Success: false,
			}},
		}

		trackerCount := 0
		_, mc := mockForProcessorTest()
		txTracker := mc.TxTracker().(*MockTxTracker)
		txTracker.OnTxFailedFunc = func(chain, hash string, status types.TxStatus) {
			trackerCount = 1
		}

		processor := NewProcessor(nil, mc)
		err := processor.OnTxIns(txs)
		require.NoError(t, err)
		require.Equal(t, 1, trackerCount)
	})
}
