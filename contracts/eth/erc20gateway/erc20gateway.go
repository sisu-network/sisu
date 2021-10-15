// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20Gateway

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

// Erc20GatewayABI is the input ABI used to generate the binding from.
const Erc20GatewayABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferInAssetOfThisChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferOutFromContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferWithin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newChain\",\"type\":\"string\"}],\"name\":\"addAllowedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInAssetOfThisChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOutFromContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferWithin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Erc20GatewayBin is the compiled bytecode used for deploying new contracts.
var Erc20GatewayBin = "0x60806040523480156200001157600080fd5b5060405162001c6438038062001c648339818101604052810190620000379190620001b9565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600190805190602001906200008f92919062000097565b50506200032f565b828054620000a5906200029b565b90600052602060002090601f016020900481019282620000c9576000855562000115565b82601f10620000e457805160ff191683800117855562000115565b8280016001018555821562000115579182015b8281111562000114578251825591602001919060010190620000f7565b5b50905062000124919062000128565b5090565b5b808211156200014357600081600090555060010162000129565b5090565b60006200015e620001588462000232565b620001fe565b9050828152602081018484840111156200017757600080fd5b6200018484828562000265565b509392505050565b600082601f8301126200019e57600080fd5b8151620001b084826020860162000147565b91505092915050565b600060208284031215620001cc57600080fd5b600082015167ffffffffffffffff811115620001e757600080fd5b620001f5848285016200018c565b91505092915050565b6000604051905081810181811067ffffffffffffffff8211171562000228576200022762000300565b5b8060405250919050565b600067ffffffffffffffff82111562000250576200024f62000300565b5b601f19601f8301169050602081019050919050565b60005b838110156200028557808201518184015260208101905062000268565b8381111562000295576000848401525b50505050565b60006002820490506001821680620002b457607f821691505b60208210811415620002cb57620002ca620002d1565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611925806200033f6000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063893d20e811610066578063893d20e81461010857806389f4b317146101265780639a12bb6014610156578063a6f9dae114610172578063bf6022be1461018e57610093565b80630832f134146100985780636c189db8146100b45780637dd1f364146100d0578063810b9f12146100ec575b600080fd5b6100b260048036038101906100ad9190611027565b6101aa565b005b6100ce60048036038101906100c99190610eb0565b6102ea565b005b6100ea60048036038101906100e59190611027565b61068b565b005b6101066004803603810190610101919061108e565b610890565b005b610110610a28565b60405161011d91906113fe565b60405180910390f35b610140600480360381019061013b9190610fd3565b610a51565b60405161014d9190611685565b60405180910390f35b610170600480360381019061016b9190610f92565b610ab7565b005b61018c60048036038101906101879190610e87565b610b48565b005b6101a860048036038101906101a39190610eff565b610c19565b005b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610238576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161022f90611625565b60405180910390fd5b8060028460405161024991906113e7565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546102a3919061173e565b925050819055507fbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea0948383836040516102dd93929190611587565b60405180910390a1505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610378576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161036f90611625565b60405180910390fd5b6000808473ffffffffffffffffffffffffffffffffffffffff16306040516024016103a391906113fe565b6040516020818303038152906040527f70a08231000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161042d91906113d0565b6000604051808303816000865af19150503d806000811461046a576040519150601f19603f3d011682016040523d82523d6000602084013e61046f565b606091505b5091509150816104b4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104ab90611645565b60405180910390fd5b82818060200190518101906104c99190611139565b101561050a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050190611665565b60405180910390fd5b60008573ffffffffffffffffffffffffffffffffffffffff1685856040516024016105369291906114f6565b6040516020818303038152906040527fa9059cbb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516105c091906113d0565b6000604051808303816000865af19150503d80600081146105fd576040519150601f19603f3d011682016040523d82523d6000602084013e610602565b606091505b5050905080610646576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161063d90611605565b60405180910390fd5b7fc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c38686868460405161067b9493929190611450565b60405180910390a1505050505050565b600081116106ce576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106c5906115e5565b60405180910390fd5b806002846040516106df91906113e7565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561076c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610763906115c5565b60405180910390fd5b8060028460405161077d91906113e7565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546107d79190611794565b92505081905550806002846040516107ef91906113e7565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610849919061173e565b925050819055507f8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df483838360405161088393929190611587565b60405180910390a1505050565b600081116108d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108ca906115e5565b60405180910390fd5b806002856040516108e491906113e7565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541015610971576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161096890611665565b60405180910390fd5b8060028560405161098291906113e7565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546109dc9190611794565b925050819055507f265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f8433858585604051610a1a95949392919061151f565b60405180910390a150505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600283604051610a6391906113e7565b908152602001604051809103902060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610b45576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b3c90611625565b60405180910390fd5b50565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610bd6576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bcd90611625565b60405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008111610c5c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c53906115e5565b60405180910390fd5b60008473ffffffffffffffffffffffffffffffffffffffff16333084604051602401610c8a93929190611419565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610d1491906113d0565b6000604051808303816000865af19150503d8060008114610d51576040519150601f19603f3d011682016040523d82523d6000602084013e610d56565b606091505b5050905080610d9a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d9190611605565b60405180910390fd5b7f889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e5418585858585604051610dd1959493929190611495565b60405180910390a15050505050565b6000610df3610dee846116d1565b6116a0565b905082815260208101848484011115610e0b57600080fd5b610e16848285611810565b509392505050565b600081359050610e2d816118c1565b92915050565b600082601f830112610e4457600080fd5b8135610e54848260208601610de0565b91505092915050565b600081359050610e6c816118d8565b92915050565b600081519050610e81816118d8565b92915050565b600060208284031215610e9957600080fd5b6000610ea784828501610e1e565b91505092915050565b600080600060608486031215610ec557600080fd5b6000610ed386828701610e1e565b9350506020610ee486828701610e1e565b9250506040610ef586828701610e5d565b9150509250925092565b60008060008060808587031215610f1557600080fd5b6000610f2387828801610e1e565b945050602085013567ffffffffffffffff811115610f4057600080fd5b610f4c87828801610e33565b935050604085013567ffffffffffffffff811115610f6957600080fd5b610f7587828801610e33565b9250506060610f8687828801610e5d565b91505092959194509250565b600060208284031215610fa457600080fd5b600082013567ffffffffffffffff811115610fbe57600080fd5b610fca84828501610e33565b91505092915050565b60008060408385031215610fe657600080fd5b600083013567ffffffffffffffff81111561100057600080fd5b61100c85828601610e33565b925050602061101d85828601610e1e565b9150509250929050565b60008060006060848603121561103c57600080fd5b600084013567ffffffffffffffff81111561105657600080fd5b61106286828701610e33565b935050602061107386828701610e1e565b925050604061108486828701610e5d565b9150509250925092565b600080600080608085870312156110a457600080fd5b600085013567ffffffffffffffff8111156110be57600080fd5b6110ca87828801610e33565b945050602085013567ffffffffffffffff8111156110e757600080fd5b6110f387828801610e33565b935050604085013567ffffffffffffffff81111561111057600080fd5b61111c87828801610e33565b925050606061112d87828801610e5d565b91505092959194509250565b60006020828403121561114b57600080fd5b600061115984828501610e72565b91505092915050565b61116b816117c8565b82525050565b61117a816117da565b82525050565b600061118b82611701565b6111958185611717565b93506111a581856020860161181f565b80840191505092915050565b60006111bc8261170c565b6111c68185611722565b93506111d681856020860161181f565b6111df816118b0565b840191505092915050565b60006111f58261170c565b6111ff8185611733565b935061120f81856020860161181f565b80840191505092915050565b6000611228602a83611722565b91507f42616c616e6365206c657373207468616e20616d6f756e74206265696e67207460008301527f72616e73666572726564000000000000000000000000000000000000000000006020830152604082019050919050565b600061128e601d83611722565b91507f416d6f756e74206d7573742062652067726561746572207468616e20300000006000830152602082019050919050565b60006112ce601283611722565b91507f4661696c656420746f207472616e7366657200000000000000000000000000006000830152602082019050919050565b600061130e600d83611722565b91507f4d757374206265206f776e6572000000000000000000000000000000000000006000830152602082019050919050565b600061134e601583611722565b91507f4661696c656420746f206765742062616c616e636500000000000000000000006000830152602082019050919050565b600061138e601183611722565b91507f4e6f7420656e6f75676820746f6b656e730000000000000000000000000000006000830152602082019050919050565b6113ca81611806565b82525050565b60006113dc8284611180565b915081905092915050565b60006113f382846111ea565b915081905092915050565b60006020820190506114136000830184611162565b92915050565b600060608201905061142e6000830186611162565b61143b6020830185611162565b61144860408301846113c1565b949350505050565b60006080820190506114656000830187611162565b6114726020830186611162565b61147f60408301856113c1565b61148c6060830184611171565b95945050505050565b600060a0820190506114aa6000830188611162565b81810360208301526114bc81876111b1565b905081810360408301526114d081866111b1565b90506114df60608301856113c1565b6114ec6080830184611171565b9695505050505050565b600060408201905061150b6000830185611162565b61151860208301846113c1565b9392505050565b600060a082019050818103600083015261153981886111b1565b90506115486020830187611162565b818103604083015261155a81866111b1565b9050818103606083015261156e81856111b1565b905061157d60808301846113c1565b9695505050505050565b600060608201905081810360008301526115a181866111b1565b90506115b06020830185611162565b6115bd60408301846113c1565b949350505050565b600060208201905081810360008301526115de8161121b565b9050919050565b600060208201905081810360008301526115fe81611281565b9050919050565b6000602082019050818103600083015261161e816112c1565b9050919050565b6000602082019050818103600083015261163e81611301565b9050919050565b6000602082019050818103600083015261165e81611341565b9050919050565b6000602082019050818103600083015261167e81611381565b9050919050565b600060208201905061169a60008301846113c1565b92915050565b6000604051905081810181811067ffffffffffffffff821117156116c7576116c6611881565b5b8060405250919050565b600067ffffffffffffffff8211156116ec576116eb611881565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600081519050919050565b600081905092915050565b600082825260208201905092915050565b600081905092915050565b600061174982611806565b915061175483611806565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561178957611788611852565b5b828201905092915050565b600061179f82611806565b91506117aa83611806565b9250828210156117bd576117bc611852565b5b828203905092915050565b60006117d3826117e6565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b8381101561183d578082015181840152602081019050611822565b8381111561184c576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b6118ca816117c8565b81146118d557600080fd5b50565b6118e181611806565b81146118ec57600080fd5b5056fea2646970667358221220ddaf44d97620e79765a9597f9cefd900250b252eb2f4e9b87ae5b14abff9100764736f6c63430008000033"

