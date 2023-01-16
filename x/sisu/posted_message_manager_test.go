package sisu

import (
	"math/big"
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
	db "github.com/tendermint/tm-db"
)

func mockForPostedMessageManager() (sdk.Context, ManagerContainer) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)
	globalData := &common.MockGlobalData{}
	dheartClient := &external.MockDheartClient{}
	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}
	valsMgr := NewValidatorManager(k)
	valsMgr.AddValidator(ctx, &types.Node{
		ValPubkey: &types.ValPubkey{
			Type:  "ed25519",
			Bytes: []byte("some_key"),
		},
	})
	txOutProducer := &MockTxOutputProducer{}
	mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient, valsMgr,
		&MockTransferQueue{}, &MockTxOutQueue{}, keeper.NewPrivateDb(".", db.MemDBBackend))

	return ctx, mc
}

func mockTxOutWithSignerForPostedMessageManager() *types.TxOutMsg {
	ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: big.NewInt(100),
		Gas:      uint64(100),
		To:       &ecommon.Address{},
		Value:    big.NewInt(100),
	})
	binary, _ := ethTx.MarshalBinary()

	txOutWithSigner := &types.TxOutMsg{
		Signer: "signer",
		Data: &types.TxOut{
			Content: &types.TxOutContent{
				OutChain: "ganache1",
				OutBytes: binary,
			},
		},
	}

	return txOutWithSigner
}

func TestPostedMessageManager(t *testing.T) {
	t.Run("keygen_with_signer", func(t *testing.T) {
		ctx, mc := mockForPostedMessageManager()
		pmm := mc.PostedMessageManager()

		msg := &types.KeygenWithSigner{
			Signer: "signer",
			Data:   &types.Keygen{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerKeygen(mc)
		_, err := h.doKeygen(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("keygen_result_with_signer", func(t *testing.T) {
		ctx, mc := mockForPostedMessageManager()
		pmm := mc.PostedMessageManager()

		msg := &types.KeygenResultWithSigner{
			Signer: "signer",
			Keygen: &types.Keygen{},
			Data:   &types.KeygenResult{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerKeygenResult(mc)
		_, err := h.doKeygenResult(ctx, msg.Keygen, []*types.KeygenResultWithSigner{msg})
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})
}
