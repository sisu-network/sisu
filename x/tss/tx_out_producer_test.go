package tss

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang/mock/gomock"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/tests/mock"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"

	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	t.Parallel()

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayMetaData.Bin)
	contract := &types.Contract{
		Chain: "eth",
		Hash:  hash,
	}

	worldState := NewWorldState(config.TssConfig{}, nil, nil)
	txOutProducer := DefaultTxOutputProducer{
		worldState: worldState,
		tssConfig: config.TssConfig{
			Enable: true,
			SupportedChains: map[string]config.TssChainConfig{
				"ganache": {
					Symbol:   "ganache",
					DeyesUrl: "http://0.0.0.0:1234",
				},
			},
		},
	}

	tx := txOutProducer.getContractTx(contract, 100)
	require.NotNil(t, tx)
	require.EqualValues(t, 100, tx.Nonce())
	require.EqualValues(t, *big.NewInt(400000000000), *tx.GasPrice())
	require.EqualValues(t, *big.NewInt(400000000000), *tx.GasFeeCap())
}

func TestTxOutProducer_getEthResponse(t *testing.T) {
	t.Parallel()

	t.Run("transaction_send_to_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		privKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		pubkeyBytes := crypto.FromECDSAPub(&privKey.PublicKey)
		mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
		mockPrivateDb.EXPECT().GetKeygenPubkey(libchain.KEY_TYPE_ECDSA).Return(pubkeyBytes).Times(1)

		keeper := mocktss.NewMockKeeper(ctrl)
		keeper.EXPECT().IsKeygenAddress(gomock.Any(), libchain.KEY_TYPE_ECDSA, gomock.Any()).Return(true).Times(1)
		keeper.EXPECT().GetPendingContracts(gomock.Any(), "eth").Return([]*types.Contract{
			{
				Chain: "eth",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}).Times(1)

		mockAppKeys := mock.NewMockAppKeys(ctrl)
		accAddress := []byte{1, 2, 3}
		mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).Times(1)

		amount := big.NewInt(100)
		gasLimit := uint64(100)
		gasPrice := big.NewInt(100)
		ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &common.Address{},
			Value:    amount,
		})
		binary, err := ethTransaction.MarshalBinary()
		require.NoError(t, err)

		observedTx := types.TxIn{
			BlockHeight: 1,
			Serialized:  binary,
			Chain:       "eth",
		}

		worldState := DefaultWorldState{
			privateDb: mockPrivateDb,
			tssConfig: config.TssConfig{},
			nonces: map[string]int64{
				"eth": 100,
			},
			deyesClients: nil,
		}
		txOutProducer := DefaultTxOutputProducer{
			worldState: &worldState,
			keeper:     keeper,
			tssConfig: config.TssConfig{
				Enable: true,
				SupportedChains: map[string]config.TssChainConfig{
					"ganache1": {
						Symbol:   "ganache",
						DeyesUrl: "http://0.0.0.0:1234",
					},
				},
			},
			appKeys: mockAppKeys,
		}

		ctx := sdk.Context{}
		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)

		// TODO Check the output of txOut to make sure that they are correct.
	})

	t.Run("transaction_send_to_contract", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		privKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		pubkeyBytes := crypto.FromECDSAPub(&privKey.PublicKey)
		mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
		mockPrivateDb.EXPECT().GetKeygenPubkey(libchain.KEY_TYPE_ECDSA).Return(pubkeyBytes).Times(1)

		keeper := mocktss.NewMockKeeper(ctrl)
		keeper.EXPECT().IsKeygenAddress(gomock.Any(), libchain.KEY_TYPE_ECDSA, gomock.Any()).Return(false).Times(1)
		keeper.EXPECT().IsContractExistedAtAddress(gomock.Any(), "eth", gomock.Any()).Return(true).Times(1)
		keeper.EXPECT().GetLatestContractAddressByName(gomock.Any(), gomock.Any(), ContractErc20Gateway).Return("0x12345").Times(1)

		mockAppKeys := mock.NewMockAppKeys(ctrl)
		accAddress := []byte{1, 2, 3}
		mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).Times(1)

		abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
		require.NoError(t, err)
		amount := big.NewInt(100)
		dummyHex := common.HexToAddress("ab1257528b3782fb40d7ed5f72e624b744dffb2f")
		data, err := abi.Pack(MethodTransferOut, "eth", dummyHex, dummyHex, dummyHex, &amount)
		require.NoError(t, err)

		gasLimit := uint64(100)
		gasPrice := big.NewInt(100)
		ethTransaction := ethTypes.NewTx(&ethTypes.LegacyTx{
			GasPrice: gasPrice,
			Gas:      gasLimit,
			To:       &common.Address{},
			Value:    big.NewInt(100),
			Data:     data,
		})
		binary, err := ethTransaction.MarshalBinary()
		require.NoError(t, err)

		observedTx := types.TxIn{
			BlockHeight: 1,
			Chain:       "eth",
			Serialized:  binary,
		}

		worldState := DefaultWorldState{
			privateDb: mockPrivateDb,
			tssConfig: config.TssConfig{},
			nonces: map[string]int64{
				"eth": 100,
			},
			deyesClients: nil,
		}
		txOutProducer := DefaultTxOutputProducer{
			worldState: &worldState,
			keeper:     keeper,
			tssConfig: config.TssConfig{
				Enable: true,
				SupportedChains: map[string]config.TssChainConfig{
					"ganache": {
						Symbol:   "ganache",
						DeyesUrl: "http://0.0.0.0:1234",
					},
				},
			},
			privateDb: mockPrivateDb,
			appKeys:   mockAppKeys,
		}

		ctx := sdk.Context{}
		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)

		// TODO Check the output of txOut to make sure that they are correct.
	})
}
