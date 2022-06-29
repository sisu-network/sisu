package types

func NewExternalInfoMsg(signer string, gasPrice *GasPrice, block *BlockHeight) *ExternalInfoMsg {
	return &ExternalInfoMsg{
		Signer: signer,
		Data: &ExternalInfoData{
			GasPrice:     gasPrice,
			BlockHeights: []*BlockHeight{block},
		},
	}
}

func NewExternalInfoObservedBlockHeight(signer string, block *BlockHeight) *ExternalInfoMsg {
	return NewExternalInfoMsg(signer, nil, block)
}
