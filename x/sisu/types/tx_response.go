package types

import (
	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// A wrapper around a tx that responds to a incoming tx.
type TxResponse struct {
	OutChain string
	EthTx    *ethtypes.Transaction
	RawBytes []byte
}
