package cardano

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTestDefaultTxOutputProducer() (sdk.Context, keeper.Keeper, *MockCardanoClient) {
	ctx := testmock.TestContext()
	k := testmock.KeeperTestAfterContractDeployed(ctx)

	client := &MockCardanoClient{}
	client.ProtocolParamsFunc = func() (*cardano.ProtocolParams, error) {
		return &cardano.ProtocolParams{
			MinFeeA:          5,
			MinFeeB:          10,
			CoinsPerUTXOWord: 20,
		}, nil
	}

	client.TipFunc = func() (*cardano.NodeTip, error) {
		return &cardano.NodeTip{
			Block: 1,
			Epoch: 2,
			Slot:  20,
		}, nil
	}

	return ctx, k, client
}

func getMultiAssetWithBalance(assetName string, balance int) *cardano.MultiAsset {
	policyID := DummyPolicyId()
	cAssetName := cardano.NewAssetName(assetName)

	asset := cardano.NewAssets().Set(cAssetName, cardano.BigNum(balance)) // 1000 Multi asset Token
	multiAsset := cardano.NewMultiAsset().Set(policyID, asset)

	return multiAsset
}

func mockUtxos(hash cardano.Hash32, sender cardano.Address, balance *cardano.Value) []cardano.UTxO {
	return []cardano.UTxO{
		{
			TxHash:  hash,
			Spender: sender,
			Amount:  balance,
			Index:   0,
		},
	}
}

func mockClient(client *MockCardanoClient, utxos []cardano.UTxO, balance *cardano.Value) {
	client.UTxOsFunc = func(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error) {
		return utxos, nil
	}

	client.BalanceFunc = func(address cardano.Address) (*cardano.Value, error) {
		return balance, nil
	}

}

func TestBridge_getCardanoTx(t *testing.T) {
	ctx, k, client := mockForTestDefaultTxOutputProducer()
	bridge := NewBridge("cardano", "signer", k, client).(*bridge)

	// Set the commission fee
	params := k.GetParams(ctx)
	params.CommissionRate = 10 // 0.1%. One unit is 1/10,000
	k.SaveParams(ctx, params)

	var balance *cardano.Value
	// Mock UTXOs
	hash, err := cardano.NewHash32("dd92bb91ac05247d21665a89fbac5958dc0d490605255571a5cc82cbcf9f2fba")
	require.Nil(t, err)
	sender, err := cardano.NewAddress("addr_test1vp9uhllavnhwc6m6422szvrtq3eerhleer4eyu00rmx8u6c42z3v8")
	require.Nil(t, err)
	recipient1, err := cardano.NewAddress("addr_test1vq0qus2tc5g2463xkng2g584gynlxs58t64d973dpu9gmqc2rrjv5")
	require.Nil(t, err)
	recipient2, err := cardano.NewAddress("addr_test1vpjdtcsa7kq9l9l3ahmfkvg6fn03k7ky87luhggt2hl4mhg7u9ly6")

	t.Run("get_tx_success", func(t *testing.T) {
		balance = cardano.NewValueWithAssets(
			cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()*10),
			getMultiAssetWithBalance("uSISU", 100*1_000_000),
		)

		utxos := mockUtxos(hash, sender, balance)
		mockClient(client, utxos, balance)
		transfers := []*types.Transfer{
			{
				Id:          "ganache1_hash1",
				ToRecipient: recipient1.String(),
				Token:       "SISU",
				Amount:      utils.EthToWei.String(), // Transfer 1 Sisu
			},
		}

		// Get tx
		tx, err := bridge.getCardanoTx(ctx, "cardano-testnet", transfers, utxos, uint64(100))
		require.Nil(t, err)
		require.Equal(t, 2, len(tx.Body.Outputs))
		require.Equal(t, recipient1, tx.Body.Outputs[1].Address)

		// ADA price: 0.4 USD. Sisu price: 4 USD. Transaction fee = 1.6 ADA = 0.16 Sisu
		// Total input: 1 Sisu
		// Commission fee: 0.001 Sisu
		// Transaction fee: 0.16 Sisu
		// Total output: 0.839
		expectedAmount := cardano.NewValueWithAssets(1_600_000, getMultiAssetWithBalance("uSISU", 839_000))
		require.Nil(t, err)
		require.Equal(t,
			0,
			tx.Body.Outputs[1].Amount.Cmp(expectedAmount),
		)
	})

	t.Run("get_tx_not_enough_balance", func(t *testing.T) {
		balance = cardano.NewValueWithAssets(
			cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()*10),
			getMultiAssetWithBalance("uSISU", 100*1_000_000), // 100 Sisu tokens
		)
		utxos := mockUtxos(hash, sender, balance)
		mockClient(client, utxos, balance)

		transfers := []*types.Transfer{
			{
				Id:          "ganache1_hash1",
				ToRecipient: recipient1.String(),
				Token:       "SISU",
				Amount:      (new(big.Int).Mul(utils.EthToWei, big.NewInt(200))).String(), // 200 Sisu
			},
		}

		// Transfer 200 tokens (exceed balanaced of 100 tokens)
		_, err := bridge.getCardanoTx(ctx, "cardano-testnet", transfers, utxos, uint64(100))
		require.NotNil(t, err)
	})

	t.Run("transfer_multiple_assets", func(t *testing.T) {
		balance = cardano.NewValueWithAssets(
			cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()*10),
			getMultiAssetWithBalance("uSISU", 100*1_000_000), // 100 Sisu tokens
		)
		utxos := mockUtxos(hash, sender, balance)
		mockClient(client, utxos, balance)

		transfers := []*types.Transfer{
			{
				Id:          "ganache1_hash1",
				ToRecipient: recipient1.String(),
				Token:       "SISU",
				Amount:      (new(big.Int).Mul(utils.EthToWei, big.NewInt(5))).String(),
			},
			{
				Id:          "ganache1_hash2",
				ToRecipient: recipient2.String(),
				Token:       "SISU",
				Amount:      (new(big.Int).Mul(utils.EthToWei, big.NewInt(10))).String(),
			},
		}

		// Get tx
		tx, err := bridge.getCardanoTx(ctx, "cardano-testnet", transfers, utxos, uint64(100))
		require.Nil(t, err)
		require.Equal(t, 3, len(tx.Body.Outputs))

		expectedAmount := cardano.NewValueWithAssets(1_600_000, getMultiAssetWithBalance("uSISU", 4835000))
		require.Nil(t, err)
		require.Equal(t,
			0,
			tx.Body.Outputs[1].Amount.Cmp(expectedAmount),
		)
		expectedAmount = cardano.NewValueWithAssets(1_600_000, getMultiAssetWithBalance("uSISU", 9830000))
		require.Nil(t, err)
		require.Equal(t,
			0,
			tx.Body.Outputs[2].Amount.Cmp(expectedAmount),
		)
	})
}
