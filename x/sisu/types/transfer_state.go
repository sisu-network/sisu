package types

type TransferState int

const (
	TransferState_Unknown      TransferState = 0
	TransferState_Confirmed    TransferState = 1
	TransferState_WaitForTxOut TransferState = 2
	TransferState_Failure      TransferState = 3
)
