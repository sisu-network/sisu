package tss

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper, txSubmit common.TxSubmit, processor *Processor) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.KeygenProposal:
			return handleKeygenProposal(msg, processor)
		case *types.KeygenResult:
			return handleKeygenResult(ctx, msg, processor)
		case *types.ObservedTxs:
			return handleObservedTxs(ctx, msg, processor)
		case *types.TxOut:
			return handleTxOut(ctx, msg, processor)
		case *types.KeysignResult:
			return handleKeysignResult(ctx, msg, processor)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleKeygenProposal(msg *types.KeygenProposal, processor *Processor) (*sdk.Result, error) {
	data, err := processor.DeliverKeyGenProposal(msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleKeygenResult(ctx sdk.Context, msg *types.KeygenResult, processor *Processor) (*sdk.Result, error) {
	utils.LogDebug("Handling TSS Keygen result")
	data, err := processor.DeliverKeygenResult(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleObservedTxs(ctx sdk.Context, msg *types.ObservedTxs, processor *Processor) (*sdk.Result, error) {
	// Update the count for all txs.
	utils.LogVerbose("Handling ObservedTxs")
	data, err := processor.DeliverObservedTxs(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxOut(ctx sdk.Context, msg *types.TxOut, processor *Processor) (*sdk.Result, error) {
	utils.LogVerbose("Handling Txout")
	data, err := processor.DeliverTxOut(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleKeysignResult(ctx sdk.Context, msg *types.KeysignResult, processor *Processor) (*sdk.Result, error) {
	utils.LogVerbose("Handling Keysign Result")
	data, err := processor.DeliverKeysignResult(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}
