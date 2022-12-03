package helper

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	econfig "github.com/sisu-network/deyes/config"

	cardanogo "github.com/echovl/cardano-go"
)

type CmdSolanaConfig struct {
	Enable          bool   `toml:"enable" json:"enable"`
	Chain           string `toml:"chain" json:"chain"`
	Rpc             string `toml:"rpc" json:"rpc"`
	Ws              string `toml:"rpc" json:"ws"`
	BlockTime       int    `toml:"block_time" json:"block_time"`
	AdjustTime      int    `toml:"adjust_time" json:"adjust_time"`
	BridgeProgramId string `toml:"bridge_program_id" json:"bridge_program_id"`
	BridgePda       string `toml:"bridge_pda" json:"bridge_pda"`
}

func ReadCmdSolanaConfig(filePath string) (CmdSolanaConfig, error) {
	cfg := CmdSolanaConfig{}

	dat, err := os.ReadFile(filePath)
	if err != nil {
		return cfg, err
	}

	if err := json.Unmarshal(dat, &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

type CardanoConfig struct {
	Enable bool   `toml:"enable" json:"enable"`
	Secret string `toml:"secret" json:"secret"`
	Chain  string `toml:"chain" json:"chain"`
}

func (c *CardanoConfig) GetCardanoNetwork() cardanogo.Network {
	switch c.Chain {
	case "cardano-testnet":
		return cardanogo.Preprod
	}

	return cardanogo.Mainnet
}

func ReadCardanoConfig(genesisFolder string) CardanoConfig {
	cfg := CardanoConfig{}

	dat, err := os.ReadFile(filepath.Join(genesisFolder, "cardano.json"))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(dat, &cfg); err != nil {
		panic(err)
	}

	return cfg
}

func ReadDeyesChainConfigs(path string) []econfig.Chain {
	deyesChains := make([]econfig.Chain, 0)
	file, _ := ioutil.ReadFile(path)
	err := json.Unmarshal([]byte(file), &deyesChains)
	if err != nil {
		panic(err)
	}

	return deyesChains
}
