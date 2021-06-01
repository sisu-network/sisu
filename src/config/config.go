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
	ConfigDir      string
	EnableTss      bool
	KeyringBackend string
	ChainId        string
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

type TssConfig struct {
	Enable bool
	Host   string
	Port   int
}

func NewLocalConfig() Config {
	return &LocalConfig{}
}

func NewTestnetConfig() Config {
	return &TestnetConfig{}
}
