package sisu

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/utils"
)

const (
	ContractErc20Gateway = "erc20gateway"

	MethodTransferIn    = "transferIn"
	MethodTransferOut   = "transferOut"
	MethodPauseGateway  = "pauseGateway"
	MethodResumeGateway = "resumeGateway"
)

// TODO: export a command line to set liquid pool address for the gateway
var liquidPoolAddrs = map[string]ecommon.Address{
	"ganache1":            ecommon.HexToAddress("0x18eD078Bf666049f02FF8193e0d6B4D45B50329f"),
	"ganache2":            ecommon.HexToAddress("0x18eD078Bf666049f02FF8193e0d6B4D45B50329f"),
	"eth-ropsten":         ecommon.HexToAddress("0xbc77D44223a75194eab5006De96cd1EBa95dA374"),
	"eth-binance-testnet": ecommon.HexToAddress("0xAF234785e8c283129968a3D0219aEDe7E9B83953"),
	"fantom-testnet":      ecommon.HexToAddress("0x1a6766342142A9BCa7151b7714f811999B7CfeCA"),
	"polygon-testnet":     ecommon.HexToAddress("0xbc77D44223a75194eab5006De96cd1EBa95dA374"),
}

var (
	SupportedContracts = map[string]struct {
		AbiString, Bin, AbiHash string
		Abi                     abi.ABI
		MethodNames             []string
	}{
		ContractErc20Gateway: {
			AbiString: erc20gateway.Erc20gatewayABI,
			Bin:       erc20gateway.Erc20gatewayBin,
			AbiHash:   utils.KeccakHash32(erc20gateway.Erc20gatewayBin),
		},
	}
)

// init initializes variables used throughout this package.
func init() {
	// 1. Initializes abi fields for SupportedContracts
	if entry, ok := SupportedContracts[ContractErc20Gateway]; ok {
		entry.Abi, _ = abi.JSON(strings.NewReader(SupportedContracts[ContractErc20Gateway].AbiString))
		SupportedContracts[ContractErc20Gateway] = entry
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
