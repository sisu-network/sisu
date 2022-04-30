package main

import (
	"os"
	"path/filepath"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/joho/godotenv"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/cmd/sisud/cmd"
)

func loadConfig() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Set main app home dir
	app.SisuHome = os.Getenv("SISU_HOME")
	if app.SisuHome == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		app.SisuHome = filepath.Join(userHomeDir, "."+app.Name)
	}

	app.MainAppHome = filepath.Join(app.SisuHome, "/main")
}

func main() {
	loadConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.MainAppHome); err != nil {
		os.Exit(1)
	}
}
