package app

import (
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sisu-network/cosmos-sdk/client"
	"github.com/sisu-network/cosmos-sdk/codec/types"
	"github.com/sisu-network/sisu/db"
	"github.com/sisu-network/sisu/utils"
	"github.com/spf13/cast"

	abci "github.com/sisu-network/tendermint/abci/types"
	"github.com/sisu-network/tendermint/libs/log"
	tmos "github.com/sisu-network/tendermint/libs/os"
	"github.com/sisu-network/tendermint/node"
	dbm "github.com/tendermint/tm-db"

	ethRpc "github.com/ethereum/go-ethereum/rpc"
	"github.com/sisu-network/cosmos-sdk/baseapp"
	"github.com/sisu-network/cosmos-sdk/client/grpc/tmservice"
	"github.com/sisu-network/cosmos-sdk/client/rpc"
	"github.com/sisu-network/cosmos-sdk/codec"
	"github.com/sisu-network/cosmos-sdk/server/api"
	cConfig "github.com/sisu-network/cosmos-sdk/server/config"
	servertypes "github.com/sisu-network/cosmos-sdk/server/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	"github.com/sisu-network/cosmos-sdk/types/module"
	"github.com/sisu-network/cosmos-sdk/version"

	authrest "github.com/sisu-network/cosmos-sdk/x/auth/client/rest"
	authkeeper "github.com/sisu-network/cosmos-sdk/x/auth/keeper"
	authtx "github.com/sisu-network/cosmos-sdk/x/auth/tx"
	authtypes "github.com/sisu-network/cosmos-sdk/x/auth/types"
	"github.com/sisu-network/cosmos-sdk/x/auth/vesting"
	"github.com/sisu-network/cosmos-sdk/x/bank"
	bankkeeper "github.com/sisu-network/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/sisu-network/cosmos-sdk/x/bank/types"
	"github.com/sisu-network/cosmos-sdk/x/capability"
	capabilitykeeper "github.com/sisu-network/cosmos-sdk/x/capability/keeper"
	capabilitytypes "github.com/sisu-network/cosmos-sdk/x/capability/types"
	"github.com/sisu-network/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/sisu-network/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/sisu-network/cosmos-sdk/x/crisis/types"
	distr "github.com/sisu-network/cosmos-sdk/x/distribution"
	distrclient "github.com/sisu-network/cosmos-sdk/x/distribution/client"
	distrkeeper "github.com/sisu-network/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/sisu-network/cosmos-sdk/x/distribution/types"
	"github.com/sisu-network/cosmos-sdk/x/evidence"
	evidencekeeper "github.com/sisu-network/cosmos-sdk/x/evidence/keeper"
	evidencetypes "github.com/sisu-network/cosmos-sdk/x/evidence/types"
	"github.com/sisu-network/cosmos-sdk/x/genutil"
	genutiltypes "github.com/sisu-network/cosmos-sdk/x/genutil/types"
	"github.com/sisu-network/cosmos-sdk/x/gov"
	govclient "github.com/sisu-network/cosmos-sdk/x/gov/client"
	govkeeper "github.com/sisu-network/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/sisu-network/cosmos-sdk/x/gov/types"
	transfer "github.com/sisu-network/cosmos-sdk/x/ibc/applications/transfer"
	ibctransferkeeper "github.com/sisu-network/cosmos-sdk/x/ibc/applications/transfer/keeper"
	ibctransfertypes "github.com/sisu-network/cosmos-sdk/x/ibc/applications/transfer/types"
	ibc "github.com/sisu-network/cosmos-sdk/x/ibc/core"
	ibcclient "github.com/sisu-network/cosmos-sdk/x/ibc/core/02-client"
	porttypes "github.com/sisu-network/cosmos-sdk/x/ibc/core/05-port/types"
	ibchost "github.com/sisu-network/cosmos-sdk/x/ibc/core/24-host"
	ibckeeper "github.com/sisu-network/cosmos-sdk/x/ibc/core/keeper"
	"github.com/sisu-network/cosmos-sdk/x/mint"
	mintkeeper "github.com/sisu-network/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/sisu-network/cosmos-sdk/x/mint/types"
	"github.com/sisu-network/cosmos-sdk/x/params"
	paramsclient "github.com/sisu-network/cosmos-sdk/x/params/client"
	paramskeeper "github.com/sisu-network/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/sisu-network/cosmos-sdk/x/params/types"
	paramproposal "github.com/sisu-network/cosmos-sdk/x/params/types/proposal"
	"github.com/sisu-network/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/sisu-network/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/sisu-network/cosmos-sdk/x/slashing/types"
	"github.com/sisu-network/cosmos-sdk/x/staking"
	stakingkeeper "github.com/sisu-network/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/sisu-network/cosmos-sdk/x/staking/types"
	"github.com/sisu-network/cosmos-sdk/x/upgrade"
	upgradeclient "github.com/sisu-network/cosmos-sdk/x/upgrade/client"
	upgradekeeper "github.com/sisu-network/cosmos-sdk/x/upgrade/keeper"
	upgradetypes "github.com/sisu-network/cosmos-sdk/x/upgrade/types"
	appparams "github.com/sisu-network/sisu/app/params"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/server"
	"github.com/sisu-network/sisu/x/auth"
	"github.com/sisu-network/sisu/x/auth/ante"
	"github.com/sisu-network/sisu/x/evm"
	evmKeeper "github.com/sisu-network/sisu/x/evm/keeper"
	evmtypes "github.com/sisu-network/sisu/x/evm/types"
	"github.com/sisu-network/sisu/x/sisu"
	sisukeeper "github.com/sisu-network/sisu/x/sisu/keeper"
	sisutypes "github.com/sisu-network/sisu/x/sisu/types"
	tss "github.com/sisu-network/sisu/x/tss"
	tssKeeper "github.com/sisu-network/sisu/x/tss/keeper"
	tsstypes "github.com/sisu-network/sisu/x/tss/types"
	tmjson "github.com/sisu-network/tendermint/libs/json"
	tmproto "github.com/sisu-network/tendermint/proto/tendermint/types"

	"github.com/sisu-network/sisu/common"
	sisuAuth "github.com/sisu-network/sisu/x/auth"
	// this line is used by starport scaffolding # stargate/app/moduleImport
)

