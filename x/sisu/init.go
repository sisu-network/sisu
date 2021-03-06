package sisu

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/sisu-network/sisu/contracts/eth/erc20gateway"
	"github.com/sisu-network/sisu/contracts/eth/liquiditypool"
	"github.com/sisu-network/sisu/utils"
)

const (
	ContractErc20Gateway  = "erc20gateway"
	ContractLiquidityPool = "liquidityPool"

	// Methods in gateway smart contract
	MethodTransferIn        = "transferIn"
	MethodTransferOut       = "transferOut"
	MethodPauseGateway      = "pauseGateway"
	MethodResumeGateway     = "resumeGateway"
	MethodTransferOwnership = "transferOwnership"
	MethodSetLiquidAddress  = "setLiquidAddress"

	// Methods in liquidity smart contract
	MethodEmergencyWithdrawFund = "emergencyWithdrawFunds"

	UpdateGasPriceFrequency = 100
)

type ContractInfo struct {
	AbiString, Bin, AbiHash string
	Abi                     abi.ABI
	MethodNames             []string
	IsDeployBySisu          bool
}

var (
	SupportedContracts = map[string]*ContractInfo{
		ContractErc20Gateway: {
			AbiString:      erc20gateway.Erc20gatewayABI,
			Bin:            erc20gateway.Erc20gatewayBin,
			AbiHash:        utils.KeccakHash32(erc20gateway.Erc20gatewayBin),
			IsDeployBySisu: true,
		},
		ContractLiquidityPool: {
			AbiString:      liquiditypool.LiquiditypoolABI,
			Bin:            liquiditypool.LiquiditypoolBin,
			AbiHash:        utils.KeccakHash32(liquiditypool.LiquiditypoolBin),
			IsDeployBySisu: false,
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
