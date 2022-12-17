package cmd

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/sisu-network/lib/log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/sisu-network/sisu/cmd/sisud/cmd/flags"
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
			root, _ := cmd.Flags().GetString(flags.RootFolder)

			if err := deleteMainAppData(root); err != nil {
				return err
			}

			if err := resetValidatorState(root); err != nil {
				return err
			}

			if err := deleteSql(); err != nil {
				return err
			}

			return nil
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	cmd.Flags().String(flags.RootFolder, home+"/.sisu/main", "Relative path to the root folder of Sisu.")

	return cmd
}

func deleteMainAppData(root string) error {
	dataDir := path.Join(root, "data")
	configDir := path.Join(root, "config")

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
		dataDir + "/private",
	}

	filepath.Walk(configDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasPrefix(info.Name(), "write-file-atomic") {
			toDelete = append(toDelete, configDir+"/"+info.Name())
		}

		return nil
	})

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

func resetValidatorState(root string) error {
	content := `{
		"height": "0",
		"round": 0,
		"step": 0
	}`

	path := path.Join(root, "data/priv_validator_state.json")
	return os.WriteFile(path, []byte(content), 0644)
}

func deleteSql() error {
	username := "root"
	password := "password"
	host := "localhost"
	port := 3306
	schema := "sisu"

	database, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, host, port, schema))
	if err != nil {
		return err
	}

	defer database.Close()

	log.Info("Deleting sql tables...")

	database.Exec("TRUNCATE TABLE deyes.watch_address")
	database.Exec("TRUNCATE TABLE deyes.latest_block_height")

	return nil
}
