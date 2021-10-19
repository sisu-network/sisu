// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20gateway

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

// Erc20gatewayABI is the input ABI used to generate the binding from.
const Erc20gatewayABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"},{\"internalType\":\"string[]\",\"name\":\"chains\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferInAssetOfThisChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferOutFromContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferWithin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newChain\",\"type\":\"string\"}],\"name\":\"addAllowedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"isChainAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInAssetOfThisChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOutFromContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferWithin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Erc20gatewayBin is the compiled bytecode used for deploying new contracts.
var Erc20gatewayBin = "0x60806040523480156200001157600080fd5b50604051620022d3380380620022d38339818101604052810190620000379190620002ea565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555081600190805190602001906200008f92919062000136565b5060005b81518110156200012d5760016003838381518110620000db577f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b6020026020010151604051620000f2919062000394565b908152602001604051809103902060006101000a81548160ff02191690831515021790555080806200012490620004cf565b91505062000093565b505050620005aa565b828054620001449062000499565b90600052602060002090601f016020900481019282620001685760008555620001b4565b82601f106200018357805160ff1916838001178555620001b4565b82800160010185558215620001b4579182015b82811115620001b357825182559160200191906001019062000196565b5b509050620001c39190620001c7565b5090565b5b80821115620001e2576000816000905550600101620001c8565b5090565b6000620001fd620001f784620003e1565b620003ad565b9050808382526020820190508260005b85811015620002415781518501620002268882620002bd565b8452602084019350602083019250506001810190506200020d565b5050509392505050565b6000620002626200025c8462000410565b620003ad565b9050828152602081018484840111156200027b57600080fd5b6200028884828562000463565b509392505050565b600082601f830112620002a257600080fd5b8151620002b4848260208601620001e6565b91505092915050565b600082601f830112620002cf57600080fd5b8151620002e18482602086016200024b565b91505092915050565b60008060408385031215620002fe57600080fd5b600083015167ffffffffffffffff8111156200031957600080fd5b6200032785828601620002bd565b925050602083015167ffffffffffffffff8111156200034557600080fd5b620003538582860162000290565b9150509250929050565b60006200036a8262000443565b6200037681856200044e565b93506200038881856020860162000463565b80840191505092915050565b6000620003a282846200035d565b915081905092915050565b6000604051905081810181811067ffffffffffffffff82111715620003d757620003d66200057b565b5b8060405250919050565b600067ffffffffffffffff821115620003ff57620003fe6200057b565b5b602082029050602081019050919050565b600067ffffffffffffffff8211156200042e576200042d6200057b565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600081905092915050565b6000819050919050565b60005b838110156200048357808201518184015260208101905062000466565b8381111562000493576000848401525b50505050565b60006002820490506001821680620004b257607f821691505b60208210811415620004c957620004c86200054c565b5b50919050565b6000620004dc8262000459565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8214156200051257620005116200051d565b5b600182019050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611d1980620005ba6000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806389f4b3171161006657806389f4b317146101315780639a12bb6014610161578063a6f9dae11461017d578063bf6022be14610199578063d29fde24146101b55761009e565b80630832f134146100a35780636c189db8146100bf5780637dd1f364146100db578063810b9f12146100f7578063893d20e814610113575b600080fd5b6100bd60048036038101906100b89190611234565b6101e5565b005b6100d960048036038101906100d491906110bd565b610325565b005b6100f560048036038101906100f09190611234565b6106c6565b005b610111600480360381019061010c919061129b565b6108cb565b005b61011b610ace565b6040516101289190611721565b60405180910390f35b61014b600480360381019061014691906111e0565b610af7565b6040516101589190611a03565b60405180910390f35b61017b6004803603810190610176919061119f565b610b5d565b005b61019760048036038101906101929190611094565b610cb5565b005b6101b360048036038101906101ae919061110c565b610d86565b005b6101cf60048036038101906101ca919061119f565b610fb8565b6040516101dc9190611842565b60405180910390f35b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610273576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161026a906119a3565b60405180910390fd5b8060028460405161028491906116f3565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546102de9190611ad1565b925050819055507fbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094838383604051610318939291906118c5565b60405180910390a1505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103b3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103aa906119a3565b60405180910390fd5b6000808473ffffffffffffffffffffffffffffffffffffffff16306040516024016103de9190611721565b6040516020818303038152906040527f70a08231000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505060405161046891906116dc565b6000604051808303816000865af19150503d80600081146104a5576040519150601f19603f3d011682016040523d82523d6000602084013e6104aa565b606091505b5091509150816104ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104e6906119c3565b60405180910390fd5b82818060200190518101906105049190611346565b1015610545576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161053c906119e3565b60405180910390fd5b60008573ffffffffffffffffffffffffffffffffffffffff168585604051602401610571929190611819565b6040516020818303038152906040527fa9059cbb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516105fb91906116dc565b6000604051808303816000865af19150503d8060008114610638576040519150601f19603f3d011682016040523d82523d6000602084013e61063d565b606091505b5050905080610681576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161067890611983565b60405180910390fd5b7fc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3868686846040516106b69493929190611773565b60405180910390a1505050505050565b60008111610709576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161070090611943565b60405180910390fd5b8060028460405161071a91906116f3565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156107a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079e90611923565b60405180910390fd5b806002846040516107b891906116f3565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108129190611b27565b925050819055508060028460405161082a91906116f3565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108849190611ad1565b925050819055507f8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df48383836040516108be939291906118c5565b60405180910390a1505050565b6000811161090e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161090590611943565b60405180910390fd5b8060028560405161091f91906116f3565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156109ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a3906119e3565b60405180910390fd5b6003836040516109bc91906116f3565b908152602001604051809103902060009054906101000a900460ff16610a17576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a0e90611963565b60405180910390fd5b80600285604051610a2891906116f3565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610a829190611b27565b925050819055507f265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f8433858585604051610ac095949392919061185d565b60405180910390a150505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600283604051610b0991906116f3565b908152602001604051809103902060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610beb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610be2906119a3565b60405180910390fd5b6001604051602001610bfd919061170a565b6040516020818303038152906040528051906020012081604051602001610c2491906116f3565b604051602081830303815290604052805190602001201415610c7b576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c7290611903565b60405180910390fd5b6001600382604051610c8d91906116f3565b908152602001604051809103902060006101000a81548160ff02191690831515021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610d43576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d3a906119a3565b60405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008111610dc9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dc090611943565b60405180910390fd5b600383604051610dd991906116f3565b908152602001604051809103902060009054906101000a900460ff16610e34576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e2b90611963565b60405180910390fd5b60008473ffffffffffffffffffffffffffffffffffffffff16333084604051602401610e629392919061173c565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610eec91906116dc565b6000604051808303816000865af19150503d8060008114610f29576040519150601f19603f3d011682016040523d82523d6000602084013e610f2e565b606091505b5050905080610f72576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f6990611983565b60405180910390fd5b7f889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e5418585858585604051610fa99594939291906117b8565b60405180910390a15050505050565b6000600382604051610fca91906116f3565b908152602001604051809103902060009054906101000a900460ff169050919050565b6000611000610ffb84611a4f565b611a1e565b90508281526020810184848401111561101857600080fd5b611023848285611ba3565b509392505050565b60008135905061103a81611cb5565b92915050565b600082601f83011261105157600080fd5b8135611061848260208601610fed565b91505092915050565b60008135905061107981611ccc565b92915050565b60008151905061108e81611ccc565b92915050565b6000602082840312156110a657600080fd5b60006110b48482850161102b565b91505092915050565b6000806000606084860312156110d257600080fd5b60006110e08682870161102b565b93505060206110f18682870161102b565b92505060406111028682870161106a565b9150509250925092565b6000806000806080858703121561112257600080fd5b60006111308782880161102b565b945050602085013567ffffffffffffffff81111561114d57600080fd5b61115987828801611040565b935050604085013567ffffffffffffffff81111561117657600080fd5b61118287828801611040565b92505060606111938782880161106a565b91505092959194509250565b6000602082840312156111b157600080fd5b600082013567ffffffffffffffff8111156111cb57600080fd5b6111d784828501611040565b91505092915050565b600080604083850312156111f357600080fd5b600083013567ffffffffffffffff81111561120d57600080fd5b61121985828601611040565b925050602061122a8582860161102b565b9150509250929050565b60008060006060848603121561124957600080fd5b600084013567ffffffffffffffff81111561126357600080fd5b61126f86828701611040565b93505060206112808682870161102b565b92505060406112918682870161106a565b9150509250925092565b600080600080608085870312156112b157600080fd5b600085013567ffffffffffffffff8111156112cb57600080fd5b6112d787828801611040565b945050602085013567ffffffffffffffff8111156112f457600080fd5b61130087828801611040565b935050604085013567ffffffffffffffff81111561131d57600080fd5b61132987828801611040565b925050606061133a8782880161106a565b91505092959194509250565b60006020828403121561135857600080fd5b60006113668482850161107f565b91505092915050565b61137881611b5b565b82525050565b61138781611b6d565b82525050565b600061139882611a94565b6113a28185611aaa565b93506113b2818560208601611bb2565b80840191505092915050565b60006113c982611a9f565b6113d38185611ab5565b93506113e3818560208601611bb2565b6113ec81611ca4565b840191505092915050565b600061140282611a9f565b61140c8185611ac6565b935061141c818560208601611bb2565b80840191505092915050565b6000815461143581611be5565b61143f8186611ac6565b9450600182166000811461145a576001811461146b5761149e565b60ff1983168652818601935061149e565b61147485611a7f565b60005b8381101561149657815481890152600182019150602081019050611477565b838801955050505b50505092915050565b60006114b4602083611ab5565b91507f4e657720636861696e206d757374206e6f74206265207468697320636861696e6000830152602082019050919050565b60006114f4602a83611ab5565b91507f42616c616e6365206c657373207468616e20616d6f756e74206265696e67207460008301527f72616e73666572726564000000000000000000000000000000000000000000006020830152604082019050919050565b600061155a601d83611ab5565b91507f416d6f756e74206d7573742062652067726561746572207468616e20300000006000830152602082019050919050565b600061159a601183611ab5565b91507f436861696e206e6f7420616c6c6f7765640000000000000000000000000000006000830152602082019050919050565b60006115da601283611ab5565b91507f4661696c656420746f207472616e7366657200000000000000000000000000006000830152602082019050919050565b600061161a600d83611ab5565b91507f4d757374206265206f776e6572000000000000000000000000000000000000006000830152602082019050919050565b600061165a601583611ab5565b91507f4661696c656420746f206765742062616c616e636500000000000000000000006000830152602082019050919050565b600061169a601183611ab5565b91507f4e6f7420656e6f75676820746f6b656e730000000000000000000000000000006000830152602082019050919050565b6116d681611b99565b82525050565b60006116e8828461138d565b915081905092915050565b60006116ff82846113f7565b915081905092915050565b60006117168284611428565b915081905092915050565b6000602082019050611736600083018461136f565b92915050565b6000606082019050611751600083018661136f565b61175e602083018561136f565b61176b60408301846116cd565b949350505050565b6000608082019050611788600083018761136f565b611795602083018661136f565b6117a260408301856116cd565b6117af606083018461137e565b95945050505050565b600060a0820190506117cd600083018861136f565b81810360208301526117df81876113be565b905081810360408301526117f381866113be565b905061180260608301856116cd565b61180f608083018461137e565b9695505050505050565b600060408201905061182e600083018561136f565b61183b60208301846116cd565b9392505050565b6000602082019050611857600083018461137e565b92915050565b600060a082019050818103600083015261187781886113be565b9050611886602083018761136f565b818103604083015261189881866113be565b905081810360608301526118ac81856113be565b90506118bb60808301846116cd565b9695505050505050565b600060608201905081810360008301526118df81866113be565b90506118ee602083018561136f565b6118fb60408301846116cd565b949350505050565b6000602082019050818103600083015261191c816114a7565b9050919050565b6000602082019050818103600083015261193c816114e7565b9050919050565b6000602082019050818103600083015261195c8161154d565b9050919050565b6000602082019050818103600083015261197c8161158d565b9050919050565b6000602082019050818103600083015261199c816115cd565b9050919050565b600060208201905081810360008301526119bc8161160d565b9050919050565b600060208201905081810360008301526119dc8161164d565b9050919050565b600060208201905081810360008301526119fc8161168d565b9050919050565b6000602082019050611a1860008301846116cd565b92915050565b6000604051905081810181811067ffffffffffffffff82111715611a4557611a44611c75565b5b8060405250919050565b600067ffffffffffffffff821115611a6a57611a69611c75565b5b601f19601f8301169050602081019050919050565b60008190508160005260206000209050919050565b600081519050919050565b600081519050919050565b600081905092915050565b600082825260208201905092915050565b600081905092915050565b6000611adc82611b99565b9150611ae783611b99565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611b1c57611b1b611c17565b5b828201905092915050565b6000611b3282611b99565b9150611b3d83611b99565b925082821015611b5057611b4f611c17565b5b828203905092915050565b6000611b6682611b79565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015611bd0578082015181840152602081019050611bb5565b83811115611bdf576000848401525b50505050565b60006002820490506001821680611bfd57607f821691505b60208210811415611c1157611c10611c46565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b611cbe81611b5b565b8114611cc957600080fd5b50565b611cd581611b99565b8114611ce057600080fd5b5056fea2646970667358221220b1f535ec92c49aca9d5327804e7c558872d85820b6caccb4d57e6fd35e099cda64736f6c63430008000033"

