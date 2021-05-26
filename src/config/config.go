package config

import (
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/node"
)

const (
	SISU_CONFIG = "sisu_config"
	ETH_CONFIG  = "eth_config"
	TSS_CONFIG  = "tss_config"
)

type AppConfig struct {
	ConfigDir string
	EnableTss bool
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
	Host string
	Port int
}
