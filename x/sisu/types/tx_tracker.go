package types

type TxStatus int64

const (
	// Any update to this status enum should update the StatusStrings as well.
	TxStatusUnknown TxStatus = iota
	TxStatusCreated
	TxStatusDelivered
	TxStatusSigned
	TxStatusSignFailed
	TxStatusDepoyed   // transaction has been sent to blockchain but not confirmed yet.
	TxStatusConfirmed // Tx is successfully executed in a block.
	TxStatusReverted  // Tx is included in the blockchain but reverted/failed during execution.
)

var (
	StatusStrings = []string{
		"TxStatusUnknown",
		"TxStatusCreated",
		"TxStatusDelivered",
		"TxStatusSigned",
		"TxStatusSignFailed",
		"TxStatusDepoyed",
		"TxStatusConfirmed",
	}
)
