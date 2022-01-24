package gen

import (
	"encoding/json"
	"path/filepath"

	cfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/genutil/types"
)

// GenAppStateFromConfig gets the genesis app state from the config
func GenAppStateFromConfig(cdc codec.JSONMarshaler,
	config *cfg.Config, genDoc tmtypes.GenesisDoc,
	persistentPeers string,
) (appState json.RawMessage, err error) {
	if err != nil {
		return appState, err
	}

	config.P2P.PersistentPeers = persistentPeers
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	// create the app state
	appGenesisState, err := types.GenesisStateFromGenDoc(genDoc)
	if err != nil {
		return appState, err
	}

	appState, err = json.MarshalIndent(appGenesisState, "", "  ")
	if err != nil {
		return appState, err
	}

	genDoc.AppState = appState
	err = genutil.ExportGenesisFile(&genDoc, config.GenesisFile())

	return appState, err
}
