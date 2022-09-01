package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForHandlerTxIn() (sdk.Context, ManagerContainer) {
	ctx := testContext()
	k := keeperTestAfterContractDeployed(ctx)
	pmm := NewPostedMessageManager(k)
	k.SaveParams(ctx, &types.Params{
		MajorityThreshold: 1,
	})

	mc := MockManagerContainer(ctx, k, pmm, &MockTransferQueue{})
	return ctx, mc
}

func TestHandlerTxIn_HappyCase(t *testing.T) {
	t.Run("transfer_is_saved", func(t *testing.T) {
		ctx, mc := mockForHandlerTxIn()
		srcChain := "ganache1"
		destChain := "ganache2"
		recipient := "0x8095f5b69F2970f38DC6eBD2682ed71E4939f988"
		token := "SISU"
		hash1 := "123"
		amount := "10000"

		handler := NewHandlerTxIn(mc)
		msg := types.NewTxsInMsg("signer", &types.TxsIn{
			Chain:  srcChain,
			Height: 10,
			Requests: []*types.TransferOut{
				{
					ToChain:   destChain,
					Token:     token,
					Hash:      hash1,
					Recipient: recipient,
					Amount:    amount,
				},
			},
		})

		_, err := handler.DeliverMsg(ctx, msg)
		require.Nil(t, err)

		keeper := mc.Keeper()
		queue := keeper.GetTransferQueue(ctx, destChain)
		require.Equal(t, []*types.Transfer{
			{
				Id:        fmt.Sprintf("%s__%s", srcChain, hash1),
				Recipient: recipient,
				Token:     token,
				Amount:    amount,
			},
		}, queue)

		// Add the second request
		hash2 := "456"
		recipient2 := "0x98Fa8Ab1dd59389138B286d0BeB26bfa4808EC80"
		token2 := "ADA"
		handler = NewHandlerTxIn(mc)
		msg = types.NewTxsInMsg("signer", &types.TxsIn{
			Chain:  srcChain,
			Height: 11,
			Requests: []*types.TransferOut{
				{
					ToChain:   destChain,
					Token:     token2,
					Hash:      hash2,
					Recipient: recipient2,
					Amount:    amount,
				},
			},
		})
		_, err = handler.DeliverMsg(ctx, msg)
		require.Nil(t, err)
		queue = keeper.GetTransferQueue(ctx, destChain)
		require.Equal(t, []*types.Transfer{
			{
				Id:        fmt.Sprintf("%s__%s", srcChain, hash1),
				Recipient: recipient,
				Token:     token,
				Amount:    amount,
			},
			{
				Id:        fmt.Sprintf("%s__%s", srcChain, hash2),
				Recipient: recipient2,
				Token:     token2,
				Amount:    amount,
			},
		}, queue)
	})
}
