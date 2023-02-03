package types

import (
	"fmt"

	"github.com/sisu-network/lib/log"
)

// GetId returns a unique id of the txOut. This is the hash of tx out type, its content and list
// of transfer id inputs.
func (t *TxOut) GetId() string {
	switch t.TxType {
	case TxOutType_TRANSFER_OUT:
		return fmt.Sprintf("%s__%s", t.Content.OutChain, t.Content.OutHash)
	default:
		log.Errorf("TxOut GetId is not implemented for type %s", t.TxType.String())
		return ""
	}
}

// GetValidatorId returns an id string that could be used to find the assigned validator.
func (t *TxOut) GetValidatorId() string {
	switch t.TxType {
	case TxOutType_TRANSFER_OUT:
		if len(t.Input.TransferIds) == 0 {
			log.Errorf("TxOut transfer out does not have associated transfer input")
			return ""
		}

		return t.Input.TransferIds[0]
	default:
		log.Errorf("TxOut GetValidatorId is not implemented for type %s", t.TxType.String())
		return ""
	}
}
