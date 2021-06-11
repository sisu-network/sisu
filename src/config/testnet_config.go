package config

import (
	"math/big"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/sisu/utils"
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
		SignerName:     "owner",
		ConfigDir:      appDir,
		Home:           appDir + "/main",
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
		Host:          "0.0.0.0",
		Port:          1234,
		UseInMemDb:    false,
		DbPath:        home + "leveldb",
		Node:          getLocalEthNodeConfig(home),
		ImportAccount: false,
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
	config := &TssConfig{}
	// Read the config fiel from tss folder.
	_, err := toml.DecodeFile(tomlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}
