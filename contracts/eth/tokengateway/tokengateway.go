// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package TokenGateway

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

// TokenGatewayABI is the input ABI used to generate the binding from.
const TokenGatewayABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferInAssetOfThisChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferOutFromContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferWithin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newChain\",\"type\":\"string\"}],\"name\":\"addAllowedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"isChainAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInAssetOfThisChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOutFromContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferWithin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// TokenGatewayBin is the compiled bytecode used for deploying new contracts.
var TokenGatewayBin = "0x60806040523480156200001157600080fd5b50604051620017d1380380620017d18339818101604052810190620000379190620001b9565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600190805190602001906200008f92919062000097565b50506200032f565b828054620000a5906200029b565b90600052602060002090601f016020900481019282620000c9576000855562000115565b82601f10620000e457805160ff191683800117855562000115565b8280016001018555821562000115579182015b8281111562000114578251825591602001919060010190620000f7565b5b50905062000124919062000128565b5090565b5b808211156200014357600081600090555060010162000129565b5090565b60006200015e620001588462000232565b620001fe565b9050828152602081018484840111156200017757600080fd5b6200018484828562000265565b509392505050565b600082601f8301126200019e57600080fd5b8151620001b084826020860162000147565b91505092915050565b600060208284031215620001cc57600080fd5b600082015167ffffffffffffffff811115620001e757600080fd5b620001f5848285016200018c565b91505092915050565b6000604051905081810181811067ffffffffffffffff8211171562000228576200022762000300565b5b8060405250919050565b600067ffffffffffffffff82111562000250576200024f62000300565b5b601f19601f8301169050602081019050919050565b60005b838110156200028557808201518184015260208101905062000268565b8381111562000295576000848401525b50505050565b60006002820490506001821680620002b457607f821691505b60208210811415620002cb57620002ca620002d1565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611492806200033f6000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806389f4b3171161006657806389f4b317146101315780639a12bb6014610161578063a6f9dae11461017d578063bf6022be14610199578063d29fde24146101b55761009e565b80630832f134146100a35780636c189db8146100bf5780637dd1f364146100db578063810b9f12146100f7578063893d20e814610113575b600080fd5b6100bd60048036038101906100b89190610e31565b6101e5565b005b6100d960048036038101906100d49190610cba565b6102ef565b005b6100f560048036038101906100f09190610e31565b610531565b005b610111600480360381019061010c9190610e98565b6106ca565b005b61011b61082b565b6040516101289190611039565b60405180910390f35b61014b60048036038101906101469190610ddd565b610854565b60405161015891906111f2565b60405180910390f35b61017b60048036038101906101769190610d9c565b6108ba565b005b61019760048036038101906101929190610c91565b61094c565b005b6101b360048036038101906101ae9190610d09565b6109e7565b005b6101cf60048036038101906101ca9190610d9c565b610bca565b6040516101dc9190611131565b60405180910390f35b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461023d57600080fd5b8060028460405161024e9190611022565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546102a891906112ab565b925050819055507fbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea0948383836040516102e2939291906111b4565b60405180910390a1505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461034757600080fd5b80600460008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561039357600080fd5b60008373ffffffffffffffffffffffffffffffffffffffff163084846040516024016103c193929190611054565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161044b919061100b565b6000604051808303816000865af19150503d8060008114610488576040519150601f19603f3d011682016040523d82523d6000602084013e61048d565b606091505b5050905080156104ee5781600460008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546104e69190611301565b925050819055505b7fc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c384848484604051610523949392919061108b565b60405180910390a150505050565b6000811161053e57600080fd5b8060028460405161054f9190611022565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156105a657600080fd5b806002846040516105b79190611022565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546106119190611301565b92505081905550806002846040516106299190611022565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020600082825461068391906112ab565b925050819055507f8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df48383836040516106bd939291906111b4565b60405180910390a1505050565b600081116106d757600080fd5b806002856040516106e89190611022565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561073f57600080fd5b60038360405161074f9190611022565b908152602001604051809103902060009054906101000a900460ff1661077457600080fd5b806002856040516107859190611022565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107df9190611301565b925050819055507f265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f843385858560405161081d95949392919061114c565b60405180910390a150505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b60006002836040516108669190611022565b908152602001604051809103902060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461091257600080fd5b60016003826040516109249190611022565b908152602001604051809103902060006101000a81548160ff02191690831515021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146109a457600080fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b600081116109f457600080fd5b600383604051610a049190611022565b908152602001604051809103902060009054906101000a900460ff16610a2957600080fd5b60008473ffffffffffffffffffffffffffffffffffffffff16333084604051602401610a5793929190611054565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610ae1919061100b565b6000604051808303816000865af19150503d8060008114610b1e576040519150601f19603f3d011682016040523d82523d6000602084013e610b23565b606091505b505090508015610b845781600460008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610b7c91906112ab565b925050819055505b7f889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e5418585858585604051610bbb9594939291906110d0565b60405180910390a15050505050565b6000600382604051610bdc9190611022565b908152602001604051809103902060009054906101000a900460ff169050919050565b6000610c12610c0d8461123e565b61120d565b905082815260208101848484011115610c2a57600080fd5b610c3584828561137d565b509392505050565b600081359050610c4c8161142e565b92915050565b600082601f830112610c6357600080fd5b8135610c73848260208601610bff565b91505092915050565b600081359050610c8b81611445565b92915050565b600060208284031215610ca357600080fd5b6000610cb184828501610c3d565b91505092915050565b600080600060608486031215610ccf57600080fd5b6000610cdd86828701610c3d565b9350506020610cee86828701610c3d565b9250506040610cff86828701610c7c565b9150509250925092565b60008060008060808587031215610d1f57600080fd5b6000610d2d87828801610c3d565b945050602085013567ffffffffffffffff811115610d4a57600080fd5b610d5687828801610c52565b935050604085013567ffffffffffffffff811115610d7357600080fd5b610d7f87828801610c52565b9250506060610d9087828801610c7c565b91505092959194509250565b600060208284031215610dae57600080fd5b600082013567ffffffffffffffff811115610dc857600080fd5b610dd484828501610c52565b91505092915050565b60008060408385031215610df057600080fd5b600083013567ffffffffffffffff811115610e0a57600080fd5b610e1685828601610c52565b9250506020610e2785828601610c3d565b9150509250929050565b600080600060608486031215610e4657600080fd5b600084013567ffffffffffffffff811115610e6057600080fd5b610e6c86828701610c52565b9350506020610e7d86828701610c3d565b9250506040610e8e86828701610c7c565b9150509250925092565b60008060008060808587031215610eae57600080fd5b600085013567ffffffffffffffff811115610ec857600080fd5b610ed487828801610c52565b945050602085013567ffffffffffffffff811115610ef157600080fd5b610efd87828801610c52565b935050604085013567ffffffffffffffff811115610f1a57600080fd5b610f2687828801610c52565b9250506060610f3787828801610c7c565b91505092959194509250565b610f4c81611335565b82525050565b610f5b81611347565b82525050565b6000610f6c8261126e565b610f768185611284565b9350610f8681856020860161138c565b80840191505092915050565b6000610f9d82611279565b610fa7818561128f565b9350610fb781856020860161138c565b610fc08161141d565b840191505092915050565b6000610fd682611279565b610fe081856112a0565b9350610ff081856020860161138c565b80840191505092915050565b61100581611373565b82525050565b60006110178284610f61565b915081905092915050565b600061102e8284610fcb565b915081905092915050565b600060208201905061104e6000830184610f43565b92915050565b60006060820190506110696000830186610f43565b6110766020830185610f43565b6110836040830184610ffc565b949350505050565b60006080820190506110a06000830187610f43565b6110ad6020830186610f43565b6110ba6040830185610ffc565b6110c76060830184610f52565b95945050505050565b600060a0820190506110e56000830188610f43565b81810360208301526110f78187610f92565b9050818103604083015261110b8186610f92565b905061111a6060830185610ffc565b6111276080830184610f52565b9695505050505050565b60006020820190506111466000830184610f52565b92915050565b600060a08201905081810360008301526111668188610f92565b90506111756020830187610f43565b81810360408301526111878186610f92565b9050818103606083015261119b8185610f92565b90506111aa6080830184610ffc565b9695505050505050565b600060608201905081810360008301526111ce8186610f92565b90506111dd6020830185610f43565b6111ea6040830184610ffc565b949350505050565b60006020820190506112076000830184610ffc565b92915050565b6000604051905081810181811067ffffffffffffffff82111715611234576112336113ee565b5b8060405250919050565b600067ffffffffffffffff821115611259576112586113ee565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600081519050919050565b600081905092915050565b600082825260208201905092915050565b600081905092915050565b60006112b682611373565b91506112c183611373565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff038211156112f6576112f56113bf565b5b828201905092915050565b600061130c82611373565b915061131783611373565b92508282101561132a576113296113bf565b5b828203905092915050565b600061134082611353565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b838110156113aa57808201518184015260208101905061138f565b838111156113b9576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b61143781611335565b811461144257600080fd5b50565b61144e81611373565b811461145957600080fd5b5056fea2646970667358221220122c28f2d40c8b5a843c8c0f9aaa77feda005ed08368658b86747d573b119d3364736f6c63430008000033"

