package types

func ConvertBlockHeightsMapToArray(m map[string]*BlockHeight) ([]string, []*BlockHeight) {
	signers := make([]string, 0, len(m))
	blockHeights := make([]*BlockHeight, 0, len(m))

	for signer, blockHeight := range m {
		signers = append(signers, signer)
		blockHeights = append(blockHeights, blockHeight)
	}

	return signers, blockHeights
}
