package keeper

import (
	"fmt"
	"os"

	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	dcore "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	"github.com/sisu-network/sisu/x/evm/types"
)

type Keeper struct {
	cdc       codec.Marshaler
	ethConfig *config.ETHConfig
	chain     *ethchain.ETHChain
}

func NewKeeper(cdc codec.Marshaler) *Keeper {
	keeper := &Keeper{
		cdc: cdc,
	}

	baseDir := os.Getenv("HOME") + "/.sisu"
	keeper.ethConfig = config.LocalETHConfig(baseDir)

	// 	// Setting log level
	ethLog.Root().SetHandler(ethLog.LvlFilterHandler(
		ethLog.LvlDebug, ethLog.StreamHandler(os.Stderr, ethLog.TerminalFormat(false))))

	err := keeper.createChain()
	if err != nil {
		panic(err)
	}

	keeper.chain.Start()

	go func() {
		// Import account if needed. Used in dev mode only
		if keeper.ethConfig.ImportAccount {
			keeper.chain.ImportAccounts()
		}

	}()

	return keeper
}

func (k *Keeper) createChain() error {
	k.chain = ethchain.NewETHChain(k.ethConfig, eth.DefaultSettings, k.onTxSubmitted)

	err := k.chain.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) onTxSubmitted(*dcore.Transaction) {

}