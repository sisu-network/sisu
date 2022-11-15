package helper

import "path/filepath"

func IsSolanaEnabled(genesisFolder string) bool {
	solanaConfig, err := ReadCmdSolanaConfig(filepath.Join(genesisFolder, "solana.json"))
	if err != nil {
		panic(err)
	}

	return solanaConfig.Enable
}

func IsCardanoEnabled(genesisFolder string) bool {
	solanaConfig := ReadCardanoConfig(filepath.Join(genesisFolder, "cardano.json"))
	return solanaConfig.Enable
}