// DeployTokenGateway deploys a new Ethereum contract, binding an instance of TokenGateway to it.
func DeployTokenGateway(auth *bind.TransactOpts, backend bind.ContractBackend, chain string) (common.Address, *types.Transaction, *TokenGateway, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenGatewayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(TokenGatewayBin), backend, chain)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TokenGateway{TokenGatewayCaller: TokenGatewayCaller{contract: contract}, TokenGatewayTransactor: TokenGatewayTransactor{contract: contract}, TokenGatewayFilterer: TokenGatewayFilterer{contract: contract}}, nil
}

// TokenGateway is an auto generated Go binding around an Ethereum contract.
type TokenGateway struct {
	TokenGatewayCaller     // Read-only binding to the contract
	TokenGatewayTransactor // Write-only binding to the contract
	TokenGatewayFilterer   // Log filterer for contract events
}

// TokenGatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type TokenGatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TokenGatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TokenGatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TokenGatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TokenGatewaySession struct {
	Contract     *TokenGateway     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TokenGatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TokenGatewayCallerSession struct {
	Contract *TokenGatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// TokenGatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TokenGatewayTransactorSession struct {
	Contract     *TokenGatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// TokenGatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type TokenGatewayRaw struct {
	Contract *TokenGateway // Generic contract binding to access the raw methods on
}

// TokenGatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TokenGatewayCallerRaw struct {
	Contract *TokenGatewayCaller // Generic read-only contract binding to access the raw methods on
}

// TokenGatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TokenGatewayTransactorRaw struct {
	Contract *TokenGatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTokenGateway creates a new instance of TokenGateway, bound to a specific deployed contract.
func NewTokenGateway(address common.Address, backend bind.ContractBackend) (*TokenGateway, error) {
	contract, err := bindTokenGateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TokenGateway{TokenGatewayCaller: TokenGatewayCaller{contract: contract}, TokenGatewayTransactor: TokenGatewayTransactor{contract: contract}, TokenGatewayFilterer: TokenGatewayFilterer{contract: contract}}, nil
}

// NewTokenGatewayCaller creates a new read-only instance of TokenGateway, bound to a specific deployed contract.
func NewTokenGatewayCaller(address common.Address, caller bind.ContractCaller) (*TokenGatewayCaller, error) {
	contract, err := bindTokenGateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGatewayCaller{contract: contract}, nil
}

// NewTokenGatewayTransactor creates a new write-only instance of TokenGateway, bound to a specific deployed contract.
func NewTokenGatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*TokenGatewayTransactor, error) {
	contract, err := bindTokenGateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransactor{contract: contract}, nil
}

// NewTokenGatewayFilterer creates a new log filterer instance of TokenGateway, bound to a specific deployed contract.
func NewTokenGatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*TokenGatewayFilterer, error) {
	contract, err := bindTokenGateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TokenGatewayFilterer{contract: contract}, nil
}

