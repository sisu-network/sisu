package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/node"
)

type Config struct {
	Mode string

	Sisu SisuConfig
	Eth  ETHConfig
	Tss  TssConfig
}

type SisuConfig struct {
	Dir            string
	KeyringBackend string `toml:"keyring-backend"`
	ChainId        string `toml:"chain-id"`
	ApiHost        string `toml:"api-host"`
	ApiPort        uint16 `toml:"api-port"`
	Sql            SqlConfig
}

type SqlConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"` // TODO: Move this sensitive data into separate place.
	Schema   string `toml:"schema"`
}

type ETHConfig struct {
	Dir           string
	Eth           *eth.Config
	Host          string `toml:"host"`
	Port          int    `toml:"port"`
	DbPath        string
	Node          *node.Config
	ImportAccount bool `toml:"import-account"`
}

type TssChainConfig struct {
	Symbol   string `toml:"symbol"`
	DeyesUrl string `toml:"deyes-url"`
}

// Example of supported chains in the toml config file.
// [supported_chains]
// [supported_chains.eth]
//   symbol = "eth"
// 	 id = 1
// 	 deyes_url = "http://localhost:31001"
type TssConfig struct {
	Enable          bool                      `toml:"enable"`
	DheartHost      string                    `toml:"dheart-host"`
	DheartPort      int                       `toml:"dheart-port"`
	SupportedChains map[string]TssChainConfig `toml:"supported-chains"`

	Dir string
}

// Overrides some config values since we don't want users to changes these values. They should be
// fixed and consistent throughout all the nodes.
func OverrideConfigValues(config *Config) {
	switch config.Mode {
	case "dev":
		overrideDevConfig(config)
	case "testnet":
		overrideTestnetConfig(config)
	default:
		panic(fmt.Errorf("unknown config mode %s", config.Mode))
	}
}

func ReadConfig() (Config, error) {
	cfg := Config{}

	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	cfg.Sisu.Dir = appDir + "/main"
	cfg.Eth.Dir = appDir + "/eth"
	cfg.Tss.Dir = appDir + "/tss"

	configFile := cfg.Sisu.Dir + "/config/sisu.toml"
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}

	OverrideConfigValues(&cfg)

	return cfg, nil
}
