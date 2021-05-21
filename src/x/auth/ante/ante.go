package ante

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAnte "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	evmKeeper "github.com/sisu-network/sisu/x/evm/keeper"
	evmTypes "github.com/sisu-network/sisu/x/evm/types"
)

func NewAnteHandler(
	ak AccountKeeper,
	bankKeeper BankKeeper,
	evmKeeper evmKeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler

		fmt.Println("Running ante. checkTx = ", ctx.IsCheckTx())

		// If the tx contains a non-evm message, we calculate the transaction like normal with
		// gas fee taken into account. Otherwise, use EvmAnteHandler which does not subtract gas fee
		// from the sender.
		if hasNonEvmTx(tx) {
			anteHandler = cosmosAnte.NewAnteHandler(ak, bankKeeper, sigGasConsumer, signModeHandler)
		} else {
			anteHandler = EvmAnteHandler(ak, bankKeeper, evmKeeper, sigGasConsumer, signModeHandler)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func hasNonEvmTx(tx sdk.Tx) bool {
	msgs := tx.GetMsgs()
	for _, msg := range msgs {
		_, ok := msg.(*evmTypes.EthTx)
		if !ok {
			return true
		}
	}

	return false
}

func EvmAnteHandler(ak AccountKeeper,
	bankKeeper BankKeeper,
	evmKeeper evmKeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		// TODO: Add EVM mempool ante handler
		// NewMempoolFeeDecorator(), // No cosmos mempool
		NewValidateBasicDecorator(),
		TxTimeoutHeightDecorator{},
		NewValidateMemoDecorator(ak),
		NewConsumeGasForTxSizeDecorator(ak),
		NewRejectFeeGranterDecorator(),
		NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		NewValidateSigCountDecorator(ak),
		NewDeductFeeDecorator(ak, bankKeeper),
		NewSigGasConsumeDecorator(ak, sigGasConsumer),
		NewSigVerificationDecorator(ak, signModeHandler),
		NewIncrementSequenceDecorator(ak),
	)
}
