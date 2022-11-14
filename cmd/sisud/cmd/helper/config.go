package helper

import (
	"encoding/json"
	"os"
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
