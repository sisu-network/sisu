package utils

import (
	"fmt"

	eTypes "github.com/ethereum/go-ethereum/core/types"
)

func IsETHBasedChain(chain string) bool {
	switch chain {
	case "eth":
		return true
	}

	return false
}

func GetTxHash(chain string, serialized []byte) (string, error) {
	if IsETHBasedChain(chain) {
		tx := &eTypes.Transaction{}
		err := tx.UnmarshalBinary(serialized)
		if err != nil {
			return "", err
		}

		bz, err := tx.MarshalBinary()
		if err != nil {
			return "", err
		}

		return KeccakHash32(string(bz))
	}

	// TODO: Support more chain other than ETH family.
	return "", fmt.Errorf("Unknwon chain: %s", chain)
}