// bindTokenGateway binds a generic wrapper to an already deployed contract.
func bindTokenGateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(TokenGatewayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGateway *TokenGatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenGateway.Contract.TokenGatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGateway *TokenGatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGateway.Contract.TokenGatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGateway *TokenGatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGateway.Contract.TokenGatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TokenGateway *TokenGatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TokenGateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TokenGateway *TokenGatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TokenGateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TokenGateway *TokenGatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TokenGateway.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_TokenGateway *TokenGatewayCaller) GetBalance(opts *bind.CallOpts, assetId string, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TokenGateway.contract.Call(opts, &out, "getBalance", assetId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_TokenGateway *TokenGatewaySession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _TokenGateway.Contract.GetBalance(&_TokenGateway.CallOpts, assetId, account)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_TokenGateway *TokenGatewayCallerSession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _TokenGateway.Contract.GetBalance(&_TokenGateway.CallOpts, assetId, account)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_TokenGateway *TokenGatewayCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TokenGateway.contract.Call(opts, &out, "getOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_TokenGateway *TokenGatewaySession) GetOwner() (common.Address, error) {
	return _TokenGateway.Contract.GetOwner(&_TokenGateway.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_TokenGateway *TokenGatewayCallerSession) GetOwner() (common.Address, error) {
	return _TokenGateway.Contract.GetOwner(&_TokenGateway.CallOpts)
}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_TokenGateway *TokenGatewayCaller) IsChainAllowed(opts *bind.CallOpts, chain string) (bool, error) {
	var out []interface{}
	err := _TokenGateway.contract.Call(opts, &out, "isChainAllowed", chain)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_TokenGateway *TokenGatewaySession) IsChainAllowed(chain string) (bool, error) {
	return _TokenGateway.Contract.IsChainAllowed(&_TokenGateway.CallOpts, chain)
}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_TokenGateway *TokenGatewayCallerSession) IsChainAllowed(chain string) (bool, error) {
	return _TokenGateway.Contract.IsChainAllowed(&_TokenGateway.CallOpts, chain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_TokenGateway *TokenGatewayTransactor) AddAllowedChain(opts *bind.TransactOpts, newChain string) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "addAllowedChain", newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_TokenGateway *TokenGatewaySession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _TokenGateway.Contract.AddAllowedChain(&_TokenGateway.TransactOpts, newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_TokenGateway *TokenGatewayTransactorSession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _TokenGateway.Contract.AddAllowedChain(&_TokenGateway.TransactOpts, newChain)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_TokenGateway *TokenGatewayTransactor) ChangeOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "changeOwner", newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_TokenGateway *TokenGatewaySession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _TokenGateway.Contract.ChangeOwner(&_TokenGateway.TransactOpts, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_TokenGateway *TokenGatewayTransactorSession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _TokenGateway.Contract.ChangeOwner(&_TokenGateway.TransactOpts, newOwner)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactor) TransferIn(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "transferIn", assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewaySession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferIn(&_TokenGateway.TransactOpts, assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactorSession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferIn(&_TokenGateway.TransactOpts, assetId, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactor) TransferInAssetOfThisChain(opts *bind.TransactOpts, assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "transferInAssetOfThisChain", assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewaySession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferInAssetOfThisChain(&_TokenGateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactorSession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferInAssetOfThisChain(&_TokenGateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactor) TransferOut(opts *bind.TransactOpts, assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "transferOut", assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewaySession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferOut(&_TokenGateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactorSession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferOut(&_TokenGateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactor) TransferOutFromContract(opts *bind.TransactOpts, assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "transferOutFromContract", assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewaySession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferOutFromContract(&_TokenGateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactorSession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferOutFromContract(&_TokenGateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactor) TransferWithin(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.contract.Transact(opts, "transferWithin", assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewaySession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferWithin(&_TokenGateway.TransactOpts, assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_TokenGateway *TokenGatewayTransactorSession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _TokenGateway.Contract.TransferWithin(&_TokenGateway.TransactOpts, assetId, recipient, amount)
}

// TokenGatewayTransferInIterator is returned from FilterTransferIn and is used to iterate over the raw logs and unpacked data for TransferIn events raised by the TokenGateway contract.
type TokenGatewayTransferInIterator struct {
	Event *TokenGatewayTransferIn // Event containing the contract specifics and raw log

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
func (it *TokenGatewayTransferInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGatewayTransferIn)
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
		it.Event = new(TokenGatewayTransferIn)
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
func (it *TokenGatewayTransferInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGatewayTransferInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGatewayTransferIn represents a TransferIn event raised by the TokenGateway contract.
type TokenGatewayTransferIn struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferIn is a free log retrieval operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) FilterTransferIn(opts *bind.FilterOpts) (*TokenGatewayTransferInIterator, error) {

	logs, sub, err := _TokenGateway.contract.FilterLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransferInIterator{contract: _TokenGateway.contract, event: "TransferIn", logs: logs, sub: sub}, nil
}

// WatchTransferIn is a free log subscription operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) WatchTransferIn(opts *bind.WatchOpts, sink chan<- *TokenGatewayTransferIn) (event.Subscription, error) {

	logs, sub, err := _TokenGateway.contract.WatchLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGatewayTransferIn)
				if err := _TokenGateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
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

// ParseTransferIn is a log parse operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) ParseTransferIn(log types.Log) (*TokenGatewayTransferIn, error) {
	event := new(TokenGatewayTransferIn)
	if err := _TokenGateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenGatewayTransferInAssetOfThisChainIterator is returned from FilterTransferInAssetOfThisChain and is used to iterate over the raw logs and unpacked data for TransferInAssetOfThisChain events raised by the TokenGateway contract.
type TokenGatewayTransferInAssetOfThisChainIterator struct {
	Event *TokenGatewayTransferInAssetOfThisChain // Event containing the contract specifics and raw log

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
func (it *TokenGatewayTransferInAssetOfThisChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGatewayTransferInAssetOfThisChain)
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
		it.Event = new(TokenGatewayTransferInAssetOfThisChain)
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
func (it *TokenGatewayTransferInAssetOfThisChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGatewayTransferInAssetOfThisChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGatewayTransferInAssetOfThisChain represents a TransferInAssetOfThisChain event raised by the TokenGateway contract.
type TokenGatewayTransferInAssetOfThisChain struct {
	AssetAddr common.Address
	Recipient common.Address
	Amount    *big.Int
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInAssetOfThisChain is a free log retrieval operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) FilterTransferInAssetOfThisChain(opts *bind.FilterOpts) (*TokenGatewayTransferInAssetOfThisChainIterator, error) {

	logs, sub, err := _TokenGateway.contract.FilterLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransferInAssetOfThisChainIterator{contract: _TokenGateway.contract, event: "TransferInAssetOfThisChain", logs: logs, sub: sub}, nil
}

// WatchTransferInAssetOfThisChain is a free log subscription operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) WatchTransferInAssetOfThisChain(opts *bind.WatchOpts, sink chan<- *TokenGatewayTransferInAssetOfThisChain) (event.Subscription, error) {

	logs, sub, err := _TokenGateway.contract.WatchLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGatewayTransferInAssetOfThisChain)
				if err := _TokenGateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
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

// ParseTransferInAssetOfThisChain is a log parse operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) ParseTransferInAssetOfThisChain(log types.Log) (*TokenGatewayTransferInAssetOfThisChain, error) {
	event := new(TokenGatewayTransferInAssetOfThisChain)
	if err := _TokenGateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenGatewayTransferOutIterator is returned from FilterTransferOut and is used to iterate over the raw logs and unpacked data for TransferOut events raised by the TokenGateway contract.
type TokenGatewayTransferOutIterator struct {
	Event *TokenGatewayTransferOut // Event containing the contract specifics and raw log

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
func (it *TokenGatewayTransferOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGatewayTransferOut)
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
		it.Event = new(TokenGatewayTransferOut)
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
func (it *TokenGatewayTransferOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGatewayTransferOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGatewayTransferOut represents a TransferOut event raised by the TokenGateway contract.
type TokenGatewayTransferOut struct {
	AssetId   string
	From      common.Address
	ToChain   string
	Recipient string
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferOut is a free log retrieval operation binding the contract event 0x265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f.
//
// Solidity: event TransferOut(string assetId, address from, string toChain, string recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) FilterTransferOut(opts *bind.FilterOpts) (*TokenGatewayTransferOutIterator, error) {

	logs, sub, err := _TokenGateway.contract.FilterLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransferOutIterator{contract: _TokenGateway.contract, event: "TransferOut", logs: logs, sub: sub}, nil
}

// WatchTransferOut is a free log subscription operation binding the contract event 0x265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f.
//
// Solidity: event TransferOut(string assetId, address from, string toChain, string recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) WatchTransferOut(opts *bind.WatchOpts, sink chan<- *TokenGatewayTransferOut) (event.Subscription, error) {

	logs, sub, err := _TokenGateway.contract.WatchLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGatewayTransferOut)
				if err := _TokenGateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
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

// ParseTransferOut is a log parse operation binding the contract event 0x265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f.
//
// Solidity: event TransferOut(string assetId, address from, string toChain, string recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) ParseTransferOut(log types.Log) (*TokenGatewayTransferOut, error) {
	event := new(TokenGatewayTransferOut)
	if err := _TokenGateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenGatewayTransferOutFromContractIterator is returned from FilterTransferOutFromContract and is used to iterate over the raw logs and unpacked data for TransferOutFromContract events raised by the TokenGateway contract.
type TokenGatewayTransferOutFromContractIterator struct {
	Event *TokenGatewayTransferOutFromContract // Event containing the contract specifics and raw log

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
func (it *TokenGatewayTransferOutFromContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGatewayTransferOutFromContract)
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
		it.Event = new(TokenGatewayTransferOutFromContract)
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
func (it *TokenGatewayTransferOutFromContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGatewayTransferOutFromContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGatewayTransferOutFromContract represents a TransferOutFromContract event raised by the TokenGateway contract.
type TokenGatewayTransferOutFromContract struct {
	AssetAddr common.Address
	ToChain   string
	Recipient string
	Amount    *big.Int
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferOutFromContract is a free log retrieval operation binding the contract event 0x889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e541.
//
// Solidity: event TransferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) FilterTransferOutFromContract(opts *bind.FilterOpts) (*TokenGatewayTransferOutFromContractIterator, error) {

	logs, sub, err := _TokenGateway.contract.FilterLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransferOutFromContractIterator{contract: _TokenGateway.contract, event: "TransferOutFromContract", logs: logs, sub: sub}, nil
}

// WatchTransferOutFromContract is a free log subscription operation binding the contract event 0x889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e541.
//
// Solidity: event TransferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) WatchTransferOutFromContract(opts *bind.WatchOpts, sink chan<- *TokenGatewayTransferOutFromContract) (event.Subscription, error) {

	logs, sub, err := _TokenGateway.contract.WatchLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGatewayTransferOutFromContract)
				if err := _TokenGateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
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

// ParseTransferOutFromContract is a log parse operation binding the contract event 0x889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e541.
//
// Solidity: event TransferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount, bool success)
func (_TokenGateway *TokenGatewayFilterer) ParseTransferOutFromContract(log types.Log) (*TokenGatewayTransferOutFromContract, error) {
	event := new(TokenGatewayTransferOutFromContract)
	if err := _TokenGateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TokenGatewayTransferWithinIterator is returned from FilterTransferWithin and is used to iterate over the raw logs and unpacked data for TransferWithin events raised by the TokenGateway contract.
type TokenGatewayTransferWithinIterator struct {
	Event *TokenGatewayTransferWithin // Event containing the contract specifics and raw log

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
func (it *TokenGatewayTransferWithinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TokenGatewayTransferWithin)
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
		it.Event = new(TokenGatewayTransferWithin)
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
func (it *TokenGatewayTransferWithinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TokenGatewayTransferWithinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TokenGatewayTransferWithin represents a TransferWithin event raised by the TokenGateway contract.
type TokenGatewayTransferWithin struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferWithin is a free log retrieval operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) FilterTransferWithin(opts *bind.FilterOpts) (*TokenGatewayTransferWithinIterator, error) {

	logs, sub, err := _TokenGateway.contract.FilterLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return &TokenGatewayTransferWithinIterator{contract: _TokenGateway.contract, event: "TransferWithin", logs: logs, sub: sub}, nil
}

// WatchTransferWithin is a free log subscription operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) WatchTransferWithin(opts *bind.WatchOpts, sink chan<- *TokenGatewayTransferWithin) (event.Subscription, error) {

	logs, sub, err := _TokenGateway.contract.WatchLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TokenGatewayTransferWithin)
				if err := _TokenGateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
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

// ParseTransferWithin is a log parse operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_TokenGateway *TokenGatewayFilterer) ParseTransferWithin(log types.Log) (*TokenGatewayTransferWithin, error) {
	event := new(TokenGatewayTransferWithin)
	if err := _TokenGateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
