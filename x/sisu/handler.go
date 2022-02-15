package sisu

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

// Handler is an interface that all handlers in this app should implement.
type Handler interface {
	// Deliver and execute a transaction message.
	DeliverMsg(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error)
}

type SisuHandler struct {
}

// NewHandler ...
func NewHandler(processor *Processor, valsManager ValidatorManager) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, fmt.Errorf("incorrect signers length: %d", len(signers))
		}

		if !valsManager.IsValidator(signers[0].String()) {
			log.Verbose("sender is not a validator", signers[0].String())
			return nil, fmt.Errorf("sender is not a validator: %s", signers[0].String())
		}

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
		case *types.TxOutContractConfirmWithSigner:
			return handleTxOutContractConfirm(ctx, msg, processor)
		case *types.KeysignResult:
			return handleKeysignResult(ctx, msg, processor)
		case *types.ContractsWithSigner:
			return handleContractWithSigner(ctx, msg, processor)
		case *types.GasPriceMsg:
			return handleGasPriceMsg(ctx, msg, processor)
		case *types.UpdateTokenPrice:
			return handleUpdateTokenPrice(ctx, msg, processor)

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
	data, err := processor.deliverKeygenResult(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxIn(ctx sdk.Context, msg *types.TxInWithSigner, processor *Processor) (*sdk.Result, error) {
	// Update the count for all txs.
	data, err := processor.deliverTxIn(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxOut(ctx sdk.Context, msg *types.TxOutWithSigner, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverTxOut(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleTxOutContractConfirm(ctx sdk.Context, msg *types.TxOutContractConfirmWithSigner, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverTxOutContractConfirm(ctx, msg)
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

func handleGasPriceMsg(ctx sdk.Context, msg *types.GasPriceMsg, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverGasPriceMsg(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}

func handleUpdateTokenPrice(ctx sdk.Context, msg *types.UpdateTokenPrice, processor *Processor) (*sdk.Result, error) {
	data, err := processor.deliverUpdateTokenPrice(ctx, msg)
	return &sdk.Result{
		Data: data,
	}, err
}
