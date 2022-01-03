// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package erc20gateway

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
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_supportedChains\",\"type\":\"string[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"AddSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"RemoveSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferInEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"destChain\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOutEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"AddSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"PauseGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"RemoveSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ResumeGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"TransferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_destChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"TransferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"supportedChains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001e3938038062001e39833981810160405281019062000037919062000456565b620000576200004b620000f160201b60201c565b620000f960201b60201c565b60008060146101000a81548160ff02191690831515021790555060005b8151811015620000e957600180838381518110620000975762000096620004a7565b5b6020026020010151604051620000ae919062000523565b908152602001604051809103902060006101000a81548160ff0219169083151502179055508080620000e09062000575565b91505062000074565b5050620005c3565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200022182620001d6565b810181811067ffffffffffffffff82111715620002435762000242620001e7565b5b80604052505050565b600062000258620001bd565b905062000266828262000216565b919050565b600067ffffffffffffffff821115620002895762000288620001e7565b5b602082029050602081019050919050565b600080fd5b600080fd5b600067ffffffffffffffff821115620002c257620002c1620001e7565b5b620002cd82620001d6565b9050602081019050919050565b60005b83811015620002fa578082015181840152602081019050620002dd565b838111156200030a576000848401525b50505050565b6000620003276200032184620002a4565b6200024c565b9050828152602081018484840111156200034657620003456200029f565b5b62000353848285620002da565b509392505050565b600082601f830112620003735762000372620001d1565b5b81516200038584826020860162000310565b91505092915050565b6000620003a56200039f846200026b565b6200024c565b90508083825260208201905060208402830185811115620003cb57620003ca6200029a565b5b835b818110156200041957805167ffffffffffffffff811115620003f457620003f3620001d1565b5b8086016200040389826200035b565b85526020850194505050602081019050620003cd565b5050509392505050565b600082601f8301126200043b576200043a620001d1565b5b81516200044d8482602086016200038e565b91505092915050565b6000602082840312156200046f576200046e620001c7565b5b600082015167ffffffffffffffff81111562000490576200048f620001cc565b5b6200049e8482850162000423565b91505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081519050919050565b600081905092915050565b6000620004f982620004d6565b620005058185620004e1565b935062000517818560208601620002da565b80840191505092915050565b6000620005318284620004ec565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600062000582826200056b565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620005b857620005b76200053c565b5b600182019050919050565b61186680620005d36000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c80638da5cb5b116100715780638da5cb5b1461013e578063ca569dbf1461015c578063ecf62f8214610166578063f1897eb514610170578063f2fde38b1461018c578063fc69a67a146101a8576100a9565b806358b67fe1146100ae5780636c30aaa2146100ca578063715018a6146100fa5780638456cb59146101045780638ab008cb14610122575b600080fd5b6100c860048036038101906100c39190610f79565b6101c4565b005b6100e460048036038101906100df9190610f79565b6102bc565b6040516100f19190610fdd565b60405180910390f35b6101026102f2565b005b61010c61037a565b6040516101199190610fdd565b60405180910390f35b61013c6004803603810190610137919061108c565b61038d565b005b610146610596565b60405161015391906110ee565b60405180910390f35b6101646105bf565b005b61016e6106ae565b005b61018a60048036038101906101859190611109565b61079c565b005b6101a660048036038101906101a191906111a0565b6108f5565b005b6101c260048036038101906101bd9190610f79565b6109ed565b005b6101cc610ae4565b73ffffffffffffffffffffffffffffffffffffffff166101ea610596565b73ffffffffffffffffffffffffffffffffffffffff1614610240576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102379061122a565b60405180910390fd5b600060018260405161025291906112c4565b908152602001604051809103902060006101000a81548160ff0219169083151502179055508060405161028591906112c4565b60405180910390207ff300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca63560405160405180910390a250565b6001818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900460ff1681565b6102fa610ae4565b73ffffffffffffffffffffffffffffffffffffffff16610318610596565b73ffffffffffffffffffffffffffffffffffffffff161461036e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103659061122a565b60405180910390fd5b6103786000610aec565b565b600060149054906101000a900460ff1681565b610395610ae4565b73ffffffffffffffffffffffffffffffffffffffff166103b3610596565b73ffffffffffffffffffffffffffffffffffffffff1614610409576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016104009061122a565b60405180910390fd5b60001515600060149054906101000a900460ff1615151461045f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161045690611327565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161049a91906110ee565b602060405180830381865afa1580156104b7573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906104db919061135c565b905081811015610520576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610517906113fb565b60405180910390fd5b61052b848484610bb0565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c84604051610588919061142a565b60405180910390a350505050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6105c7610ae4565b73ffffffffffffffffffffffffffffffffffffffff166105e5610596565b73ffffffffffffffffffffffffffffffffffffffff161461063b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106329061122a565b60405180910390fd5b60001515600060149054906101000a900460ff16151514610691576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161068890611491565b60405180910390fd5b6001600060146101000a81548160ff021916908315150217905550565b6106b6610ae4565b73ffffffffffffffffffffffffffffffffffffffff166106d4610596565b73ffffffffffffffffffffffffffffffffffffffff161461072a576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107219061122a565b60405180910390fd5b60011515600060149054906101000a900460ff16151514610780576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610777906114fd565b60405180910390fd5b60008060146101000a81548160ff021916908315150217905550565b60001515600060149054906101000a900460ff161515146107f2576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107e990611327565b60405180910390fd5b6001151560018660405161080691906112c4565b908152602001604051809103902060009054906101000a900460ff16151514610864576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161085b90611569565b60405180910390fd5b61087083333084610ce6565b8273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff16866040516108ac91906112c4565b60405180910390207e6b0e4d260e96ab50544d327c9b2747d2c9032870e6c00d5479ac75d06635188533866040516108e693929190611589565b60405180910390a45050505050565b6108fd610ae4565b73ffffffffffffffffffffffffffffffffffffffff1661091b610596565b73ffffffffffffffffffffffffffffffffffffffff1614610971576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109689061122a565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614156109e1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016109d890611632565b60405180910390fd5b6109ea81610aec565b50565b6109f5610ae4565b73ffffffffffffffffffffffffffffffffffffffff16610a13610596565b73ffffffffffffffffffffffffffffffffffffffff1614610a69576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a609061122a565b60405180910390fd5b60018082604051610a7a91906112c4565b908152602001604051809103902060006101000a81548160ff02191690831515021790555080604051610aad91906112c4565b60405180910390207f7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a760405160405180910390a250565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000808473ffffffffffffffffffffffffffffffffffffffff1663a9059cbb8585604051602401610be2929190611652565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610c3091906116c2565b6000604051808303816000865af19150503d8060008114610c6d576040519150601f19603f3d011682016040523d82523d6000602084013e610c72565b606091505b5091509150818015610ca05750600081511480610c9f575080806020019051810190610c9e9190611705565b5b5b610cdf576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cd69061177e565b60405180910390fd5b5050505050565b6000808573ffffffffffffffffffffffffffffffffffffffff166323b872dd868686604051602401610d1a93929190611589565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610d6891906116c2565b6000604051808303816000865af19150503d8060008114610da5576040519150601f19603f3d011682016040523d82523d6000602084013e610daa565b606091505b5091509150818015610dd85750600081511480610dd7575080806020019051810190610dd69190611705565b5b5b610e17576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e0e90611810565b60405180910390fd5b505050505050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610e8682610e3d565b810181811067ffffffffffffffff82111715610ea557610ea4610e4e565b5b80604052505050565b6000610eb8610e1f565b9050610ec48282610e7d565b919050565b600067ffffffffffffffff821115610ee457610ee3610e4e565b5b610eed82610e3d565b9050602081019050919050565b82818337600083830152505050565b6000610f1c610f1784610ec9565b610eae565b905082815260208101848484011115610f3857610f37610e38565b5b610f43848285610efa565b509392505050565b600082601f830112610f6057610f5f610e33565b5b8135610f70848260208601610f09565b91505092915050565b600060208284031215610f8f57610f8e610e29565b5b600082013567ffffffffffffffff811115610fad57610fac610e2e565b5b610fb984828501610f4b565b91505092915050565b60008115159050919050565b610fd781610fc2565b82525050565b6000602082019050610ff26000830184610fce565b92915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061102382610ff8565b9050919050565b61103381611018565b811461103e57600080fd5b50565b6000813590506110508161102a565b92915050565b6000819050919050565b61106981611056565b811461107457600080fd5b50565b60008135905061108681611060565b92915050565b6000806000606084860312156110a5576110a4610e29565b5b60006110b386828701611041565b93505060206110c486828701611041565b92505060406110d586828701611077565b9150509250925092565b6110e881611018565b82525050565b600060208201905061110360008301846110df565b92915050565b600080600080600060a0868803121561112557611124610e29565b5b600086013567ffffffffffffffff81111561114357611142610e2e565b5b61114f88828901610f4b565b955050602061116088828901611041565b945050604061117188828901611041565b935050606061118288828901611041565b925050608061119388828901611077565b9150509295509295909350565b6000602082840312156111b6576111b5610e29565b5b60006111c484828501611041565b91505092915050565b600082825260208201905092915050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006112146020836111cd565b915061121f826111de565b602082019050919050565b6000602082019050818103600083015261124381611207565b9050919050565b600081519050919050565b600081905092915050565b60005b8381101561127e578082015181840152602081019050611263565b8381111561128d576000848401525b50505050565b600061129e8261124a565b6112a88185611255565b93506112b8818560208601611260565b80840191505092915050565b60006112d08284611293565b915081905092915050565b7f4761746577617920697320706175736564000000000000000000000000000000600082015250565b60006113116011836111cd565b915061131c826112db565b602082019050919050565b6000602082019050818103600083015261134081611304565b9050919050565b60008151905061135681611060565b92915050565b60006020828403121561137257611371610e29565b5b600061138084828501611347565b91505092915050565b7f476174657761792062616c616e6365206973206c657373207468616e2072657160008201527f756972656420616d6f756e740000000000000000000000000000000000000000602082015250565b60006113e5602c836111cd565b91506113f082611389565b604082019050919050565b60006020820190508181036000830152611414816113d8565b9050919050565b61142481611056565b82525050565b600060208201905061143f600083018461141b565b92915050565b7f476174657761792069732070617573656420616c726561647900000000000000600082015250565b600061147b6019836111cd565b915061148682611445565b602082019050919050565b600060208201905081810360008301526114aa8161146e565b9050919050565b7f47617465776179206973206e6f742070617573656420616c7265616479000000600082015250565b60006114e7601d836111cd565b91506114f2826114b1565b602082019050919050565b60006020820190508181036000830152611516816114da565b9050919050565b7f64657374436861696e206973206e6f7420737570706f72746564000000000000600082015250565b6000611553601a836111cd565b915061155e8261151d565b602082019050919050565b6000602082019050818103600083015261158281611546565b9050919050565b600060608201905061159e60008301866110df565b6115ab60208301856110df565b6115b8604083018461141b565b949350505050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061161c6026836111cd565b9150611627826115c0565b604082019050919050565b6000602082019050818103600083015261164b8161160f565b9050919050565b600060408201905061166760008301856110df565b611674602083018461141b565b9392505050565b600081519050919050565b600081905092915050565b600061169c8261167b565b6116a68185611686565b93506116b6818560208601611260565b80840191505092915050565b60006116ce8284611691565b915081905092915050565b6116e281610fc2565b81146116ed57600080fd5b50565b6000815190506116ff816116d9565b92915050565b60006020828403121561171b5761171a610e29565b5b6000611729848285016116f0565b91505092915050565b7f5472616e7366657248656c7065723a205452414e534645525f4641494c454400600082015250565b6000611768601f836111cd565b915061177382611732565b602082019050919050565b600060208201905081810360008301526117978161175b565b9050919050565b7f5472616e7366657248656c7065723a205452414e534645525f46524f4d5f464160008201527f494c454400000000000000000000000000000000000000000000000000000000602082015250565b60006117fa6024836111cd565b91506118058261179e565b604082019050919050565b60006020820190508181036000830152611829816117ed565b905091905056fea26469706673582212202bca9acb412b528f97bc4eab2bb72a5840c86de6bcf4f037d8dcd4b06bf1b77664736f6c634300080a0033",
}

