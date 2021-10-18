package tss

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	erc20gateway "github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
)

const (
	ContractErc20 = "erc20"

	MethodTransferOutFromContract = "transferOutFromContract"
	MethodTransferIn              = "transferIn"
)

var (
	SupportedContracts = map[string]struct {
		AbiString, Bin, AbiHash string
		Abi                     abi.ABI
		MethodNames             []string
	}{
		ContractErc20: {
			AbiString: erc20gateway.Erc20gatewayABI,
			Bin:       erc20gateway.Erc20gatewayBin,
			AbiHash:   utils.KeccakHash32(erc20gateway.Erc20gatewayBin),
			MethodNames: []string{
				MethodTransferOutFromContract,
				MethodTransferIn,
			},
		},
	}
)

// init initializes variables used throughout this package.
func init() {
	// 1. Initializes abi fields for SupportedContracts
	if entry, ok := SupportedContracts[ContractErc20]; ok {
		entry.Abi, _ = abi.JSON(strings.NewReader(SupportedContracts[ContractErc20].AbiString))
		SupportedContracts[ContractErc20] = entry
	}

	// 2. Make sure that all the method names in our struct are present in the ABI methods. This is
	// to make sure that the data in the contract is consistent with our constants.
	for _, contractData := range SupportedContracts {
		for _, methodName := range contractData.MethodNames {
			_, ok := contractData.Abi.Methods[methodName]
			if !ok {
				panic(fmt.Errorf("cannot find method name '%s' in the generated abi", methodName))
			}
		}
	}
}
