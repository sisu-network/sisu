package cmd

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/config"

	"github.com/sisu-network/sisu/app"
	"github.com/spf13/cobra"
)

func resetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Resets sisu's database.",
		Long: `This command resets Sisu app's database (including the ethereum's data).
Next time when user runs the app, it will sync with network from scratch. The sync
requires more than 1 validator nodes in the network and will not work for local testing.
This command only deletes data files but not config files.
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := deleteMainAppData(); err != nil {
				return err
			}

			if err := deleteEthData(); err != nil {
				return err
			}

			if err := deleteTssData(); err != nil {
				return err
			}

			if err := resetValidatorState(); err != nil {
				return err
			}

			if err := deleteSql(); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func deleteMainAppData() error {
	dataDir := app.MainAppHome + "/data"
	configDir := app.MainAppHome + "/config"

	toDelete := []string{
		// config folders.
		configDir + "/addrbook.json",

		// Data folder
		dataDir + "/application.db",
		dataDir + "/blockstore.db",
		dataDir + "/cs.wal",
		dataDir + "/evidence.db",
		dataDir + "/state.db",
		dataDir + "/tx_index.db",
		dataDir + "/snapshots",
	}

	filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "write-file-atomic") {
			toDelete = append(toDelete, configDir+"/"+info.Name())
		}

		return nil
	})

	return deleteFiles(toDelete)
}

func deleteEthData() error {
	toDelete := []string{
		app.SisuHome + "/eth",
		app.SisuHome + "/ethleveldb",
	}

	return deleteFiles(toDelete)
}

func deleteTssData() error {
	toDelete := []string{
		app.SisuHome + "/tss/processor.db",
	}

	return deleteFiles(toDelete)
}

func deleteFiles(toDelete []string) error {
	for _, file := range toDelete {
		log.Info("Deleting file/directory", file)
		if err := os.RemoveAll(file); err != nil {
			log.Error("Cannot delete", file, ". Error =", err)
			return err
		}
	}

	return nil
}

func resetValidatorState() error {
	content := `{
		"height": "0",
		"round": 0,
		"step": 0
	}`

	path := app.MainAppHome + "/data/priv_validator_state.json"
	return ioutil.WriteFile(path, []byte(content), 0644)
}

func deleteSql() error {
	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	cfg := config.Config{}
	cfg.Sisu.Dir = appDir + "/main"
	configFile := cfg.Sisu.Dir + "/config/sisu.toml"
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return err
	}

	sqlConfig := cfg.Sisu.Sql

	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", sqlConfig.Username, sqlConfig.Password, sqlConfig.Host, sqlConfig.Port, sqlConfig.Schema))
	if err != nil {
		return err
	}

	defer database.Close()

	log.Info("Deleting sql tables...")

	database.Exec("DROP TABLE contract")
	database.Exec("DROP TABLE tx_in")
	database.Exec("DROP TABLE tx_out")
	database.Exec("DROP TABLE schema_migrations")
	database.Exec("DROP TABLE keygen")
	database.Exec("DROP TABLE keygen_result")
	database.Exec("DROP TABLE mempool_tx")

	database.Exec("TRUNCATE TABLE deyes.watch_address")
	database.Exec("TRUNCATE TABLE deyes.latest_block_height")

	return nil
}
