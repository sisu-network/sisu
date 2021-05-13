package ante

import (
	ante "github.com/cosmos/cosmos-sdk/x/auth/ante"
)

var (
	NewAnteHandler                    = ante.NewAnteHandler
	GetSignerAcc                      = ante.GetSignerAcc
	DefaultSigVerificationGasConsumer = ante.DefaultSigVerificationGasConsumer
	DeductFees                        = ante.DeductFees
	SetGasMeter                       = ante.SetGasMeter
)