// Erc20gwABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc20gwMetaData.ABI instead.
var Erc20gwABI = Erc20gwMetaData.ABI

// Erc20gwBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Erc20gwMetaData.Bin instead.
var Erc20gwBin = Erc20gwMetaData.Bin

// DeployErc20gw deploys a new Ethereum contract, binding an instance of Erc20gw to it.
func DeployErc20gw(auth *bind.TransactOpts, backend bind.ContractBackend, _supportedChains []string) (common.Address, *types.Transaction, *Erc20gw, error) {
	parsed, err := Erc20gwMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Erc20gwBin), backend, _supportedChains)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Erc20gw{Erc20gwCaller: Erc20gwCaller{contract: contract}, Erc20gwTransactor: Erc20gwTransactor{contract: contract}, Erc20gwFilterer: Erc20gwFilterer{contract: contract}}, nil
}

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

// TransferIn is a paid mutator transaction binding the contract method 0x8ab008cb.
//
// Solidity: function TransferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactor) TransferIn(opts *bind.TransactOpts, _token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "TransferIn", _token, _recipient, _amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x8ab008cb.
//
// Solidity: function TransferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwSession) TransferIn(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferIn(&_Erc20gw.TransactOpts, _token, _recipient, _amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0x8ab008cb.
//
// Solidity: function TransferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactorSession) TransferIn(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferIn(&_Erc20gw.TransactOpts, _token, _recipient, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0xf1897eb5.
//
// Solidity: function TransferOut(string _destChain, address _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactor) TransferOut(opts *bind.TransactOpts, _destChain string, _recipient common.Address, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.contract.Transact(opts, "TransferOut", _destChain, _recipient, _tokenOut, _tokenIn, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0xf1897eb5.
//
// Solidity: function TransferOut(string _destChain, address _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gw *Erc20gwSession) TransferOut(_destChain string, _recipient common.Address, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut(&_Erc20gw.TransactOpts, _destChain, _recipient, _tokenOut, _tokenIn, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0xf1897eb5.
//
// Solidity: function TransferOut(string _destChain, address _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gw *Erc20gwTransactorSession) TransferOut(_destChain string, _recipient common.Address, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gw.Contract.TransferOut(&_Erc20gw.TransactOpts, _destChain, _recipient, _tokenOut, _tokenIn, _amount)
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
	Token     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInEvent is a free log retrieval operation binding the contract event 0xb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c.