const Name = "sisu"

// this line is used by starport scaffolding # stargate/wasm/app/enabledProposals

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		distrclient.ProposalHandler,
		upgradeclient.ProposalHandler,
		upgradeclient.CancelProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

var (
	SisuHome string

	// MainAppHome default home directories for the application daemon
	MainAppHome string

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		genutil.AppModuleBasic{},
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()...),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		ibc.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		vesting.AppModuleBasic{},
		sisu.AppModuleBasic{},
		evm.AppModuleBasic{},
		tss.AppModuleBasic{},
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
	}
)

var (
	_ CosmosApp               = (*App)(nil)
	_ servertypes.Application = (*App)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type App struct {
	txSubmitter       *common.TxSubmitter
	appKeys           *common.AppKeys
	globalData        common.GlobalData
	internalApiServer server.Server
	tssProcessor      *tss.Processor

	///////////////////////////////////////////////////////////////

	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Marshaler
	interfaceRegistry types.InterfaceRegistry

	invCheckPeriod uint

	// keys to access the substores
	keys    map[string]*sdk.KVStoreKey
	tkeys   map[string]*sdk.TransientStoreKey
	memKeys map[string]*sdk.MemoryStoreKey

	// keepers
	AccountKeeper    authkeeper.AccountKeeper
	BankKeeper       bankkeeper.Keeper
	CapabilityKeeper *capabilitykeeper.Keeper
	StakingKeeper    stakingkeeper.Keeper
	SlashingKeeper   slashingkeeper.Keeper
	MintKeeper       mintkeeper.Keeper
	DistrKeeper      distrkeeper.Keeper
	GovKeeper        govkeeper.Keeper
	CrisisKeeper     crisiskeeper.Keeper
	UpgradeKeeper    upgradekeeper.Keeper
	ParamsKeeper     paramskeeper.Keeper
	IBCKeeper        *ibckeeper.Keeper // IBC Keeper must be a pointer in the app, so we can SetRouter on it correctly
	EvidenceKeeper   evidencekeeper.Keeper
	TransferKeeper   ibctransferkeeper.Keeper

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper

	db         db.Database
	sisuKeeper sisukeeper.Keeper
	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	evmKeeper evmKeeper.Keeper
	tssKeeper tssKeeper.Keeper

	// the module manager
	mm *module.Manager
}

// New returns a reference to an initialized Gaia.
// NewSimApp returns a reference to an initialized SimApp.
func New(
	logger log.Logger, tdb dbm.DB, traceStore io.Writer, loadLatest bool, skipUpgradeHeights map[int64]bool,
	homePath string, invCheckPeriod uint, encodingConfig appparams.EncodingConfig,
	// this line is used by starport scaffolding # stargate/app/newArgument
	appOpts servertypes.AppOptions, baseAppOptions ...func(*baseapp.BaseApp),
) *App {
	_ = sisuAuth.AppModule{}

	appCodec := encodingConfig.Marshaler
	cdc := encodingConfig.Amino
	interfaceRegistry := encodingConfig.InterfaceRegistry

	bApp := baseapp.NewBaseApp(Name, logger, tdb, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetAppVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)

	keys := sdk.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey,
		govtypes.StoreKey, paramstypes.StoreKey, ibchost.StoreKey, upgradetypes.StoreKey,
		evidencetypes.StoreKey, ibctransfertypes.StoreKey, capabilitytypes.StoreKey,
		sisutypes.StoreKey, tsstypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
	)
	tkeys := sdk.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := sdk.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	app := &App{
		////////////////

		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
	}

	app.setupDefaultKeepers(homePath, bApp, skipUpgradeHeights)
	////////////// Sisu related keeper //////////////

	cfg, err := app.ReadConfig()
	if err != nil {
		panic(err)
	}

	app.db = db.NewDatabase(cfg.Sisu.Sql)
	err = app.db.Init()
	if err != nil {
		panic(err)
	}

	app.sisuKeeper = *sisukeeper.NewKeeper(
		appCodec, keys[sisutypes.StoreKey], keys[sisutypes.MemStoreKey],
	)

	app.appKeys = common.NewAppKeys(cfg.Sisu)
	app.appKeys.Init()

	app.globalData = common.NewGlobalData(cfg)
	app.globalData.Init()

	app.txSubmitter = common.NewTxSubmitter(cfg, app.appKeys)
	go app.txSubmitter.Start()

	// EVM keeper
	app.evmKeeper = *evmKeeper.NewKeeper(appCodec, app.txSubmitter, &cfg.Eth)
	app.evmKeeper.Initialize()

	// TSS keeper
	tssConfig := cfg.Tss
	utils.LogInfo("tssConfig = ", tssConfig)

	app.tssKeeper = *tssKeeper.NewKeeper(keys[tsstypes.StoreKey])

	//////////////////////////////////////////////////////////////////////

	transferModule := transfer.NewAppModule(app.TransferKeeper)
	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	// this line is used by starport scaffolding # ibc/app/router
	app.IBCKeeper.SetRouter(ibcRouter)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	var skipGenesisInvariants = cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	modules := []module.AppModule{
		genutil.NewAppModule(
			app.AccountKeeper, app.StakingKeeper, app.BaseApp.DeliverTx,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.AccountKeeper, nil),
		vesting.NewAppModule(app.AccountKeeper, app.BankKeeper),
		bank.NewAppModule(appCodec, app.BankKeeper, app.AccountKeeper),
		capability.NewAppModule(appCodec, *app.CapabilityKeeper),
		crisis.NewAppModule(&app.CrisisKeeper, skipGenesisInvariants),
		gov.NewAppModule(appCodec, app.GovKeeper, app.AccountKeeper, app.BankKeeper),
		mint.NewAppModule(appCodec, app.MintKeeper, app.AccountKeeper),
		slashing.NewAppModule(appCodec, app.SlashingKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		distr.NewAppModule(appCodec, app.DistrKeeper, app.AccountKeeper, app.BankKeeper, app.StakingKeeper),
		staking.NewAppModule(appCodec, app.StakingKeeper, app.AccountKeeper, app.BankKeeper),
		upgrade.NewAppModule(app.UpgradeKeeper),
		evidence.NewAppModule(app.EvidenceKeeper),
		ibc.NewAppModule(app.IBCKeeper),
		params.NewAppModule(app.ParamsKeeper),
		transferModule,

		sisu.NewAppModule(appCodec, app.sisuKeeper),
		evm.NewAppModule(appCodec, app.evmKeeper),
	}

	tssProcessor := tss.NewProcessor(app.tssKeeper, tssConfig, app.appKeys, app.txSubmitter, app.globalData)
	if tssConfig.Enable {
		utils.LogInfo("TSS is enabled")
		tssProcessor.Init()
		modules = append(modules, tss.NewAppModule(appCodec, app.tssKeeper, app.appKeys, app.txSubmitter, tssProcessor, app.globalData))
	}
	app.tssProcessor = tssProcessor

	app.mm = module.NewManager(modules...)

	// Set module begin
	beginBlockers := []string{upgradetypes.ModuleName, minttypes.ModuleName, distrtypes.ModuleName, slashingtypes.ModuleName,
		evidencetypes.ModuleName, stakingtypes.ModuleName, ibchost.ModuleName,
		evmtypes.ModuleName}

	if tssConfig.Enable {
		beginBlockers = append(beginBlockers, tsstypes.ModuleName)
	}
	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(beginBlockers...)

	// Set module end blocker
	endBlockers := []string{evmtypes.ModuleName,
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName}
	if tssConfig.Enable {
		endBlockers = append(endBlockers, tsstypes.ModuleName)
	}
	app.mm.SetOrderEndBlockers(endBlockers...)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	initGenesisModules := []string{
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		ibchost.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		ibctransfertypes.ModuleName,
		sisutypes.ModuleName,
		evmtypes.ModuleName,
	}

	if tssConfig.Enable {
		initGenesisModules = append(initGenesisModules, tsstypes.ModuleName)
	}
	app.mm.SetOrderInitGenesis(initGenesisModules...)

	app.mm.RegisterInvariants(&app.CrisisKeeper)
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter(), encodingConfig.Amino)
	app.mm.RegisterServices(module.NewConfigurator(app.MsgServiceRouter(), app.GRPCQueryRouter()))

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetAnteHandler(
		ante.NewAnteHandler(
			tssConfig,
			app.AccountKeeper, app.BankKeeper, app.evmKeeper,
			ante.DefaultSigVerificationGasConsumer, encodingConfig.TxConfig.SignModeHandler(),
			app.evmKeeper.GetEthValidator(), tssProcessor,
		),
	)
	app.SetEndBlocker(app.EndBlocker)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}

		// Initialize and seal the capability keeper so all persistent capabilities
		// are loaded in-memory and prevent any further modules from creating scoped
		// sub-keepers.
		// This must be done during creation of baseapp rather than in InitChain so
		// that in-memory capabilities get regenerated on app restart.
		// Note that since this reads from the store, we can only perform it when
		// `loadLatest` is set to true.
		ctx := app.BaseApp.NewUncachedContext(true, tmproto.Header{})
		app.CapabilityKeeper.InitializeAndSeal(ctx)
	}

	app.setupApiServer(cfg)

	return app
}

