package config

import (
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/eth/ethconfig"
	"github.com/sisu-network/dcore/node"
	"github.com/sisu-network/dcore/params"
	"github.com/sisu-network/sisu/utils"
)

var (
	basicTxGasLimit = 21000
	initialBalance  = new(big.Int).Mul(big.NewInt(100), utils.ONE_ETHER_IN_WEI) // 100 ETHER
	chainID         = big.NewInt(1)
)

func LocalAppConfig() *AppConfig {
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	appConfig := &AppConfig{
		ConfigDir: appDir,
	}

	return appConfig
}

func LocalETHConfig(baseDir string) *ETHConfig {
	home := baseDir + "/eth"

	return &ETHConfig{
		Home:          home,
		Eth:           getLocalEthConfig(),
		Host:          "localhost",
		Port:          1234,
		UseInMemDb:    false,
		DbPath:        home + "leveldb",
		Node:          getLocalEthNodeConfig(home),
		ImportAccount: true,
	}
}

// getLocalEthConfig returns ETH configuration used for localhost or testing.
func getLocalEthConfig() *ethconfig.Config {
	config := ethconfig.NewDefaultConfig()
	chainConfig := &params.ChainConfig{
		ChainID:                     chainID,
		HomesteadBlock:              big.NewInt(0),
		DAOForkBlock:                big.NewInt(0),
		DAOForkSupport:              true,
		EIP150Block:                 big.NewInt(0),
		EIP150Hash:                  common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:                 big.NewInt(0),
		EIP158Block:                 big.NewInt(0),
		ByzantiumBlock:              big.NewInt(0),
		ConstantinopleBlock:         big.NewInt(0),
		PetersburgBlock:             big.NewInt(0),
		IstanbulBlock:               big.NewInt(0),
		ApricotPhase1BlockTimestamp: big.NewInt(0),
		ApricotPhase2BlockTimestamp: big.NewInt(0),
	}

	accounts := utils.GetLocalAccounts()
	alloc := make(map[common.Address]core.GenesisAccount)

	for _, account := range accounts {
		alloc[account.Address] = core.GenesisAccount{
			Balance: initialBalance,
		}
	}

	config.Genesis = &core.Genesis{
		Config:     chainConfig,
		Nonce:      0,
		Number:     0,
		ExtraData:  hexutil.MustDecode("0x00"),
		GasLimit:   1000000000000,
		Difficulty: big.NewInt(0),
		Alloc:      alloc,
	}

	config.TxPool = core.TxPoolConfig{
		PriceLimit: 50,
	}

	return &config
}

func getLocalEthNodeConfig(ethHome string) *node.Config {
	ksDir := ethHome + "/keystore"

	return &node.Config{
		KeyStoreDir: ksDir,
	}
}
