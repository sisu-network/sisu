package sisu

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTxOutPauseResumeContract(t *testing.T) {
	t.Parallel()

	t.Run("can_pause_or_resume", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
			signer: "signer",
		}

		chain := "ganache1"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		txOutWithSigner, err := txOutputProducer.PauseOrResumeEthContract(ctx, chain, hash, false)
		require.NoError(t, err)
		require.NotNil(t, txOutWithSigner)

		txOutWithSigner, err = txOutputProducer.PauseOrResumeEthContract(ctx, chain, hash, true)
		require.NoError(t, err)
		require.NotNil(t, txOutWithSigner)
	})

	t.Run("unsupported_chain", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
			signer: "signer",
		}

		chain := "chain"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		txOutWithSigner, err := txOutputProducer.PauseOrResumeEthContract(ctx, chain, hash, false)
		require.Error(t, err)
		require.Nil(t, txOutWithSigner)
	})

	t.Run("can_not_find_gateway_address", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterKeygen(ctx)
		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
			signer: "signer",
		}

		chain := "ganache1"
		hash := SupportedContracts[ContractErc20Gateway].AbiHash
		txOutWithSigner, err := txOutputProducer.PauseOrResumeEthContract(ctx, chain, hash, false)
		require.Error(t, err)
		require.Nil(t, txOutWithSigner)
	})
}
