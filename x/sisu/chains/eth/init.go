package eth

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sisu-network/sisu/contracts/eth/vault"
	"github.com/sisu-network/sisu/utils"
)

const (
	ContractVault = "vault"

	// Methods in gateway smart contract
	MethodTransferIn            = "transferIn"
	MethodTransferInMultiple    = "transferInMultiple"
	MethodTransferOut           = "transferOut"
	MethodSetNotPausedChain     = "setNotPausedChain"
	MethodRemoteExecute         = "remoteExecute"
	MethodRemoteExecuteMultiple = "remoteExecuteMultiple"
)

type ContractInfo struct {
	AbiString, Bin, AbiHash string
	Abi                     abi.ABI
	MethodNames             []string
	IsDeployBySisu          bool
}

var (
	SupportedContracts = map[string]*ContractInfo{
		ContractVault: {
			AbiString:      vault.VaultABI,
			Bin:            vault.VaultBin,
			AbiHash:        utils.KeccakHash32(vault.VaultBin),
			IsDeployBySisu: true,
		},
	}
)

// init initializes variables used throughout this package.
func init() {
	// 1. Initializes abi fields for SupportedContracts
	for _, contractInfo := range SupportedContracts {
		var err error
		contractInfo.Abi, err = abi.JSON(strings.NewReader(contractInfo.AbiString))
		if err != nil {
			panic("error when read abi json")
		}
	}

	// // 2. Make sure that all the method names in our struct are present in the ABI methods. This is
	// // to make sure that the data in the contract is consistent with our constants.
	// for _, contractData := range SupportedContracts {
	// 	for _, methodName := range contractData.MethodNames {
	// 		_, ok := contractData.Abi.Methods[methodName]
	// 		if !ok {
	// 			panic(fmt.Errorf("cannot find method name '%s' in the generated abi", methodName))
	// 		}
	// 	}
	// }
}
