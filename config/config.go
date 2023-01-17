package config

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/lib/log"
)

type Config struct {
	Mode string

	Sisu   SisuConfig       `toml:"sisu"`
	Tss    TssConfig        `toml:"tss"`
	LogDNA log.LogDNAConfig `toml:"log_dna"`
	Solana SolanaConfig     `toml:"solana"`
	Lisk   LiskConfig       `toml:"lisk"`
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

type SolanaConfig struct {
	BridgeProgramId string `toml:"bridge_program_id" json:"bridge_program_id"`
	BridgePda       string `toml:"bridge_pda" json:"bridge_pda"`
}

type LiskConfig struct {
	RPC     string            `toml:"rpc" json:"rpc"`
	Network map[string]string `toml:"network" json:"network"`
}

// Example of supported chains in the toml config file.
// [supported_chains]
// [supported_chains.eth]
//
//	  symbol = "eth"
//		 id = 1
//		 deyes_url = "http://localhost:31001"
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
