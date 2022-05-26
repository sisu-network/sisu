package sisu

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/tssclients"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/sisu-network/sisu/x/sisu/world"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	testKeyStore          = sdk.NewKVStoreKey("TestContext")
	testSisuEthAddr       = "0x743E1388AAd8EC7c47Df39AFbAEd58EBc1f43901"
	testSisuEthPubkeyHex  = "04b3cb1c95782b1793e3102d2ba493c34456f11ce471ca7e1ec1a731275b72bb2ba93e45069dab6d2b84815baeb3824f39c344bb9cf03d62cca9504724a808cc42"
	testContractAddr      = "0x50cc7ceDe8532d5f431EfC3e3EF167423Bc1807a"
	testErc20TokenAddress = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
)

func testContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(testKeyStore, sdk.StoreTypeIAVL, db)
	cms.LoadVersion(0)
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, tlog.NewNopLogger())
	return ctx
}

// Default keeper
func keeperTestGenesis(ctx sdk.Context) keeper.Keeper {
	keeper := keeper.NewKeeper(testKeyStore)
	keeper.SaveChain(ctx, &types.Chain{
		Id:       "ganache1",
		GasPrice: int64(5_000_000_000),
	})
	keeper.SaveChain(ctx, &types.Chain{
		Id:       "ganache2",
		GasPrice: int64(10_000_000_000),
	})
	liquidities := map[string]*types.Liquidity{
		"ganache1": {
			Id:      "ganache1",
			Address: "0xf0D676183dD5ae6b370adDdbE770235F23546f9d",
		},
		"ganache2": {
			Id:      "ganache2",
			Address: "0xf0D676183dD5ae6b370adDdbE770235F23546f9d",
		},
	}
	keeper.SetLiquidities(ctx, liquidities)
	keeper.SetTokens(ctx, map[string]*types.Token{
		"NATIVE_GANACHE1": {
			Id:       "NATIVE_GANACHE1",
			Price:    2000000000,
			Decimals: 18,
		},
		"NATIVE_GANACHE2": {
			Id:       "NATIVE_GANACHE1",
			Price:    2000000000,
			Decimals: 18,
		},
		"SISU": {
			Id:        "SISU",
			Price:     4000000000,
			Decimals:  18,
			Chains:    []string{"ganache1", "ganache2"},
			Addresses: []string{testErc20TokenAddress, testErc20TokenAddress},
		},
	})
	keeper.SaveParams(ctx, &types.Params{
		MajorityThreshold: 1,
	})
	return keeper
}

// Keeper after keygen has been saved
func keeperTestAfterKeygen(ctx sdk.Context) keeper.Keeper {
	ethTx := defaultTestEthTx(0)
	keeper := keeperTestGenesis(ctx)

	keeper.SaveKeygen(ctx, &types.Keygen{
		KeyType:     libchain.KEY_TYPE_ECDSA,
		Address:     ethTx.To().String(),
		PubKeyBytes: defaultTestEthPubkeyBytes(),
	})
	keeper.SaveContract(ctx, &types.Contract{
		Chain: "ganache1",
		Name:  "erc20gateway",
		Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
	}, false)
	keeper.SaveContract(ctx, &types.Contract{
		Chain: "ganache2",
		Name:  "erc20gateway",
		Hash:  SupportedContracts[ContractErc20Gateway].AbiHash,
	}, false)

	return keeper
}

func keeperTestAfterContractDeployed(ctx sdk.Context) keeper.Keeper {
	keeper := keeperTestAfterKeygen(ctx)

	keeper.SaveContract(ctx, &types.Contract{
		Chain:   "ganache1",
		Name:    "erc20gateway",
		Address: testContractAddr,
		Hash:    SupportedContracts[ContractErc20Gateway].AbiHash,
	}, false)
	keeper.SaveContract(ctx, &types.Contract{
		Chain:   "ganache2",
		Name:    "erc20gateway",
		Address: testContractAddr,
		Hash:    SupportedContracts[ContractErc20Gateway].AbiHash,
	}, false)

	return keeper
}

func defaultTestEthPubkeyBytes() []byte {
	bz, err := hex.DecodeString(testSisuEthPubkeyHex)
	if err != nil {
		panic(err)
	}

	return bz
}

func defaultTestEthPubkey() *ecdsa.PublicKey {
	bz := defaultTestEthPubkeyBytes()

	pubkey, err := crypto.UnmarshalPubkey(bz)
	if err != nil {
		panic(err)
	}
	return pubkey
}

func defaultTestEthTx(nonce uint64) *ethTypes.Transaction {
	amount := big.NewInt(100)
	gasLimit := uint64(100)
	gasPrice := big.NewInt(100)

	return ethTypes.NewTransaction(nonce,
		ecommon.HexToAddress("0x743E1388AAd8EC7c47Df39AFbAEd58EBc1f43901"), amount, gasLimit, gasPrice, nil)
}

func defaultWorldStateTest(ctx sdk.Context, keeper keeper.Keeper, deyesClients tssclients.DeyesClient) world.WorldState {
	ws := world.NewWorldState(keeper, deyesClients)
	ws.InitData(ctx)

	return ws
}
