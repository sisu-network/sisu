package ante

import (
	"fmt"

	cosmosTypes "github.com/sisu-network/cosmos-sdk/types"
	sdk "github.com/sisu-network/cosmos-sdk/types"
	cosmosAnte "github.com/sisu-network/cosmos-sdk/x/auth/ante"
	"github.com/sisu-network/cosmos-sdk/x/auth/signing"
	"github.com/sisu-network/sisu/config"
	"github.com/sisu-network/sisu/utils"
	"github.com/sisu-network/sisu/x/evm/ethchain"
	evmKeeper "github.com/sisu-network/sisu/x/evm/keeper"
	evmTypes "github.com/sisu-network/sisu/x/evm/types"
	"github.com/sisu-network/sisu/x/tss"
	tssTypes "github.com/sisu-network/sisu/x/tss/types"
)

type SisuTxType int

const (
	TYPE_TX_COSMOS SisuTxType = iota
	TYPE_TX_ETH
	TYPE_TX_TSS
)

func NewAnteHandler(
	tssConfig *config.TssConfig,
	ak AccountKeeper,
	bankKeeper BankKeeper,
	evmKeeper evmKeeper.Keeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	ethValidator ethchain.EthValidator,
	tssValidator tss.TssValidator,
) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, sim bool) (sdk.Context, error) {
		var anteHandler sdk.AnteHandler
		utils.LogDebug("Running ante. checkTx & recheck = ", ctx.IsCheckTx(), ctx.IsReCheckTx())

		txType, err := getTxType(tx)
		if err != nil {
			// TODO: Handle when there are errors here.
			return ctx, err
		}

		switch txType {
		case TYPE_TX_ETH:
			anteHandler = EvmAnteHandler(ctx, tx, ak, bankKeeper, evmKeeper, sigGasConsumer, signModeHandler, ethValidator)
		case TYPE_TX_TSS:
			if tssConfig.Enable {
				// Add TSS AnteHandler here.
				anteHandler = TssAnteHandler(ctx, tx, ak, bankKeeper, sigGasConsumer, signModeHandler, tssValidator)
			} else {
				return ctx, fmt.Errorf("Tss is not enabled")
			}
		case TYPE_TX_COSMOS:
			anteHandler = cosmosAnte.NewAnteHandler(ak, bankKeeper, sigGasConsumer, signModeHandler)
		}

		return anteHandler(ctx, tx, sim)
	}
}

func getTxType(tx sdk.Tx) (SisuTxType, error) {
	var txType, lastType SisuTxType

	msgs := tx.GetMsgs()

	for i, msg := range msgs {
		switch msg.Route() {
		case evmTypes.ModuleName:
			txType = TYPE_TX_ETH
		case tssTypes.ModuleName:
			txType = TYPE_TX_TSS
		default:
			txType = TYPE_TX_COSMOS
		}
		if i > 0 && txType != lastType {
			return txType, fmt.Errorf("There are more than 1 message types in a transaction")
		}

		lastType = txType
	}

	return txType, nil
}

// TODO: Clean up EVM Ante handler.
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

// TODO: Clean up TSS Ante handler.
func TssAnteHandler(
	ctx sdk.Context,
	tx sdk.Tx,
	ak AccountKeeper,
	bankKeeper BankKeeper,
	sigGasConsumer SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	validator tss.TssValidator,
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
		NewTssDecorator(validator),
		NewIncrementSequenceDecorator(ak),
	}

	return sdk.ChainAnteDecorators(
		decors...,
	)
}
