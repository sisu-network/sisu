package config

import (
	"math/big"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
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
	initialBalance  = new(big.Int).Mul(big.NewInt(10000), utils.ONE_ETHER_IN_WEI) // 10,000 ETHER
	ethChainId      = big.NewInt(1)
)

type LocalConfig struct {
	sisuConfig *SisuConfig
	ethConfig  *ETHConfig
	tssConfig  *TssConfig
}

func (c *LocalConfig) GetSisuConfig() *SisuConfig {
	if c.sisuConfig == nil {
		c.sisuConfig = localSisuConfig()
	}
	return c.sisuConfig
}

func (c *LocalConfig) GetETHConfig() *ETHConfig {
	sisuConfig := c.GetSisuConfig()

	if c.ethConfig == nil {
		c.ethConfig = localETHConfig(sisuConfig.ConfigDir)
	}
	return c.ethConfig
}

func (c *LocalConfig) GetTssConfig() *TssConfig {
	sisuConfig := c.GetSisuConfig()

	if c.tssConfig == nil {
		c.tssConfig = localTssConfig(sisuConfig.ConfigDir)
	}
	return c.tssConfig
}

func localSisuConfig() *SisuConfig {
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	sisuConfig := &SisuConfig{
		ConfigDir:      appDir,
		KeyringBackend: keyring.BackendTest,
		ChainId:        "sisu-dev",
	}

	return sisuConfig
}

func localETHConfig(baseDir string) *ETHConfig {
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
	config := ethconfig.Defaults
	// config := ethconfig.NewDefaultConfig()
	chainConfig := &params.ChainConfig{
		ChainID:             ethChainId,
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
		KeyStoreDir:         ksDir,
		AllowUnprotectedTxs: true,
	}
}

func localTssConfig(baseDir string) *TssConfig {
	return &TssConfig{
		// Enable: true,
		Host: "localhost",
		Port: 5678,
	}
}
