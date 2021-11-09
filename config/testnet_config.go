package config

import (
	"math/big"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sisu-network/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/eth/ethconfig"
	"github.com/sisu-network/dcore/miner"
	"github.com/sisu-network/dcore/node"
	"github.com/sisu-network/dcore/params"
	"github.com/sisu-network/sisu/utils"
)

var (
	testnetEthChainId = big.NewInt(34567)
)

type TestnetConfig struct {
	sisuConfig *SisuConfig
	ethConfig  *ETHConfig
	tssConfig  *TssConfig
}

func (c *TestnetConfig) GetSisuConfig() *SisuConfig {
	if c.sisuConfig == nil {
		c.sisuConfig = testnetSisuConfig()
	}

	return c.sisuConfig
}

func (c *TestnetConfig) GetETHConfig() *ETHConfig {
	// sisuConfig := c.GetSisuConfig()

	// if c.ethConfig == nil {
	// 	c.ethConfig = testnetETHConfig(sisuConfig.ConfigDir)
	// }

	return c.ethConfig
}

func (c *TestnetConfig) GetTssConfig() *TssConfig {
	// sisuConfig := c.GetSisuConfig()

	// if c.tssConfig == nil {
	// 	c.tssConfig = testnetTssConfig(sisuConfig.ConfigDir)
	// }

	return c.tssConfig
}

func testnetSisuConfig() *SisuConfig {
	appDir := os.Getenv("HOME") + "/.sisu"

	sisuConfig := &SisuConfig{
		Dir:            appDir + "/main",
		ChainId:        "talon-1",
		KeyringBackend: keyring.BackendFile,
		ApiHost:        "0.0.0.0",
		ApiPort:        25456,
	}

	return sisuConfig
}

func testnetETHConfig(baseDir string) *ETHConfig {
	home := baseDir + "/eth"

	return &ETHConfig{
		Eth:           testTestnetEthConfig(),
		Host:          "0.0.0.0",
		Port:          1234,
		DbPath:        home + "leveldb",
		Node:          getTestnetEthNodeConfig(home),
		ImportAccount: false,
	}
}

func testTestnetEthConfig() *ethconfig.Config {
	config := ethconfig.Defaults
	chainConfig := &params.ChainConfig{
		ChainID:             testnetEthChainId,
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

	blockGasLimit := uint64(15000000)
	alloc := make(map[common.Address]core.GenesisAccount)

	addrs := []common.Address{common.HexToAddress("0x018309Ce82ED587F568B3ae04549897d88066eE1")}
	for _, addr := range addrs {
		alloc[addr] = core.GenesisAccount{
			Balance: initialBalance,
		}
	}

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

func getTestnetEthNodeConfig(ethHome string) *node.Config {
	ksDir := ethHome + "/keystore"

	return &node.Config{
		KeyStoreDir:         ksDir,
		AllowUnprotectedTxs: false,
	}
}

func testnetTssConfig(baseDir string) *TssConfig {
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

func overrideTestnetConfig(config *Config) {
	config.Eth.Eth = testTestnetEthConfig()
	config.Eth.Node = &node.Config{
		KeyStoreDir:         filepath.Join(config.Eth.Dir, "keystore"),
		AllowUnprotectedTxs: false,
	}
}
