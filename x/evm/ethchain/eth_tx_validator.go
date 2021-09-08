package ethchain

import (
	etypes "github.com/sisu-network/dcore/core/types"
)

type EthValidator interface {
	CheckTx(txs []*etypes.Transaction) error
}
