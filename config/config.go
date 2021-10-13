package config

import (
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
	Id       int    `toml:"id"`
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
	if config.Mode == "dev" {
		overrideDevConfig(config)
	}
}

func overrideDevConfig(config *Config) {
	config.Eth.Eth = getLocalEthConfig()
	config.Eth.Node = getLocalEthNodeConfig(config.Eth.Dir)
}
