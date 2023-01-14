package types

type TxOutState int

const (
	TxOutState_Unknown   TxOutState = 0
	TxOutState_Confirmed TxOutState = 1
	TxOutState_Signing   TxOutState = 2
	TxOutState_Done      TxOutState = 3
)
