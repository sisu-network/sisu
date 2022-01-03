package cmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sisu-network/lib/log"

	_ "github.com/go-sql-driver/mysql"
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
		dataDir + "/private.db",
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
