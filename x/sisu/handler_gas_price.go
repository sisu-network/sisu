package sisu

import (
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/sisu/keeper"
	"github.com/sisu-network/sisu/x/sisu/types"

	libchain "github.com/sisu-network/lib/chain"
)

type HandlerGasPrice struct {
	keeper     keeper.Keeper
	globalData common.GlobalData
}

func NewHandlerGasPrice(mc ManagerContainer) *HandlerGasPrice {
	return &HandlerGasPrice{
		keeper:     mc.Keeper(),
		globalData: mc.GlobalData(),
	}
}

func (h *HandlerGasPrice) DeliverMsg(ctx sdk.Context, msg *types.GasPriceMsg) (*sdk.Result, error) {
	h.keeper.SetGasPrice(ctx, msg)

	params := h.keeper.GetParams(ctx)
	if params == nil {
		return nil, fmt.Errorf("Cannot find tss params")
	}

	h.keeper.SetGasPrice(ctx, msg)
	savedRecord := h.keeper.GetGasPriceRecord(ctx, msg.BlockHeight)

	allChains := make(map[string]bool)
	for _, record := range savedRecord.Messages {
		for _, chain := range record.Chains {
			allChains[chain] = true
		}
	}

	for chain := range allChains {
		prices := make([]int64, 0)
		for _, record := range savedRecord.Messages {
			for i, c := range record.Chains {
				if c == chain {
					prices = append(prices, record.Prices[i])
					break
				}
			}
		}

		if len(prices) >= int(params.MajorityThreshold) {
			// Calculate the median
			sort.SliceStable(prices, func(i, j int) bool {
				return prices[i] < prices[j]
			})

			median := prices[len(prices)/2]
			if median == 0 && libchain.IsETHBasedChain(chain) {
				log.Warn("Median gas price for chain ", chain, " is ", median)
			} else {
				log.Verbose("Median gas price for chain ", chain, " is ", median)
			}

			if median == 0 {
				continue
			}

			// Save to db
			savedChain := h.keeper.GetChain(ctx, chain)
			if savedChain == nil {
				savedChain = &types.Chain{
					Id: chain,
				}
			}
			savedChain.GasPrice = median
			h.keeper.SaveChain(ctx, savedChain)
		}
	}

	return &sdk.Result{}, nil
}
