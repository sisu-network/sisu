package sisu

import (
	"testing"

	ctypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerContractPauseResume() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	pmm := NewPostedMessageManager(k)
	globalData := &common.MockGlobalData{}
	dheartClient := &tssclients.MockDheartClient{}
	partyManager := &MockPartyManager{}
	partyManager.GetActivePartyPubkeysFunc = func() []ctypes.PubKey {
		return []ctypes.PubKey{}
	}
	txOutProducer := &MockTxOutputProducer{}
	txOutProducer.PauseContractFunc = func(ctx sdk.Context, chain, hash string) (*types.TxOutMsg, error) {
		ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{})
		binary, _ := ethTx.MarshalBinary()

		txOutWithSigner := &types.TxOutMsg{
			Signer: "signer",
			Data: &types.TxOut{
				OutChain: "ganache1",
				OutBytes: binary,
			},
		}

		return txOutWithSigner, nil
	}

	txOutProducer.ResumeContractFunc = func(ctx sdk.Context, chain, hash string) (*types.TxOutMsg, error) {
		ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{})
		binary, _ := ethTx.MarshalBinary()

		txOutWithSigner := &types.TxOutMsg{
			Signer: "signer",
			Data: &types.TxOut{
				OutChain: "ganache1",
				OutBytes: binary,
			},
		}

		return txOutWithSigner, nil
	}

	mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)

	return ctx, mc
}

func TestHandlerContractPauseResume_doPauseOrResume(t *testing.T) {
	t.Parallel()

	t.Run("is_catching_up", func(t *testing.T) {
		ctx, mc := mockForHandlerContractPauseResume()

		globalData := mc.GlobalData().(*common.MockGlobalData)
		globalData.IsCatchingUpFunc = func() bool {
			return true
		}

		h := newHandlerPauseResumeContract(mc)

		chain := "ganache1"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		_, err := h.doPauseOrResume(ctx, chain, hash, true)
		require.NoError(t, err)
	})

	t.Run("contract_not_found", func(t *testing.T) {
		ctx, mc := mockForHandlerContractPauseResume()
		h := newHandlerPauseResumeContract(mc)

		chain := "ganache1"
		hash := "hash"
		_, err := h.doPauseOrResume(ctx, chain, hash, true)
		require.Error(t, err)
	})

	t.Run("pause_contract", func(t *testing.T) {
		ctx, mc := mockForHandlerContractPauseResume()
		h := newHandlerPauseResumeContract(mc)

		chain := "ganache1"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		_, err := h.doPauseOrResume(ctx, chain, hash, true)
		require.NoError(t, err)
	})

	t.Run("resume_contract", func(t *testing.T) {
		ctx, mc := mockForHandlerContractPauseResume()

		h := newHandlerPauseResumeContract(mc)

		chain := "ganache1"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		_, err := h.doPauseOrResume(ctx, chain, hash, false)
		require.NoError(t, err)
	})
}
