package types

import (
	"fmt"
	"strconv"
	"strings"
)

func GetTransferId(chain, inHash string) string {
	return fmt.Sprintf("%s__%s", chain, inHash)
}

func (t *TransferDetails) GetRetryId() string {
	return fmt.Sprintf("%s___%d", t.Id, t.RetryNum)
}

func GetIdFromRetryId(uniq string) (string, int) {
	parts := strings.SplitN(uniq, "___", 2)
	num, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}

	return parts[0], num
}

func GetIdsFromRetryIds(uniqs []string) []string {
	ids := make([]string, len(uniqs))
	for i, uniq := range uniqs {
		ids[i], _ = GetIdFromRetryId(uniq)
	}
	return ids
}
