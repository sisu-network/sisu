package config

import (
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/lib/log"
)

type Config struct {
	Mode string

	Sisu   SisuConfig       `toml:"sisu"`
	Tss    TssConfig        `toml:"tss"`
	LogDNA log.LogDNAConfig `toml:"log_dna"`
}

type SisuConfig struct {
	Dir            string
	KeyringBackend string `toml:"keyring-backend"`
	ChainId        string `toml:"chain-id"`
	ApiHost        string `toml:"api-host"`
	ApiPort        uint16 `toml:"api-port"`
}

type TssChainConfig struct {
	Id    string `toml:"id"`
	Token string `toml:"token"`
}

// Example of supported chains in the toml config file.
// [supported_chains]
// [supported_chains.eth]
//   symbol = "eth"
// 	 id = 1
// 	 deyes_url = "http://localhost:31001"
type TssConfig struct {
	DheartHost        string                    `toml:"dheart-host"`
	DheartPort        int                       `toml:"dheart-port"`
	SupportedChains   map[string]TssChainConfig `toml:"supported-chains"`
	MajorityThreshold int                       `toml:"majority-threshold"`

	DeyesUrl string `toml:"deyes-url"`

	Dir string
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
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