//
// Solidity: event TransferInEvent(address indexed token, address indexed recipient, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) FilterTransferInEvent(opts *bind.FilterOpts, token []common.Address, recipient []common.Address) (*Erc20gwTransferInEventIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "TransferInEvent", tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwTransferInEventIterator{contract: _Erc20gw.contract, event: "TransferInEvent", logs: logs, sub: sub}, nil
}

// WatchTransferInEvent is a free log subscription operation binding the contract event 0xb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c.
//
// Solidity: event TransferInEvent(address indexed token, address indexed recipient, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) WatchTransferInEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwTransferInEvent, token []common.Address, recipient []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "TransferInEvent", tokenRule, recipientRule)
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

// ParseTransferInEvent is a log parse operation binding the contract event 0xb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c.
//
// Solidity: event TransferInEvent(address indexed token, address indexed recipient, uint256 amount)
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
	DestChain common.Hash
	Recipient common.Address
	TokenOut  common.Address
	TokenIn   common.Address
	Sender    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferOutEvent is a free log retrieval operation binding the contract event 0x006b0e4d260e96ab50544d327c9b2747d2c9032870e6c00d5479ac75d0663518.
//
// Solidity: event TransferOutEvent(string indexed destChain, address indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) FilterTransferOutEvent(opts *bind.FilterOpts, destChain []string, recipient []common.Address, tokenOut []common.Address) (*Erc20gwTransferOutEventIterator, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var tokenOutRule []interface{}
	for _, tokenOutItem := range tokenOut {
		tokenOutRule = append(tokenOutRule, tokenOutItem)
	}

	logs, sub, err := _Erc20gw.contract.FilterLogs(opts, "TransferOutEvent", destChainRule, recipientRule, tokenOutRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gwTransferOutEventIterator{contract: _Erc20gw.contract, event: "TransferOutEvent", logs: logs, sub: sub}, nil
}

