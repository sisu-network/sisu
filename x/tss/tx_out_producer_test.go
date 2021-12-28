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

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
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
	require.EqualValues(t, *big.NewInt(10000000000), *tx.GasPrice())
	require.EqualValues(t, *big.NewInt(10000000000), *tx.GasFeeCap())
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
		mockPrivateDb.EXPECT().GetKeygenPubkey("ecdsa").Return(pubkeyBytes).Times(1)

		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(true).Times(1)
		mockDb.EXPECT().GetPubKey("ecdsa").Return(pubkeyBytes).Times(1)

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
			tssConfig: config.TssConfig{
				Enable: true,
				SupportedChains: map[string]config.TssChainConfig{
					"ganache": {
						Symbol:   "ganache",
						DeyesUrl: "http://0.0.0.0:1234",
					},
				},
			},
			db:      mockDb,
			appKeys: mockAppKeys,
		}

		ctx := sdk.Context{}
		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)
	})

	t.Run("transaction_send_to_contract", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		contractEntity := &types.ContractEntity{
			Chain: "eth",
			Hash:  SupportedContracts[ContractErc20].AbiHash,
		}

		privKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		pubkeyBytes := crypto.FromECDSAPub(&privKey.PublicKey)
		mockPrivateDb := mocktss.NewMockPrivateDb(ctrl)
		mockPrivateDb.EXPECT().GetKeygenPubkey("ecdsa").Return(pubkeyBytes).Times(1)

		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(false).Times(1)
		mockDb.EXPECT().GetContractFromAddress(gomock.Any(), gomock.Any()).Return(contractEntity).Times(1)
		mockDb.EXPECT().GetContractFromHash(gomock.Any(), gomock.Any()).Return(contractEntity).Times(1)
		mockDb.EXPECT().GetPubKey("ecdsa").Return(pubkeyBytes).Times(1)

		mockAppKeys := mock.NewMockAppKeys(ctrl)
		accAddress := []byte{1, 2, 3}
		mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).Times(1)

		abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20].AbiString))
		require.NoError(t, err)
		amount := big.NewInt(100)
		hex := "ab1257528b3782fb40d7ed5f72e624b744dffb2f"
		data, err := abi.Pack(MethodTransferOutFromContract, common.HexToAddress(hex), "eth", hex, &amount)
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
			tssConfig: config.TssConfig{
				Enable: true,
				SupportedChains: map[string]config.TssChainConfig{
					"ganache": {
						Symbol:   "ganache",
						DeyesUrl: "http://0.0.0.0:1234",
					},
				},
			},
			db:      mockDb,
			appKeys: mockAppKeys,
		}

		ctx := sdk.Context{}
		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)
	})
}
