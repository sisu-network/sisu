package types

type TxOutEntity struct {
	OutChain       string
	HashWithoutSig string
	HashWithSig    string
	InChain        string
	InHash         string
	Outbytes       []byte

	ContractHash string // optional field for Eth contracts
}

func TxOutToEntity(txOut *TxOut) *TxOutEntity {
	return &TxOutEntity{
		OutChain:       txOut.OutChain,
		HashWithoutSig: txOut.GetHash(),
		InChain:        txOut.InChain,
		InHash:         txOut.InHash,
		Outbytes:       txOut.OutBytes,
	}
}