func (app *App) setupDefaultKeepers(homePath string, bApp *baseapp.BaseApp, skipUpgradeHeights map[int64]bool) {
	appCodec := app.appCodec
	cdc := app.cdc
	keys := app.keys
	tkeys := app.tkeys
	memKeys := app.memKeys
	invCheckPeriod := app.invCheckPeriod

	app.ParamsKeeper = initParamsKeeper(appCodec, cdc, keys[paramstypes.StoreKey], tkeys[paramstypes.TStoreKey])
	// set the BaseApp's parameter store
	bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramskeeper.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.CapabilityKeeper = capabilitykeeper.NewKeeper(appCodec, keys[capabilitytypes.StoreKey], memKeys[capabilitytypes.MemStoreKey])

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibchost.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	// add keepers
	app.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec, keys[authtypes.StoreKey], app.GetSubspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, maccPerms,
	)
	app.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec, keys[banktypes.StoreKey], app.AccountKeeper, app.GetSubspace(banktypes.ModuleName), app.ModuleAccountAddrs(),
	)
	stakingKeeper := stakingkeeper.NewKeeper(
		appCodec, keys[stakingtypes.StoreKey], app.AccountKeeper, app.BankKeeper, app.GetSubspace(stakingtypes.ModuleName),
	)
	app.MintKeeper = mintkeeper.NewKeeper(
		appCodec, keys[minttypes.StoreKey], app.GetSubspace(minttypes.ModuleName), &stakingKeeper,
		app.AccountKeeper, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.DistrKeeper = distrkeeper.NewKeeper(
		appCodec, keys[distrtypes.StoreKey], app.GetSubspace(distrtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&stakingKeeper, authtypes.FeeCollectorName, app.ModuleAccountAddrs(),
	)
	app.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec, keys[slashingtypes.StoreKey], &stakingKeeper, app.GetSubspace(slashingtypes.ModuleName),
	)
	app.CrisisKeeper = crisiskeeper.NewKeeper(
		app.GetSubspace(crisistypes.ModuleName), invCheckPeriod, app.BankKeeper, authtypes.FeeCollectorName,
	)
	app.UpgradeKeeper = upgradekeeper.NewKeeper(skipUpgradeHeights, keys[upgradetypes.StoreKey], appCodec, homePath)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.StakingKeeper = *stakingKeeper.SetHooks(
		stakingtypes.NewMultiStakingHooks(app.DistrKeeper.Hooks(), app.SlashingKeeper.Hooks()),
	)

	// ... other modules keepers

	// Create IBC Keeper
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec, keys[ibchost.StoreKey], app.GetSubspace(ibchost.ModuleName), app.StakingKeeper, scopedIBCKeeper,
	)

	// register the proposal types
	govRouter := govtypes.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govtypes.ProposalHandler).
		AddRoute(paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.ParamsKeeper)).
		AddRoute(distrtypes.RouterKey, distr.NewCommunityPoolSpendProposalHandler(app.DistrKeeper)).
		AddRoute(upgradetypes.RouterKey, upgrade.NewSoftwareUpgradeProposalHandler(app.UpgradeKeeper)).
		AddRoute(ibchost.RouterKey, ibcclient.NewClientUpdateProposalHandler(app.IBCKeeper.ClientKeeper))

	app.GovKeeper = govkeeper.NewKeeper(
		appCodec, keys[govtypes.StoreKey], app.GetSubspace(govtypes.ModuleName), app.AccountKeeper, app.BankKeeper,
		&app.StakingKeeper, govRouter,
	)

	// Create Transfer Keepers
	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec, keys[ibctransfertypes.StoreKey], app.GetSubspace(ibctransfertypes.ModuleName),
		app.IBCKeeper.ChannelKeeper, &app.IBCKeeper.PortKeeper,
		app.AccountKeeper, app.BankKeeper, app.ScopedTransferKeeper,
	)
	// Create evidence Keeper for to register the IBC light client misbehaviour evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec, keys[evidencetypes.StoreKey], &app.StakingKeeper, app.SlashingKeeper,
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.EvidenceKeeper = *evidenceKeeper
}