// DeployErc20Gateway deploys a new Ethereum contract, binding an instance of Erc20Gateway to it.
func DeployErc20Gateway(auth *bind.TransactOpts, backend bind.ContractBackend, chain string) (common.Address, *types.Transaction, *Erc20Gateway, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20GatewayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Erc20GatewayBin), backend, chain)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc20Gateway{Erc20GatewayCaller: Erc20GatewayCaller{contract: contract}, Erc20GatewayTransactor: Erc20GatewayTransactor{contract: contract}, Erc20GatewayFilterer: Erc20GatewayFilterer{contract: contract}}, nil
}

// Erc20Gateway is an auto generated Go binding around an Ethereum contract.
type Erc20Gateway struct {
	Erc20GatewayCaller     // Read-only binding to the contract
	Erc20GatewayTransactor // Write-only binding to the contract
	Erc20GatewayFilterer   // Log filterer for contract events
}

// Erc20GatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20GatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20GatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20GatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20GatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20GatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20GatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20GatewaySession struct {
	Contract     *Erc20Gateway     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20GatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20GatewayCallerSession struct {
	Contract *Erc20GatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// Erc20GatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20GatewayTransactorSession struct {
	Contract     *Erc20GatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// Erc20GatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20GatewayRaw struct {
	Contract *Erc20Gateway // Generic contract binding to access the raw methods on
}

