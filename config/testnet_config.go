package config

import (
	"math/big"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/sisu-network/dcore/core"
	"github.com/sisu-network/dcore/eth/ethconfig"
	"github.com/sisu-network/dcore/miner"
	"github.com/sisu-network/dcore/node"
	"github.com/sisu-network/dcore/params"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/utils"
)

var (
	// testnetEthChainId = big.NewInt(34567)
	testnetEthChainId = libchain.GetChainIntFromId("eth-sisu-testnet")
)

// 10 accounts from the mnemmonic
// '0xF0cB0b1eAe36dBFE049ef4eDBeB5FDa9f1e05E50',
// '0x53FE06ffc5E89588d7CBE7415aB678ac683e6E98',
// '0xC1BB7762b9c1095A84528DD13A5771b34C5F06B1',
// '0xfeFdEa69eCCbA627ed014fB6C9eaE0ed50eAfDd3',
// '0xe59e1792966Bfe6AF001760DE15Ace60DD31848D',
// '0x2A98Ac0C5Dde849940a9790f7a52817e0CF5F924',
// '0xe3B1ba8De7f74De201Dad5978e3F207aA6B70aB4',
// '0xE11Dd1de3f74031EF82013C7C1F435CED2ae56a3',
// '0x4193cfF48b41076228621E25fb61868ec8563534',
// '0xd6c11d904b0F4eebf6102d3161ecD2c60170576E'

type TestnetConfig struct {
	sisuConfig *SisuConfig
	ethConfig  *ETHConfig
	tssConfig  *TssConfig
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

	addrs := []common.Address{
		common.HexToAddress("0xF0cB0b1eAe36dBFE049ef4eDBeB5FDa9f1e05E50"),
		common.HexToAddress("0x53FE06ffc5E89588d7CBE7415aB678ac683e6E98"),
	}
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
