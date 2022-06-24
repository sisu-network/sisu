package sisu

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"

	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	ctx := testContext()
	keeper := keeperTestGenesis(ctx)

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
	contract := &types.Contract{
		Chain: "ganache1",
		Hash:  hash,
	}

	txOutProducer := DefaultTxOutputProducer{
		keeper: keeper,
		tssConfig: config.TssConfig{
			DeyesUrl: "http://0.0.0.0:1234",
			SupportedChains: map[string]config.TssChainConfig{
				"eth": {
					Id: "eth",
				},
			},
		},
	}

	tx := txOutProducer.getContractTx(ctx, contract, 100)
	require.NotNil(t, tx)
	require.EqualValues(t, 100, tx.Nonce())
	require.EqualValues(t, *big.NewInt(5_000_000_000), *tx.GasPrice())
	require.EqualValues(t, *big.NewInt(5_000_000_000), *tx.GasFeeCap())
}

func TestTxOutProducer_getEthResponse2(t *testing.T) {
	t.Parallel()

	t.Run("transaction_send_to_key", func(t *testing.T) {
		t.Parallel()

		ethTx := defaultTestEthTx(0)
		ctx := testContext()
		keeper := keeperTestAfterKeygen(ctx)

		appKeys := common.NewMockAppKeys()

		worldState := world.NewWorldState(keeper, &tssclients.MockDeyesClient{})

		binary, err := ethTx.MarshalBinary()
		require.NoError(t, err)

		txIn := types.TxIn{
			BlockHeight: 1,
			Serialized:  binary,
			Chain:       "ganache1",
		}

		txOutProducer := NewTxOutputProducer(worldState, appKeys,
			keeper,
			config.TssConfig{
				DeyesUrl: "http://0.0.0.0:1234",
				SupportedChains: map[string]config.TssChainConfig{
					"ganache1": {
						Id: "ganache",
					},
				},
			},
			config.CardanoConfig{},
			&MockCardanoNode{},
		).(*DefaultTxOutputProducer)

		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &txIn)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)

		// TODO Check the output of txOut to make sure that they are correct.
	})

	t.Run("transaction_send_to_contract", func(t *testing.T) {
		t.Parallel()

		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		txOutHash := "someTxOutHash"
		contractAddress := ecommon.HexToAddress("0x08BAB502c5e7125fD558B19a98D14907CF7f7E93")
		keeper.SaveTxOut(ctx, &types.TxOut{
			OutChain:     "ganache1",
			OutHash:      txOutHash,
			ContractHash: "contractHash",
		})
		keeper.CreateContractAddress(ctx, "ganache1", txOutHash, contractAddress.String())

		worldState := defaultWorldStateTest(ctx, keeper, &tssclients.MockDeyesClient{})
		appKeys := common.NewMockAppKeys()

		// Create transfer tx
		abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
		require.NoError(t, err)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		data, err := abi.Pack(MethodTransferOut, "ganache2", contractAddress.String(), tokenAddr, tokenAddr, amount)
		require.NoError(t, err)

		gasLimit := uint64(100)
		gasPrice := big.NewInt(100)
		ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &contractAddress,
			Value:    big.NewInt(100),
			Data:     data,
		})
		binary, err := ethTransaction.MarshalBinary()
		require.NoError(t, err)

		observedTx := types.TxIn{
			BlockHeight: 1,
			Chain:       "ganache1",
			Serialized:  binary,
		}

		txOutProducer := DefaultTxOutputProducer{
			worldState: worldState,
			tssConfig: config.TssConfig{
				DeyesUrl: "http://0.0.0.0:1234",
				SupportedChains: map[string]config.TssChainConfig{
					"ganache": {
						Id: "ganache",
					},
				},
			},
			keeper:  keeper,
			appKeys: appKeys,
		}

		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)

		// TODO Check the output of txOut to make sure that they are correct.
	})
}
