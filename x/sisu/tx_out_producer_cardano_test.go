package sisu

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	"github.com/echovl/cardano-go/crypto"
	"github.com/sisu-network/sisu/utils"
	scardano "github.com/sisu-network/sisu/x/sisu/cardano"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForTestDefaultTxOutputProducer() (sdk.Context, keeper.Keeper, *scardano.MockCardanoClient) {
	ctx := testContext()
	k := keeperTestAfterContractDeployed(ctx)

	client := &scardano.MockCardanoClient{}
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

func DummyPolicyId() cardano.PolicyID {
	policyKey := crypto.NewXPrvKeyFromEntropy([]byte("policy"), "")
	policyScript, _ := cardano.NewScriptPubKey(policyKey.PubKey())
	policyID, _ := cardano.NewPolicyID(policyScript)

	return policyID
}

func getMultiAssetWithBalance(assetName string, balance int) *cardano.MultiAsset {
	policyID := DummyPolicyId()
	cAssetName := cardano.NewAssetName(assetName)

	asset := cardano.NewAssets().Set(cAssetName, cardano.BigNum(1_000_000_000*balance)) // 1000 Multi asset Token
	multiAsset := cardano.NewMultiAsset().Set(policyID, asset)

	return multiAsset
}

func TestDefaultTxOutputProducer_getCardanoTx(t *testing.T) {
	ctx, k, client := mockForTestDefaultTxOutputProducer()

	txOutProducer := &DefaultTxOutputProducer{
		keeper:        k,
		cardanoClient: client,
	}

	transfers := []*types.Transfer{
		{
			Id:        "ganache1_hash1",
			Recipient: "addr_test1vq0qus2tc5g2463xkng2g584gynlxs58t64d973dpu9gmqc2rrjv5",
			Token:     "SISU",
			Amount:    utils.EthToWei.String(),
		},
	}

	// Mock UTXOs
	hash, err := cardano.NewHash32("dd92bb91ac05247d21665a89fbac5958dc0d490605255571a5cc82cbcf9f2fba")
	if err != nil {
		require.Error(t, err)
	}
	sender, err := cardano.NewAddress("addr_test1vp9uhllavnhwc6m6422szvrtq3eerhleer4eyu00rmx8u6c42z3v8")
	balance := cardano.NewValueWithAssets(
		cardano.Coin(utils.ONE_ADA_IN_LOVELACE.Uint64()*10),
		getMultiAssetWithBalance("uSISU", 1),
	)
	utxos := []cardano.UTxO{
		{
			TxHash:  hash,
			Spender: sender,
			Amount:  balance,
			Index:   0,
		},
	}

	client.UTxOsFunc = func(addr cardano.Address, maxBlock uint64) ([]cardano.UTxO, error) {
		return utxos, nil
	}

	client.BalanceFunc = func(address cardano.Address) (*cardano.Value, error) {
		return balance, nil
	}

	// Get tx
	tx, err := txOutProducer.getCardanoTx(ctx, "cardano-testnet", transfers, utxos, uint64(100))
	require.Nil(t, err)

	fmt.Println("tx.Body.Outputs = ", tx.Body.Outputs)

	require.Equal(t, 2, len(tx.Body.Outputs))
}
