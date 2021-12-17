// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20gw

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// Erc20gwMetaData contains all meta data concerning the Erc20gw contract.
var Erc20gwMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_supportedChains\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"AddSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"RemoveSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"destChain\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferInEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"reipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOutEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"AddSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PauseGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"RemoveSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ResumeGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"destChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"supportedChains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// Erc20gwABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc20gwMetaData.ABI instead.
var Erc20gwABI = Erc20gwMetaData.ABI

// Erc20gw is an auto generated Go binding around an Ethereum contract.
type Erc20gw struct {
	Erc20gwCaller     // Read-only binding to the contract
	Erc20gwTransactor // Write-only binding to the contract
	Erc20gwFilterer   // Log filterer for contract events
}

// Erc20gwCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20gwCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gwTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20gwTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gwFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20gwFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gwSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20gwSession struct {
	Contract     *Erc20gw          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20gwCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20gwCallerSession struct {
	Contract *Erc20gwCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// Erc20gwTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20gwTransactorSession struct {
	Contract     *Erc20gwTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// Erc20gwRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20gwRaw struct {
	Contract *Erc20gw // Generic contract binding to access the raw methods on
}

// Erc20gwCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20gwCallerRaw struct {
	Contract *Erc20gwCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20gwTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20gwTransactorRaw struct {
	Contract *Erc20gwTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20gw creates a new instance of Erc20gw, bound to a specific deployed contract.
func NewErc20gw(address common.Address, backend bind.ContractBackend) (*Erc20gw, error) {
	contract, err := bindErc20gw(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20gw{Erc20gwCaller: Erc20gwCaller{contract: contract}, Erc20gwTransactor: Erc20gwTransactor{contract: contract}, Erc20gwFilterer: Erc20gwFilterer{contract: contract}}, nil
}

// NewErc20gwCaller creates a new read-only instance of Erc20gw, bound to a specific deployed contract.
func NewErc20gwCaller(address common.Address, caller bind.ContractCaller) (*Erc20gwCaller, error) {
	contract, err := bindErc20gw(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20gwCaller{contract: contract}, nil
}

// NewErc20gwTransactor creates a new write-only instance of Erc20gw, bound to a specific deployed contract.
func NewErc20gwTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20gwTransactor, error) {
	contract, err := bindErc20gw(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20gwTransactor{contract: contract}, nil
}

// NewErc20gwFilterer creates a new log filterer instance of Erc20gw, bound to a specific deployed contract.
func NewErc20gwFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc20gwFilterer, error) {
	contract, err := bindErc20gw(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20gwFilterer{contract: contract}, nil
}

// bindErc20gw binds a generic wrapper to an already deployed contract.
func bindErc20gw(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20gwABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20gw *Erc20gwRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20gw.Contract.Erc20gwCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20gw *Erc20gwRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gw.Contract.Erc20gwTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20gw *Erc20gwRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20gw.Contract.Erc20gwTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20gw *Erc20gwCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20gw.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20gw *Erc20gwTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gw.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20gw *Erc20gwTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20gw.Contract.contract.Transact(opts, method, params...)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gw *Erc20gwCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc20gw.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gw *Erc20gwSession) Owner() (common.Address, error) {
	return _Erc20gw.Contract.Owner(&_Erc20gw.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gw *Erc20gwCallerSession) Owner() (common.Address, error) {
	return _Erc20gw.Contract.Owner(&_Erc20gw.CallOpts)
}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gw *Erc20gwCaller) Pause(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Erc20gw.contract.Call(opts, &out, "pause")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gw *Erc20gwSession) Pause() (bool, error) {
	return _Erc20gw.Contract.Pause(&_Erc20gw.CallOpts)
}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gw *Erc20gwCallerSession) Pause() (bool, error) {
	return _Erc20gw.Contract.Pause(&_Erc20gw.CallOpts)
}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gw *Erc20gwCaller) SupportedChains(opts *bind.CallOpts, arg0 string) (bool, error) {
	var out []interface{}
	err := _Erc20gw.contract.Call(opts, &out, "supportedChains", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gw *Erc20gwSession) SupportedChains(arg0 string) (bool, error) {
	return _Erc20gw.Contract.SupportedChains(&_Erc20gw.CallOpts, arg0)
}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gw *Erc20gwCallerSession) SupportedChains(arg0 string) (bool, error) {
	return _Erc20gw.Contract.SupportedChains(&_Erc20gw.CallOpts, arg0)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xfc69a67a.
//
// Solidity: function AddSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwTransactor) AddSupportedChain(opts *bind.TransactOpts, chain string) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "AddSupportedChain", chain)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xfc69a67a.
//
// Solidity: function AddSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwSession) AddSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gw.Contract.AddSupportedChain(&_Erc20gw.TransactOpts, chain)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xfc69a67a.
//
// Solidity: function AddSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwTransactorSession) AddSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gw.Contract.AddSupportedChain(&_Erc20gw.TransactOpts, chain)
}

// PauseGateway is a paid mutator transaction binding the contract method 0xca569dbf.
//
// Solidity: function PauseGateway() returns()
func (_Erc20gw *Erc20gwTransactor) PauseGateway(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "PauseGateway")
}

