package ante

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

func NewAnteHandler(
	ak AccountKeeper,
	bankKeeper BankKeeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	decors := []cosmosTypes.AnteDecorator{
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),

		// TODO: Check signature of the sender. Only valdiator can submit tss tx.
		// NewMempoolFeeDecorator(), // No cosmos mempool
		NewValidateBasicDecorator(),
		TxTimeoutHeightDecorator{},
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewRejectFeeGranterDecorator(),
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(ak),
		NewSigVerificationDecorator(ak, signModeHandler),
		NewIncrementSequenceDecorator(ak),
	}

	return sdk.ChainAnteDecorators(
		decors...,
	)
}
