package types

// A wrapper around a tx that responds to a incoming tx.
type TxResponse struct {
	OutChain string
	RawBytes []byte
}
