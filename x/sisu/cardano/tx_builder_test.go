package cardano

import (
	"fmt"
	"testing"

	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/crypto"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestTxBuilder_Fee(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		node := &MockCardanoNode{}
		sender, err := cardano.NewAddress("addr_test1vp9uhllavnhwc6m6422szvrtq3eerhleer4eyu00rmx8u6c42z3v8")
		require.NoError(t, err)
		receiver, err := cardano.NewAddress("addr_test1vqyqp03az6w8xuknzpfup3h7ghjwu26z7xa6gk7l9j7j2gs8zfwcy")
		require.NoError(t, err)

		policyID := NewDummyPolicyID(t)
		assetName := "WRAP_ANIMAL"
		cAssetName := cardano.NewAssetName(assetName)
		// sender has 5 ADA and 1_000_000_000 ANIMAL
		node.UTxOsFunc = func(_ cardano.Address) ([]cardano.UTxO, error) {
			asset := cardano.NewAssets().Set(cAssetName, 1_000_000_000)
			multiAsset := cardano.NewMultiAsset().Set(policyID, asset)
			utxos := []cardano.UTxO{}
			utxos = append(utxos, cardano.UTxO{
				TxHash:  nil,
				Spender: sender,
				Amount:  cardano.NewValueWithAssets(cardano.Coin(5*utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset),
				Index:   0,
			})

			return utxos, nil
		}

		node.TipFunc = func() (*cardano.NodeTip, error) {
			return &cardano.NodeTip{
				Block: 1,
				Epoch: 2,
				Slot:  20,
			}, nil
		}

		node.ProtocolParamsFunc = func() (*cardano.ProtocolParams, error) {
			return &cardano.ProtocolParams{
				MinFeeA:          5,
				MinFeeB:          10,
				CoinsPerUTXOWord: 20,
			}, nil
		}

		// send 1 ADA and 1_000_000 ANIMAL
		assetAmount := uint64(10_000_000)
		asset := cardano.NewAssets().Set(cAssetName, cardano.BigNum(assetAmount))
		multiAsset := cardano.NewMultiAsset().Set(policyID, asset)
		amount := cardano.NewValueWithAssets(cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()), multiAsset)

		adaPrice := 10
		destChain := "cardano-testnet"
		token := &types.Token{
			Id:        "ANIMAL",
			Price:     5,
			Chains:    []string{destChain},
			Addresses: []string{fmt.Sprintf("%s:%s", policyID.String(), assetName)},
		}
		tx, err := BuildTx(node, cardano.Testnet, sender, receiver, amount, nil, int64(adaPrice), token, destChain, assetAmount)
		require.NoError(t, err)
		require.Len(t, tx.Body.Outputs, 2)
		output := tx.Body.Outputs[0]
		require.Greater(t, uint64(output.Amount.Coin), uint64(1_300_000)) // must not less than 1,3 ADA

		requiredAdaFeeInToken := 2_600_000
		// User must receive less than <Origin amount - required ADA fee>
		require.Less(t, uint64(10_000_000-requiredAdaFeeInToken), uint64(output.Amount.MultiAsset.Get(policyID).Get(cAssetName)))
	})
}

func NewDummyPolicyID(t *testing.T) cardano.PolicyID {
	policyKey := crypto.NewXPrvKeyFromEntropy([]byte("policy"), "")
	policyScript, err := cardano.NewScriptPubKey(policyKey.PubKey())
	require.NoError(t, err)
	policyID, err := cardano.NewPolicyID(policyScript)
	require.NoError(t, err)
	return policyID
}
