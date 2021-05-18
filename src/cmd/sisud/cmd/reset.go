package cmd

import (
	"io/ioutil"
	"os"

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
`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			var err error

			err = deleteMainAppFolder()
			if err != nil {
				return err
			}

			err = deleteEthFolder()
			if err != nil {
				return err
			}

			err = resetValidatorState()
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func deleteMainAppFolder() error {
	dataDir := app.MainAppHome + "/data"
	toDelete := []string{
		dataDir + "/application.db",
		dataDir + "/blockstore.db",
		dataDir + "/cs.wal",
		dataDir + "/evidence.db",
		dataDir + "/state.db",
		dataDir + "/tx_index.db",
	}

	for _, dir := range toDelete {
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
	}

	return nil
}

func deleteEthFolder() error {
	sisuHome := os.Getenv("SISU_HOME")
	toDelete := []string{
		sisuHome + "/eth",
		sisuHome + "/ethleveldb",
	}

	for _, dir := range toDelete {
		if err := os.RemoveAll(dir); err != nil {
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
