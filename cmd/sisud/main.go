package main

import (
	"os"
	"path/filepath"

	"github.com/sisu-network/sisu/utils"

	svrcmd "github.com/sisu-network/cosmos-sdk/server/cmd"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/cmd/sisud/cmd"

	"github.com/joho/godotenv"
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
		app.MainAppHome = filepath.Join(app.SisuHome, "/main")
	}

	utils.LogInfo("Sisu home = ", app.SisuHome)
	utils.LogInfo("Main App home = ", app.MainAppHome)

	// Print IP address
	// _, err = server.ExternalIP()
	// if err != nil {
	// 	panic(err)
	// }
}

func main() {
	loadConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.MainAppHome); err != nil {
		os.Exit(1)
	}
}
