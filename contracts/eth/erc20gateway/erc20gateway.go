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
const Erc20GatewayABI = "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferIn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferInAssetOfThisChain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"success\",\"type\":\"bool\"}],\"name\":\"TransferOutFromContract\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferWithin\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newChain\",\"type\":\"string\"}],\"name\":\"addAllowedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"changeOwner\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"getBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getOwner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"isChainAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInAssetOfThisChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"assetAddr\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"toChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOutFromContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"assetId\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferWithin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Erc20GatewayBin is the compiled bytecode used for deploying new contracts.
var Erc20GatewayBin = "0x60806040523480156200001157600080fd5b5060405162001e5c38038062001e5c8339818101604052810190620000379190620001b9565b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600190805190602001906200008f92919062000097565b50506200032f565b828054620000a5906200029b565b90600052602060002090601f016020900481019282620000c9576000855562000115565b82601f10620000e457805160ff191683800117855562000115565b8280016001018555821562000115579182015b8281111562000114578251825591602001919060010190620000f7565b5b50905062000124919062000128565b5090565b5b808211156200014357600081600090555060010162000129565b5090565b60006200015e620001588462000232565b620001fe565b9050828152602081018484840111156200017757600080fd5b6200018484828562000265565b509392505050565b600082601f8301126200019e57600080fd5b8151620001b084826020860162000147565b91505092915050565b600060208284031215620001cc57600080fd5b600082015167ffffffffffffffff811115620001e757600080fd5b620001f5848285016200018c565b91505092915050565b6000604051905081810181811067ffffffffffffffff8211171562000228576200022762000300565b5b8060405250919050565b600067ffffffffffffffff82111562000250576200024f62000300565b5b601f19601f8301169050602081019050919050565b60005b838110156200028557808201518184015260208101905062000268565b8381111562000295576000848401525b50505050565b60006002820490506001821680620002b457607f821691505b60208210811415620002cb57620002ca620002d1565b5b50919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611b1d806200033f6000396000f3fe608060405234801561001057600080fd5b506004361061009e5760003560e01c806389f4b3171161006657806389f4b317146101315780639a12bb6014610161578063a6f9dae11461017d578063bf6022be14610199578063d29fde24146101b55761009e565b80630832f134146100a35780636c189db8146100bf5780637dd1f364146100db578063810b9f12146100f7578063893d20e814610113575b600080fd5b6100bd60048036038101906100b891906111a4565b6101e5565b005b6100d960048036038101906100d4919061102d565b610325565b005b6100f560048036038101906100f091906111a4565b6106c6565b005b610111600480360381019061010c919061120b565b6108cb565b005b61011b610ace565b60405161012891906115bb565b60405180910390f35b61014b60048036038101906101469190611150565b610af7565b604051610158919061187d565b60405180910390f35b61017b6004803603810190610176919061110f565b610b5d565b005b61019760048036038101906101929190611004565b610c25565b005b6101b360048036038101906101ae919061107c565b610cf6565b005b6101cf60048036038101906101ca919061110f565b610f28565b6040516101dc91906116dc565b60405180910390f35b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610273576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161026a9061181d565b60405180910390fd5b8060028460405161028491906115a4565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546102de9190611936565b925050819055507fbd5741014816dfb8e25ad23c75773cd9d3609adb1f701eb45b44452a2c3ea0948383836040516103189392919061175f565b60405180910390a1505050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16146103b3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103aa9061181d565b60405180910390fd5b6000808473ffffffffffffffffffffffffffffffffffffffff16306040516024016103de91906115bb565b6040516020818303038152906040527f70a08231000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610468919061158d565b6000604051808303816000865af19150503d80600081146104a5576040519150601f19603f3d011682016040523d82523d6000602084013e6104aa565b606091505b5091509150816104ef576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104e69061183d565b60405180910390fd5b828180602001905181019061050491906112b6565b1015610545576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161053c9061185d565b60405180910390fd5b60008573ffffffffffffffffffffffffffffffffffffffff1685856040516024016105719291906116b3565b6040516020818303038152906040527fa9059cbb000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506040516105fb919061158d565b6000604051808303816000865af19150503d8060008114610638576040519150601f19603f3d011682016040523d82523d6000602084013e61063d565b606091505b5050905080610681576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610678906117fd565b60405180910390fd5b7fc9c5da88430bb17ba62a65a8e9e0089011884793e67cfe276d85052c0aa711c3868686846040516106b6949392919061160d565b60405180910390a1505050505050565b60008111610709576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610700906117bd565b60405180910390fd5b8060028460405161071a91906115a4565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156107a7576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161079e9061179d565b60405180910390fd5b806002846040516107b891906115a4565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610812919061198c565b925050819055508060028460405161082a91906115a4565b908152602001604051809103902060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546108849190611936565b925050819055507f8622a0377434878052d354ceb356637058303a6c8c91314a0d2a82e3a4a55df48383836040516108be9392919061175f565b60405180910390a1505050565b6000811161090e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610905906117bd565b60405180910390fd5b8060028560405161091f91906115a4565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410156109ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109a39061185d565b60405180910390fd5b6003836040516109bc91906115a4565b908152602001604051809103902060009054906101000a900460ff16610a17576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a0e906117dd565b60405180910390fd5b80600285604051610a2891906115a4565b908152602001604051809103902060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000206000828254610a82919061198c565b925050819055507f265e08033322072b661d58b72599cd7c92f4fcf44da508cf9030e0f6b295f57f8433858585604051610ac09594939291906116f7565b60405180910390a150505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6000600283604051610b0991906115a4565b908152602001604051809103902060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610beb576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610be29061181d565b60405180910390fd5b6001600382604051610bfd91906115a4565b908152602001604051809103902060006101000a81548160ff02191690831515021790555050565b60008054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610cb3576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610caa9061181d565b60405180910390fd5b806000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008111610d39576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d30906117bd565b60405180910390fd5b600383604051610d4991906115a4565b908152602001604051809103902060009054906101000a900460ff16610da4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d9b906117dd565b60405180910390fd5b60008473ffffffffffffffffffffffffffffffffffffffff16333084604051602401610dd2939291906115d6565b6040516020818303038152906040527f23b872dd000000000000000000000000000000000000000000000000000000007bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610e5c919061158d565b6000604051808303816000865af19150503d8060008114610e99576040519150601f19603f3d011682016040523d82523d6000602084013e610e9e565b606091505b5050905080610ee2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ed9906117fd565b60405180910390fd5b7f889398e2ef5cc790a1063ee34969e83ec5cb197819219fae1eceed6544c8e5418585858585604051610f19959493929190611652565b60405180910390a15050505050565b6000600382604051610f3a91906115a4565b908152602001604051809103902060009054906101000a900460ff169050919050565b6000610f70610f6b846118c9565b611898565b905082815260208101848484011115610f8857600080fd5b610f93848285611a08565b509392505050565b600081359050610faa81611ab9565b92915050565b600082601f830112610fc157600080fd5b8135610fd1848260208601610f5d565b91505092915050565b600081359050610fe981611ad0565b92915050565b600081519050610ffe81611ad0565b92915050565b60006020828403121561101657600080fd5b600061102484828501610f9b565b91505092915050565b60008060006060848603121561104257600080fd5b600061105086828701610f9b565b935050602061106186828701610f9b565b925050604061107286828701610fda565b9150509250925092565b6000806000806080858703121561109257600080fd5b60006110a087828801610f9b565b945050602085013567ffffffffffffffff8111156110bd57600080fd5b6110c987828801610fb0565b935050604085013567ffffffffffffffff8111156110e657600080fd5b6110f287828801610fb0565b925050606061110387828801610fda565b91505092959194509250565b60006020828403121561112157600080fd5b600082013567ffffffffffffffff81111561113b57600080fd5b61114784828501610fb0565b91505092915050565b6000806040838503121561116357600080fd5b600083013567ffffffffffffffff81111561117d57600080fd5b61118985828601610fb0565b925050602061119a85828601610f9b565b9150509250929050565b6000806000606084860312156111b957600080fd5b600084013567ffffffffffffffff8111156111d357600080fd5b6111df86828701610fb0565b93505060206111f086828701610f9b565b925050604061120186828701610fda565b9150509250925092565b6000806000806080858703121561122157600080fd5b600085013567ffffffffffffffff81111561123b57600080fd5b61124787828801610fb0565b945050602085013567ffffffffffffffff81111561126457600080fd5b61127087828801610fb0565b935050604085013567ffffffffffffffff81111561128d57600080fd5b61129987828801610fb0565b92505060606112aa87828801610fda565b91505092959194509250565b6000602082840312156112c857600080fd5b60006112d684828501610fef565b91505092915050565b6112e8816119c0565b82525050565b6112f7816119d2565b82525050565b6000611308826118f9565b611312818561190f565b9350611322818560208601611a17565b80840191505092915050565b600061133982611904565b611343818561191a565b9350611353818560208601611a17565b61135c81611aa8565b840191505092915050565b600061137282611904565b61137c818561192b565b935061138c818560208601611a17565b80840191505092915050565b60006113a5602a8361191a565b91507f42616c616e6365206c657373207468616e20616d6f756e74206265696e67207460008301527f72616e73666572726564000000000000000000000000000000000000000000006020830152604082019050919050565b600061140b601d8361191a565b91507f416d6f756e74206d7573742062652067726561746572207468616e20300000006000830152602082019050919050565b600061144b60118361191a565b91507f436861696e206e6f7420616c6c6f7765640000000000000000000000000000006000830152602082019050919050565b600061148b60128361191a565b91507f4661696c656420746f207472616e7366657200000000000000000000000000006000830152602082019050919050565b60006114cb600d8361191a565b91507f4d757374206265206f776e6572000000000000000000000000000000000000006000830152602082019050919050565b600061150b60158361191a565b91507f4661696c656420746f206765742062616c616e636500000000000000000000006000830152602082019050919050565b600061154b60118361191a565b91507f4e6f7420656e6f75676820746f6b656e730000000000000000000000000000006000830152602082019050919050565b611587816119fe565b82525050565b600061159982846112fd565b915081905092915050565b60006115b08284611367565b915081905092915050565b60006020820190506115d060008301846112df565b92915050565b60006060820190506115eb60008301866112df565b6115f860208301856112df565b611605604083018461157e565b949350505050565b600060808201905061162260008301876112df565b61162f60208301866112df565b61163c604083018561157e565b61164960608301846112ee565b95945050505050565b600060a08201905061166760008301886112df565b8181036020830152611679818761132e565b9050818103604083015261168d818661132e565b905061169c606083018561157e565b6116a960808301846112ee565b9695505050505050565b60006040820190506116c860008301856112df565b6116d5602083018461157e565b9392505050565b60006020820190506116f160008301846112ee565b92915050565b600060a0820190508181036000830152611711818861132e565b905061172060208301876112df565b8181036040830152611732818661132e565b90508181036060830152611746818561132e565b9050611755608083018461157e565b9695505050505050565b60006060820190508181036000830152611779818661132e565b905061178860208301856112df565b611795604083018461157e565b949350505050565b600060208201905081810360008301526117b681611398565b9050919050565b600060208201905081810360008301526117d6816113fe565b9050919050565b600060208201905081810360008301526117f68161143e565b9050919050565b600060208201905081810360008301526118168161147e565b9050919050565b60006020820190508181036000830152611836816114be565b9050919050565b60006020820190508181036000830152611856816114fe565b9050919050565b600060208201905081810360008301526118768161153e565b9050919050565b6000602082019050611892600083018461157e565b92915050565b6000604051905081810181811067ffffffffffffffff821117156118bf576118be611a79565b5b8060405250919050565b600067ffffffffffffffff8211156118e4576118e3611a79565b5b601f19601f8301169050602081019050919050565b600081519050919050565b600081519050919050565b600081905092915050565b600082825260208201905092915050565b600081905092915050565b6000611941826119fe565b915061194c836119fe565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561198157611980611a4a565b5b828201905092915050565b6000611997826119fe565b91506119a2836119fe565b9250828210156119b5576119b4611a4a565b5b828203905092915050565b60006119cb826119de565b9050919050565b60008115159050919050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b82818337600083830152505050565b60005b83811015611a35578082015181840152602081019050611a1a565b83811115611a44576000848401525b50505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6000601f19601f8301169050919050565b611ac2816119c0565b8114611acd57600080fd5b50565b611ad9816119fe565b8114611ae457600080fd5b5056fea2646970667358221220a0f3cffde07e3568f7bd395a70a5f177d9367cd9806bf6584a142465411532f064736f6c63430008000033"

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

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20Gateway *Erc20GatewayCaller) IsChainAllowed(opts *bind.CallOpts, chain string) (bool, error) {
	var out []interface{}
	err := _Erc20Gateway.contract.Call(opts, &out, "isChainAllowed", chain)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20Gateway *Erc20GatewaySession) IsChainAllowed(chain string) (bool, error) {
	return _Erc20Gateway.Contract.IsChainAllowed(&_Erc20Gateway.CallOpts, chain)
}

// IsChainAllowed is a free data retrieval call binding the contract method 0xd29fde24.
//
// Solidity: function isChainAllowed(string chain) view returns(bool)
func (_Erc20Gateway *Erc20GatewayCallerSession) IsChainAllowed(chain string) (bool, error) {
	return _Erc20Gateway.Contract.IsChainAllowed(&_Erc20Gateway.CallOpts, chain)
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
