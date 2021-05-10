package app

import (
	"fmt"
	"os"
	"time"

	cServer "github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/ethdb"
	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/sisu-network/dcore/core/rawdb"
	dcore "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/dcore/extra"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/ethchain"
	"github.com/sisu-network/sisu/utils"
	abci "github.com/tendermint/tendermint/abci/types"
)

func (app *App) InitializeEth() error {
	app.Logger().Info("Init chain....")

	// Use local host config for dev.
	// TODO: Moev this to a separate config file.
	app.appConfigs = config.LocalAppConfig()
	app.ethConfig = config.LocalETHConfig(app.appConfigs.ConfigDir)

	// Setting log level
	ethLog.Root().SetHandler(ethLog.LvlFilterHandler(
		ethLog.LvlDebug, ethLog.StreamHandler(os.Stderr, ethLog.TerminalFormat(false))))

	app.createConfigDir()
	err := app.createChain()
	if err != nil {
		panic(err)
	}

	// Start the ETH chain
	go func() {
		utils.LogDebug("Starting the ETH chain....")

		app.chain.Start()

		app.chain.BlockChain().Accept(app.chain.GetGenesisBlock())

		app.listenKillSignal()

		// if app.ethConfig.ImportAccount {
		// 	app.chain.ImportAccounts()
		// }
	}()

	// Start ETH API server
	go app.startApiServer()

	return nil
}

func (app *App) createConfigDir() {
	if _, err := os.Stat(app.appConfigs.ConfigDir); os.IsNotExist(err) {
		utils.LogInfo("Creating app configuration directory:", app.appConfigs.ConfigDir)
		os.Mkdir(app.appConfigs.ConfigDir, os.ModeDir|0755)
	}
}

// startApiServer starts an ETH RPC api server
func (app *App) startApiServer() {
	chain := app.chain
	s := &ethchain.Server{}

	handler := chain.NewRPCHandler(time.Second * 10)
	handler.RegisterName("web3", &extra.Web3API{})
	handler.RegisterName("net", &extra.NetAPI{NetworkId: "1"})
	handler.RegisterName("evm", &extra.EvmApi{})

	chain.AttachEthService(handler, []string{"eth", "personal", "txpool", "debug"})

	s.Initialize(app.ethConfig.Host, uint16(app.ethConfig.Port), []string{}, handler)

	go s.Dispatch()
}

func (app *App) createChain() error {
	db, err := app.getChainDb()
	if err != nil {
		return err
	}
	chain := ethchain.NewETHChain(app.ethConfig, db, eth.DefaultSettings, true,
		app.broadcastSubmittedTx)

	err = chain.Initialize()
	if err != nil {
		return err
	}

	app.chain = chain

	return nil
}

func (app *App) getChainDb() (ethdb.Database, error) {
	var db ethdb.Database
	var err error

	if app.ethConfig.UseInMemDb {
		utils.LogInfo("Use In memory for ETH")
		db = rawdb.NewMemoryDatabase()
	} else {
		utils.LogInfo("Use real DB for ETH")
		// Use level DB.
		// TODO: Create new configs.
		db, err = rawdb.NewLevelDBDatabase(app.ethConfig.DbPath, 1024, 500, "metrics_")
	}

	return db, err
}

func (app *App) Shutdown() error {
	app.chain.Stop()

	return nil
}

func (app *App) listenKillSignal() {
	cServer.TrapSignal(func() {
		utils.LogDebug("Shutting down sisu app...")
		app.Shutdown()
	})
}

func (app *App) broadcastSubmittedTx(tx *dcore.Transaction) {
	// TODO: Broadcast to cosmos network here.
}

func (app *App) DeliverTx(req abci.RequestDeliverTx) abci.ResponseDeliverTx {
	fmt.Println("Delivering tx.....")
	app.BaseApp.DeliverTx(req)

	return abci.ResponseDeliverTx{Code: 0}
}

// BeginBlocker application updates every begin block
func (app *App) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	fmt.Println("Begin blocker....")
	app.mm.BeginBlock(ctx, req)

	// ETH chain
	app.chain.BeginBlock()

	return abci.ResponseBeginBlock{}
}

// EndBlocker application updates every end block
func (app *App) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	utils.LogDebug("End blocker....")
	app.chain.EndBlock()

	res := app.mm.EndBlock(ctx, req)

	return res
}
