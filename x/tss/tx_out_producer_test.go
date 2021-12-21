package tss

import (
	"context"
	"github.com/ethereum/go-ethereum/ethclient"
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

	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/types"
	"github.com/stretchr/testify/require"
)

func TestTxOutProducer_getContractTx(t *testing.T) {
	t.Parallel()

	hash := utils.KeccakHash32(erc20gateway.Erc20gatewayBin)
	contractEntity := &types.ContractEntity{
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

	tx := txOutProducer.getContractTx(contractEntity, 100)
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

		contractEntities := []*types.ContractEntity{
			{
				Chain: "eth",
				Hash:  SupportedContracts[ContractErc20].AbiHash,
			},
		}

		privKey, err := crypto.GenerateKey()
		require.NoError(t, err)

		pubkeyBytes := crypto.FromECDSAPub(&privKey.PublicKey)
		mockDb := mock.NewMockDatabase(ctrl)
		mockDb.EXPECT().IsChainKeyAddress(gomock.Any(), gomock.Any()).Return(true).Times(1)
		mockDb.EXPECT().GetPendingDeployContracts(gomock.Any()).Return(contractEntities).Times(1)
		mockDb.EXPECT().GetPubKey("ecdsa").Return(pubkeyBytes).Times(1)
		mockDb.EXPECT().UpdateContractsStatus(gomock.Any(), gomock.Any()).Return(nil).Times(1)

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

		observedTx := types.ObservedTx{
			BlockHeight: 1,
			Serialized:  binary,
			Chain:       "eth",
		}

		worldState := DefaultWorldState{
			db:        mockDb,
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
		txOuts, txOutEntities, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)
		require.Len(t, txOutEntities, 1)
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

		observedTx := types.ObservedTx{
			BlockHeight: 1,
			Serialized:  binary,
		}

		worldState := DefaultWorldState{
			db:        mockDb,
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
		txOuts, txOutEntities, err := txOutProducer.getEthResponse(ctx, 1, &observedTx)
		require.NoError(t, err)
		require.Len(t, txOuts, 1)
		require.Len(t, txOutEntities, 1)
	})
}

func TestTxOutProducer_createERC20TransferIn(t *testing.T) {
	// Uncomment t.Skip() to run this test
	t.Skip()

	// Please run ganache and deploy ERC20 gateway + ERC20 token before running this test. See repo smart-contracts to get instructions
	txOutProducer := DefaultTxOutputProducer{}
	gatewayAddr := "0x5FbDB2315678afecb367f032d93F642f64180aa3"
	tokenAddr := "0xe7f1725E7734CE288F8367e1Bb143E90bb3F0512"
	recipient := "0xbcd4042de499d14e55001ccbb24a551f3b954096"
	amount := big.NewInt(9999)
	txResponse, err := txOutProducer.createERC20TransferIn(gatewayAddr, tokenAddr, recipient, amount, "eth")
	require.NoError(t, err)

	signer := ethTypes.NewEIP2930Signer(big.NewInt(31337))
	txHash := signer.Hash(txResponse.EthTx)

	privKeyHex := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	privKey, err := crypto.HexToECDSA(privKeyHex)
	require.NoError(t, err)

	sig, err := crypto.Sign(txHash.Bytes(), privKey)
	require.NoError(t, err)

	signedTx, err := txResponse.EthTx.WithSignature(signer, sig)
	require.NoError(t, err)

	ethClient := initETHClient(t, "http://localhost:8545")
	err = ethClient.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
}

func initETHClient(t *testing.T, rawURL string) *ethclient.Client {
	ethClient, err := ethclient.Dial(rawURL)
	require.NoError(t, err)
	return ethClient
}
