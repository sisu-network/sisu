package types

import (
	"encoding/binary"
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// GetId returns a unique id of the txOut. This is the hash of tx out type, its content and list
// of transfer id inputs.
func (t *TxOutOld) GetId() (string, error) {
	hash := sha3.NewLegacyKeccak256()

	_, err := hash.Write([]byte(t.TxType.String()))
	if err != nil {
		return "", err
	}

	_, err = hash.Write(t.Content.OutBytes)
	if err != nil {
		return "", err
	}

	_, err = hash.Write([]byte(t.Content.OutChain))
	if err != nil {
		return "", err
	}

	_, err = hash.Write([]byte(t.Content.OutHash))
	if err != nil {
		return "", err
	}

	bz := make([]byte, 4)
	binary.BigEndian.PutUint32(bz, uint32(t.Content.RetryNum))
	_, err = hash.Write(bz)
	if err != nil {
		return "", err
	}

	for _, id := range t.Input.TransferIds {
		_, err = hash.Write([]byte(id))
		if err != nil {
			return "", err
		}
	}

	var buf []byte
	buf = hash.Sum(nil)

	encoded := hex.EncodeToString(buf)
	if len(encoded) > 32 {
		encoded = encoded[:32]
	}

	return encoded, nil
}
