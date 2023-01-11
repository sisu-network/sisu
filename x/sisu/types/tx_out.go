package types

import (
	"fmt"
)

// GetId returns a unique id of the txOut. This is the hash of tx out type, its content and list
// of transfer id inputs.
func (t *TxOutOld) GetId() string {
	return fmt.Sprintf("%s__%s", t.Content.OutChain, t.Content.OutHash)
}
