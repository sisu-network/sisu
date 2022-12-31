package service

import (
	"sync"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	libchain "github.com/sisu-network/lib/chain"

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

	lastUpdateTime map[string]int64
	lock           *sync.RWMutex
}

func NewAuxiliaryDataTracker(
	deyesClient external.DeyesClient,
	appKeys common.AppKeys,
	k keeper.Keeper,
) AuxiliaryDataTracker {
	return &defaultAuxiliaryDataTracker{
		deyesClient: deyesClient,
		appKeys:     appKeys,
		keeper:      k,
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

			chainCfg := tracker.keeper.GetChain(ctx, chain)
			if chainCfg.EthConfig.UseEip_1559 {
				// Get base fee and tip
				baseFee := gasInfo.BaseFee
				tip := gasInfo.Tip

				tracker.updateGasPriceIfNeeded(ctx, chain, chainCfg.EthConfig, baseFee, tip)
			}
		}

		// TODO: Move the solana's recent block to here.
	}
}

func (tracker *defaultAuxiliaryDataTracker) updateGasPriceIfNeeded(ctx sdk.Context, chain string,
	ethCfg *types.ChainEthConfig, baseFee, tip int64) {
	update := false
	// if Tip changes by more than 20% (either downward or upward, up the tip)
	ratio := ethCfg.Tip * 100 / tip
	if ratio < 80 || ratio > 120 {
		update = true
	}

	if !update {
		if baseFee+tip > 2*ethCfg.BaseFee+ethCfg.Tip {
			update = true
		}
	}

	if update {
		tracker.lock.RLock()
		lastUpdateTime := tracker.lastUpdateTime[chain]
		tracker.lock.RUnlock()

		if time.Now().UnixMilli()-lastUpdateTime > AuxiliaryDataUpdateInterval {
			// Post update message now
		}
	}
}
