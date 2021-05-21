package keeper

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ethLog "github.com/ethereum/go-ethereum/log"
	"github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authKeepr "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	etypes "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/dcore/eth"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/config"
	evmCodec "github.com/sisu-network/sisu/x/evm/codec"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	"github.com/sisu-network/sisu/x/evm/types"
)

type Keeper struct {
	txSubmitter common.TxSubmit

	client    *grpc.ClientConn
	cdc       codec.Marshaler
	ethConfig *config.ETHConfig
	chain     *ethchain.ETHChain
	ak        *authKeepr.AccountKeeper
}

func NewKeeper(cdc codec.Marshaler, txSubmitter common.TxSubmit) *Keeper {
	keeper := &Keeper{
		cdc:         cdc,
		txSubmitter: txSubmitter,
	}

	// TODO: Put this in the config file.
	baseDir := os.Getenv("HOME") + "/.sisu"
	keeper.ethConfig = config.LocalETHConfig(baseDir)

	// 	// Setting log level
	ethLog.Root().SetHandler(ethLog.LvlFilterHandler(
		ethLog.LvlDebug, ethLog.StreamHandler(os.Stderr, ethLog.TerminalFormat(false))))

	return keeper
}

func (k *Keeper) Initialize() {
	err := k.createChain()
	if err != nil {
		panic(err)
	}

	k.chain.Start()
	k.listenSignalKill()

	go func() {
		// Import account if needed. Used in dev mode only
		if k.ethConfig.ImportAccount {
			k.chain.ImportAccounts()
		}
	}()
}

func (k *Keeper) listenSignalKill() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c

		k.chain.Stop()
	}()
}

func (k *Keeper) createChain() error {
	k.chain = ethchain.NewETHChain(k.ethConfig, eth.DefaultSettings, k.txSubmitter)

	err := k.chain.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func (k *Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) DeliverTx(etx *etypes.Transaction) ([]byte, error) {
	receipt, _, err := k.chain.DeliverTx(etx)

	if err != nil {
		return []byte{}, nil
	}

	// RLP encode of the receipt
	var buff bytes.Buffer
	err = receipt.EncodeRLP(bufio.NewWriter(&buff))
	if err != nil {
		return []byte{}, nil
	}

	// Prefixed length encoded
	prefixedData, err := evmCodec.EncodePrefixedLength(buff.Bytes())
	if err != nil {
		return []byte{}, nil
	}

	return prefixedData, err
}

func (k *Keeper) BeginBlock() error {
	k.chain.BeginBlock()

	return nil
}

func (k *Keeper) EndBlock() error {
	k.chain.EndBlock()

	return nil
}