// Erc20GatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20GatewayCallerRaw struct {
	Contract *Erc20GatewayCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20GatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20GatewayTransactorRaw struct {
	Contract *Erc20GatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20Gateway creates a new instance of Erc20Gateway, bound to a specific deployed contract.
func NewErc20Gateway(address common.Address, backend bind.ContractBackend) (*Erc20Gateway, error) {
	contract, err := bindErc20Gateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20Gateway{Erc20GatewayCaller: Erc20GatewayCaller{contract: contract}, Erc20GatewayTransactor: Erc20GatewayTransactor{contract: contract}, Erc20GatewayFilterer: Erc20GatewayFilterer{contract: contract}}, nil
}

// NewErc20GatewayCaller creates a new read-only instance of Erc20Gateway, bound to a specific deployed contract.
func NewErc20GatewayCaller(address common.Address, caller bind.ContractCaller) (*Erc20GatewayCaller, error) {
	contract, err := bindErc20Gateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayCaller{contract: contract}, nil
}

// NewErc20GatewayTransactor creates a new write-only instance of Erc20Gateway, bound to a specific deployed contract.
func NewErc20GatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20GatewayTransactor, error) {
	contract, err := bindErc20Gateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransactor{contract: contract}, nil
}

// NewErc20GatewayFilterer creates a new log filterer instance of Erc20Gateway, bound to a specific deployed contract.
func NewErc20GatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc20GatewayFilterer, error) {
	contract, err := bindErc20Gateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayFilterer{contract: contract}, nil
}

