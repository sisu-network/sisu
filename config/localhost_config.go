package config

import (
	"math/big"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/sisu-network/sisu/utils"
)

var (
	initialBalance = new(big.Int).Mul(big.NewInt(10000), big.NewInt(utils.OneEtherInWei)) // 10,000 ETHER
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
