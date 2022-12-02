package service

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/deyes/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestQueryRecentSolanBlock(t *testing.T) {
	chain := "solana-devnet"
	hash := "hash"
	height := int64(123)
	deyesClient := external.MockDeyesClient{
		SolanaQueryRecentBlockFunc: func(chain string) (*etypes.SolanaQueryRecentBlockResult, error) {
			return &etypes.SolanaQueryRecentBlockResult{
				Hash:   hash,
				Height: height,
			}, nil
		},
	}
	txSubmit := common.MockTxSubmit{
		SubmitMessageAsyncFunc: func(msg sdk.Msg) error {
			updateMsg, ok := msg.(*types.UpdateSolanaRecentHashMsg)
			require.True(t, ok)
			require.Equal(t, chain, updateMsg.Data.Chain)
			require.Equal(t, hash, updateMsg.Data.Hash)
			require.Equal(t, height, updateMsg.Data.Height)
			return nil
		},
	}

	chainPolling := NewChainPolling("signer", &deyesClient, &txSubmit)
	chainPolling.QueryRecentSolanBlock(chain)
}
