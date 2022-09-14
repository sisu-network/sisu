package helper

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
	dbm "github.com/tendermint/tm-db"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

var (
	testKeyStore = sdk.NewKVStoreKey("TestContext")
)

// TODO: Refactor test context out of the sisu package
func testmock.TestContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(testKeyStore, sdk.StoreTypeIAVL, db)
	cms.LoadVersion(0)
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, tlog.NewNopLogger())
	return ctx
}

func TestGasCostInToken(t *testing.T) {
	ctx := testmock.TestContext()
	k := keeper.NewKeeper(testKeyStore)

	chain := "ganache1"
	k.SaveChain(ctx, &types.Chain{
		Id:          chain,
		GasPrice:    10 * 1_000_000_000,
		NativeToken: "NATIVE_GANACHE1",
	})
	k.SetTokens(ctx, map[string]*types.Token{
		"NATIVE_GANACHE1": {
			Id:       "NATIVE_GANACHE1",
			Price:    new(big.Int).Mul(big.NewInt(2), utils.EthToWei).String(), // $2
			Decimals: 18,
		},
		"SISU": {
			Id:        "SISU",
			Price:     new(big.Int).Mul(big.NewInt(4), utils.EthToWei).String(), // $4
			Decimals:  18,
			Chains:    []string{"ganache1", "ganache2"},
			Addresses: []string{"", ""},
		},
	})

	gas := big.NewInt(8_000_000)
	amount, err := GetChainGasCostInToken(ctx, k, "SISU", chain, gas)

	require.Equal(t, nil, err)

	// amount = 0.008 * 10 * 2 / 4 ~ 0.04. Since 1 ETH = 10^18 wei, 0.04 ETH is 40_000_000_000_000_000 wei.
	require.Equal(t, big.NewInt(40_000_000_000_000_000), amount)
}
