package helper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	econfig "github.com/sisu-network/deyes/config"
	"github.com/sisu-network/sisu/x/sisu/types"

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

type LiskConfig struct {
	Enable  bool   `toml:"enable" json:"enable"`
	Chain   string `toml:"chain" json:"chain"`
	RPC     string `toml:"chain" json:"rpc"`
	Network string `toml:"network" json:"network"`
}

func ReadLiskConfig(genesisFolder string) LiskConfig {
	cfg := LiskConfig{}

	dat, err := os.ReadFile(filepath.Join(genesisFolder, "lisk.json"))
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
	err := json.Unmarshal(file, &deyesChains)
	if err != nil {
		panic(err)
	}

	return deyesChains
}

func ReadVaults(genesisFolder string, chains []string) []string {
	dat, err := os.ReadFile(filepath.Join(genesisFolder, "vault.json"))
	if err != nil {
		panic(err)
	}

	vaults := make([]*types.Vault, 0)
	err = json.Unmarshal(dat, &vaults)
	if err != nil {
		panic(err)
	}

	ret := make([]string, 0)
	for _, chain := range chains {
		found := false
		for _, vault := range vaults {
			if vault.Chain == chain {
				ret = append(ret, vault.Address)
				found = true
				break
			}
		}

		if !found {
			panic(fmt.Errorf("Cannot find vault in chain %s", chain))
		}
	}

	return ret
}

func ReadToken(genesisFolder string) []*types.Token {
	bz, err := os.ReadFile(filepath.Join(genesisFolder, "tokens.json"))
	if err != nil {
		panic(err)
	}

	tokens := make([]*types.Token, 0)
	err = json.Unmarshal(bz, &tokens)
	if err != nil {
		panic(err)
	}

	return tokens
}
