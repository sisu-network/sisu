package main

import (
	"os"
	"path/filepath"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/cmd/sisud/cmd"

	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	app.SisuHome = os.Getenv("SISU_HOME")

	if app.SisuHome == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		app.SisuHome = filepath.Join(userHomeDir, "."+app.Name, "main")
	}

	app.KeyringBackend = os.Getenv("KEYRING_BACKEND")
}

func main() {
	loadConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.SisuHome); err != nil {
		os.Exit(1)
	}
}
