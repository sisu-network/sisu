package utils

import (
	"strings"

	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func GetTokenOnChain(allTokens map[string]*types.Token, tokenAddr, targetChain string) *types.Token {
	for _, t := range allTokens {
		if len(t.Chains) != len(t.Addresses) {
			log.Error("Chains length is not the same as address length ")
			log.Error("t.Chains = ", t.Chains)
			log.Error("t.Addresses = ", t.Addresses)
			continue
		}

		for j, chain := range t.Chains {
			if chain == targetChain && strings.EqualFold(t.Addresses[j], tokenAddr) {
				return t
			}
		}
	}

	return nil
}
