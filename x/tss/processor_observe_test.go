package tss

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
	"testing"

	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	"github.com/sisu-network/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	eyesTypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/tests/mock"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestProcessor_OnObservedTxs(t *testing.T) {
	t.Parallel()

	t.Run("empty_tx", func(t *testing.T) {
		t.Parallel()

		processor := &Processor{}
		require.NoError(t, processor.OnObservedTxs(&eyesTypes.Txs{}))
	})

	t.Run("success_from_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		observedChain := "eth"
		keygenAddress := utils.RandomHeximalString(64)

		// Define mock db
		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(true).MinTimes(1)
		mockDb.EXPECT().UpdateTxOutStatus(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)
		mockDb.EXPECT().GetTxOutWithHash(gomock.Any(), gomock.Any(), gomock.Any()).Return(
			&types.TxOutEntity{
				ContractHash: utils.RandomHeximalString(64),
			}).MinTimes(1)
		mockDb.EXPECT().UpdateContractsStatus(gomock.Any(), gomock.Any()).Return(nil).MinTimes(1)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				From:       keygenAddress,
			}},
		}
		processor := &Processor{}
		processor.db = mockDb

		require.NoError(t, processor.OnObservedTxs(txs))
	})

	t.Run("success_to_our_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup( func() {
			ctrl.Finish()
		})

		txSubmitterMock := mock.NewMockTxSubmit(ctrl)
		txSubmitterMock.EXPECT().SubmitMessage(gomock.Any()).Return(nil).AnyTimes()

		priv := ed25519.GenPrivKey()
		addr := sdk.AccAddress(priv.PubKey().Address())
		appKeysMock := mock.NewMockAppKeys(ctrl)
		appKeysMock.EXPECT().GetSignerAddress().Return(addr).MinTimes(1)

		observedChain := "eth"
		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(false).MinTimes(1)
		mockDb.EXPECT().UpdateTxOutStatus(observedChain, gomock.Any(), types.TxOutStatusPreBroadcast, gomock.Any()).Return(nil).AnyTimes()
		keygenAddress := utils.RandomHeximalString(64)

		txs := &eyesTypes.Txs{
			Chain: observedChain,
			Block: int64(utils.RandomNaturalNumber(1000)),
			Arr: []*eyesTypes.Tx{{
				Hash:       utils.RandomHeximalString(64),
				Serialized: []byte{},
				To:         keygenAddress,
				From:       utils.RandomHeximalString(64),
			}},
		}

		// Init processor with mocks
		processor := &Processor{
			db:       mockDb,
			appKeys:  appKeysMock,
			txSubmit: txSubmitterMock,
		}

		require.NoError(t, processor.OnObservedTxs(txs))
	})
}

func TestProcessor_createAndBroadcastTxOuts(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		processor := &Processor{}
		signer    := ethTypes.NewEIP2930Signer(common.Big1)
		recipient := common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
		addr      := common.HexToAddress("0x0000000000000000000000000000000000000001")
		accesses  := ethTypes.AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}}

		txdata := &ethTypes.AccessListTx{
			ChainID:    big.NewInt(1),
			Nonce:      10,
			To:         &recipient,
			Gas:        123457,
			GasPrice:   big.NewInt(10),
			AccessList: accesses,
			Data:       []byte("abcdef"),
		}

		key, err := crypto.GenerateKey()
		require.NoError(t, err)

		tx, err := ethTypes.SignNewTx(key, signer, txdata)
		require.NoError(t, err)

		outBytes, err := tx.MarshalBinary()
		require.NoError(t, err)

		ctx := sdk.Context{}
		observedTx := types.ObservedTx{
			Chain:       "ganache1",
			BlockHeight: 10,
			Serialized:  outBytes,
		}
		txOuts := processor.createAndBroadcastTxOuts(ctx, &observedTx)
		require.NotNil(t, txOuts)
	})

}