func (app *App) ReadConfig() (config.Config, error) {
	cfg := config.Config{}

	appDir := os.Getenv("APP_DIR")
	if appDir == "" {
		appDir = os.Getenv("HOME") + "/.sisu"
	}

	cfg.Sisu.Dir = appDir + "/main"
	cfg.Eth.Dir = appDir + "/eth"
	cfg.Tss.Dir = appDir + "/tss"

	configFile := cfg.Sisu.Dir + "/config/sisu.toml"
	_, err := toml.DecodeFile(configFile, &cfg)
	if err != nil {
		return cfg, err
	}

	config.OverrideConfigValues(&cfg)

	return cfg, nil
}

// Name returns the name of the App
func (app *App) Name() string { return app.BaseApp.Name() }

// InitChainer application update at chain initialization
func (app *App) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState
	if err := tmjson.Unmarshal(req.AppStateBytes, &genesisState); err != nil {
		panic(err)
	}
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// LoadHeight loads a particular height
func (app *App) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *App) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns Gaia's app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *App) AppCodec() codec.Marshaler {
	return app.appCodec
}

// InterfaceRegistry returns Gaia's InterfaceRegistry
func (app *App) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetKey(storeKey string) *sdk.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetTKey(storeKey string) *sdk.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *App) GetMemKey(storeKey string) *sdk.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *App) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *App) RegisterAPIRoutes(apiSvr *api.Server, apiConfig cConfig.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	rpc.RegisterRoutes(clientCtx, apiSvr.Router)
	// Register legacy tx routes.
	authrest.RegisterTxRoutes(clientCtx, apiSvr.Router)
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new tendermint queries routes from grpc-gateway.
	tmservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register legacy and grpc-gateway routes for all modules.
	ModuleBasics.RegisterRESTRoutes(clientCtx, apiSvr.Router)
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *App) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *App) RegisterTendermintService(clientCtx client.Context) {
	tmservice.RegisterTendermintService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.interfaceRegistry)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryMarshaler, legacyAmino *codec.LegacyAmino, key, tkey sdk.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName).WithKeyTable(govtypes.ParamKeyTable())
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName)
	paramsKeeper.Subspace(ibchost.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	return paramsKeeper
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	app.txSubmitter.SyncBlockSequence(ctx, app.AccountKeeper)

	app.globalData.UpdateCatchingUp()
	app.globalData.UpdateValidatorSets()

	return app.mm.BeginBlock(ctx, req)
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

// This is internal server for Sisu to communicate with other components (e.g. TSS or watchers).
// These apis are not exposed to public.
// TODO: Make this separate from tssProcessor to avoid other components to invoke restricted functions
// of the processor.
func (app *App) setupApiServer(c config.Config) {
	handler := ethRpc.NewServer()
	handler.RegisterName("tss", tss.NewApi(app.tssProcessor, &app.tssKeeper))

	appConfig := c.Sisu
	s := server.NewServer(handler, appConfig.ApiHost, appConfig.ApiPort)

	utils.LogInfo("Starting Internal API server")
	go s.Run()
}

func (app *App) GetTendermintOptions() []node.Option {
	options := make([]node.Option, 0)
	options = append(options, func(n *node.Node) {
		n.SetPreAddTxFunc(app.tssProcessor.PreAddTxToMempoolFunc)
	})

	return options
}
