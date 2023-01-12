package lisk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/testmock"
	"github.com/sisu-network/sisu/x/sisu/types"
	"math/big"
	"testing"
)

func mockKeeperForBridge(ctx sdk.Context, tokenPrice *big.Int) keeper.Keeper {
	k := testmock.KeeperTestAfterContractDeployed(ctx)
	k.SetMpcNonce(ctx, &types.MpcNonce{
		Chain: "lisk-testnet",
		Nonce: 1,
	})

	token := &types.Token{
		Id:        "SISU",
		Price:     tokenPrice.String(),
		Chains:    []string{"lisk-testnet"},
		Addresses: []string{testmock.TestErc20TokenAddress, testmock.TestErc20TokenAddress},
	}
	k.SetTokens(ctx, map[string]*types.Token{"SISU": token})

	return k
}

func TestProcessTransfer(t *testing.T) {
	ctx := testmock.TestContext()

	amount := "1000000000"
	recipientAddress := "445f5ae1342837a1231f9d36d34a79145c1cd014"
	networks := map[string]string{"mainnet": "", "testnet": "e8832331820e5ba835012106a7c807b46c8b9c8672b6217b01373773fe87daf8"}
	keeper := mockKeeperForBridge(ctx, new(big.Int).Mul(big.NewInt(10_000_000), utils.GweiToWei))
	cfg := config.Config{}
	cfg.Lisk.RPC = "https://testnet-service.lisk.com/api/v2"
	cfg.Lisk.Network = networks

	bridge := NewBridge("lisk-testnet", "445f5ae1342837a1231f9d36d34a79145c1cd014", keeper, cfg)
	transfers := []*types.Transfer{{Amount: amount, Id: "id", ToRecipient: recipientAddress}}
	bridge.ProcessTransfers(ctx, transfers)
}
