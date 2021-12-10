package types

type TxOutEntity struct {
	OutChain        string
	HashWithoutSig  string
	HashWithSig     string
	InChain         string
	InHash          string
	BytesWithoutSig []byte
	Status          string
	Signature       string

	ContractHash string // optional field for Eth contracts
}

func (txOut *TxOut) ToEntity() *TxOutEntity {
	return &TxOutEntity{
		OutChain:        txOut.OutChain,
		HashWithoutSig:  txOut.GetHash(),
		InChain:         txOut.InChain,
		InHash:          txOut.InHash,
		BytesWithoutSig: txOut.OutBytes,
	}
}
