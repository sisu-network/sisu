package ante

import (
	cosmosTypes "github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	cosmosAnte "github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	evmKeeper "github.com/sisu-network/sisu/x/evm/keeper"
	evmTypes "github.com/sisu-network/sisu/x/evm/types"
)

func NewAnteHandler(
	ak AccountKeeper,
	bankKeeper BankKeeper,
	evmKeeper evmKeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	validator ethchain.EthValidator,
) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (newCtx sdk.Context, err error) {
		var anteHandler sdk.AnteHandler
		utils.LogDebug("Running ante. checkTx & recheck = ", ctx.IsCheckTx(), ctx.IsReCheckTx())

		// If the tx contains a non-evm message, we calculate the transaction like normal with
		// gas fee taken into account. Otherwise, use EvmAnteHandler which does not subtract gas fee
		// from the sender.
		if hasNonEvmTx(tx) {
			anteHandler = cosmosAnte.NewAnteHandler(ak, bankKeeper, sigGasConsumer, signModeHandler)
		} else {
			utils.LogDebug("This is an EVM transaction")
			anteHandler = EvmAnteHandler(ctx, tx, ak, bankKeeper, evmKeeper, sigGasConsumer, signModeHandler, validator)
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

func EvmAnteHandler(
	ctx sdk.Context,
	tx sdk.Tx,
	ak AccountKeeper,
	bankKeeper BankKeeper,
	evmKeeper evmKeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	validator ethchain.EthValidator,
) sdk.AnteHandler {
	decors := []cosmosTypes.AnteDecorator{
		NewSetUpContextDecorator(), // outermost AnteDecorator. SetUpContext must be called first
		NewRejectExtensionOptionsDecorator(),
		// TODO: Check signature of the sender. Only valdiator can submit evm tx.
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

	// If this is a checkTx or recheckTx, check to see if we can add the tx to the ETH mempool.
	if ctx.IsCheckTx() || ctx.IsReCheckTx() {
		decors = append(decors, NewEvmTxDecorator(validator))
	}

	return sdk.ChainAnteDecorators(
		decors...,
	)
}
