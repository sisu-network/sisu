package keeper

import (
	"fmt"
	"os"
	"sync"

	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	etypes "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	"github.com/sisu-network/sisu/x/evm/types"
)

type Keeper struct {
	txSubmitter *TxSubmitter

	client    *grpc.ClientConn
	cdc       codec.Marshaler
	ethConfig *config.ETHConfig
	chain     *ethchain.ETHChain

	softState *ethchain.SoftState
	ssLock    *sync.RWMutex
}

func NewKeeper(cdc codec.Marshaler, sisuHome string, keyRingBackend string) *Keeper {
	keeper := &Keeper{
		cdc:         cdc,
		txSubmitter: NewTxSubmitter(sisuHome, keyRingBackend),
		ssLock:      &sync.RWMutex{},
	}

	// TODO: Put this in the config file.
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
	k.chain = ethchain.NewETHChain(k.ethConfig, eth.DefaultSettings, k.txSubmitter.onTxSubmitted)

	err := k.chain.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) DeliverTx(etx *etypes.Transaction) error {
	_, err := k.chain.DeliverTx(etx)

	return err
}

func (k *Keeper) BeginBlock() error {
	// Create a new softstate for execution in this block.
	k.chain.BeginBlock()

	return nil
}

func (k *Keeper) EndBlock() error {
	utils.LogDebug("End of block")

	k.chain.EndBlock()

	return nil
}
