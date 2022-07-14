package app

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/sisu-network/lib/log"
)

// LogStackTraceOnRunTxFail log stack trace when run tx failed
// To test this, developer can add a panic in msg.ValidateBasic() of any msg
// For example: panic(sdkerrors.Wrap(sdkerrors.ErrPanic, "invalid msg"))
func LogStackTraceOnRunTxFail(recoveryObj interface{}) error {
	log.Errorf("type of recoveryObj: %T\n", recoveryObj)
	log.Error("recoveryObj: ", recoveryObj)
	err, ok := recoveryObj.(error)
	if !ok {
		return nil
	}

	if sdkerrors.ErrPanic.Is(err) {
		log.Error(fmt.Errorf("panic app with error = %w\n", err))
		return err
	}

	// Other type of error, still need to be printed
	log.Error(err)
	return err
}
