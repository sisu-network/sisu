package config

import (
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/node"
)

const (
	APP_CONFIG  = "app_config"
	SISU_CONFIG = "sisu_config"
	ETH_CONFIG  = "eth_config"
	TSS_CONFIG  = "tss_config"
)

type Config interface {
	GetSisuConfig() *SisuConfig
	GetETHConfig() *ETHConfig
	GetTssConfig() *TssConfig
}

type SisuConfig struct {
	ConfigDir string

	Home            string
	SignerName      string
	EnableTss       bool
	KeyringBackend  string
	ChainId         string
	InternalApiHost string
	InternalApiPort uint16
}

type ETHConfig struct {
	Home          string
	Eth           *eth.Config
	Host          string
	Port          int
	UseInMemDb    bool
	DbPath        string
	Node          *node.Config
	ImportAccount bool
}

type TssChainConfig struct {
	Symbol   string `toml:"symbol"`
	Id       int    `toml:"id"`
	DeyesUrl string `toml:"deyes_url"`
}

// Example of supported chains in the toml config file.
// [supported_chains]
// [supported_chains.eth]
//   symbol = "eth"
// 	 id = 1
// 	 deyes_url = "http://localhost:31001"
type TssConfig struct {
	Enable          bool                      `toml:"enable"`
	Host            string                    `toml:"host"`
	Port            int                       `toml:"port"`
	SupportedChains map[string]TssChainConfig `toml:"supported_chains"`

	Dir string

	// Keygen
	PoolSizeLowerBound  int
	PoolSizeUpperBound  int
	BlockProposalLength int
}

func NewLocalConfig() Config {
	return &LocalConfig{}
}

func NewTestnetConfig() Config {
	return &TestnetConfig{}
}
