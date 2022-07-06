package types

const (
	MsgTypeKeygenWithSigner               = "KeygenWithSigner"
	MsgTypeKeygenResultWithSigner         = "KeygenResultWithSigner"
	MsgTypeContractsWithSigner            = "ContractsWithSigner"
	MsgTypeTxOutWithSigner                = "TxOutWithSigner"
	MsgTypeTxOutConfirmMsg                = "TxOutConfirmMsg"
	MsgTypeTxsInMsg                       = "TxsInMsg"
	MsgTypeGasPriceWithSigner             = "GasPriceWithSigner"
	MsgTypeUpdateTokenPrice               = "UpdateTokenPrice"
	MsgTypePauseContract                  = "PauseContract"
	MsgTypeResumeContract                 = "ResumeContract"
	MsgTypeContractChangeOwnership        = "ContractChangeOwnership"
	MsgTypeContractChangeLiquidityAddress = "ContractChangeLiquidityAddress"
	MsgTypeContractLiquidityWithdrawFund  = "ContractLiquidityWithdrawFund"
	MsgTypeTransferRequestsMsg            = "TransferRequestsMsg"
	MsgTypeFundGatewayMsg                 = "FundGatewayMsg"
	MsgTypeBlockHeightMsg                 = "BlockHeightMsg"

	MsgTypeTxOut         = "TxOut"
	MsgTypeKeysignResult = "KeysignResult"
)
