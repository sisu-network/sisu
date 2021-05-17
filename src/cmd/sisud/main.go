package main

import (
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/cmd/sisud/cmd"
	"github.com/sisu-network/sisu/utils"

	"github.com/joho/godotenv"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Set main app home dir
	app.MainAppHome = os.Getenv("SISU_HOME")
	if app.MainAppHome == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		app.MainAppHome = filepath.Join(userHomeDir, "."+app.Name, "main")
	}

	// Keyring
	app.KeyringBackend = os.Getenv("KEYRING_BACKEND")
	utils.LogInfo("keyRingBackend = ", app.KeyringBackend)

	// Print IP address
	ip, err := server.ExternalIP()
	if err != nil {
		panic(err)
	}
	utils.LogInfo("IP address = ", ip)
}

func main() {
	loadConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.MainAppHome); err != nil {
		os.Exit(1)
	}
}
