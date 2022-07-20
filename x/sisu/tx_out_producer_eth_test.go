package sisu

import (
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func mockEthTx(t *testing.T, txOutputProducer DefaultTxOutputProducer, destChain string, tokenAddr ecommon.Address, amount *big.Int) *ethTypes.Transaction {
	contractAddress := ecommon.HexToAddress("0x08BAB502c5e7125fD558B19a98D14907CF7f7E93")
	gasLimit := uint64(100)
	gasPrice := big.NewInt(100)

	abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
	require.NoError(t, err)

	data, err := abi.Pack(MethodTransferOut, destChain, contractAddress.String(), tokenAddr, tokenAddr, amount)
	require.NoError(t, err)

	ethTx := ethTypes.NewTx(&ethTypes.LegacyTx{
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &contractAddress,
		Value:    big.NewInt(100),
		Data:     data,
	})

	return ethTx
}

func mockKeeperForTxOutProducerEth(ctx sdk.Context) keeper.Keeper {
	k := keeperTestAfterContractDeployed(ctx)
	k.AddGatewayCheckPoint(ctx, &types.GatewayCheckPoint{
		Chain: "ganache2",
		Nonce: 1,
	})

	return k
}

func TestTxOutProducerErc20_processERC20TransferOut(t *testing.T) {
	t.Parallel()

	t.Run("token_has_low_price", func(t *testing.T) {
		ctx := testContext()
		keeper := mockKeeperForTxOutProducerEth(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     new(big.Int).Mul(big.NewInt(10_000_000), utils.GweiToWei).String(), // 0.01
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})

		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
		}

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)
		data, err := parseEthTransferOut(ctx, keeper, ethTx, "ganache1")
		require.Nil(t, err)
		txResponse, err := txOutputProducer.buildERC20TransferIn(ctx, keeper, []*types.Token{data.Token},
			[]ecommon.Address{ecommon.HexToAddress(data.Recipient)}, []*big.Int{data.Amount}, data.DestChain)
		require.Nil(t, err)

		txIns, err := parseTransferInData(txResponse.EthTx)
		require.NoError(t, err)
		txIn := txIns[0]

		// gasPriceInToken = 0.00008 * 10 * 2 / 0.01 ~ 0.16. Since 1 ETH = 10^18 wei, 0.16 ETH is 160_000_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(160_000_000_000_000_000)), txIn.amount)
	})

	t.Run("token_has_high_price", func(t *testing.T) {
		ctx := testContext()
		keeper := mockKeeperForTxOutProducerEth(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     utils.EtherToWei(big.NewInt(100)).String(),
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})

		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
		}

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)
		data, err := parseEthTransferOut(ctx, keeper, ethTx, "ganache1")
		require.Nil(t, err)
		txResponse, err := txOutputProducer.buildERC20TransferIn(ctx, keeper, []*types.Token{data.Token},
			[]ecommon.Address{ecommon.HexToAddress(data.Recipient)}, []*big.Int{data.Amount}, data.DestChain)
		require.Nil(t, err)

		txIns, err := parseTransferInData(txResponse.EthTx)
		require.NoError(t, err)
		txIn := txIns[0]

		// gasPriceInToken = 0.00008 * 10 * 2 / 100 ~ 0.000016. Since 1 ETH = 10^18 wei, 0.000016 ETH is 16_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(16_000_000_000_000)), txIn.amount)
	})

	t.Run("insufficient_fund", func(t *testing.T) {
		ctx := testContext()
		keeper := mockKeeperForTxOutProducerEth(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     utils.EtherToWei(big.NewInt(8)).String(),
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})

		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
		}

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := big.NewInt(10_000_000_000)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)
		data, err := parseEthTransferOut(ctx, keeper, ethTx, "ganache1")
		require.Nil(t, err)
		txResponse, err := txOutputProducer.buildERC20TransferIn(ctx, keeper, []*types.Token{data.Token},
			[]ecommon.Address{ecommon.HexToAddress(data.Recipient)}, []*big.Int{data.Amount}, data.DestChain)

		// gasPriceInToken = 0.00008 * 10 * 2 / 8 ~ 0.0002. Since 1 ETH = 10^18 wei, 0.0002 ETH is 200_000_000_000_000 wei.
		// gasPriceInToken > amountIn
		require.Error(t, err)
		require.Nil(t, txResponse)
	})

	t.Run("token_has_zero_price", func(t *testing.T) {
		ctx := testContext()
		keeper := mockKeeperForTxOutProducerEth(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     utils.ZeroBigInt.String(),
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})

		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
		}

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)
		data, err := parseEthTransferOut(ctx, keeper, ethTx, "ganache1")
		require.Nil(t, err)
		txResponse, err := txOutputProducer.buildERC20TransferIn(ctx, keeper, []*types.Token{data.Token},
			[]ecommon.Address{ecommon.HexToAddress(data.Recipient)}, []*big.Int{data.Amount}, data.DestChain)
		require.Error(t, err)
		require.Nil(t, txResponse)
	})

	t.Run("token_has_negative_price", func(t *testing.T) {
		ctx := testContext()
		keeper := mockKeeperForTxOutProducerEth(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     utils.EtherToWei(big.NewInt(-100)).String(),
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})

		txOutputProducer := DefaultTxOutputProducer{
			keeper: keeper,
		}

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)
		data, err := parseEthTransferOut(ctx, keeper, ethTx, "ganache1")
		require.Nil(t, err)
		txResponse, err := txOutputProducer.buildERC20TransferIn(ctx, keeper, []*types.Token{data.Token},
			[]ecommon.Address{ecommon.HexToAddress(data.Recipient)}, []*big.Int{data.Amount}, data.DestChain)
		require.Error(t, err)
		require.Nil(t, txResponse)
	})
}
