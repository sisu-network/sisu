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
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

func mockTxOutputProducer(ctx sdk.Context, keeper keeper.Keeper, worldState world.WorldState) DefaultTxOutputProducer {

	txOutputProducer := DefaultTxOutputProducer{
		worldState: worldState,
		keeper:     keeper,
	}

	return txOutputProducer
}

func mockEthTx(t *testing.T, txOutputProducer DefaultTxOutputProducer, destChain string, tokenAddr ecommon.Address, amount *big.Int) *ethTypes.Transaction {
	contractAddress := ecommon.HexToAddress("0x08BAB502c5e7125fD558B19a98D14907CF7f7E93")
	gasLimit := uint64(100)
	gasPrice := big.NewInt(100)

	abi, err := abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
	require.NoError(t, err)

	data, err := abi.Pack(MethodTransferOut, destChain, contractAddress, tokenAddr, tokenAddr, amount)
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

func TestTxOutProducerErc20_getGasCostInToken(t *testing.T) {
	ctx := testContext()
	k := keeperTestGenesis(ctx)
	deyesClient := &tssclients.MockDeyesClient{}

	worldState := defaultWorldStateTest(ctx, k, deyesClient)

	chain := "ganache1"
	token := &types.Token{
		Id:    "SISU",
		Price: int64(4 * utils.DecinmalUnit),
	}
	worldState.SetTokens(map[string]*types.Token{
		"SISU": token,
	})

	gas := big.NewInt(8_000_000)
	gasPrice := big.NewInt(10 * 1_000_000_000) // 10 gwei
	nativeTokenPrice, err := worldState.GetNativeTokenPriceForChain(chain)
	require.NoError(t, err)
	amount, err := helper.GetGasCostInToken(gas, gasPrice, big.NewInt(token.Price), big.NewInt(nativeTokenPrice))

	require.Equal(t, nil, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}

func TestTxOutProducerErc20_processERC20TransferOut(t *testing.T) {
	t.Parallel()

	t.Run("token_has_low_price", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)
		worldState.SetTokens(map[string]*types.Token{
			"SISU": {
				Id:    "SISU",
				Price: int64(0.01 * utils.DecinmalUnit),
			},
		})

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.NoError(t, err)

		txIn, err := parseTransferInData(txResponse.EthTx, worldState)
		require.NoError(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 0.01 ~ 0.16. Since 1 ETH = 10^18 wei, 0.16 ETH is 160_000_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(160_000_000_000_000_000)), txIn.amount)
	})

	t.Run("token_has_high_price", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)
		worldState.SetTokens(map[string]*types.Token{
			"SISU": {
				Id:    "SISU",
				Price: int64(100 * utils.DecinmalUnit),
			},
		})

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.NoError(t, err)

		txIn, err := parseTransferInData(txResponse.EthTx, worldState)
		require.NoError(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 100 ~ 0.000016. Since 1 ETH = 10^18 wei, 0.000016 ETH is 16_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(16_000_000_000_000)), txIn.amount)
	})

	t.Run("insufficient_fund", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)
		worldState.SetTokens(map[string]*types.Token{
			"SISU": {
				Id:    "SISU",
				Price: int64(8 * utils.DecinmalUnit),
			},
		})

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := big.NewInt(10_000_000_000)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)

		// gasPriceInToken = 0.00008 * 10 * 2 / 8 ~ 0.0002. Since 1 ETH = 10^18 wei, 0.0002 ETH is 200_000_000_000_000 wei.
		// gasPriceInToken > amountIn
		require.Error(t, err)
		require.Nil(t, txResponse)
	})

	t.Run("token_has_zero_price", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)
		worldState.SetTokens(map[string]*types.Token{
			"SISU": {
				Id:    "SISU",
				Price: 0,
			},
		})

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.Error(t, err)
		require.Nil(t, txResponse)
	})

	t.Run("token_has_negative_price", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)
		worldState.SetTokens(map[string]*types.Token{
			"SISU": {
				Id:    "SISU",
				Price: -1000,
			},
		})

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.Error(t, err)
		require.Nil(t, txResponse)
	})
}
