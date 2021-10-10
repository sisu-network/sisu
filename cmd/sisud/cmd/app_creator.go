package cmd

import (
	"errors"
	"io"
	"path/filepath"

	"github.com/sisu-network/cosmos-sdk/baseapp"
	"github.com/sisu-network/cosmos-sdk/client/flags"
	"github.com/sisu-network/cosmos-sdk/server"
	servertypes "github.com/sisu-network/cosmos-sdk/server/types"
	"github.com/sisu-network/cosmos-sdk/snapshots"
	"github.com/sisu-network/cosmos-sdk/store"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/sisu/app"
	"github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/tendermint/libs/log"
	"github.com/spf13/cast"
	dbm "github.com/tendermint/tm-db"
)

type appCreator struct {
	encCfg params.EncodingConfig
	cfg    config.Config
}

type AppOptionWrapper struct {
	appOpts servertypes.AppOptions
	cfg     config.Config
}

func (wrapper *AppOptionWrapper) Get(key string) interface{} {
	switch key {
	case config.APP_CONFIG:
		return wrapper.cfg
	case config.SISU_CONFIG:
		return wrapper.cfg.GetSisuConfig()
	case config.ETH_CONFIG:
		return wrapper.cfg.GetETHConfig()
	case config.TSS_CONFIG:
		return wrapper.cfg.GetTssConfig()
	default:
		return wrapper.appOpts.Get(key)
	}
}

// newApp is an AppCreator
func (a appCreator) newApp(logger log.Logger, db dbm.DB, traceStore io.Writer, appOpts servertypes.AppOptions) servertypes.Application {
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

	return app.New(
		logger, db, traceStore, true, skipUpgradeHeights,
		cast.ToString(appOpts.Get(flags.FlagHome)),
		cast.ToUint(appOpts.Get(server.FlagInvCheckPeriod)),
		a.encCfg,
		// this line is used by starport scaffolding # stargate/root/appArgument
		a.getAppOptionsWrapper(appOpts),
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
		cfg:     a.cfg,
	}
}