// PauseGateway is a paid mutator transaction binding the contract method 0xca569dbf.
//
// Solidity: function PauseGateway() returns()
func (_Erc20gw *Erc20gwSession) PauseGateway() (*types.Transaction, error) {
	return _Erc20gw.Contract.PauseGateway(&_Erc20gw.TransactOpts)
}

// PauseGateway is a paid mutator transaction binding the contract method 0xca569dbf.
//
// Solidity: function PauseGateway() returns()
func (_Erc20gw *Erc20gwTransactorSession) PauseGateway() (*types.Transaction, error) {
	return _Erc20gw.Contract.PauseGateway(&_Erc20gw.TransactOpts)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x58b67fe1.
//
// Solidity: function RemoveSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwTransactor) RemoveSupportedChain(opts *bind.TransactOpts, chain string) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "RemoveSupportedChain", chain)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x58b67fe1.
//
// Solidity: function RemoveSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwSession) RemoveSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gw.Contract.RemoveSupportedChain(&_Erc20gw.TransactOpts, chain)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x58b67fe1.
//
// Solidity: function RemoveSupportedChain(string chain) returns()
func (_Erc20gw *Erc20gwTransactorSession) RemoveSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gw.Contract.RemoveSupportedChain(&_Erc20gw.TransactOpts, chain)
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xecf62f82.
//
// Solidity: function ResumeGateway() returns()
func (_Erc20gw *Erc20gwTransactor) ResumeGateway(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "ResumeGateway")
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xecf62f82.
//
// Solidity: function ResumeGateway() returns()
func (_Erc20gw *Erc20gwSession) ResumeGateway() (*types.Transaction, error) {
	return _Erc20gw.Contract.ResumeGateway(&_Erc20gw.TransactOpts)
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xecf62f82.
//
// Solidity: function ResumeGateway() returns()
func (_Erc20gw *Erc20gwTransactorSession) ResumeGateway() (*types.Transaction, error) {
	return _Erc20gw.Contract.ResumeGateway(&_Erc20gw.TransactOpts)
}

// TransferOut is a paid mutator transaction binding the contract method 0x5d2c285d.
//
// Solidity: function TransferOut(address _token, address recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactor) TransferOut(opts *bind.TransactOpts, _token common.Address, recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "TransferOut", _token, recipient, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x5d2c285d.
//
// Solidity: function TransferOut(address _token, address recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwSession) TransferOut(_token common.Address, recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut(&_Erc20gw.TransactOpts, _token, recipient, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x5d2c285d.
//
// Solidity: function TransferOut(address _token, address recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactorSession) TransferOut(_token common.Address, recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut(&_Erc20gw.TransactOpts, _token, recipient, _amount)
}

// TransferOut0 is a paid mutator transaction binding the contract method 0xaa1e756e.
//
// Solidity: function TransferOut(string destChain, address _token, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactor) TransferOut0(opts *bind.TransactOpts, destChain string, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "TransferOut0", destChain, _token, _amount)
}

// TransferOut0 is a paid mutator transaction binding the contract method 0xaa1e756e.
//
// Solidity: function TransferOut(string destChain, address _token, uint256 _amount) returns()
func (_Erc20gw *Erc20gwSession) TransferOut0(destChain string, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut0(&_Erc20gw.TransactOpts, destChain, _token, _amount)
}

