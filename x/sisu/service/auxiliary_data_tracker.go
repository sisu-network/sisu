package service

import (
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"

	deyesethtypes "github.com/sisu-network/deyes/chains/eth/types"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"
)

var (
	AuxiliaryDataUpdateInterval = int64(1000 * 30) // 30s
)

// AuxiliaryDataTracker tracks auxiliary data for transaction (e.g. gas price in ETH or recent block
// hash in Solana).
type AuxiliaryDataTracker interface {
	UpdateData(ctx sdk.Context, chains []string)
}

type defaultAuxiliaryDataTracker struct {
	deyesClient external.DeyesClient
	appKeys     common.AppKeys
	keeper      keeper.Keeper
	txSubmit    common.TxSubmit

	lastUpdateTime map[string]int64
	lock           *sync.RWMutex
}

func NewAuxiliaryDataTracker(
	deyesClient external.DeyesClient,
	appKeys common.AppKeys,
	k keeper.Keeper,
	txSubmit common.TxSubmit,
) AuxiliaryDataTracker {
	return &defaultAuxiliaryDataTracker{
		deyesClient: deyesClient,
		appKeys:     appKeys,
		keeper:      k,
		txSubmit:    txSubmit,
		lock:        &sync.RWMutex{},
	}
}

// UpdateData checks chain's auxiliary data (like gas price, recent hash block) and sends update
// message to Sisu chain if needed.
func (tracker *defaultAuxiliaryDataTracker) UpdateData(ctx sdk.Context, chains []string) {
	for _, chain := range chains {
		// ETH
		if libchain.IsETHBasedChain(chain) {
			gasInfo, err := tracker.deyesClient.GetGasInfo(chain)
			if err != nil {
				continue
			}

			log.Verbosef("gasInfo for chain %s = %v\n", chain, gasInfo)

			chainCfg := tracker.keeper.GetChain(ctx, chain)
			tracker.updateGasPriceIfNeeded(ctx, chain, chainCfg.EthConfig, gasInfo)
		}

		// TODO: Move the solana's recent block to here.
	}
}

func (tracker *defaultAuxiliaryDataTracker) updateGasPriceIfNeeded(ctx sdk.Context, chain string,
	ethCfg *types.ChainEthConfig, gasInfo *deyesethtypes.GasInfo) {
	update := false

	if ethCfg.UseEip_1559 {
		baseFee := gasInfo.BaseFee
		tip := gasInfo.Tip

		if tip != 0 {
			// if Tip changes by more than 20% (either downward or upward, update the tip)
			ratio := ethCfg.Tip * 100 / tip
			if ratio < 80 || ratio > 120 {
				update = true
			}

			// if baseFee changes by more than 20% (either downward or upward, update the tip)
			ratio = ethCfg.BaseFee * 100 / baseFee
			if ratio < 80 || ratio > 120 {
				update = true
			}
		}
	} else {
		gasPrice := gasInfo.GasPrice
		if gasPrice != 0 {
			ratio := ethCfg.GasPrice * 100 / gasPrice
			if ratio < 90 || ratio > 110 {
				update = true
			}
		}
	}

	if update {
		tracker.lock.RLock()
		lastUpdateTime := tracker.lastUpdateTime[chain]
		tracker.lock.RUnlock()

		if time.Now().UnixMilli()-lastUpdateTime > AuxiliaryDataUpdateInterval {
			// Post update message now
			msg := types.NewGasPriceMsg(
				tracker.appKeys.GetSignerAddress().String(),
				[]string{chain},
				[]int64{gasInfo.GasPrice},
				[]int64{gasInfo.BaseFee},
				[]int64{gasInfo.Tip},
			)

			if err := tracker.txSubmit.SubmitMessageAsync(msg); err != nil {
				log.Errorf("Failed to submit tx, err = %s", err)
			}
		}
	}
}
