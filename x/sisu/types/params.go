package types

func (params *Params) GetMaxTransferOutBatch(chain string) int {
	for _, transferParam := range params.TransferOutParams {
		if transferParam.Chain == chain {
			return int(transferParam.MaxBatching)
		}
	}

	return 10
}
