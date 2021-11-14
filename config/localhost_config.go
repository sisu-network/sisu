package config

import (
	"math/big"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/eth/ethconfig"
	"github.com/sisu-network/dcore/miner"
	"github.com/sisu-network/dcore/node"
	"github.com/sisu-network/dcore/params"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/utils"
)

var (
	basicTxGasLimit = 21000
	initialBalance  = new(big.Int).Mul(big.NewInt(10000), utils.ONE_ETHER_IN_WEI) // 10,000 ETHER
	localEthChainId = libchain.GetChainIntFromId("eth-sisu-local")
)

func localSisuConfig() *SisuConfig {
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	sisuConfig := &SisuConfig{
		Dir:            appDir + "/main",
		KeyringBackend: keyring.BackendTest,
		ChainId:        "eth-sisu-local",
		ApiHost:        "0.0.0.0",
		ApiPort:        25456,
	}

	return sisuConfig
}

func localETHConfig(baseDir string) *ETHConfig {
	home := baseDir + "/eth"

	return &ETHConfig{
		Eth:    getLocalEthConfig(),
		Host:   "0.0.0.0",
		Port:   1234,
		DbPath: home + "leveldb",
		Node:   getLocalEthNodeConfig(home),
		// ImportAccount: true,
		ImportAccount: false,
	}
}

// getLocalEthConfig returns ETH configuration used for localhost or testing.
func getLocalEthConfig() *ethconfig.Config {
	config := ethconfig.Defaults
	// config := ethconfig.NewDefaultConfig()
	chainConfig := &params.ChainConfig{
		ChainID:             localEthChainId,
		HomesteadBlock:      big.NewInt(0),
		DAOForkBlock:        big.NewInt(0),
		DAOForkSupport:      true,
		EIP150Block:         big.NewInt(0),
		EIP150Hash:          common.HexToHash("0x2086799aeebeae135c246c65021c82b4e15a2c451340993aacfd2751886514f0"),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
	}

	accounts := utils.GetLocalAccounts()
	alloc := make(map[common.Address]core.GenesisAccount)

	for _, account := range accounts {
		alloc[account.Address] = core.GenesisAccount{
			Balance: initialBalance,
		}
	}

	blockGasLimit := uint64(15000000)

	config.Genesis = &core.Genesis{
		Config:     chainConfig,
		Nonce:      0,
		Number:     0,
		ExtraData:  hexutil.MustDecode("0x00"),
		GasLimit:   blockGasLimit,
		Difficulty: big.NewInt(0),
		Alloc:      alloc,
	}

	config.Miner = miner.Config{
		BlockGasLimit: blockGasLimit,
	}

	config.TxPool = core.TxPoolConfig{
		PriceLimit: 50,
	}

	return &config
}

func getLocalEthNodeConfig(ethHome string) *node.Config {
	ksDir := ethHome + "/keystore"

	return &node.Config{
		KeyStoreDir:         ksDir,
		AllowUnprotectedTxs: true,
	}
}

func localTssConfig(baseDir string) *TssConfig {
	// 1. Check Tss home directory
	home := baseDir + "/tss"
	if !utils.IsFileExisted(home) {
		err := os.Mkdir(home, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// 2. Check toml file existence. Create a new one if needed.
	tomlFile := home + "/tss.toml"
	if !utils.IsFileExisted(tomlFile) {
		err := writeDefaultTss(tomlFile)
		if err != nil {
			panic(err)
		}
	}

	// 3. Read the toml config file
	config := &TssConfig{
		Dir: home,
	}
	// Read the config fiel from tss folder.
	_, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}

func overrideDevConfig(config *Config) {
	config.Eth.Eth = getLocalEthConfig()
	config.Eth.Node = getLocalEthNodeConfig(config.Eth.Dir)
}
