package components

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForPostedMessageManager() (sdk.Context, keeper.Keeper) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestGenesis(ctx)

	return ctx, k
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
		ctx, k := mockForPostedMessageManager()
		pmm := NewPostedMessageManager(k)

		msg := &types.KeygenWithSigner{
			Signer: "signer",
			Data:   &types.Keygen{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		k.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})

	t.Run("keygen_result_with_signer", func(t *testing.T) {
		ctx, k := mockForPostedMessageManager()
		pmm := NewPostedMessageManager(k)

		msg := &types.KeygenResultWithSigner{
			Signer: "signer",
			Keygen: &types.Keygen{},
			Data:   &types.KeygenResult{},
		}

		process, hash := pmm.ShouldProcessMsg(ctx, msg)
		require.True(t, process)

		k.ProcessTxRecord(ctx, hash)
		process, _ = pmm.ShouldProcessMsg(ctx, msg)
		require.False(t, process)
	})
}