// DeployErc20gateway deploys a new Ethereum contract, binding an instance of Erc20gateway to it.
func DeployErc20gateway(auth *bind.TransactOpts, backend bind.ContractBackend, chain string, chains []string) (common.Address, *types.Transaction, *Erc20gateway, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20gatewayABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(Erc20gatewayBin), backend, chain, chains)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc20gateway{Erc20gatewayCaller: Erc20gatewayCaller{contract: contract}, Erc20gatewayTransactor: Erc20gatewayTransactor{contract: contract}, Erc20gatewayFilterer: Erc20gatewayFilterer{contract: contract}}, nil
}

// Erc20gateway is an auto generated Go binding around an Ethereum contract.
type Erc20gateway struct {
	Erc20gatewayCaller     // Read-only binding to the contract
	Erc20gatewayTransactor // Write-only binding to the contract
	Erc20gatewayFilterer   // Log filterer for contract events
}

// Erc20gatewayCaller is an auto generated read-only Go binding around an Ethereum contract.
type Erc20gatewayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gatewayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Erc20gatewayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gatewayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Erc20gatewayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Erc20gatewaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Erc20gatewaySession struct {
	Contract     *Erc20gateway     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Erc20gatewayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Erc20gatewayCallerSession struct {
	Contract *Erc20gatewayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// Erc20gatewayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Erc20gatewayTransactorSession struct {
	Contract     *Erc20gatewayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// Erc20gatewayRaw is an auto generated low-level Go binding around an Ethereum contract.
type Erc20gatewayRaw struct {
	Contract *Erc20gateway // Generic contract binding to access the raw methods on
}

// Erc20gatewayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Erc20gatewayCallerRaw struct {
	Contract *Erc20gatewayCaller // Generic read-only contract binding to access the raw methods on
}

// Erc20gatewayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Erc20gatewayTransactorRaw struct {
	Contract *Erc20gatewayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewErc20gateway creates a new instance of Erc20gateway, bound to a specific deployed contract.
func NewErc20gateway(address common.Address, backend bind.ContractBackend) (*Erc20gateway, error) {
	contract, err := bindErc20gateway(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Erc20gateway{Erc20gatewayCaller: Erc20gatewayCaller{contract: contract}, Erc20gatewayTransactor: Erc20gatewayTransactor{contract: contract}, Erc20gatewayFilterer: Erc20gatewayFilterer{contract: contract}}, nil
}

// NewErc20gatewayCaller creates a new read-only instance of Erc20gateway, bound to a specific deployed contract.
func NewErc20gatewayCaller(address common.Address, caller bind.ContractCaller) (*Erc20gatewayCaller, error) {
	contract, err := bindErc20gateway(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayCaller{contract: contract}, nil
}

// NewErc20gatewayTransactor creates a new write-only instance of Erc20gateway, bound to a specific deployed contract.
func NewErc20gatewayTransactor(address common.Address, transactor bind.ContractTransactor) (*Erc20gatewayTransactor, error) {
	contract, err := bindErc20gateway(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransactor{contract: contract}, nil
}

// NewErc20gatewayFilterer creates a new log filterer instance of Erc20gateway, bound to a specific deployed contract.
func NewErc20gatewayFilterer(address common.Address, filterer bind.ContractFilterer) (*Erc20gatewayFilterer, error) {
	contract, err := bindErc20gateway(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayFilterer{contract: contract}, nil
}

// bindErc20gateway binds a generic wrapper to an already deployed contract.
func bindErc20gateway(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Erc20gatewayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20gateway *Erc20gatewayRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20gateway.Contract.Erc20gatewayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20gateway *Erc20gatewayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gateway.Contract.Erc20gatewayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20gateway *Erc20gatewayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20gateway.Contract.Erc20gatewayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Erc20gateway *Erc20gatewayCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Erc20gateway.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Erc20gateway *Erc20gatewayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gateway.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Erc20gateway *Erc20gatewayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Erc20gateway.Contract.contract.Transact(opts, method, params...)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20gateway *Erc20gatewayCaller) GetBalance(opts *bind.CallOpts, assetId string, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "getBalance", assetId, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20gateway *Erc20gatewaySession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _Erc20gateway.Contract.GetBalance(&_Erc20gateway.CallOpts, assetId, account)
}

// GetBalance is a free data retrieval call binding the contract method 0x89f4b317.
//
// Solidity: function getBalance(string assetId, address account) view returns(uint256)
func (_Erc20gateway *Erc20gatewayCallerSession) GetBalance(assetId string, account common.Address) (*big.Int, error) {
	return _Erc20gateway.Contract.GetBalance(&_Erc20gateway.CallOpts, assetId, account)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20gateway *Erc20gatewayCaller) GetOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "getOwner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20gateway *Erc20gatewaySession) GetOwner() (common.Address, error) {
	return _Erc20gateway.Contract.GetOwner(&_Erc20gateway.CallOpts)
}

// GetOwner is a free data retrieval call binding the contract method 0x893d20e8.
//
// Solidity: function getOwner() view returns(address)
func (_Erc20gateway *Erc20gatewayCallerSession) GetOwner() (common.Address, error) {
	return _Erc20gateway.Contract.GetOwner(&_Erc20gateway.CallOpts)
}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20gateway *Erc20gatewayCaller) IsChainAllowed(opts *bind.CallOpts, chain string) (bool, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "isChainAllowed", chain)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20gateway *Erc20gatewaySession) IsChainAllowed(chain string) (bool, error) {
	return _Erc20gateway.Contract.IsChainAllowed(&_Erc20gateway.CallOpts, chain)
}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20gateway *Erc20gatewayCallerSession) IsChainAllowed(chain string) (bool, error) {
	return _Erc20gateway.Contract.IsChainAllowed(&_Erc20gateway.CallOpts, chain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20gateway *Erc20gatewayTransactor) AddAllowedChain(opts *bind.TransactOpts, newChain string) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "addAllowedChain", newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20gateway *Erc20gatewaySession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.AddAllowedChain(&_Erc20gateway.TransactOpts, newChain)
}

// AddAllowedChain is a paid mutator transaction binding the contract method 0x9a12bb60.
//
// Solidity: function addAllowedChain(string newChain) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) AddAllowedChain(newChain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.AddAllowedChain(&_Erc20gateway.TransactOpts, newChain)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20gateway *Erc20gatewayTransactor) ChangeOwner(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "changeOwner", newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20gateway *Erc20gatewaySession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.ChangeOwner(&_Erc20gateway.TransactOpts, newOwner)
}

