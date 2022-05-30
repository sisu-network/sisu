package sisu

import (
	"fmt"
	"math/big"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"

	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/helper"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
)

func mockTxOutputProducer(ctx sdk.Context, k keeper.Keeper, worldState world.WorldState) DefaultTxOutputProducer {
	appKeys := common.NewMockAppKeys()

	txOutputProducer := DefaultTxOutputProducer{
		worldState: worldState,
		tssConfig: config.TssConfig{
			DeyesUrl: "https://0.0.0.0:1234",
			SupportedChains: map[string]config.TssChainConfig{
				"ganache1": {
					Id: "ganache1",
				},
			},
		},
		keeper:  k,
		appKeys: appKeys,
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

	t.Run("call_transfer_in_successfully", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.NoError(t, err)

		token := worldState.GetTokenFromAddress(destChain, tokenAddr.String())
		gasPriceInToken, err := worldState.GetGasCostInToken(token.Id, destChain)
		require.NoError(t, err)

		txIn, err := parseTransferInData(txResponse.EthTx, worldState)
		require.NoError(t, err)
		require.Equal(t, txIn.amount, amount.Sub(amount, big.NewInt(gasPriceInToken)))
	})

	t.Run("insufficient_fund", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := big.NewInt(10000000000)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.EqualError(t, err, fmt.Sprint(world.ErrInsufficientFund))
		require.Nil(t, txResponse)
	})

	t.Run("token_has_price_0", func(t *testing.T) {
		ctx := testContext()
		keeper := keeperTestAfterContractDeployed(ctx)
		keeper.SetTokens(ctx, map[string]*types.Token{
			"SISU": {
				Id:        "SISU",
				Price:     0,
				Decimals:  18,
				Chains:    []string{"ganache1", "ganache2"},
				Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
			},
		})
		deyesClient := &tssclients.MockDeyesClient{}
		worldState := defaultWorldStateTest(ctx, keeper, deyesClient)

		txOutputProducer := mockTxOutputProducer(ctx, keeper, worldState)

		destChain := "ganache2"
		tokenAddr := ecommon.HexToAddress(testErc20TokenAddress)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)

		ethTx := mockEthTx(t, txOutputProducer, destChain, tokenAddr, amount)

		token := worldState.GetTokenFromAddress(destChain, tokenAddr.String())
		txResponse, err := txOutputProducer.processERC20TransferOut(ctx, ethTx)
		require.EqualError(t, err, fmt.Sprintf("token %s has price 0", token.Id))
		require.Nil(t, txResponse)
	})
}