// TransferOut0 is a paid mutator transaction binding the contract method 0xaa1e756e.
//
// Solidity: function TransferOut(string destChain, address _token, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactorSession) TransferOut0(destChain string, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut0(&_Erc20gw.TransactOpts, destChain, _token, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gw *Erc20gwTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gw *Erc20gwSession) RenounceOwnership() (*types.Transaction, error) {
	return _Erc20gw.Contract.RenounceOwnership(&_Erc20gw.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gw *Erc20gwTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Erc20gw.Contract.RenounceOwnership(&_Erc20gw.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gw *Erc20gwTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gw *Erc20gwSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOwnership(&_Erc20gw.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gw *Erc20gwTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOwnership(&_Erc20gw.TransactOpts, newOwner)
}

// Erc20gwAddSupportedChainEventIterator is returned from FilterAddSupportedChainEvent and is used to iterate over the raw logs and unpacked data for AddSupportedChainEvent events raised by the Erc20gw contract.
type Erc20gwAddSupportedChainEventIterator struct {
	Event *Erc20gwAddSupportedChainEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gwAddSupportedChainEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gwAddSupportedChainEvent)
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
		it.Event = new(Erc20gwAddSupportedChainEvent)
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
func (it *Erc20gwAddSupportedChainEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gwAddSupportedChainEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gwAddSupportedChainEvent represents a AddSupportedChainEvent event raised by the Erc20gw contract.
type Erc20gwAddSupportedChainEvent struct {
	Chain common.Hash
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddSupportedChainEvent is a free log retrieval operation binding the contract event 0x7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a7.
//
// Solidity: event AddSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) FilterAddSupportedChainEvent(opts *bind.FilterOpts, chain []string) (*Erc20gwAddSupportedChainEventIterator, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "AddSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwAddSupportedChainEventIterator{contract: _Erc20gw.contract, event: "AddSupportedChainEvent", logs: logs, sub: sub}, nil
}

// WatchAddSupportedChainEvent is a free log subscription operation binding the contract event 0x7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a7.
//
// Solidity: event AddSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) WatchAddSupportedChainEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwAddSupportedChainEvent, chain []string) (event.Subscription, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "AddSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gwAddSupportedChainEvent)
				if err := _Erc20gw.contract.UnpackLog(event, "AddSupportedChainEvent", log); err != nil {
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

// ParseAddSupportedChainEvent is a log parse operation binding the contract event 0x7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a7.
//
// Solidity: event AddSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) ParseAddSupportedChainEvent(log types.Log) (*Erc20gwAddSupportedChainEvent, error) {
	event := new(Erc20gwAddSupportedChainEvent)
	if err := _Erc20gw.contract.UnpackLog(event, "AddSupportedChainEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gwOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Erc20gw contract.
type Erc20gwOwnershipTransferredIterator struct {
	Event *Erc20gwOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Erc20gwOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gwOwnershipTransferred)
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
		it.Event = new(Erc20gwOwnershipTransferred)
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
func (it *Erc20gwOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gwOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gwOwnershipTransferred represents a OwnershipTransferred event raised by the Erc20gw contract.
type Erc20gwOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Erc20gw *Erc20gwFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Erc20gwOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwOwnershipTransferredIterator{contract: _Erc20gw.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Erc20gw *Erc20gwFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Erc20gwOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gwOwnershipTransferred)
				if err := _Erc20gw.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Erc20gw *Erc20gwFilterer) ParseOwnershipTransferred(log types.Log) (*Erc20gwOwnershipTransferred, error) {
	event := new(Erc20gwOwnershipTransferred)
	if err := _Erc20gw.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gwRemoveSupportedChainEventIterator is returned from FilterRemoveSupportedChainEvent and is used to iterate over the raw logs and unpacked data for RemoveSupportedChainEvent events raised by the Erc20gw contract.
type Erc20gwRemoveSupportedChainEventIterator struct {
	Event *Erc20gwRemoveSupportedChainEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gwRemoveSupportedChainEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gwRemoveSupportedChainEvent)
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
		it.Event = new(Erc20gwRemoveSupportedChainEvent)
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
func (it *Erc20gwRemoveSupportedChainEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gwRemoveSupportedChainEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gwRemoveSupportedChainEvent represents a RemoveSupportedChainEvent event raised by the Erc20gw contract.
type Erc20gwRemoveSupportedChainEvent struct {
	Chain common.Hash
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRemoveSupportedChainEvent is a free log retrieval operation binding the contract event 0xf300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca635.
//
// Solidity: event RemoveSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) FilterRemoveSupportedChainEvent(opts *bind.FilterOpts, chain []string) (*Erc20gwRemoveSupportedChainEventIterator, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "RemoveSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwRemoveSupportedChainEventIterator{contract: _Erc20gw.contract, event: "RemoveSupportedChainEvent", logs: logs, sub: sub}, nil
}

// WatchRemoveSupportedChainEvent is a free log subscription operation binding the contract event 0xf300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca635.
//
// Solidity: event RemoveSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) WatchRemoveSupportedChainEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwRemoveSupportedChainEvent, chain []string) (event.Subscription, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "RemoveSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gwRemoveSupportedChainEvent)
				if err := _Erc20gw.contract.UnpackLog(event, "RemoveSupportedChainEvent", log); err != nil {
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

// ParseRemoveSupportedChainEvent is a log parse operation binding the contract event 0xf300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca635.
//
// Solidity: event RemoveSupportedChainEvent(string indexed chain)
func (_Erc20gw *Erc20gwFilterer) ParseRemoveSupportedChainEvent(log types.Log) (*Erc20gwRemoveSupportedChainEvent, error) {
	event := new(Erc20gwRemoveSupportedChainEvent)
	if err := _Erc20gw.contract.UnpackLog(event, "RemoveSupportedChainEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gwTransferInEventIterator is returned from FilterTransferInEvent and is used to iterate over the raw logs and unpacked data for TransferInEvent events raised by the Erc20gw contract.
type Erc20gwTransferInEventIterator struct {
	Event *Erc20gwTransferInEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gwTransferInEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gwTransferInEvent)
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
		it.Event = new(Erc20gwTransferInEvent)
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
func (it *Erc20gwTransferInEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gwTransferInEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gwTransferInEvent represents a TransferInEvent event raised by the Erc20gw contract.
type Erc20gwTransferInEvent struct {
	DestChain common.Hash
	Token     common.Address
	Sender    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInEvent is a free log retrieval operation binding the contract event 0x28b8d4ba55a81ea694c0446c5914f4b2e4762e943d9fb9bae9e525ac9e788850.
//
// Solidity: event TransferInEvent(string indexed destChain, address indexed token, address indexed sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) FilterTransferInEvent(opts *bind.FilterOpts, destChain []string, token []common.Address, sender []common.Address) (*Erc20gwTransferInEventIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "TransferInEvent", destChainRule, tokenRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwTransferInEventIterator{contract: _Erc20gw.contract, event: "TransferInEvent", logs: logs, sub: sub}, nil
}

// WatchTransferInEvent is a free log subscription operation binding the contract event 0x28b8d4ba55a81ea694c0446c5914f4b2e4762e943d9fb9bae9e525ac9e788850.
//
// Solidity: event TransferInEvent(string indexed destChain, address indexed token, address indexed sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) WatchTransferInEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwTransferInEvent, destChain []string, token []common.Address, sender []common.Address) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "TransferInEvent", destChainRule, tokenRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gwTransferInEvent)
				if err := _Erc20gw.contract.UnpackLog(event, "TransferInEvent", log); err != nil {
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

// ParseTransferInEvent is a log parse operation binding the contract event 0x28b8d4ba55a81ea694c0446c5914f4b2e4762e943d9fb9bae9e525ac9e788850.
//
// Solidity: event TransferInEvent(string indexed destChain, address indexed token, address indexed sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) ParseTransferInEvent(log types.Log) (*Erc20gwTransferInEvent, error) {
	event := new(Erc20gwTransferInEvent)
	if err := _Erc20gw.contract.UnpackLog(event, "TransferInEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gwTransferOutEventIterator is returned from FilterTransferOutEvent and is used to iterate over the raw logs and unpacked data for TransferOutEvent events raised by the Erc20gw contract.
type Erc20gwTransferOutEventIterator struct {
	Event *Erc20gwTransferOutEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gwTransferOutEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gwTransferOutEvent)
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
		it.Event = new(Erc20gwTransferOutEvent)
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
func (it *Erc20gwTransferOutEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gwTransferOutEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gwTransferOutEvent represents a TransferOutEvent event raised by the Erc20gw contract.
type Erc20gwTransferOutEvent struct {
	Token    common.Address
	Reipient common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferOutEvent is a free log retrieval operation binding the contract event 0x7ef6d9ca986c0567356edd27e245c895a7c749cc39a887bdc6e2520486dffbb7.
//
// Solidity: event TransferOutEvent(address indexed token, address indexed reipient, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) FilterTransferOutEvent(opts *bind.FilterOpts, token []common.Address, reipient []common.Address) (*Erc20gwTransferOutEventIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var reipientRule []interface{}
	for _, reipientItem := range reipient {
		reipientRule = append(reipientRule, reipientItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "TransferOutEvent", tokenRule, reipientRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwTransferOutEventIterator{contract: _Erc20gw.contract, event: "TransferOutEvent", logs: logs, sub: sub}, nil
}

// WatchTransferOutEvent is a free log subscription operation binding the contract event 0x7ef6d9ca986c0567356edd27e245c895a7c749cc39a887bdc6e2520486dffbb7.
//
// Solidity: event TransferOutEvent(address indexed token, address indexed reipient, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) WatchTransferOutEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwTransferOutEvent, token []common.Address, reipient []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var reipientRule []interface{}
	for _, reipientItem := range reipient {
		reipientRule = append(reipientRule, reipientItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "TransferOutEvent", tokenRule, reipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gwTransferOutEvent)
				if err := _Erc20gw.contract.UnpackLog(event, "TransferOutEvent", log); err != nil {
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

// ParseTransferOutEvent is a log parse operation binding the contract event 0x7ef6d9ca986c0567356edd27e245c895a7c749cc39a887bdc6e2520486dffbb7.
//
// Solidity: event TransferOutEvent(address indexed token, address indexed reipient, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) ParseTransferOutEvent(log types.Log) (*Erc20gwTransferOutEvent, error) {
	event := new(Erc20gwTransferOutEvent)
	if err := _Erc20gw.contract.UnpackLog(event, "TransferOutEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
