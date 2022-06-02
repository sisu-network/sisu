package sisu

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestPostedMessageManager(t *testing.T) {
	t.Parallel()

	t.Run("keygen_with_signer", func(t *testing.T) {
		ctx, mc := mockForHandlerKeygen()
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
		ctx, mc := mockForHandlerKeygen()
		pmm := mc.PostedMessageManager()

		msg := &types.KeygenResultWithSigner{
			Signer: "signer",
			Keygen: &types.Keygen{},
			Data:   &types.KeygenResult{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerKeygenResult(mc)
		_, err := h.doKeygenResult(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("tx_in_with_signer", func(t *testing.T) {
		ctx, mc := mockForHandlerTxIn()
		pmm := mc.PostedMessageManager()

		msg := &types.TxInWithSigner{
			Signer: "signer",
			Data:   &types.TxIn{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerTxIn(mc)
		_, err := h.doTxIn(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("tx_out_with_signer", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		pmm := mc.PostedMessageManager()

		msg := &types.TxOutWithSigner{
			Signer: "signer",
			Data:   &types.TxOut{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerTxOut(mc)
		_, err := h.doTxOut(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("tx_out_contract_confirm_with_signer", func(t *testing.T) {
		ctx, mc := mockForHandlerTxOut()
		pmm := mc.PostedMessageManager()

		msg := &types.TxOutContractConfirmWithSigner{
			Signer: "signer",
			Data:   &types.TxOutContractConfirm{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerTxOutContractConfirmation(mc)
		_, err := h.doTxOutContractConfirm(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("contract_with_signer", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		mc := MockManagerContainer(k, pmm)

		msg := &types.ContractsWithSigner{
			Signer: "signer",
			Data:   &types.Contracts{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerContract(mc)
		_, err := h.doContracts(ctx, msg)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("pause_contract_msg", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		globalData := &common.MockGlobalData{}
		dheartClient := &tssclients.MockDheartClient{}
		partyManager := &MockPartyManager{}
		txOutProducer := &MockTxOutputProducer{}
		txOutProducer.PauseContractFunc = func(ctx sdk.Context, chain, hash string) (*types.TxOutWithSigner, error) {
			ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				GasPrice: big.NewInt(100),
				Gas:      uint64(100),
				To:       &ecommon.Address{},
				Value:    big.NewInt(100),
			})
			binary, err := ethTx.MarshalBinary()
			require.NoError(t, err)

			txOutWithSigner := &types.TxOutWithSigner{
				Signer: "signer",
				Data: &types.TxOut{
					OutChain: "ganache1",
					OutBytes: binary,
				},
			}

			return txOutWithSigner, nil
		}

		msg := &types.PauseContractMsg{
			Signer: "signer",
			Data: &types.PauseContract{
				Chain: "ganache1",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}

		mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)
		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerPauseContract(mc)
		_, err := newHandlerPauseResumeContract(h.mc).doPauseOrResume(ctx, msg.Data.Chain, msg.Data.Hash, true)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("resume_contract_msg", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		globalData := &common.MockGlobalData{}
		dheartClient := &tssclients.MockDheartClient{}
		partyManager := &MockPartyManager{}
		txOutProducer := &MockTxOutputProducer{}
		txOutProducer.ResumeContractFunc = func(ctx sdk.Context, chain, hash string) (*types.TxOutWithSigner, error) {
			ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				GasPrice: big.NewInt(100),
				Gas:      uint64(100),
				To:       &ecommon.Address{},
				Value:    big.NewInt(100),
			})
			binary, err := ethTx.MarshalBinary()
			require.NoError(t, err)

			txOutWithSigner := &types.TxOutWithSigner{
				Signer: "signer",
				Data: &types.TxOut{
					OutChain: "ganache1",
					OutBytes: binary,
				},
			}

			return txOutWithSigner, nil
		}

		msg := &types.ResumeContractMsg{
			Signer: "signer",
			Data: &types.ResumeContract{
				Chain: "ganache1",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}

		mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)
		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerResumeContract(mc)
		_, err := newHandlerPauseResumeContract(h.mc).doPauseOrResume(ctx, msg.Data.Chain, msg.Data.Hash, false)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("change_ownership_contract_msg", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		globalData := &common.MockGlobalData{}
		dheartClient := &tssclients.MockDheartClient{}
		partyManager := &MockPartyManager{}
		txOutProducer := &MockTxOutputProducer{}
		txOutProducer.ContractChangeOwnershipFunc = func(ctx sdk.Context, chain, contractHash, newOwner string) (*types.TxOutWithSigner, error) {
			ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				GasPrice: big.NewInt(100),
				Gas:      uint64(100),
				To:       &ecommon.Address{},
				Value:    big.NewInt(100),
			})
			binary, err := ethTx.MarshalBinary()
			require.NoError(t, err)

			txOutWithSigner := &types.TxOutWithSigner{
				Signer: "signer",
				Data: &types.TxOut{
					OutChain: "ganache1",
					OutBytes: binary,
				},
			}

			return txOutWithSigner, nil
		}

		msg := &types.ChangeOwnershipContractMsg{
			Signer: "signer",
			Data: &types.ChangeOwnership{
				Chain: "ganache1",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}

		mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)
		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerContractChangeOwnership(mc)
		_, err := newHandlerContractChangeOwnership(h.mc).doChangeOwner(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.NewOwner)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("change_liquid_pool_address_msg", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		globalData := &common.MockGlobalData{}
		dheartClient := &tssclients.MockDheartClient{}
		partyManager := &MockPartyManager{}
		txOutProducer := &MockTxOutputProducer{}
		txOutProducer.ContractSetLiquidPoolAddressFunc = func(ctx sdk.Context, chain, contractHash, newAddress string) (*types.TxOutWithSigner, error) {
			ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				GasPrice: big.NewInt(100),
				Gas:      uint64(100),
				To:       &ecommon.Address{},
				Value:    big.NewInt(100),
			})
			binary, err := ethTx.MarshalBinary()
			require.NoError(t, err)

			txOutWithSigner := &types.TxOutWithSigner{
				Signer: "signer",
				Data: &types.TxOut{
					OutChain: "ganache1",
					OutBytes: binary,
				},
			}

			return txOutWithSigner, nil
		}

		msg := &types.ChangeLiquidPoolAddressMsg{
			Signer: "signer",
			Data: &types.ChangeLiquidAddress{
				Chain: "ganache1",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}

		mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)
		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerContractSetLiquidityAddress(mc)
		_, err := newHandlerContractSetLiquidityAddress(h.mc).doSetLiquidityAddress(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.NewLiquidAddress)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("liquidity_withdraw_fund_msg", func(t *testing.T) {
		ctx := testContext()
		k := keeperTestGenesis(ctx)
		pmm := NewPostedMessageManager(k)
		globalData := &common.MockGlobalData{}
		dheartClient := &tssclients.MockDheartClient{}
		partyManager := &MockPartyManager{}
		txOutProducer := &MockTxOutputProducer{}
		txOutProducer.ContractEmergencyWithdrawFundFunc = func(ctx sdk.Context, chain, contractHash string, tokens []string, newOwner string) (*types.TxOutWithSigner, error) {
			ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
				GasPrice: big.NewInt(100),
				Gas:      uint64(100),
				To:       &ecommon.Address{},
				Value:    big.NewInt(100),
			})
			binary, err := ethTx.MarshalBinary()
			require.NoError(t, err)

			txOutWithSigner := &types.TxOutWithSigner{
				Signer: "signer",
				Data: &types.TxOut{
					OutChain: "ganache1",
					OutBytes: binary,
				},
			}

			return txOutWithSigner, nil
		}

		msg := &types.LiquidityWithdrawFundMsg{
			Signer: "signer",
			Data: &types.LiquidityWithdrawFund{
				Chain: "ganache1",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}

		mc := MockManagerContainer(k, pmm, globalData, txOutProducer, partyManager, dheartClient)
		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		h := NewHandlerContractLiquidityWithdrawFund(mc)
		_, err := h.doWithdrawFund(ctx, msg.Data.Chain, msg.Data.Hash, msg.Data.TokenAddresses, msg.Data.NewOwner)
		require.NoError(t, err)

		h.keeper.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})
}