// ChangeOwner is a paid mutator transaction binding the contract method 0xa6f9dae1.
//
// Solidity: function changeOwner(address newOwner) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) ChangeOwner(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.ChangeOwner(&_Erc20gateway.TransactOpts, newOwner)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferIn(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferIn", assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferIn(&_Erc20gateway.TransactOpts, assetId, recipient, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x0832f134.
//
// Solidity: function transferIn(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferIn(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferIn(&_Erc20gateway.TransactOpts, assetId, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferInAssetOfThisChain(opts *bind.TransactOpts, assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferInAssetOfThisChain", assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferInAssetOfThisChain(&_Erc20gateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferInAssetOfThisChain is a paid mutator transaction binding the contract method 0x6c189db8.
//
// Solidity: function transferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferInAssetOfThisChain(assetAddr common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferInAssetOfThisChain(&_Erc20gateway.TransactOpts, assetAddr, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferOut(opts *bind.TransactOpts, assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferOut", assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOut(&_Erc20gateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x810b9f12.
//
// Solidity: function transferOut(string assetId, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferOut(assetId string, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOut(&_Erc20gateway.TransactOpts, assetId, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferOutFromContract(opts *bind.TransactOpts, assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferOutFromContract", assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOutFromContract(&_Erc20gateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferOutFromContract is a paid mutator transaction binding the contract method 0xbf6022be.
//
// Solidity: function transferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferOutFromContract(assetAddr common.Address, toChain string, recipient string, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOutFromContract(&_Erc20gateway.TransactOpts, assetAddr, toChain, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferWithin(opts *bind.TransactOpts, assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferWithin", assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferWithin(&_Erc20gateway.TransactOpts, assetId, recipient, amount)
}