// WatchTransferOutEvent is a free log subscription operation binding the contract event 0x006b0e4d260e96ab50544d327c9b2747d2c9032870e6c00d5479ac75d0663518.
//
// Solidity: event TransferOutEvent(string indexed destChain, address indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) WatchTransferOutEvent(opts *bind.WatchOpts, sink chan<- *Erc20gwTransferOutEvent, destChain []string, recipient []common.Address, tokenOut []common.Address) (event.Subscription, error) {

	var destChainRule []interface{}
	for _, destChainItem := range destChain {
		destChainRule = append(destChainRule, destChainItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}
	var tokenOutRule []interface{}
	for _, tokenOutItem := range tokenOut {
		tokenOutRule = append(tokenOutRule, tokenOutItem)
	}

	logs, sub, err := _Erc20gw.contract.WatchLogs(opts, "TransferOutEvent", destChainRule, recipientRule, tokenOutRule)
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

// ParseTransferOutEvent is a log parse operation binding the contract event 0x006b0e4d260e96ab50544d327c9b2747d2c9032870e6c00d5479ac75d0663518.
//
// Solidity: event TransferOutEvent(string indexed destChain, address indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gw *Erc20gwFilterer) ParseTransferOutEvent(log types.Log) (*Erc20gwTransferOutEvent, error) {
	event := new(Erc20gwTransferOutEvent)
	if err := _Erc20gw.contract.UnpackLog(event, "TransferOutEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
