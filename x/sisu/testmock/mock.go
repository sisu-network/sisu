package testmock

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"strings"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ecommon "github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

var (
	TestKeyStore          = sdk.NewKVStoreKey("TestContext")
	testSisuEthAddr       = "0x743E1388AAd8EC7c47Df39AFbAEd58EBc1f43901"
	testSisuEthPubkeyHex  = "04b3cb1c95782b1793e3102d2ba493c34456f11ce471ca7e1ec1a731275b72bb2ba93e45069dab6d2b84815baeb3824f39c344bb9cf03d62cca9504724a808cc42"
	testContractAddr      = "0x50cc7ceDe8532d5f431EfC3e3EF167423Bc1807a"
	TestErc20TokenAddress = "0x3A84fBbeFD21D6a5ce79D54d348344EE11EBd45C"
)

func TestContext() sdk.Context {
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)
	cms.MountStoreWithDB(TestKeyStore, sdk.StoreTypeIAVL, db)
	cms.LoadVersion(0)
	ctx := sdk.NewContext(cms, tmproto.Header{}, false, tlog.NewNopLogger())
	return ctx
}

// Default keeper
func KeeperTestGenesis(ctx sdk.Context) keeper.Keeper {
	keeper := keeper.NewKeeper(TestKeyStore)
	keeper.SaveChain(ctx, &types.Chain{
		Id:          "ganache1",
		NativeToken: "NATIVE_GANACHE1",
		EthConfig: &types.ChainEthConfig{
			GasPrice: int64(5_000_000_000),
		},
	})
	keeper.SaveChain(ctx, &types.Chain{
		Id:          "ganache2",
		NativeToken: "NATIVE_GANACHE2",
		EthConfig: &types.ChainEthConfig{
			GasPrice: int64(10_000_000_000),
		},
	})

	vaults := []*types.Vault{
		{
			Id:      "ganache1_v0",
			Chain:   "ganache1",
			Address: "0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c",
		},
		{
			Id:      "ganache2_v0",
			Chain:   "ganache2",
			Address: "0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c",
		},
	}
	keeper.SetVaults(ctx, vaults)
	keeper.SetTokens(ctx, map[string]*types.Token{
		"NATIVE_GANACHE1": {
			Id:       "NATIVE_GANACHE1",
			Price:    new(big.Int).Mul(big.NewInt(2), utils.EthToWei).String(),
			Chains:   []string{"ganache1"},
			Decimals: []uint32{18},
		},
		"NATIVE_GANACHE2": {
			Id:       "NATIVE_GANACHE1",
			Price:    new(big.Int).Mul(big.NewInt(2), utils.EthToWei).String(),
			Chains:   []string{"ganache2"},
			Decimals: []uint32{18},
		},
		"SISU": {
			Id:       "SISU",
			Price:    new(big.Int).Mul(big.NewInt(4), utils.EthToWei).String(),
			Chains:   []string{"ganache1", "ganache2", "cardano-testnet", "solana-devnet"},
			Decimals: []uint32{18, 18, 6, 8},
			Addresses: []string{
				strings.ToLower(TestErc20TokenAddress),
				strings.ToLower(TestErc20TokenAddress),
				"ccf1a53e157a7277e717045578a6e9834405730be0b778fd0daab794:uSISU",
				"8a6Kn1uwFAuePztJSBkLjUvJiD6YWZ33JMuSaXErKPCX",
			},
		},
		"ADA": {
			Id:       "ADA",
			Price:    new(big.Int).Mul(big.NewInt(400_000_000), utils.GweiToWei).String(),
			Chains:   []string{"ganache1", "ganache2", "cardano-testnet", "solana-devnet"},
			Decimals: []uint32{18, 18, 6, 8},
			Addresses: []string{
				"0xf0D676183dD5ae6b370adDdbE770235F23546f9d",
				"0xf0D676183dD5ae6b370adDdbE770235F23546f9d",
				"ccf1a53e157a7277e717045578a6e9834405730be0b778fd0daab794:uADA",
				"BJ9ArHvbeUhVLChS2yksw8xqvoRpWYLtGkg7CVHNa31a",
			},
		},
	})
	keeper.SaveParams(ctx, &types.Params{
		MajorityThreshold: 1,
		SupportedChains:   []string{"ganache1", "ganache2"},
	})
	return keeper
}

// Keeper after keygen has been saved
func KeeperTestAfterKeygen(ctx sdk.Context) keeper.Keeper {
	ethTx := defaultTestEthTx(0)
	keeper := KeeperTestGenesis(ctx)

	keeper.SaveKeygen(ctx, &types.Keygen{
		KeyType:     libchain.KEY_TYPE_ECDSA,
		Address:     ethTx.To().String(),
		PubKeyBytes: defaultTestEthPubkeyBytes(),
	})

	return keeper
}

func KeeperTestAfterContractDeployed(ctx sdk.Context) keeper.Keeper {
	keeper := KeeperTestAfterKeygen(ctx)

	keeper.SetMpcAddress(ctx, "ganache1", testContractAddr)
	keeper.SetMpcAddress(ctx, "ganache2", testContractAddr)

	keeper.SetVaults(ctx, []*types.Vault{
		{
			Id:      "ganache1_v0",
			Chain:   "ganache1",
			Address: "0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c",
		},
		{
			Id:      "ganache2_v0",
			Chain:   "ganache2",
			Address: "0x3a84fbbefd21d6a5ce79d54d348344ee11ebd45c",
		},
	})

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
