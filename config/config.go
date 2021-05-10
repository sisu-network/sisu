package config

import (
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/node"
)

type AppConfig struct {
	ConfigDir string
}

type ETHConfig struct {
	Home          string
	Eth           *eth.Config
	Port          int
	UseInMemDb    bool
	DbPath        string
	Node          *node.Config
	ImportAccount bool
}
