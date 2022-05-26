package sisu_test

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForProcessorTest() (sdk.Context, sisu.ManagerContainer) {
	k, ctx := keeper.GetTestKeeperAndContext()

	globalData := &common.MockGlobalData{}
	pmm := sisu.NewPostedMessageManager(k)
	txSubmit := &common.MockTxSubmit{}

	partyManager := &sisu.MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}
	appKeys := common.NewMockAppKeys()

	mc := sisu.MockManagerContainer(k, pmm, globalData, partyManager, dheartClient, txSubmit, appKeys, ctx)
	return ctx, mc
}

func TestProcessor_OnTxIns(t *testing.T) {
	t.Parallel()

	t.Run("empty_tx", func(t *testing.T) {
		t.Parallel()

		_, mc := mockForProcessorTest()
		processor := sisu.NewProcessor(mc.Keeper(), nil, config.TssConfig{}, nil, nil, nil, nil, nil, nil, nil, mc)

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

		// Init processor with mocks
		processor := sisu.NewProcessor(k, nil, config.TssConfig{}, mc.AppKeys(), mc.TxSubmit(), nil, nil, nil, nil, nil, mc)

		err := processor.OnTxIns(txs)
		// <-done

		require.NoError(t, err)
	})
}
