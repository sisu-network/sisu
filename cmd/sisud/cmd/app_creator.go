package cmd

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/server"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/snapshots"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"
	"github.com/spf13/cast"
	tlog "github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
)

type appCreator struct {
	encCfg           params.EncodingConfig
	appConfig        config.Config
	tendermintLogger tlog.Logger
}

type AppOptionWrapper struct {
	appOpts servertypes.AppOptions
}

// newApp is an AppCreator
func (a appCreator) newApp(tLogger tlog.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
	var cache sdk.MultiStorePersistentCache

	if cast.ToBool(appOpts.Get(server.FlagInterBlockCache)) {
		cache = store.NewCommitKVStoreCacheManager()
	}

	skipUpgradeHeights := make(map[int64]bool)
	for _, h := range cast.ToIntSlice(appOpts.Get(server.FlagUnsafeSkipUpgrades)) {
		skipUpgradeHeights[int64(h)] = true
	}

	pruningOpts, err := server.GetPruningOptionsFromFlags(appOpts)
	if err != nil {
		panic(err)
	}

	snapshotDir := filepath.Join(cast.ToString(appOpts.Get(flags.FlagHome)), "data", "snapshots")
	snapshotDB, err := sdk.NewLevelDB("metadata", snapshotDir)
	if err != nil {
		panic(err)
	}
	snapshotStore, err := snapshots.NewStore(snapshotDB, snapshotDir)
	if err != nil {
		panic(err)
	}

	return app.New(a.appConfig, tLogger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		// this line is used by starport scaffolding # stargate/root/appArgument
		appOpts,
		baseapp.SetPruning(pruningOpts),
		baseapp.SetMinGasPrices(cast.ToString(appOpts.Get(server.FlagMinGasPrices))),
		baseapp.SetMinRetainBlocks(cast.ToUint64(appOpts.Get(server.FlagMinRetainBlocks))),
		baseapp.SetHaltHeight(cast.ToUint64(appOpts.Get(server.FlagHaltHeight))),
		baseapp.SetHaltTime(cast.ToUint64(appOpts.Get(server.FlagHaltTime))),
		baseapp.SetInterBlockCache(cache),
		baseapp.SetTrace(cast.ToBool(appOpts.Get(server.FlagTrace))),
		baseapp.SetIndexEvents(cast.ToStringSlice(appOpts.Get(server.FlagIndexEvents))),
		baseapp.SetSnapshotStore(snapshotStore),
		baseapp.SetSnapshotInterval(cast.ToUint64(appOpts.Get(server.FlagStateSyncSnapshotInterval))),
		baseapp.SetSnapshotKeepRecent(cast.ToUint32(appOpts.Get(server.FlagStateSyncSnapshotKeepRecent))),
	)
}

// appExport creates a new simapp (optionally at a given height)
func (a appCreator) appExport(
	logger log.Logger, db dbm.DB, traceStore io.Writer, height int64, forZeroHeight bool, jailAllowedAddrs []string,
	appOpts servertypes.AppOptions) (servertypes.ExportedApp, error) {

	return servertypes.ExportedApp{}, errors.New("not supported")

	// var anApp *app.App

	// homePath, ok := appOpts.Get(flags.FlagHome).(string)
	// if !ok || homePath == "" {
	// 	return servertypes.ExportedApp{}, errors.New("application home not set")
	// }

	// if height != -1 {
	// 	anApp = app.New(
	// 		logger,
	// 		db,
	// 		traceStore,
	// 		false,
	// 		map[int64]bool{},
	// 		homePath,
	// 		uint(1),
	// 		a.encCfg,
	// 		// this line is used by starport scaffolding # stargate/root/exportArgument
	// 		a.getAppOptionsWrapper(appOpts),
	// 	)

	// 	if err := anApp.LoadHeight(height); err != nil {
	// 		return servertypes.ExportedApp{}, err
	// 	}
	// } else {
	// 	anApp = app.New(
	// 		logger,
	// 		db,
	// 		traceStore,
	// 		true,
	// 		map[int64]bool{},
	// 		homePath,
	// 		uint(1),
	// 		a.encCfg,
	// 		// this line is used by starport scaffolding # stargate/root/noHeightExportArgument
	// 		a.getAppOptionsWrapper(appOpts),
	// 	)
	// }

	// return anApp.ExportAppStateAndValidators(forZeroHeight, jailAllowedAddrs)
}

func (a appCreator) getAppOptionsWrapper(appOpts servertypes.AppOptions) *AppOptionWrapper {
	return &AppOptionWrapper{
		appOpts: appOpts,
	}
}
