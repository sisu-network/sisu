// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dummy

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

// DummyABI is the input ABI used to generate the binding from.
const DummyABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"getName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newName\",\"type\":\"string\"}],\"name\":\"setName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// DummyBin is the compiled bytecode used for deploying new contracts.
var DummyBin = "0x608060405234801561001057600080fd5b506040518060400160405280600581526020017f48656c6c6f0000000000000000000000000000000000000000000000000000008152506000908051906020019061005c929190610062565b50610166565b82805461006e90610105565b90600052602060002090601f01602090048101928261009057600085556100d7565b82601f106100a957805160ff19168380011785556100d7565b828001600101855582156100d7579182015b828111156100d65782518255916020019190600101906100bb565b5b5090506100e491906100e8565b5090565b5b808211156101015760008160009055506001016100e9565b5090565b6000600282049050600182168061011d57607f821691505b6020821081141561013157610130610137565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b61045e806101756000396000f3fe608060405234801561001057600080fd5b50600436106100365760003560e01c806317d7de7c1461003b578063c47f002714610059575b600080fd5b610043610075565b60405161005091906102a6565b60405180910390f35b610073600480360381019061006e919061022c565b610107565b005b60606000805461008490610387565b80601f01602080910402602001604051908101604052809291908181526020018280546100b090610387565b80156100fd5780601f106100d2576101008083540402835291602001916100fd565b820191906000526020600020905b8154815290600101906020018083116100e057829003601f168201915b5050505050905090565b806000908051906020019061011d929190610121565b5050565b82805461012d90610387565b90600052602060002090601f01602090048101928261014f5760008555610196565b82601f1061016857805160ff1916838001178555610196565b82800160010185558215610196579182015b8281111561019557825182559160200191906001019061017a565b5b5090506101a391906101a7565b5090565b5b808211156101c05760008160009055506001016101a8565b5090565b60006101d76101d2846102f9565b6102c8565b9050828152602081018484840111156101ef57600080fd5b6101fa848285610345565b509392505050565b600082601f83011261021357600080fd5b81356102238482602086016101c4565b91505092915050565b60006020828403121561023e57600080fd5b600082013567ffffffffffffffff81111561025857600080fd5b61026484828501610202565b91505092915050565b600061027882610329565b6102828185610334565b9350610292818560208601610354565b61029b81610417565b840191505092915050565b600060208201905081810360008301526102c0818461026d565b905092915050565b6000604051905081810181811067ffffffffffffffff821117156102ef576102ee6103e8565b5b8060405250919050565b600067ffffffffffffffff821115610314576103136103e8565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600082825260208201905092915050565b82818337600083830152505050565b60005b83811015610372578082015181840152602081019050610357565b83811115610381576000848401525b50505050565b6000600282049050600182168061039f57607f821691505b602082108114156103b3576103b26103b9565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f830116905091905056fea264697066735822122046c141509fcbce3b156fe2ab841feaa1ed909568bcec463a4beb7518a05beb3764736f6c63430008000033"

// DeployDummy deploys a new Ethereum contract, binding an instance of Dummy to it.
func DeployDummy(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Dummy, error) {
	parsed, err := abi.JSON(strings.NewReader(DummyABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DummyBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Dummy{DummyCaller: DummyCaller{contract: contract}, DummyTransactor: DummyTransactor{contract: contract}, DummyFilterer: DummyFilterer{contract: contract}}, nil
}

// Dummy is an auto generated Go binding around an Ethereum contract.
type Dummy struct {
	DummyCaller     // Read-only binding to the contract
	DummyTransactor // Write-only binding to the contract
	DummyFilterer   // Log filterer for contract events
}

// DummyCaller is an auto generated read-only Go binding around an Ethereum contract.
type DummyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DummyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DummyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DummySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DummySession struct {
	Contract     *Dummy            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DummyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DummyCallerSession struct {
	Contract *DummyCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DummyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DummyTransactorSession struct {
	Contract     *DummyTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DummyRaw is an auto generated low-level Go binding around an Ethereum contract.
type DummyRaw struct {
	Contract *Dummy // Generic contract binding to access the raw methods on
}

// DummyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DummyCallerRaw struct {
	Contract *DummyCaller // Generic read-only contract binding to access the raw methods on
}

// DummyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DummyTransactorRaw struct {
	Contract *DummyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDummy creates a new instance of Dummy, bound to a specific deployed contract.
func NewDummy(address common.Address, backend bind.ContractBackend) (*Dummy, error) {
	contract, err := bindDummy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dummy{DummyCaller: DummyCaller{contract: contract}, DummyTransactor: DummyTransactor{contract: contract}, DummyFilterer: DummyFilterer{contract: contract}}, nil
}

// NewDummyCaller creates a new read-only instance of Dummy, bound to a specific deployed contract.
func NewDummyCaller(address common.Address, caller bind.ContractCaller) (*DummyCaller, error) {
	contract, err := bindDummy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DummyCaller{contract: contract}, nil
}

// NewDummyTransactor creates a new write-only instance of Dummy, bound to a specific deployed contract.
func NewDummyTransactor(address common.Address, transactor bind.ContractTransactor) (*DummyTransactor, error) {
	contract, err := bindDummy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DummyTransactor{contract: contract}, nil
}

// NewDummyFilterer creates a new log filterer instance of Dummy, bound to a specific deployed contract.
func NewDummyFilterer(address common.Address, filterer bind.ContractFilterer) (*DummyFilterer, error) {
	contract, err := bindDummy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DummyFilterer{contract: contract}, nil
}

// bindDummy binds a generic wrapper to an already deployed contract.
func bindDummy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DummyABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dummy *DummyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dummy.Contract.DummyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dummy *DummyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dummy.Contract.DummyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dummy *DummyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dummy.Contract.DummyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dummy *DummyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dummy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dummy *DummyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dummy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dummy *DummyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dummy.Contract.contract.Transact(opts, method, params...)
}

// GetName is a free data retrieval call binding the contract method 0x17d7de7c.
//
// Solidity: function getName() view returns(string)
func (_Dummy *DummyCaller) GetName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Dummy.contract.Call(opts, &out, "getName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetName is a free data retrieval call binding the contract method 0x17d7de7c.
//
// Solidity: function getName() view returns(string)
func (_Dummy *DummySession) GetName() (string, error) {
	return _Dummy.Contract.GetName(&_Dummy.CallOpts)
}

// GetName is a free data retrieval call binding the contract method 0x17d7de7c.
//
// Solidity: function getName() view returns(string)
func (_Dummy *DummyCallerSession) GetName() (string, error) {
	return _Dummy.Contract.GetName(&_Dummy.CallOpts)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string newName) returns()
func (_Dummy *DummyTransactor) SetName(opts *bind.TransactOpts, newName string) (*types.Transaction, error) {
	return _Dummy.contract.Transact(opts, "setName", newName)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string newName) returns()
func (_Dummy *DummySession) SetName(newName string) (*types.Transaction, error) {
	return _Dummy.Contract.SetName(&_Dummy.TransactOpts, newName)
}

// SetName is a paid mutator transaction binding the contract method 0xc47f0027.
//
// Solidity: function setName(string newName) returns()
func (_Dummy *DummyTransactorSession) SetName(newName string) (*types.Transaction, error) {
	return _Dummy.Contract.SetName(&_Dummy.TransactOpts, newName)
}