// bindErc20Gateway binds a generic wrapper to an already deployed contract.
func bindErc20Gateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20GatewayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Gateway *Erc20GatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Gateway.Contract.Erc20GatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Gateway *Erc20GatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.Erc20GatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Gateway *Erc20GatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.Erc20GatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20Gateway *Erc20GatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20Gateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20Gateway *Erc20GatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20Gateway *Erc20GatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20Gateway *Erc20GatewayCaller) GetBalance(opts *bind.CallOpts, assetId string, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc20Gateway.contract.Call(opts, &out, "getBalance", assetId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20Gateway *Erc20GatewaySession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _Erc20Gateway.Contract.GetBalance(&_Erc20Gateway.CallOpts, assetId, account)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20Gateway *Erc20GatewayCallerSession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _Erc20Gateway.Contract.GetBalance(&_Erc20Gateway.CallOpts, assetId, account)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20Gateway *Erc20GatewayCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc20Gateway.contract.Call(opts, &out, "getOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20Gateway *Erc20GatewaySession) GetOwner() (common.Address, error) {
	return _Erc20Gateway.Contract.GetOwner(&_Erc20Gateway.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20Gateway *Erc20GatewayCallerSession) GetOwner() (common.Address, error) {
	return _Erc20Gateway.Contract.GetOwner(&_Erc20Gateway.CallOpts)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) AddAllowedChain(opts *bind.TransactOpts, newChain string) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "addAllowedChain", newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20Gateway *Erc20GatewaySession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.AddAllowedChain(&_Erc20Gateway.TransactOpts, newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.AddAllowedChain(&_Erc20Gateway.TransactOpts, newChain)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) ChangeOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "changeOwner", newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20Gateway *Erc20GatewaySession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.ChangeOwner(&_Erc20Gateway.TransactOpts, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.ChangeOwner(&_Erc20Gateway.TransactOpts, newOwner)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) TransferIn(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "transferIn", assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewaySession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferIn(&_Erc20Gateway.TransactOpts, assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferIn(&_Erc20Gateway.TransactOpts, assetId, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) TransferInAssetOfThisChain(opts *bind.TransactOpts, assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "transferInAssetOfThisChain", assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewaySession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferInAssetOfThisChain(&_Erc20Gateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferInAssetOfThisChain(&_Erc20Gateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) TransferOut(opts *bind.TransactOpts, assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "transferOut", assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewaySession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferOut(&_Erc20Gateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferOut(&_Erc20Gateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) TransferOutFromContract(opts *bind.TransactOpts, assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "transferOutFromContract", assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewaySession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferOutFromContract(&_Erc20Gateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferOutFromContract(&_Erc20Gateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactor) TransferWithin(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.contract.Transact(opts, "transferWithin", assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewaySession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferWithin(&_Erc20Gateway.TransactOpts, assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20Gateway *Erc20GatewayTransactorSession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20Gateway.Contract.TransferWithin(&_Erc20Gateway.TransactOpts, assetId, recipient, amount)
}

// Erc20GatewayTransferInIterator is returned from FilterTransferIn and is used to iterate over the raw logs and unpacked data for TransferIn events raised by the Erc20Gateway contract.
type Erc20GatewayTransferInIterator struct {
	Event *Erc20GatewayTransferIn // Event containing the contract specifics and raw log

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
func (it *Erc20GatewayTransferInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20GatewayTransferIn)
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
		it.Event = new(Erc20GatewayTransferIn)
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
func (it *Erc20GatewayTransferInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20GatewayTransferInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20GatewayTransferIn represents a TransferIn event raised by the Erc20Gateway contract.
type Erc20GatewayTransferIn struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferIn is a free log retrieval operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_Erc20Gateway *Erc20GatewayFilterer) FilterTransferIn(opts *bind.FilterOpts) (*Erc20GatewayTransferInIterator, error) {

	logs, sub, err := _Erc20Gateway.contract.FilterLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransferInIterator{contract: _Erc20Gateway.contract, event: "TransferIn", logs: logs, sub: sub}, nil
}

// WatchTransferIn is a free log subscription operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_Erc20Gateway *Erc20GatewayFilterer) WatchTransferIn(opts *bind.WatchOpts, sink chan<- *Erc20GatewayTransferIn) (event.Subscription, error) {

	logs, sub, err := _Erc20Gateway.contract.WatchLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20GatewayTransferIn)
				if err := _Erc20Gateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
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
func (_Erc20Gateway *Erc20GatewayFilterer) ParseTransferIn(log types.Log) (*Erc20GatewayTransferIn, error) {
	event := new(Erc20GatewayTransferIn)
	if err := _Erc20Gateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20GatewayTransferInAssetOfThisChainIterator is returned from FilterTransferInAssetOfThisChain and is used to iterate over the raw logs and unpacked data for TransferInAssetOfThisChain events raised by the Erc20Gateway contract.
type Erc20GatewayTransferInAssetOfThisChainIterator struct {
	Event *Erc20GatewayTransferInAssetOfThisChain // Event containing the contract specifics and raw log

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
func (it *Erc20GatewayTransferInAssetOfThisChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20GatewayTransferInAssetOfThisChain)
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
		it.Event = new(Erc20GatewayTransferInAssetOfThisChain)
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
func (it *Erc20GatewayTransferInAssetOfThisChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20GatewayTransferInAssetOfThisChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20GatewayTransferInAssetOfThisChain represents a TransferInAssetOfThisChain event raised by the Erc20Gateway contract.
type Erc20GatewayTransferInAssetOfThisChain struct {
	AssetAddr common.Address
	Recipient common.Address
	Amount    *big.Int
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInAssetOfThisChain is a free log retrieval operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_Erc20Gateway *Erc20GatewayFilterer) FilterTransferInAssetOfThisChain(opts *bind.FilterOpts) (*Erc20GatewayTransferInAssetOfThisChainIterator, error) {

	logs, sub, err := _Erc20Gateway.contract.FilterLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransferInAssetOfThisChainIterator{contract: _Erc20Gateway.contract, event: "TransferInAssetOfThisChain", logs: logs, sub: sub}, nil
}

// WatchTransferInAssetOfThisChain is a free log subscription operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_Erc20Gateway *Erc20GatewayFilterer) WatchTransferInAssetOfThisChain(opts *bind.WatchOpts, sink chan<- *Erc20GatewayTransferInAssetOfThisChain) (event.Subscription, error) {

	logs, sub, err := _Erc20Gateway.contract.WatchLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20GatewayTransferInAssetOfThisChain)
				if err := _Erc20Gateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
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
func (_Erc20Gateway *Erc20GatewayFilterer) ParseTransferInAssetOfThisChain(log types.Log) (*Erc20GatewayTransferInAssetOfThisChain, error) {
	event := new(Erc20GatewayTransferInAssetOfThisChain)
	if err := _Erc20Gateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20GatewayTransferOutIterator is returned from FilterTransferOut and is used to iterate over the raw logs and unpacked data for TransferOut events raised by the Erc20Gateway contract.
type Erc20GatewayTransferOutIterator struct {
	Event *Erc20GatewayTransferOut // Event containing the contract specifics and raw log

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
func (it *Erc20GatewayTransferOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20GatewayTransferOut)
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
		it.Event = new(Erc20GatewayTransferOut)
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
func (it *Erc20GatewayTransferOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20GatewayTransferOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20GatewayTransferOut represents a TransferOut event raised by the Erc20Gateway contract.
type Erc20GatewayTransferOut struct {
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
func (_Erc20Gateway *Erc20GatewayFilterer) FilterTransferOut(opts *bind.FilterOpts) (*Erc20GatewayTransferOutIterator, error) {

	logs, sub, err := _Erc20Gateway.contract.FilterLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransferOutIterator{contract: _Erc20Gateway.contract, event: "TransferOut", logs: logs, sub: sub}, nil
}

// WatchTransferOut is a free log subscription operation binding the contract event 0x265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f.
//
// Solidity: event TransferOut(string assetId, address from, string toChain, string recipient, uint256 amount)
func (_Erc20Gateway *Erc20GatewayFilterer) WatchTransferOut(opts *bind.WatchOpts, sink chan<- *Erc20GatewayTransferOut) (event.Subscription, error) {

	logs, sub, err := _Erc20Gateway.contract.WatchLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20GatewayTransferOut)
				if err := _Erc20Gateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
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
func (_Erc20Gateway *Erc20GatewayFilterer) ParseTransferOut(log types.Log) (*Erc20GatewayTransferOut, error) {
	event := new(Erc20GatewayTransferOut)
	if err := _Erc20Gateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20GatewayTransferOutFromContractIterator is returned from FilterTransferOutFromContract and is used to iterate over the raw logs and unpacked data for TransferOutFromContract events raised by the Erc20Gateway contract.
type Erc20GatewayTransferOutFromContractIterator struct {
	Event *Erc20GatewayTransferOutFromContract // Event containing the contract specifics and raw log

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
func (it *Erc20GatewayTransferOutFromContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20GatewayTransferOutFromContract)
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
		it.Event = new(Erc20GatewayTransferOutFromContract)
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
func (it *Erc20GatewayTransferOutFromContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20GatewayTransferOutFromContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20GatewayTransferOutFromContract represents a TransferOutFromContract event raised by the Erc20Gateway contract.
type Erc20GatewayTransferOutFromContract struct {
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
func (_Erc20Gateway *Erc20GatewayFilterer) FilterTransferOutFromContract(opts *bind.FilterOpts) (*Erc20GatewayTransferOutFromContractIterator, error) {

	logs, sub, err := _Erc20Gateway.contract.FilterLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransferOutFromContractIterator{contract: _Erc20Gateway.contract, event: "TransferOutFromContract", logs: logs, sub: sub}, nil
}

// WatchTransferOutFromContract is a free log subscription operation binding the contract event 0x889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e541.
//
// Solidity: event TransferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount, bool success)
func (_Erc20Gateway *Erc20GatewayFilterer) WatchTransferOutFromContract(opts *bind.WatchOpts, sink chan<- *Erc20GatewayTransferOutFromContract) (event.Subscription, error) {

	logs, sub, err := _Erc20Gateway.contract.WatchLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20GatewayTransferOutFromContract)
				if err := _Erc20Gateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
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
func (_Erc20Gateway *Erc20GatewayFilterer) ParseTransferOutFromContract(log types.Log) (*Erc20GatewayTransferOutFromContract, error) {
	event := new(Erc20GatewayTransferOutFromContract)
	if err := _Erc20Gateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20GatewayTransferWithinIterator is returned from FilterTransferWithin and is used to iterate over the raw logs and unpacked data for TransferWithin events raised by the Erc20Gateway contract.
type Erc20GatewayTransferWithinIterator struct {
	Event *Erc20GatewayTransferWithin // Event containing the contract specifics and raw log

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
func (it *Erc20GatewayTransferWithinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20GatewayTransferWithin)
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
		it.Event = new(Erc20GatewayTransferWithin)
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
func (it *Erc20GatewayTransferWithinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20GatewayTransferWithinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20GatewayTransferWithin represents a TransferWithin event raised by the Erc20Gateway contract.
type Erc20GatewayTransferWithin struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferWithin is a free log retrieval operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_Erc20Gateway *Erc20GatewayFilterer) FilterTransferWithin(opts *bind.FilterOpts) (*Erc20GatewayTransferWithinIterator, error) {

	logs, sub, err := _Erc20Gateway.contract.FilterLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return &Erc20GatewayTransferWithinIterator{contract: _Erc20Gateway.contract, event: "TransferWithin", logs: logs, sub: sub}, nil
}

// WatchTransferWithin is a free log subscription operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_Erc20Gateway *Erc20GatewayFilterer) WatchTransferWithin(opts *bind.WatchOpts, sink chan<- *Erc20GatewayTransferWithin) (event.Subscription, error) {

	logs, sub, err := _Erc20Gateway.contract.WatchLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20GatewayTransferWithin)
				if err := _Erc20Gateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
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
func (_Erc20Gateway *Erc20GatewayFilterer) ParseTransferWithin(log types.Log) (*Erc20GatewayTransferWithin, error) {
	event := new(Erc20GatewayTransferWithin)
	if err := _Erc20Gateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
