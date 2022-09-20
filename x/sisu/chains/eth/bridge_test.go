package eth

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockKeeperForBridge(ctx sdk.Context, tokenPrice *big.Int) keeper.Keeper {
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.AddGatewayCheckPoint(ctx, &types.GatewayCheckPoint{
		Chain: "ganache2",
		Nonce: 1,
	})

	token := &types.Token{
		Id:        "SISU",
		Price:     tokenPrice.String(), // 0.01
		Chains:    []string{"ganache1", "ganache2"},
		Addresses: []string{testmock.TestErc20TokenAddress, testmock.TestErc20TokenAddress},
	}
	k.SetTokens(ctx, map[string]*types.Token{"SISU": token})

	return k
}

func TestBridge(t *testing.T) {
	recipients := []ethcommon.Address{ethcommon.HexToAddress("0x08BAB502c5e7125fD558B19a98D14907CF7f7E93")}

	t.Run("token_has_low_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		keeper := mockKeeperForBridge(ctx, new(big.Int).Mul(big.NewInt(10_000_000), utils.GweiToWei))

		bridge := NewBridge("ganache2", "", keeper).(*bridge)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		txResponse, err := bridge.buildERC20TransferIn(ctx,
			[]*types.Token{keeper.GetTokens(ctx, []string{"SISU"})["SISU"]},
			recipients,
			[]*big.Int{amount},
		)
		require.Nil(t, err)

		data, err := parseTransferIn(ctx, keeper, txResponse.EthTx)
		require.NoError(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 0.01 ~ 0.16. Since 1 ETH = 10^18 wei, 0.16 ETH is 160_000_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(160_000_000_000_000_000)), data["amount"])
	})

	t.Run("token_has_high_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		keeper := mockKeeperForBridge(ctx, utils.EtherToWei(big.NewInt(100)))

		bridge := NewBridge("ganache2", "", keeper).(*bridge)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		txResponse, err := bridge.buildERC20TransferIn(ctx,
			[]*types.Token{keeper.GetTokens(ctx, []string{"SISU"})["SISU"]},
			recipients,
			[]*big.Int{amount},
		)
		require.Nil(t, err)

		data, err := parseTransferIn(ctx, keeper, txResponse.EthTx)
		require.NoError(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 100 ~ 0.000016. Since 1 ETH = 10^18 wei, 0.000016 ETH is 16_000_000_000_000 wei.
		require.Equal(t, amount.Sub(amount, big.NewInt(16_000_000_000_000)), data["amount"])
	})

	t.Run("insufficient_fund", func(t *testing.T) {
		ctx := testmock.TestContext()
		keeper := mockKeeperForBridge(ctx, utils.EtherToWei(big.NewInt(8)))

		bridge := NewBridge("ganache2", "", keeper).(*bridge)
		amount := big.NewInt(10_000_000_000)
		txResponse, err := bridge.buildERC20TransferIn(ctx,
			[]*types.Token{keeper.GetTokens(ctx, []string{"SISU"})["SISU"]},
			recipients,
			[]*big.Int{amount},
		)

		require.Error(t, err)
		require.Nil(t, txResponse)
	})

	t.Run("token_has_zero_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		keeper := mockKeeperForBridge(ctx, big.NewInt(0))

		bridge := NewBridge("ganache2", "", keeper).(*bridge)
		amount := new(big.Int).Mul(big.NewInt(1), utils.EthToWei)
		txResponse, err := bridge.buildERC20TransferIn(ctx,
			[]*types.Token{keeper.GetTokens(ctx, []string{"SISU"})["SISU"]},
			recipients,
			[]*big.Int{amount},
		)
		require.Error(t, err)
		require.Nil(t, txResponse)
	})
}
