package sisu

import (
	"testing"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	// ctx := testContext()
	// keeper := keeperTestGenesis(ctx)

	// hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
	// contract := &types.Contract{
	// 	Chain: "ganache1",
	// 	Hash:  hash,
	// }

	// txOutProducer := DefaultTxOutputProducer{
	// 	keeper: keeper,
	// 	tssConfig: config.TssConfig{
	// 		DeyesUrl: "http://0.0.0.0:1234",
	// 		SupportedChains: map[string]config.TssChainConfig{
	// 			"eth": {
	// 				Id: "eth",
	// 			},
	// 		},
	// 	},
	// }

	// tx := txOutProducer.getContractTx(ctx, contract, 100)
	// require.NotNil(t, tx)
	// require.EqualValues(t, 100, tx.Nonce())
	// require.EqualValues(t, *big.NewInt(5_000_000_000), *tx.GasPrice())
	// require.EqualValues(t, *big.NewInt(5_000_000_000), *tx.GasFeeCap())
}

func TestTxOutProducer_getEthResponse2(t *testing.T) {
	t.Run("transaction_send_to_key", func(t *testing.T) {
		// t.Parallel()

		// ctx := testContext()
		// keeper := keeperTestAfterKeygen(ctx)
		// appKeys := common.NewMockAppKeys()

		// worldState := world.NewWorldState(keeper, &tssclients.MockDeyesClient{})
		// transfer := &types.Transfer{
		// 	Id:        "ganache1_hash1",
		// 	Recipient: "0x123",
		// 	Amount:    "1000",
		// }

		// txOutProducer := NewTxOutputProducer(worldState, appKeys,
		// 	keeper,
		// 	nil,
		// 	config.CardanoConfig{},
		// 	nil,
		// 	&MockTxTracker{},
		// ).(*DefaultTxOutputProducer)

		// txOuts, err := txOutProducer.GetTxOuts(ctx, "ganache2", []*types.Transfer{transfer})
		// require.Nil(t, err)
		// require.Len(t, txOuts, 1)

		// // TODO Check the output of txOut to make sure that they are correct.
	})

	t.Run("transaction_send_to_contract", func(t *testing.T) {
		// t.Parallel()

		// ctx := testContext()
		// keeper := keeperTestAfterContractDeployed(ctx)
		// txOutHash := "someTxOutHash"
		// contractAddress := ecommon.HexToAddress("0x08BAB502c5e7125fD558B19a98D14907CF7f7E93")
		// keeper.SaveTxOut(ctx, &types.TxOut{
		// 	OutChain:     "ganache1",
		// 	OutHash:      txOutHash,
		// 	ContractHash: "contractHash",
		// })
		// keeper.CreateContractAddress(ctx, "ganache1", txOutHash, contractAddress.String())

		// worldState := defaultWorldStateTest(ctx, keeper, &tssclients.MockDeyesClient{})
		// appKeys := common.NewMockAppKeys()

		// // Create transfer tx
		// abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
		// require.NoError(t, err)
		// amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		// tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		// data, err := abi.Pack(MethodTransferOut, "ganache2", contractAddress.String(), tokenAddr, tokenAddr, amount)
		// require.NoError(t, err)

		// gasLimit := uint64(100)
		// gasPrice := big.NewInt(100)
		// ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
		// 	GasPrice: gasPrice,
		// 	Gas:      gasLimit,
		// 	To:       &contractAddress,
		// 	Value:    big.NewInt(100),
		// 	Data:     data,
		// })
		// binary, err := ethTransaction.MarshalBinary()
		// require.NoError(t, err)

		// observedTx := &types.TxIn{
		// 	BlockHeight: 1,
		// 	Chain:       "ganache1",
		// 	Serialized:  binary,
		// }

		// txOutProducer := DefaultTxOutputProducer{
		// 	worldState: worldState,
		// 	tssConfig: config.TssConfig{
		// 		DeyesUrl: "http://0.0.0.0:1234",
		// 		SupportedChains: map[string]config.TssChainConfig{
		// 			"ganache": {
		// 				Id: "ganache",
		// 			},
		// 		},
		// 	},
		// 	keeper:    keeper,
		// 	appKeys:   appKeys,
		// 	txTracker: &MockTxTracker{},
		// }

		// txOuts, _ := txOutProducer.GetTxOuts(TxInRequest{
		// 	Ctx: ctx,
		// }, []*types.TxIn{observedTx})
		// require.Len(t, txOuts, 1)

		// // TODO Check the output of txOut to make sure that they are correct.
	})
}
