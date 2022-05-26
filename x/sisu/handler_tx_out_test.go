package sisu

import (
	"math/big"
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	htypes "github.com/sisu-network/dheart/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxOut() (sdk.Context, ManagerContainer) {
	txTracker := &MockTxTracker{}
	k, ctx := keeper.GetTestKeeperAndContext()
	k.SaveParams(ctx, &types.Params{
		MajorityThreshold: 1,
	})
	globalData := &common.MockGlobalData{}
	pmm := NewPostedMessageManager(k)

	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}

	dheartClient := &tssclients.MockDheartClient{}

	mc := MockManagerContainer(k, pmm, globalData, txTracker, partyManager, dheartClient)
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

	txOutWithSigner := &types.TxOutWithSigner{
		Signer: "signer",
		Data: &types.TxOut{
			OutChain: "eth",
			OutBytes: binary,
		},
	}

	t.Run("transfer_out_successful", func(t *testing.T) {
		keysignCount := 0
		trackerCount := 0
		ctx, mc := mockForHandlerTxOut()
		dheartClient := mc.DheartClient().(*tssclients.MockDheartClient)
		dheartClient.KeySignFunc = func(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error {
			keysignCount = 1
			return nil
		}
		txTracker := mc.TxTracker().(*MockTxTracker)
		txTracker.UpdateStatusFunc = func(chain, hash string, status types.TxStatus) {
			require.Equal(t, types.TxStatusDelivered, status)
			trackerCount = 1
		}

		handler := NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutWithSigner)
		require.NoError(t, err)
		require.Equal(t, 1, keysignCount)
		require.Equal(t, 1, trackerCount)
	})

	t.Run("node_is_catching_up", func(t *testing.T) {
		keysignCount := 0
		trackerCount := 0
		ctx, mc := mockForHandlerTxOut()
		dheartClient := mc.DheartClient().(*tssclients.MockDheartClient)
		dheartClient.KeySignFunc = func(req *htypes.KeysignRequest, pubKeys []ctypes.PubKey) error {
			keysignCount = 1
			return nil
		}
		txTracker := mc.TxTracker().(*MockTxTracker)
		txTracker.UpdateStatusFunc = func(chain, hash string, status types.TxStatus) {
			require.Equal(t, types.TxStatusDelivered, status)
			trackerCount = 1
		}
		globalData := mc.GlobalData().(*common.MockGlobalData)
		globalData.IsCatchingUpFunc = func() bool {
			return true
		}

		handler := NewHandlerTxOut(mc)
		_, err = handler.DeliverMsg(ctx, txOutWithSigner)
		require.NoError(t, err)
		require.Equal(t, 0, keysignCount)
		require.Equal(t, 0, trackerCount)
	})
}
