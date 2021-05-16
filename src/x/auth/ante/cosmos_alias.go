package ante

import (
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
	types "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var (
	GetSignerAcc                       = ante.GetSignerAcc
	DefaultSigVerificationGasConsumer  = ante.DefaultSigVerificationGasConsumer
	DeductFees                         = ante.DeductFees
	SetGasMeter                        = ante.SetGasMeter
	NewSetUpContextDecorator           = ante.NewSetUpContextDecorator
	NewRejectExtensionOptionsDecorator = ante.NewRejectExtensionOptionsDecorator
	NewValidateBasicDecorator          = ante.NewValidateBasicDecorator
	NewValidateMemoDecorator           = ante.NewValidateMemoDecorator
	NewConsumeGasForTxSizeDecorator    = ante.NewConsumeGasForTxSizeDecorator
	NewRejectFeeGranterDecorator       = ante.NewRejectFeeGranterDecorator
	NewSetPubKeyDecorator              = ante.NewSetPubKeyDecorator
	NewValidateSigCountDecorator       = ante.NewValidateSigCountDecorator
	NewDeductFeeDecorator              = ante.NewDeductFeeDecorator
	NewSigGasConsumeDecorator          = ante.NewSigGasConsumeDecorator
	NewSigVerificationDecorator        = ante.NewSigVerificationDecorator
	NewIncrementSequenceDecorator      = ante.NewIncrementSequenceDecorator
	NewMempoolFeeDecorator             = ante.NewMempoolFeeDecorator
)

type (
	AccountKeeper                    = ante.AccountKeeper
	SignatureVerificationGasConsumer = ante.SignatureVerificationGasConsumer
	BankKeeper                       = types.BankKeeper
	TxTimeoutHeightDecorator         = ante.TxTimeoutHeightDecorator
)
