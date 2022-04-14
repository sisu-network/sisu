// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package sampleerc20

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Sampleerc20ABI is the input ABI used to generate the binding from.
const Sampleerc20ABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Sampleerc20Bin is the compiled bytecode used for deploying new contracts.
var Sampleerc20Bin = "0x60806040523480156200001157600080fd5b5060405162001a7138038062001a71833981810160405281019062000037919062000461565b818181600390805190602001906200005192919062000214565b5080600490805190602001906200006a92919062000214565b505050620000893369d3c21bcecceda10000006200009160201b60201c565b505062000692565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141562000104576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401620000fb9062000547565b60405180910390fd5b62000118600083836200020a60201b60201c565b80600260008282546200012c9190620005a2565b92505081905550806000808473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254620001839190620005a2565b925050819055508173ffffffffffffffffffffffffffffffffffffffff16600073ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef83604051620001ea919062000610565b60405180910390a362000206600083836200020f60201b60201c565b5050565b505050565b505050565b82805462000222906200065c565b90600052602060002090601f01602090048101928262000246576000855562000292565b82601f106200026157805160ff191683800117855562000292565b8280016001018555821562000292579182015b828111156200029157825182559160200191906001019062000274565b5b509050620002a19190620002a5565b5090565b5b80821115620002c0576000816000905550600101620002a6565b5090565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200032d82620002e2565b810181811067ffffffffffffffff821117156200034f576200034e620002f3565b5b80604052505050565b600062000364620002c4565b905062000372828262000322565b919050565b600067ffffffffffffffff821115620003955762000394620002f3565b5b620003a082620002e2565b9050602081019050919050565b60005b83811015620003cd578082015181840152602081019050620003b0565b83811115620003dd576000848401525b50505050565b6000620003fa620003f48462000377565b62000358565b905082815260208101848484011115620004195762000418620002dd565b5b62000426848285620003ad565b509392505050565b600082601f830112620004465762000445620002d8565b5b815162000458848260208601620003e3565b91505092915050565b600080604083850312156200047b576200047a620002ce565b5b600083015167ffffffffffffffff8111156200049c576200049b620002d3565b5b620004aa858286016200042e565b925050602083015167ffffffffffffffff811115620004ce57620004cd620002d3565b5b620004dc858286016200042e565b9150509250929050565b600082825260208201905092915050565b7f45524332303a206d696e7420746f20746865207a65726f206164647265737300600082015250565b60006200052f601f83620004e6565b91506200053c82620004f7565b602082019050919050565b60006020820190508181036000830152620005628162000520565b9050919050565b6000819050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000620005af8262000569565b9150620005bc8362000569565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115620005f457620005f362000573565b5b828201905092915050565b6200060a8162000569565b82525050565b6000602082019050620006276000830184620005ff565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b600060028204905060018216806200067557607f821691505b602082108114156200068c576200068b6200062d565b5b50919050565b6113cf80620006a26000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80633950935111610071578063395093511461016857806370a082311461019857806395d89b41146101c8578063a457c2d7146101e6578063a9059cbb14610216578063dd62ed3e14610246576100a9565b806306fdde03146100ae578063095ea7b3146100cc57806318160ddd146100fc57806323b872dd1461011a578063313ce5671461014a575b600080fd5b6100b6610276565b6040516100c39190610c3e565b60405180910390f35b6100e660048036038101906100e19190610cf9565b610308565b6040516100f39190610d54565b60405180910390f35b610104610326565b6040516101119190610d7e565b60405180910390f35b610134600480360381019061012f9190610d99565b610330565b6040516101419190610d54565b60405180910390f35b610152610428565b60405161015f9190610e08565b60405180910390f35b610182600480360381019061017d9190610cf9565b610431565b60405161018f9190610d54565b60405180910390f35b6101b260048036038101906101ad9190610e23565b6104dd565b6040516101bf9190610d7e565b60405180910390f35b6101d0610525565b6040516101dd9190610c3e565b60405180910390f35b61020060048036038101906101fb9190610cf9565b6105b7565b60405161020d9190610d54565b60405180910390f35b610230600480360381019061022b9190610cf9565b6106a2565b60405161023d9190610d54565b60405180910390f35b610260600480360381019061025b9190610e50565b6106c0565b60405161026d9190610d7e565b60405180910390f35b60606003805461028590610ebf565b80601f01602080910402602001604051908101604052809291908181526020018280546102b190610ebf565b80156102fe5780601f106102d3576101008083540402835291602001916102fe565b820191906000526020600020905b8154815290600101906020018083116102e157829003601f168201915b5050505050905090565b600061031c610315610747565b848461074f565b6001905092915050565b6000600254905090565b600061033d84848461091a565b6000600160008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000610388610747565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811015610408576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103ff90610f63565b60405180910390fd5b61041c85610414610747565b85840361074f565b60019150509392505050565b60006012905090565b60006104d361043e610747565b84846001600061044c610747565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008873ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020546104ce9190610fb2565b61074f565b6001905092915050565b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020549050919050565b60606004805461053490610ebf565b80601f016020809104026020016040519081016040528092919081815260200182805461056090610ebf565b80156105ad5780601f10610582576101008083540402835291602001916105ad565b820191906000526020600020905b81548152906001019060200180831161059057829003601f168201915b5050505050905090565b600080600160006105c6610747565b73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905082811015610683576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067a9061107a565b60405180910390fd5b61069761068e610747565b8585840361074f565b600191505092915050565b60006106b66106af610747565b848461091a565b6001905092915050565b6000600160008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b600033905090565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff1614156107bf576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107b69061110c565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16141561082f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108269061119e565b60405180910390fd5b80600160008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055508173ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9258360405161090d9190610d7e565b60405180910390a3505050565b600073ffffffffffffffffffffffffffffffffffffffff168373ffffffffffffffffffffffffffffffffffffffff16141561098a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161098190611230565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156109fa576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109f1906112c2565b60405180910390fd5b610a05838383610b9b565b60008060008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905081811015610a8b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a8290611354565b60405180910390fd5b8181036000808673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002081905550816000808573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b1e9190610fb2565b925050819055508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610b829190610d7e565b60405180910390a3610b95848484610ba0565b50505050565b505050565b505050565b600081519050919050565b600082825260208201905092915050565b60005b83811015610bdf578082015181840152602081019050610bc4565b83811115610bee576000848401525b50505050565b6000601f19601f8301169050919050565b6000610c1082610ba5565b610c1a8185610bb0565b9350610c2a818560208601610bc1565b610c3381610bf4565b840191505092915050565b60006020820190508181036000830152610c588184610c05565b905092915050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610c9082610c65565b9050919050565b610ca081610c85565b8114610cab57600080fd5b50565b600081359050610cbd81610c97565b92915050565b6000819050919050565b610cd681610cc3565b8114610ce157600080fd5b50565b600081359050610cf381610ccd565b92915050565b60008060408385031215610d1057610d0f610c60565b5b6000610d1e85828601610cae565b9250506020610d2f85828601610ce4565b9150509250929050565b60008115159050919050565b610d4e81610d39565b82525050565b6000602082019050610d696000830184610d45565b92915050565b610d7881610cc3565b82525050565b6000602082019050610d936000830184610d6f565b92915050565b600080600060608486031215610db257610db1610c60565b5b6000610dc086828701610cae565b9350506020610dd186828701610cae565b9250506040610de286828701610ce4565b9150509250925092565b600060ff82169050919050565b610e0281610dec565b82525050565b6000602082019050610e1d6000830184610df9565b92915050565b600060208284031215610e3957610e38610c60565b5b6000610e4784828501610cae565b91505092915050565b60008060408385031215610e6757610e66610c60565b5b6000610e7585828601610cae565b9250506020610e8685828601610cae565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b60006002820490506001821680610ed757607f821691505b60208210811415610eeb57610eea610e90565b5b50919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206160008201527f6c6c6f77616e6365000000000000000000000000000000000000000000000000602082015250565b6000610f4d602883610bb0565b9150610f5882610ef1565b604082019050919050565b60006020820190508181036000830152610f7c81610f40565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000610fbd82610cc3565b9150610fc883610cc3565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115610ffd57610ffc610f83565b5b828201905092915050565b7f45524332303a2064656372656173656420616c6c6f77616e63652062656c6f7760008201527f207a65726f000000000000000000000000000000000000000000000000000000602082015250565b6000611064602583610bb0565b915061106f82611008565b604082019050919050565b6000602082019050818103600083015261109381611057565b9050919050565b7f45524332303a20617070726f76652066726f6d20746865207a65726f2061646460008201527f7265737300000000000000000000000000000000000000000000000000000000602082015250565b60006110f6602483610bb0565b91506111018261109a565b604082019050919050565b60006020820190508181036000830152611125816110e9565b9050919050565b7f45524332303a20617070726f766520746f20746865207a65726f20616464726560008201527f7373000000000000000000000000000000000000000000000000000000000000602082015250565b6000611188602283610bb0565b91506111938261112c565b604082019050919050565b600060208201905081810360008301526111b78161117b565b9050919050565b7f45524332303a207472616e736665722066726f6d20746865207a65726f20616460008201527f6472657373000000000000000000000000000000000000000000000000000000602082015250565b600061121a602583610bb0565b9150611225826111be565b604082019050919050565b600060208201905081810360008301526112498161120d565b9050919050565b7f45524332303a207472616e7366657220746f20746865207a65726f206164647260008201527f6573730000000000000000000000000000000000000000000000000000000000602082015250565b60006112ac602383610bb0565b91506112b782611250565b604082019050919050565b600060208201905081810360008301526112db8161129f565b9050919050565b7f45524332303a207472616e7366657220616d6f756e742065786365656473206260008201527f616c616e63650000000000000000000000000000000000000000000000000000602082015250565b600061133e602683610bb0565b9150611349826112e2565b604082019050919050565b6000602082019050818103600083015261136d81611331565b905091905056fea2646970667358221220b8d1867afc715dd63ae1829e4c1c02578a79bd14cf005b09755a424375dbdf4c64736f6c637827302e382e31322d646576656c6f702e323032322e322e382b636f6d6d69742e35633362636236630058"

