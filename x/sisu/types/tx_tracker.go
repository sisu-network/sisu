package types

type TxStatus int64

const (
	// Any update to this status enum should update the StatusStrings as well.
	TxStatusUnknown TxStatus = iota
	TxStatusCreated
	TxStatusDelivered
	TxStatusSigned
	TxStatusSignFailed
	TxStatusDepoyed // transaction has been sent to blockchain but not confirmed yet.
	TxStatusConfirmed
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
