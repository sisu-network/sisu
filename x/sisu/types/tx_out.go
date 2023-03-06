package types

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/sisu-network/lib/log"
	"golang.org/x/crypto/sha3"
)

// GetId returns a unique id of the txOut. This is the hash of tx out type, its content and list
// of transfer id inputs.
func (t *TxOut) GetId() string {
	switch t.TxType {
	case TxOutType_TRANSFER:
		hash := sha3.NewLegacyKeccak256()
		hash.Write([]byte(strings.Join(t.Input.TransferRetryIds, "")))
		return fmt.Sprintf("%s__%s__%s",
			t.Content.OutChain,
			t.Content.OutHash,
			hex.EncodeToString(hash.Sum(nil)[:8]),
		)
	default:
		log.Errorf("TxOut GetId is not implemented for type %s", t.TxType.String())
		return ""
	}
}

// GetValidatorId returns an id string that could be used to find the assigned validator.
func (t *TxOut) GetValidatorId() string {
	switch t.TxType {
	case TxOutType_TRANSFER:
		if len(t.Input.TransferRetryIds) == 0 {
			log.Errorf("TxOut transfer out does not have associated transfer input")
			return ""
		}

		return t.Input.TransferRetryIds[0]
	default:
		log.Errorf("TxOut GetValidatorId is not implemented for type %s", t.TxType.String())
		return ""
	}
}

func GetTxOutIdFromChainAndHash(chain, hash string) string {
	return fmt.Sprintf("%s__%s", chain, hash)
}