// DeploySampleerc20 deploys a new Ethereum contract, binding an instance of Sampleerc20 to it.
func DeploySampleerc20(auth *bind.TransactOpts, backend bind.ContractBackend, name string, symbol string) (common.Address, *types.Transaction, *Sampleerc20, error) {
	parsed, err := abi.JSON(strings.NewReader(Sampleerc20ABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Sampleerc20Bin), backend, name, symbol)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Sampleerc20{Sampleerc20Caller: Sampleerc20Caller{contract: contract}, Sampleerc20Transactor: Sampleerc20Transactor{contract: contract}, Sampleerc20Filterer: Sampleerc20Filterer{contract: contract}}, nil
}

// Sampleerc20 is an auto generated Go binding around an Ethereum contract.
type Sampleerc20 struct {
	Sampleerc20Caller     // Read-only binding to the contract
	Sampleerc20Transactor // Write-only binding to the contract
	Sampleerc20Filterer   // Log filterer for contract events
}

// Sampleerc20Caller is an auto generated read-only Go binding around an Ethereum contract.
type Sampleerc20Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Sampleerc20Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Sampleerc20Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Sampleerc20Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Sampleerc20Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Sampleerc20Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Sampleerc20Session struct {
	Contract     *Sampleerc20      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Sampleerc20CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Sampleerc20CallerSession struct {
	Contract *Sampleerc20Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// Sampleerc20TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Sampleerc20TransactorSession struct {
	Contract     *Sampleerc20Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// Sampleerc20Raw is an auto generated low-level Go binding around an Ethereum contract.
type Sampleerc20Raw struct {
	Contract *Sampleerc20 // Generic contract binding to access the raw methods on
}

// Sampleerc20CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Sampleerc20CallerRaw struct {
	Contract *Sampleerc20Caller // Generic read-only contract binding to access the raw methods on
}

// Sampleerc20TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Sampleerc20TransactorRaw struct {
	Contract *Sampleerc20Transactor // Generic write-only contract binding to access the raw methods on
}

