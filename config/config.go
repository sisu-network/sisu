package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/echovl/cardano-go"
	"github.com/sisu-network/lib/log"
)

type Config struct {
	Mode string

	Sisu    SisuConfig       `toml:"sisu"`
	Tss     TssConfig        `toml:"tss"`
	LogDNA  log.LogDNAConfig `toml:"log_dna"`
	Cardano CardanoConfig    `toml:"cardano"`
}

type SisuConfig struct {
	Dir               string
	KeyringBackend    string           `toml:"keyring-backend"`
	KeyringPassphrase string           `toml:"keyring-passphrase"`
	ChainId           string           `toml:"chain-id"`
	ApiHost           string           `toml:"api-host"`
	ApiPort           uint16           `toml:"api-port"`
	EmailAlert        EmailAlertConfig `toml:"email-alert"`
}

type CardanoConfig struct {
	BlockfrostSecret string `toml:"block_frost_secret"`
	Network          int    `toml:"network"`
}

func (c *CardanoConfig) GetCardanoNetwork() cardano.Network {
	switch c.Network {
	case 0:
		return cardano.Testnet
	case 1:
		return cardano.Mainnet
	}

	panic(fmt.Errorf("Unkwown network %d", c.Network))
}

// Example of supported chains in the toml config file.
// [supported_chains]
// [supported_chains.eth]
//   symbol = "eth"
// 	 id = 1
// 	 deyes_url = "http://localhost:31001"
type TssConfig struct {
	DheartHost string `toml:"dheart-host"`
	DheartPort int    `toml:"dheart-port"`

	DeyesUrl string `toml:"deyes-url"`

	Dir string
}

// Optional Email Alert System. This could be useful when we want to send missing transaction alert
// to some email.
type EmailAlertConfig struct {
	Url    string `toml:"url" json:"url"`
	Secret string `toml:"secret" json:"secret"`
	Email  string `toml:"email" json:"email"`
}

func ReadConfig() (Config, error) {
	cfg := Config{}

	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	cfg.Sisu.Dir = appDir + "/main"
	cfg.Tss.Dir = appDir + "/tss"

	configFile := cfg.Sisu.Dir + "/config/sisu.toml"

	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		return cfg, os.ErrNotExist
	}

	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
