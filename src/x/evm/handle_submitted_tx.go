package evm

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	etypes "github.com/sisu-network/dcore/core/types"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/keeper"
	"github.com/sisu-network/sisu/x/evm/types"

	"github.com/ethereum/go-ethereum/common"
)

func handleSubmittedTx(ctx sdk.Context, k keeper.Keeper, etxMsg *types.EthTx) (*sdk.Result, error) {
	etx := new(etypes.Transaction)
	err := etx.UnmarshalJSON(etxMsg.Data)
	if err != nil {
		utils.LogError("Cannot unmarshall etx", err)
		return &sdk.Result{}, err
	}

	from := common.HexToAddress(etxMsg.Author)
	err = k.DeliverTx(from, etx)
	if err != nil {
		return &sdk.Result{}, err
	}

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
