package utils

import (
	"fmt"
	"math/big"

	eTypes "github.com/ethereum/go-ethereum/core/types"
)

func GetChainIntFromId(chain string) *big.Int {
	switch chain {
	case "eth":
		return big.NewInt(1)
	default:
		LogError("unknown chain:", chain)
		return big.NewInt(0)
	}
}

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

		return KeccakHash32(string(serialized)), nil
	}

	// TODO: Support more chain other than ETH family.
	return "", fmt.Errorf("Unknwon chain: %s", chain)
}