// TransferWithin is a paid mutator transaction binding the contract method 0x7dd1f364.
//
// Solidity: function transferWithin(string assetId, address recipient, uint256 amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferWithin(assetId string, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferWithin(&_Erc20gateway.TransactOpts, assetId, recipient, amount)
}

// Erc20gatewayTransferInIterator is returned from FilterTransferIn and is used to iterate over the raw logs and unpacked data for TransferIn events raised by the Erc20gateway contract.
type Erc20gatewayTransferInIterator struct {
	Event *Erc20gatewayTransferIn // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferInIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferIn)
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
		it.Event = new(Erc20gatewayTransferIn)
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
func (it *Erc20gatewayTransferInIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferInIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferIn represents a TransferIn event raised by the Erc20gateway contract.
type Erc20gatewayTransferIn struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferIn is a free log retrieval operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferIn(opts *bind.FilterOpts) (*Erc20gatewayTransferInIterator, error) {

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferInIterator{contract: _Erc20gateway.contract, event: "TransferIn", logs: logs, sub: sub}, nil
}

// WatchTransferIn is a free log subscription operation binding the contract event 0xbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea094.
//
// Solidity: event TransferIn(string assetId, address recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferIn(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferIn) (event.Subscription, error) {

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferIn")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferIn)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferIn(log types.Log) (*Erc20gatewayTransferIn, error) {
	event := new(Erc20gatewayTransferIn)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferIn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferInAssetOfThisChainIterator is returned from FilterTransferInAssetOfThisChain and is used to iterate over the raw logs and unpacked data for TransferInAssetOfThisChain events raised by the Erc20gateway contract.
type Erc20gatewayTransferInAssetOfThisChainIterator struct {
	Event *Erc20gatewayTransferInAssetOfThisChain // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferInAssetOfThisChainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferInAssetOfThisChain)
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
		it.Event = new(Erc20gatewayTransferInAssetOfThisChain)
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
func (it *Erc20gatewayTransferInAssetOfThisChainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferInAssetOfThisChainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferInAssetOfThisChain represents a TransferInAssetOfThisChain event raised by the Erc20gateway contract.
type Erc20gatewayTransferInAssetOfThisChain struct {
	AssetAddr common.Address
	Recipient common.Address
	Amount    *big.Int
	Success   bool
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInAssetOfThisChain is a free log retrieval operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferInAssetOfThisChain(opts *bind.FilterOpts) (*Erc20gatewayTransferInAssetOfThisChainIterator, error) {

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferInAssetOfThisChainIterator{contract: _Erc20gateway.contract, event: "TransferInAssetOfThisChain", logs: logs, sub: sub}, nil
}

// WatchTransferInAssetOfThisChain is a free log subscription operation binding the contract event 0xc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3.
//
// Solidity: event TransferInAssetOfThisChain(address assetAddr, address recipient, uint256 amount, bool success)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferInAssetOfThisChain(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferInAssetOfThisChain) (event.Subscription, error) {

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferInAssetOfThisChain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferInAssetOfThisChain)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferInAssetOfThisChain(log types.Log) (*Erc20gatewayTransferInAssetOfThisChain, error) {
	event := new(Erc20gatewayTransferInAssetOfThisChain)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferInAssetOfThisChain", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferOutIterator is returned from FilterTransferOut and is used to iterate over the raw logs and unpacked data for TransferOut events raised by the Erc20gateway contract.
type Erc20gatewayTransferOutIterator struct {
	Event *Erc20gatewayTransferOut // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferOutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferOut)
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
		it.Event = new(Erc20gatewayTransferOut)
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
func (it *Erc20gatewayTransferOutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferOutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferOut represents a TransferOut event raised by the Erc20gateway contract.
type Erc20gatewayTransferOut struct {
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
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferOut(opts *bind.FilterOpts) (*Erc20gatewayTransferOutIterator, error) {

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferOutIterator{contract: _Erc20gateway.contract, event: "TransferOut", logs: logs, sub: sub}, nil
}

// WatchTransferOut is a free log subscription operation binding the contract event 0x265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f.
//
// Solidity: event TransferOut(string assetId, address from, string toChain, string recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferOut(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferOut) (event.Subscription, error) {

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferOut")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferOut)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferOut(log types.Log) (*Erc20gatewayTransferOut, error) {
	event := new(Erc20gatewayTransferOut)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferOut", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferOutFromContractIterator is returned from FilterTransferOutFromContract and is used to iterate over the raw logs and unpacked data for TransferOutFromContract events raised by the Erc20gateway contract.
type Erc20gatewayTransferOutFromContractIterator struct {
	Event *Erc20gatewayTransferOutFromContract // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferOutFromContractIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferOutFromContract)
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
		it.Event = new(Erc20gatewayTransferOutFromContract)
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
func (it *Erc20gatewayTransferOutFromContractIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferOutFromContractIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferOutFromContract represents a TransferOutFromContract event raised by the Erc20gateway contract.
type Erc20gatewayTransferOutFromContract struct {
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
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferOutFromContract(opts *bind.FilterOpts) (*Erc20gatewayTransferOutFromContractIterator, error) {

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferOutFromContractIterator{contract: _Erc20gateway.contract, event: "TransferOutFromContract", logs: logs, sub: sub}, nil
}

// WatchTransferOutFromContract is a free log subscription operation binding the contract event 0x889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e541.
//
// Solidity: event TransferOutFromContract(address assetAddr, string toChain, string recipient, uint256 amount, bool success)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferOutFromContract(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferOutFromContract) (event.Subscription, error) {

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferOutFromContract")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferOutFromContract)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferOutFromContract(log types.Log) (*Erc20gatewayTransferOutFromContract, error) {
	event := new(Erc20gatewayTransferOutFromContract)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferOutFromContract", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferWithinIterator is returned from FilterTransferWithin and is used to iterate over the raw logs and unpacked data for TransferWithin events raised by the Erc20gateway contract.
type Erc20gatewayTransferWithinIterator struct {
	Event *Erc20gatewayTransferWithin // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferWithinIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferWithin)
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
		it.Event = new(Erc20gatewayTransferWithin)
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
func (it *Erc20gatewayTransferWithinIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferWithinIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferWithin represents a TransferWithin event raised by the Erc20gateway contract.
type Erc20gatewayTransferWithin struct {
	AssetId   string
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferWithin is a free log retrieval operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferWithin(opts *bind.FilterOpts) (*Erc20gatewayTransferWithinIterator, error) {

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferWithinIterator{contract: _Erc20gateway.contract, event: "TransferWithin", logs: logs, sub: sub}, nil
}

// WatchTransferWithin is a free log subscription operation binding the contract event 0x8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df4.
//
// Solidity: event TransferWithin(string assetId, address recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferWithin(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferWithin) (event.Subscription, error) {

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferWithin")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferWithin)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferWithin(log types.Log) (*Erc20gatewayTransferWithin, error) {
	event := new(Erc20gatewayTransferWithin)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferWithin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
