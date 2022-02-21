package sisu

import (
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/golang/mock/gomock"
	libchain "github.com/sisu-network/lib/chain"
	mock "github.com/sisu-network/sisu/tests/mock/common"
	mocktss "github.com/sisu-network/sisu/tests/mock/tss"
	mocksisu "github.com/sisu-network/sisu/tests/mock/x/sisu"

	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	t.Parallel()

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
	contract := &types.Contract{
		Chain: "eth",
		Hash:  hash,
	}

	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	mockPublicDb := mocktss.NewMockStorage(ctrl)
	mockPublicDb.EXPECT().GetChain("eth").Return(&types.Chain{
		Id:       "eth",
		GasPrice: int64(400_000_000_000),
	})
	mockPublicDb.EXPECT().GetLiquidity("eth").Return(&types.Liquidity{
		Id:      "eth",
		Address: "0x1234",
	})

	txOutProducer := DefaultTxOutputProducer{
		publicDb: mockPublicDb,
		tssConfig: config.TssConfig{
			DeyesUrl: "http://0.0.0.0:1234",
			SupportedChains: map[string]config.TssChainConfig{
				"eth": {
					Id: "eth",
				},
			},
		},
	}

	tx := txOutProducer.getContractTx(contract, 100)
	require.NotNil(t, tx)
	require.EqualValues(t, 100, tx.Nonce())
	require.EqualValues(t, *big.NewInt(400_000_000_000), *tx.GasPrice())
	require.EqualValues(t, *big.NewInt(400_000_000_000), *tx.GasFeeCap())
}

func TestTxOutProducer_getEthResponse(t *testing.T) {
	t.Parallel()

	t.Run("transaction_send_to_key", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		t.Cleanup(func() {
			ctrl.Finish()
		})

		mockPublicDb := mocktss.NewMockStorage(ctrl)
		mockPublicDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, gomock.Any()).Return(true).Times(1)
		mockPublicDb.EXPECT().GetPendingContracts("eth").Return([]*types.Contract{
			{
				Chain: "eth",
				Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
			},
		}).Times(1)
		mockPublicDb.EXPECT().GetChain("eth").Return(&types.Chain{
			Id:       "eth",
			GasPrice: int64(400_000_000_000),
		})
		mockPublicDb.EXPECT().GetLiquidity("eth").Return(&types.Liquidity{
			Id:      "eth",
			Address: "0x1234",
		})

		mockAppKeys := mock.NewMockAppKeys(ctrl)
		accAddress := []byte{1, 2, 3}
		mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).AnyTimes()

		mockWorldState := mocksisu.NewMockWorldState(ctrl)
		mockWorldState.EXPECT().UseAndIncreaseNonce("eth").Return(int64(0)).Times(1)

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

		txOutProducer := NewTxOutputProducer(mockWorldState, mockAppKeys, mockPublicDb,
			config.TssConfig{
				DeyesUrl: "http://0.0.0.0:1234",
				SupportedChains: map[string]config.TssChainConfig{
					"ganache1": {
						Id: "ganache",
					},
				},
			},
		).(*DefaultTxOutputProducer)

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

		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		recipient := common.HexToAddress("0x2d532C099CA476780c7703610D807948ae47856A")
		tokenAddr := common.HexToAddress("0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C")

		mockPublicDb := mocktss.NewMockStorage(ctrl)
		mockPublicDb.EXPECT().IsKeygenAddress(libchain.KEY_TYPE_ECDSA, gomock.Any()).Return(false).Times(1)
		mockPublicDb.EXPECT().IsContractExistedAtAddress("eth", gomock.Any()).Return(true).Times(1)
		mockPublicDb.EXPECT().GetLatestContractAddressByName(gomock.Any(), ContractErc20Gateway).Return("0x12345").Times(1)

		mockWorldState := mocksisu.NewMockWorldState(ctrl)
		mockWorldState.EXPECT().GetTokenFromAddress("eth", tokenAddr.String()).Return(&types.Token{
			Id:    "SISU",
			Price: 1_000_000_000,
		}).Times(1)
		mockWorldState.EXPECT().UseAndIncreaseNonce("eth").Return(int64(0)).Times(1)
		mockWorldState.EXPECT().GetGasPrice("eth").Return(big.NewInt(1_000_000_000), nil).Times(1)
		mockWorldState.EXPECT().GetNativeTokenPriceForChain("eth").Return(int64(1_000_000_000), nil).Times(1)

		mockAppKeys := mock.NewMockAppKeys(ctrl)
		accAddress := []byte{1, 2, 3}
		mockAppKeys.EXPECT().GetSignerAddress().Return(accAddress).AnyTimes()

		abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
		require.NoError(t, err)
		data, err := abi.Pack(MethodTransferOut, "eth", recipient, tokenAddr, tokenAddr, amount)

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

		txOutProducer := DefaultTxOutputProducer{
			worldState: mockWorldState,
			tssConfig: config.TssConfig{
				DeyesUrl: "http://0.0.0.0:1234",
				SupportedChains: map[string]config.TssChainConfig{
					"ganache": {
						Id: "ganache",
					},
				},
			},
			publicDb: mockPublicDb,
			appKeys:  mockAppKeys,
		}

		ctx := sdk.Context{}
		txOuts, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)

		// TODO Check the output of txOut to make sure that they are correct.
	})
}
