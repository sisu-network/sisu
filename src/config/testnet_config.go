package config

import (
	"math/big"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
)

var (
	testnetEthChainId = big.NewInt(31416)
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
	sisuConfig := c.GetSisuConfig()

	if c.ethConfig == nil {
		c.ethConfig = testnetETHConfig(sisuConfig.ConfigDir)
	}

	return c.ethConfig
}

func (c *TestnetConfig) GetTssConfig() *TssConfig {
	sisuConfig := c.GetSisuConfig()

	if c.tssConfig == nil {
		c.tssConfig = testnetTssConfig(sisuConfig.ConfigDir)
	}

	return c.tssConfig
}

func testnetSisuConfig() *SisuConfig {
	appDir := os.Getenv("HOME") + "/.sisu"

	sisuConfig := &SisuConfig{
		ConfigDir:      appDir,
		ChainId:        "talon-1",
		KeyringBackend: keyring.BackendFile,
	}

	return sisuConfig
}

func testnetETHConfig(baseDir string) *ETHConfig {
	home := baseDir + "/eth"

	return &ETHConfig{
		Home:          home,
		Eth:           getLocalEthConfig(),
		Host:          "localhost",
		Port:          1234,
		UseInMemDb:    false,
		DbPath:        home + "leveldb",
		Node:          getLocalEthNodeConfig(home),
		ImportAccount: false,
	}
}

func testnetTssConfig(baseDir string) *TssConfig {
	return &TssConfig{
		// Enable: true,
		Host: "localhost",
		Port: 5678,
	}
}
