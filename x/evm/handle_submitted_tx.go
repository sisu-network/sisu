package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/keeper"
	"github.com/sisu-network/sisu/x/evm/types"
)

func handleSubmittedTx(ctx sdk.Context, k keeper.Keeper, etxMsg *types.EthTx) (*sdk.Result, error) {
	etx := new(etypes.Transaction)
	err := etx.UnmarshalJSON(etxMsg.Data)
	if err != nil {
		utils.LogError("Cannot unmarshall etx", err)
		return &sdk.Result{}, err
	}

	data, err := k.DeliverTx(etx)
	if err != nil {
		return &sdk.Result{
			Data: data,
		}, err
	}

	return &sdk.Result{
		Data:   data,
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
