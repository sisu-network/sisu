package keeper

import "fmt"

func GetEthNonceKey(chain string) string {
	return fmt.Sprintf("eth_nonce_%s", chain)
}
