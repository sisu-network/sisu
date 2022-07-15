package types

import fmt "fmt"

func GetTransferId(chain, inHash string) string {
	return fmt.Sprintf("%s__%s", chain, inHash)
}
