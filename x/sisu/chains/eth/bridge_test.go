package eth

import (
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func mockForGetTransferIn(ctx sdk.Context, tokenPrice *big.Int, amountIn *big.Int) (keeper.Keeper,
	external.DeyesClient, *types.TransferDetails) {
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.SetMpcNonce(ctx, &types.MpcNonce{
		Chain: "ganache2",
		Nonce: 1,
	})

	token := &types.Token{
		Id:        "SISU",
		Price:     tokenPrice.String(),
		Chains:    []string{"ganache1", "ganache2"},
		Addresses: []string{testmock.TestErc20TokenAddress, testmock.TestErc20TokenAddress},
	}
	k.SetTokens(ctx, map[string]*types.Token{"SISU": token})

	deyesClient := &external.MockDeyesClient{
		GetTokenPriceFunc: func(id string) (*big.Int, error) {
			return tokenPrice, nil
		},
	}

	transfer := &types.TransferDetails{
		Token:       "SISU",
		Amount:      amountIn.String(),
		ToRecipient: "0x08BAB502c5e7125fD558B19a98D14907CF7f7E93",
	}

	return k, deyesClient, transfer
}

func TestBridge_GetTransferIn(t *testing.T) {
	gasCost := big.NewInt(80_000 * 10_000_000_000)

	t.Run("token_has_low_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		tokenPrice := new(big.Int).Mul(big.NewInt(10_000_000), utils.GweiToWei)
		amountIn := new(big.Int).Set(utils.EthToWei)
		keeper, deyesClient, transfer := mockForGetTransferIn(ctx, tokenPrice, amountIn)

		bridge := NewBridge("ganache2", "", keeper, deyesClient).(*bridge)
		_, amount, err := bridge.getTransferIn(
			ctx,
			keeper.GetAllTokens(ctx),
			transfer,
			gasCost,
		)
		require.Nil(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 0.01 ~ 0.16. Since 1 ETH = 10^18 wei, 0.16 ETH is
		// 160_000_000_000_000_000 wei. Amount ount = 1 - 0.16 = 0.84 ETH.
		require.Equal(t, new(big.Int).Sub(amountIn, big.NewInt(160_000_000_000_000_000)), amount)
	})

	t.Run("token_has_high_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		tokenPrice := new(big.Int).Mul(big.NewInt(100), utils.EthToWei)
		amountIn := new(big.Int).Set(utils.EthToWei)
		keeper, deyesClient, transfer := mockForGetTransferIn(ctx, tokenPrice, amountIn)

		bridge := NewBridge("ganache2", "", keeper, deyesClient).(*bridge)
		_, amount, err := bridge.getTransferIn(
			ctx,
			keeper.GetAllTokens(ctx),
			transfer,
			gasCost,
		)
		require.Nil(t, err)

		// gasPriceInToken = 0.00008 * 10 * 2 / 100 ~ 0.000016. Since 1 ETH = 10^18 wei, 0.000016 ETH is
		// 16_000_000_000_000 wei.
		require.Equal(t, new(big.Int).Sub(amountIn, big.NewInt(16_000_000_000_000)), amount)
	})

	t.Run("transfer_with_commission", func(t *testing.T) {
		ctx := testmock.TestContext()
		tokenPrice := utils.EtherToWei(big.NewInt(100))
		amountIn := new(big.Int).Set(utils.EthToWei)
		keeper, deyesClient, transfer := mockForGetTransferIn(ctx, tokenPrice, amountIn)

		// Set commission rate = 0.1%
		params := keeper.GetParams(ctx)
		params.CommissionRate = 10
		keeper.SaveParams(ctx, params)

		bridge := NewBridge("ganache2", "", keeper, deyesClient).(*bridge)
		_, amount, err := bridge.getTransferIn(
			ctx,
			keeper.GetAllTokens(ctx),
			transfer,
			gasCost,
		)
		require.Nil(t, err)

		amountAfterCommission := new(big.Int).Mul(amountIn, big.NewInt(999))
		amountAfterCommission = new(big.Int).Div(amountAfterCommission, big.NewInt(1000))

		// gasPriceInToken = 0.00008 * 10 * 2 / 100 ~ 0.000016. Since 1 ETH = 10^18 wei, 0.000016 ETH is
		// 16_000_000_000_000 wei.
		require.Equal(
			t,
			new(big.Int).Sub(amountAfterCommission, big.NewInt(16_000_000_000_000)),
			amount,
		)
	})

	t.Run("insufficient_fund", func(t *testing.T) {
		ctx := testmock.TestContext()
		amountIn := big.NewInt(10_000_000_000)
		keeper, deyesClient, transfer := mockForGetTransferIn(ctx, utils.EtherToWei(big.NewInt(8)),
			amountIn)

		bridge := NewBridge("ganache2", "", keeper, deyesClient).(*bridge)
		_, _, err := bridge.getTransferIn(
			ctx,
			keeper.GetAllTokens(ctx),
			transfer,
			gasCost,
		)
		require.NotNil(t, err)
	})

	t.Run("token_has_zero_price", func(t *testing.T) {
		ctx := testmock.TestContext()
		amountIn := new(big.Int).Set(utils.EthToWei)
		keeper, deyesClient, transfer := mockForGetTransferIn(ctx, utils.EtherToWei(big.NewInt(0)),
			amountIn)

		bridge := NewBridge("ganache2", "", keeper, deyesClient).(*bridge)
		_, _, err := bridge.getTransferIn(
			ctx,
			keeper.GetAllTokens(ctx),
			transfer,
			gasCost,
		)
		require.NotNil(t, err)
	})
}
