package types

import "math/big"

type TransferOutData struct {
	BlockHeight int64
	DestChain   string
	Token       *Token
	Recipient   string
	Amount      *big.Int

	// For tx_tracker
	InChain string
	InHash  string
	TxIn    *TransferOut
}
