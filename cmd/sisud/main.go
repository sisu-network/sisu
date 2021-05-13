package main

import (
	"os"

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
}

func main() {
	loadConfig()

	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.SisuHome); err != nil {
		os.Exit(1)
	}
}
