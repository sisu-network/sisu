package tss

import (
	"fmt"

	sdk "github.com/sisu-network/cosmos-sdk/types"
	sdkerrors "github.com/sisu-network/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/common"
	"github.com/sisu-network/sisu/x/tss/keeper"
	"github.com/sisu-network/sisu/x/tss/types"
)

// NewHandler ...
func NewHandler(k keeper.DefaultKeeper, txSubmit common.TxSubmit, processor *Processor) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.KeygenWithSigner:
			return handleKeygenProposal(ctx, msg, processor)
		case *types.KeygenResultWithSigner:
			return handleKeygenResult(ctx, msg, processor)
		case *types.TxInWithSigner:
			return handleTxIn(ctx, msg, processor)
		case *types.TxOutWithSigner:
			return handleTxOut(ctx, msg, processor)
		case *types.KeysignResult:
			return handleKeysignResult(ctx, msg, processor)
		case *types.ContractsWithSigner:
			return handleContractWithSigner(ctx, msg, processor)

		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

func handleKeygenProposal(ctx sdk.Context, msg *types.KeygenWithSigner, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverKeygen(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleKeygenResult(ctx sdk.Context, msg *types.KeygenResultWithSigner, processor *Processor) (*sdk.Result, error) {
	log.Verbose("Handling TSS Keygen result")
	data, err := processor.deliverKeygenResult(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxIn(ctx sdk.Context, msg *types.TxInWithSigner, processor *Processor) (*sdk.Result, error) {
	// Update the count for all txs.
	log.Verbose("Handling TxIn for chain", msg.Data.Chain)
	data, err := processor.deliverTxIn(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxOut(ctx sdk.Context, msg *types.TxOutWithSigner, processor *Processor) (*sdk.Result, error) {
	log.Verbose("Handling Txout, hash = ", msg.Data.GetHash())
	data, err := processor.deliverTxOut(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleKeysignResult(ctx sdk.Context, msg *types.KeysignResult, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverKeysignResult(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleContractWithSigner(ctx sdk.Context, msg *types.ContractsWithSigner, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverContracts(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}