// NewSampleerc20 creates a new instance of Sampleerc20, bound to a specific deployed contract.
func NewSampleerc20(address common.Address, backend bind.ContractBackend) (*Sampleerc20, error) {
	contract, err := bindSampleerc20(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20{Sampleerc20Caller: Sampleerc20Caller{contract: contract}, Sampleerc20Transactor: Sampleerc20Transactor{contract: contract}, Sampleerc20Filterer: Sampleerc20Filterer{contract: contract}}, nil
}

// NewSampleerc20Caller creates a new read-only instance of Sampleerc20, bound to a specific deployed contract.
func NewSampleerc20Caller(address common.Address, caller bind.ContractCaller) (*Sampleerc20Caller, error) {
	contract, err := bindSampleerc20(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20Caller{contract: contract}, nil
}

// NewSampleerc20Transactor creates a new write-only instance of Sampleerc20, bound to a specific deployed contract.
func NewSampleerc20Transactor(address common.Address, transactor bind.ContractTransactor) (*Sampleerc20Transactor, error) {
	contract, err := bindSampleerc20(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20Transactor{contract: contract}, nil
}

// NewSampleerc20Filterer creates a new log filterer instance of Sampleerc20, bound to a specific deployed contract.
func NewSampleerc20Filterer(address common.Address, filterer bind.ContractFilterer) (*Sampleerc20Filterer, error) {
	contract, err := bindSampleerc20(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20Filterer{contract: contract}, nil
}

// bindSampleerc20 binds a generic wrapper to an already deployed contract.
func bindSampleerc20(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Sampleerc20ABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sampleerc20 *Sampleerc20Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sampleerc20.Contract.Sampleerc20Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sampleerc20 *Sampleerc20Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Sampleerc20Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sampleerc20 *Sampleerc20Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Sampleerc20Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sampleerc20 *Sampleerc20CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sampleerc20.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sampleerc20 *Sampleerc20TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sampleerc20.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sampleerc20 *Sampleerc20TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sampleerc20.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Sampleerc20 *Sampleerc20Caller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Sampleerc20 *Sampleerc20Session) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Sampleerc20.Contract.Allowance(&_Sampleerc20.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_Sampleerc20 *Sampleerc20CallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _Sampleerc20.Contract.Allowance(&_Sampleerc20.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Sampleerc20 *Sampleerc20Caller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Sampleerc20 *Sampleerc20Session) BalanceOf(account common.Address) (*big.Int, error) {
	return _Sampleerc20.Contract.BalanceOf(&_Sampleerc20.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_Sampleerc20 *Sampleerc20CallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _Sampleerc20.Contract.BalanceOf(&_Sampleerc20.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Sampleerc20 *Sampleerc20Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Sampleerc20 *Sampleerc20Session) Decimals() (uint8, error) {
	return _Sampleerc20.Contract.Decimals(&_Sampleerc20.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_Sampleerc20 *Sampleerc20CallerSession) Decimals() (uint8, error) {
	return _Sampleerc20.Contract.Decimals(&_Sampleerc20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Sampleerc20 *Sampleerc20Caller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Sampleerc20 *Sampleerc20Session) Name() (string, error) {
	return _Sampleerc20.Contract.Name(&_Sampleerc20.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_Sampleerc20 *Sampleerc20CallerSession) Name() (string, error) {
	return _Sampleerc20.Contract.Name(&_Sampleerc20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Sampleerc20 *Sampleerc20Caller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Sampleerc20 *Sampleerc20Session) Symbol() (string, error) {
	return _Sampleerc20.Contract.Symbol(&_Sampleerc20.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_Sampleerc20 *Sampleerc20CallerSession) Symbol() (string, error) {
	return _Sampleerc20.Contract.Symbol(&_Sampleerc20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Sampleerc20 *Sampleerc20Caller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Sampleerc20.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Sampleerc20 *Sampleerc20Session) TotalSupply() (*big.Int, error) {
	return _Sampleerc20.Contract.TotalSupply(&_Sampleerc20.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_Sampleerc20 *Sampleerc20CallerSession) TotalSupply() (*big.Int, error) {
	return _Sampleerc20.Contract.TotalSupply(&_Sampleerc20.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Transactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Session) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Approve(&_Sampleerc20.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20TransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Approve(&_Sampleerc20.TransactOpts, spender, amount)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20Transactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20Session) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.DecreaseAllowance(&_Sampleerc20.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20TransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.DecreaseAllowance(&_Sampleerc20.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20Transactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20Session) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.IncreaseAllowance(&_Sampleerc20.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_Sampleerc20 *Sampleerc20TransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.IncreaseAllowance(&_Sampleerc20.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Transactor) Transfer(opts *bind.TransactOpts, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.contract.Transact(opts, "transfer", recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Session) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Transfer(&_Sampleerc20.TransactOpts, recipient, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20TransactorSession) Transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.Transfer(&_Sampleerc20.TransactOpts, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Transactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.contract.Transact(opts, "transferFrom", sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20Session) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.TransferFrom(&_Sampleerc20.TransactOpts, sender, recipient, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address sender, address recipient, uint256 amount) returns(bool)
func (_Sampleerc20 *Sampleerc20TransactorSession) TransferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Sampleerc20.Contract.TransferFrom(&_Sampleerc20.TransactOpts, sender, recipient, amount)
}

// Sampleerc20ApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Sampleerc20 contract.
type Sampleerc20ApprovalIterator struct {
	Event *Sampleerc20Approval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Sampleerc20ApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Sampleerc20Approval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Sampleerc20Approval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Sampleerc20ApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Sampleerc20ApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Sampleerc20Approval represents a Approval event raised by the Sampleerc20 contract.
type Sampleerc20Approval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*Sampleerc20ApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Sampleerc20.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20ApprovalIterator{contract: _Sampleerc20.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *Sampleerc20Approval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _Sampleerc20.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Sampleerc20Approval)
				if err := _Sampleerc20.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) ParseApproval(log types.Log) (*Sampleerc20Approval, error) {
	event := new(Sampleerc20Approval)
	if err := _Sampleerc20.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Sampleerc20TransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Sampleerc20 contract.
type Sampleerc20TransferIterator struct {
	Event *Sampleerc20Transfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *Sampleerc20TransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Sampleerc20Transfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(Sampleerc20Transfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *Sampleerc20TransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Sampleerc20TransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Sampleerc20Transfer represents a Transfer event raised by the Sampleerc20 contract.
type Sampleerc20Transfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*Sampleerc20TransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Sampleerc20.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &Sampleerc20TransferIterator{contract: _Sampleerc20.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *Sampleerc20Transfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _Sampleerc20.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Sampleerc20Transfer)
				if err := _Sampleerc20.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_Sampleerc20 *Sampleerc20Filterer) ParseTransfer(log types.Log) (*Sampleerc20Transfer, error) {
	event := new(Sampleerc20Transfer)
	if err := _Sampleerc20.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
