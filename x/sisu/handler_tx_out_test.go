package sisu

import (
	"math/big"
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOut() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	txTracker := &MockTxTracker{}
	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}

	mc := MockManagerContainer(k, pmm, globalData, txTracker, partyManager, dheartClient, &MockTxOutQueue{})
	return ctx, mc
}

func TestHandlerTxOut_TransferOut(t *testing.T) {
	t.Parallel()

	amount := big.NewInt(100)
	gasLimit := uint64(100)
	gasPrice := big.NewInt(100)
	ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &ecommon.Address{},
		Value:    amount,
	})
	binary, err := ethTransaction.MarshalBinary()
	require.NoError(t, err)

	txOutWithSigner := &types.TxOutMsg{
		Signer: "signer",
		Data: &types.TxOut{
			OutChain: "eth",
			OutBytes: binary,
		},
	}

	t.Run("transfer_out_successful", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		addTxCount := 0
		txOutQueue := mc.TxOutQueue()
		txOutQueue.(*MockTxOutQueue).AddTxOutFunc = func(txOut *types.TxOut) {
			addTxCount = 1
		}

		handler := NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutWithSigner)
		require.NoError(t, err)
		require.Equal(t, 1, addTxCount)
	})

	t.Run("node_is_catching_up", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		addTxCount := 0
		txOutQueue := mc.TxOutQueue()
		txOutQueue.(*MockTxOutQueue).AddTxOutFunc = func(txOut *types.TxOut) {
			addTxCount = 1
		}

		globalData := mc.GlobalData().(*common.MockGlobalData)
		globalData.IsCatchingUpFunc = func() bool {
			return true
		}

		handler := NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutWithSigner)
		require.NoError(t, err)
		require.Equal(t, 0, addTxCount)
	})
}
