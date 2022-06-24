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

// Erc20gatewayMetaData contains all meta data concerning the Erc20gateway contract.
var Erc20gatewayMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string[]\",\"name\":\"_supportedChains\",\"type\":\"string[]\"},{\"internalType\":\"address\",\"name\":\"_lpPool\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"AddSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"RemoveSupportedChainEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newLpPool\",\"type\":\"address\"}],\"name\":\"SetLiquidPoolAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferInEvent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"string\",\"name\":\"destChain\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"string\",\"name\":\"recipient\",\"type\":\"string\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"tokenOut\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"tokenIn\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TransferOutEvent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"addSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lpPool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pauseGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"chain\",\"type\":\"string\"}],\"name\":\"removeSupportedChain\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"resumeGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_newLpPool\",\"type\":\"address\"}],\"name\":\"setLiquidAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"supportedChains\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_destChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_recipient\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"_tokenOut\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_tokenIn\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b5060405162001f9b38038062001f9b8339818101604052810190620000379190620004fe565b620000576200004b6200013460201b60201c565b6200013c60201b60201c565b60008060146101000a81548160ff02191690831515021790555060005b8251811015620000ea576001600284838151811062000098576200009762000564565b5b6020026020010151604051620000af9190620005e0565b908152602001604051809103902060006101000a81548160ff0219169083151502179055508080620000e19062000632565b91505062000074565b5080600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505062000680565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b620002648262000219565b810181811067ffffffffffffffff821117156200028657620002856200022a565b5b80604052505050565b60006200029b62000200565b9050620002a9828262000259565b919050565b600067ffffffffffffffff821115620002cc57620002cb6200022a565b5b602082029050602081019050919050565b600080fd5b600080fd5b600067ffffffffffffffff8211156200030557620003046200022a565b5b620003108262000219565b9050602081019050919050565b60005b838110156200033d57808201518184015260208101905062000320565b838111156200034d576000848401525b50505050565b60006200036a6200036484620002e7565b6200028f565b905082815260208101848484011115620003895762000388620002e2565b5b620003968482856200031d565b509392505050565b600082601f830112620003b657620003b562000214565b5b8151620003c884826020860162000353565b91505092915050565b6000620003e8620003e284620002ae565b6200028f565b905080838252602082019050602084028301858111156200040e576200040d620002dd565b5b835b818110156200045c57805167ffffffffffffffff81111562000437576200043662000214565b5b8086016200044689826200039e565b8552602085019450505060208101905062000410565b5050509392505050565b600082601f8301126200047e576200047d62000214565b5b815162000490848260208601620003d1565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000620004c68262000499565b9050919050565b620004d881620004b9565b8114620004e457600080fd5b50565b600081519050620004f881620004cd565b92915050565b600080604083850312156200051857620005176200020a565b5b600083015167ffffffffffffffff8111156200053957620005386200020f565b5b620005478582860162000466565b92505060206200055a85828601620004e7565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081519050919050565b600081905092915050565b6000620005b68262000593565b620005c281856200059e565b9350620005d48185602086016200031d565b80840191505092915050565b6000620005ee8284620005a9565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b60006200063f8262000628565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620006755762000674620005f9565b5b600182019050919050565b61190b80620006906000396000f3fe608060405234801561001057600080fd5b50600436106100cf5760003560e01c80638456cb591161008c578063df773a7411610066578063df773a74146101c6578063e4652f49146101d0578063e5eb8a89146101ec578063f2fde38b14610208576100cf565b80638456cb59146101805780638da5cb5b1461019e578063d0a26451146101bc576100cf565b80631dd29c89146100d45780633737bcb4146100f0578063465600231461010e5780635b8948f31461012a5780636c30aaa214610146578063715018a614610176575b600080fd5b6100ee60048036038101906100e9919061112b565b610224565b005b6100f861039f565b60405161010591906111ed565b60405180910390f35b61012860048036038101906101239190611208565b6103c5565b005b610144600480360381019061013f9190611251565b6104bd565b005b610160600480360381019061015b9190611208565b6105c0565b60405161016d9190611299565b60405180910390f35b61017e6105f6565b005b61018861067e565b6040516101959190611299565b60405180910390f35b6101a6610691565b6040516101b391906111ed565b60405180910390f35b6101c46106ba565b005b6101ce6107a8565b005b6101ea60048036038101906101e591906112b4565b610897565b005b61020660048036038101906102019190611208565b610b48565b005b610222600480360381019061021d9190611251565b610c40565b005b60001515600060149054906101000a900460ff1615151461027a576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161027190611364565b60405180910390fd5b6001151560028660405161028e91906113fe565b908152602001604051809103902060009054906101000a900460ff161515146102ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102e390611461565b60405180910390fd5b61031a8333600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1684610d38565b8273ffffffffffffffffffffffffffffffffffffffff168460405161033f91906113fe565b60405180910390208660405161035591906113fe565b60405180910390207fe1e85420cf59cbd542dddac06f75047b4763b5744628dec1e2e312807f6633b985338660405161039093929190611490565b60405180910390a45050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6103cd610e71565b73ffffffffffffffffffffffffffffffffffffffff166103eb610691565b73ffffffffffffffffffffffffffffffffffffffff1614610441576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161043890611513565b60405180910390fd5b600060028260405161045391906113fe565b908152602001604051809103902060006101000a81548160ff0219169083151502179055508060405161048691906113fe565b60405180910390207ff300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca63560405160405180910390a250565b6104c5610e71565b73ffffffffffffffffffffffffffffffffffffffff166104e3610691565b73ffffffffffffffffffffffffffffffffffffffff1614610539576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161053090611513565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508073ffffffffffffffffffffffffffffffffffffffff167f5f45b5214c6238b9219374183843360ca6fa9f8bce8f39c211452057a09ef03b60405160405180910390a250565b6002818051602081018201805184825260208301602085012081835280955050505050506000915054906101000a900460ff1681565b6105fe610e71565b73ffffffffffffffffffffffffffffffffffffffff1661061c610691565b73ffffffffffffffffffffffffffffffffffffffff1614610672576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161066990611513565b60405180910390fd5b61067c6000610e79565b565b600060149054906101000a900460ff1681565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b6106c2610e71565b73ffffffffffffffffffffffffffffffffffffffff166106e0610691565b73ffffffffffffffffffffffffffffffffffffffff1614610736576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161072d90611513565b60405180910390fd5b60011515600060149054906101000a900460ff1615151461078c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107839061157f565b60405180910390fd5b60008060146101000a81548160ff021916908315150217905550565b6107b0610e71565b73ffffffffffffffffffffffffffffffffffffffff166107ce610691565b73ffffffffffffffffffffffffffffffffffffffff1614610824576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161081b90611513565b60405180910390fd5b60001515600060149054906101000a900460ff1615151461087a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610871906115eb565b60405180910390fd5b6001600060146101000a81548160ff021916908315150217905550565b61089f610e71565b73ffffffffffffffffffffffffffffffffffffffff166108bd610691565b73ffffffffffffffffffffffffffffffffffffffff1614610913576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161090a90611513565b60405180910390fd5b60001515600060149054906101000a900460ff16151514610969576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161096090611364565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff166370a08231600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff166040518263ffffffff1660e01b81526004016109c691906111ed565b602060405180830381865afa1580156109e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a079190611620565b905081811015610a4c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a43906116bf565b60405180910390fd5b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1663beabacc88585856040518463ffffffff1660e01b8152600401610aab93929190611490565b600060405180830381600087803b158015610ac557600080fd5b505af1158015610ad9573d6000803e3d6000fd5b505050508273ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff167fb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c84604051610b3a91906116df565b60405180910390a350505050565b610b50610e71565b73ffffffffffffffffffffffffffffffffffffffff16610b6e610691565b73ffffffffffffffffffffffffffffffffffffffff1614610bc4576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610bbb90611513565b60405180910390fd5b6001600282604051610bd691906113fe565b908152602001604051809103902060006101000a81548160ff02191690831515021790555080604051610c0991906113fe565b60405180910390207f7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a760405160405180910390a250565b610c48610e71565b73ffffffffffffffffffffffffffffffffffffffff16610c66610691565b73ffffffffffffffffffffffffffffffffffffffff1614610cbc576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610cb390611513565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610d2c576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d239061176c565b60405180910390fd5b610d3581610e79565b50565b6000808573ffffffffffffffffffffffffffffffffffffffff166323b872dd868686604051602401610d6c93929190611490565b6040516020818303038152906040529060e01b6020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050604051610dba91906117d3565b6000604051808303816000865af19150503d8060008114610df7576040519150601f19603f3d011682016040523d82523d6000602084013e610dfc565b606091505b5091509150818015610e2a5750600081511480610e29575080806020019051810190610e289190611816565b5b5b610e69576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610e60906118b5565b60405180910390fd5b505050505050565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b610fa482610f5b565b810181811067ffffffffffffffff82111715610fc357610fc2610f6c565b5b80604052505050565b6000610fd6610f3d565b9050610fe28282610f9b565b919050565b600067ffffffffffffffff82111561100257611001610f6c565b5b61100b82610f5b565b9050602081019050919050565b82818337600083830152505050565b600061103a61103584610fe7565b610fcc565b90508281526020810184848401111561105657611055610f56565b5b611061848285611018565b509392505050565b600082601f83011261107e5761107d610f51565b5b813561108e848260208601611027565b91505092915050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006110c282611097565b9050919050565b6110d2816110b7565b81146110dd57600080fd5b50565b6000813590506110ef816110c9565b92915050565b6000819050919050565b611108816110f5565b811461111357600080fd5b50565b600081359050611125816110ff565b92915050565b600080600080600060a0868803121561114757611146610f47565b5b600086013567ffffffffffffffff81111561116557611164610f4c565b5b61117188828901611069565b955050602086013567ffffffffffffffff81111561119257611191610f4c565b5b61119e88828901611069565b94505060406111af888289016110e0565b93505060606111c0888289016110e0565b92505060806111d188828901611116565b9150509295509295909350565b6111e7816110b7565b82525050565b600060208201905061120260008301846111de565b92915050565b60006020828403121561121e5761121d610f47565b5b600082013567ffffffffffffffff81111561123c5761123b610f4c565b5b61124884828501611069565b91505092915050565b60006020828403121561126757611266610f47565b5b6000611275848285016110e0565b91505092915050565b60008115159050919050565b6112938161127e565b82525050565b60006020820190506112ae600083018461128a565b92915050565b6000806000606084860312156112cd576112cc610f47565b5b60006112db868287016110e0565b93505060206112ec868287016110e0565b92505060406112fd86828701611116565b9150509250925092565b600082825260208201905092915050565b7f4761746577617920697320706175736564000000000000000000000000000000600082015250565b600061134e601183611307565b915061135982611318565b602082019050919050565b6000602082019050818103600083015261137d81611341565b9050919050565b600081519050919050565b600081905092915050565b60005b838110156113b857808201518184015260208101905061139d565b838111156113c7576000848401525b50505050565b60006113d882611384565b6113e2818561138f565b93506113f281856020860161139a565b80840191505092915050565b600061140a82846113cd565b915081905092915050565b7f64657374436861696e206973206e6f7420737570706f72746564000000000000600082015250565b600061144b601a83611307565b915061145682611415565b602082019050919050565b6000602082019050818103600083015261147a8161143e565b9050919050565b61148a816110f5565b82525050565b60006060820190506114a560008301866111de565b6114b260208301856111de565b6114bf6040830184611481565b949350505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b60006114fd602083611307565b9150611508826114c7565b602082019050919050565b6000602082019050818103600083015261152c816114f0565b9050919050565b7f47617465776179206973206e6f742070617573656420616c7265616479000000600082015250565b6000611569601d83611307565b915061157482611533565b602082019050919050565b600060208201905081810360008301526115988161155c565b9050919050565b7f476174657761792069732070617573656420616c726561647900000000000000600082015250565b60006115d5601983611307565b91506115e08261159f565b602082019050919050565b60006020820190508181036000830152611604816115c8565b9050919050565b60008151905061161a816110ff565b92915050565b60006020828403121561163657611635610f47565b5b60006116448482850161160b565b91505092915050565b7f476174657761792062616c616e6365206973206c657373207468616e2072657160008201527f756972656420616d6f756e740000000000000000000000000000000000000000602082015250565b60006116a9602c83611307565b91506116b48261164d565b604082019050919050565b600060208201905081810360008301526116d88161169c565b9050919050565b60006020820190506116f46000830184611481565b92915050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000611756602683611307565b9150611761826116fa565b604082019050919050565b6000602082019050818103600083015261178581611749565b9050919050565b600081519050919050565b600081905092915050565b60006117ad8261178c565b6117b78185611797565b93506117c781856020860161139a565b80840191505092915050565b60006117df82846117a2565b915081905092915050565b6117f38161127e565b81146117fe57600080fd5b50565b600081519050611810816117ea565b92915050565b60006020828403121561182c5761182b610f47565b5b600061183a84828501611801565b91505092915050565b7f5472616e7366657248656c7065723a205452414e534645525f46524f4d5f464160008201527f494c454400000000000000000000000000000000000000000000000000000000602082015250565b600061189f602483611307565b91506118aa82611843565b604082019050919050565b600060208201905081810360008301526118ce81611892565b905091905056fea2646970667358221220220bd5f66cb163145863ce1e634de61f6a848904cb7a04b50309c125fbe6d4d564736f6c634300080a0033",
}

// Erc20gatewayABI is the input ABI used to generate the binding from.
// Deprecated: Use Erc20gatewayMetaData.ABI instead.
var Erc20gatewayABI = Erc20gatewayMetaData.ABI

// Erc20gatewayBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Erc20gatewayMetaData.Bin instead.
var Erc20gatewayBin = Erc20gatewayMetaData.Bin

// DeployErc20gateway deploys a new Ethereum contract, binding an instance of Erc20gateway to it.
func DeployErc20gateway(auth *bind.TransactOpts, backend bind.ContractBackend, _supportedChains []string, _lpPool common.Address) (common.Address, *types.Transaction, *Erc20gateway, error) {
	parsed, err := Erc20gatewayMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Erc20gatewayBin), backend, _supportedChains, _lpPool)
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

// LpPool is a free data retrieval call binding the contract method 0x3737bcb4.
//
// Solidity: function lpPool() view returns(address)
func (_Erc20gateway *Erc20gatewayCaller) LpPool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "lpPool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpPool is a free data retrieval call binding the contract method 0x3737bcb4.
//
// Solidity: function lpPool() view returns(address)
func (_Erc20gateway *Erc20gatewaySession) LpPool() (common.Address, error) {
	return _Erc20gateway.Contract.LpPool(&_Erc20gateway.CallOpts)
}

// LpPool is a free data retrieval call binding the contract method 0x3737bcb4.
//
// Solidity: function lpPool() view returns(address)
func (_Erc20gateway *Erc20gatewayCallerSession) LpPool() (common.Address, error) {
	return _Erc20gateway.Contract.LpPool(&_Erc20gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gateway *Erc20gatewayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gateway *Erc20gatewaySession) Owner() (common.Address, error) {
	return _Erc20gateway.Contract.Owner(&_Erc20gateway.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Erc20gateway *Erc20gatewayCallerSession) Owner() (common.Address, error) {
	return _Erc20gateway.Contract.Owner(&_Erc20gateway.CallOpts)
}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gateway *Erc20gatewayCaller) Pause(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "pause")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gateway *Erc20gatewaySession) Pause() (bool, error) {
	return _Erc20gateway.Contract.Pause(&_Erc20gateway.CallOpts)
}

// Pause is a free data retrieval call binding the contract method 0x8456cb59.
//
// Solidity: function pause() view returns(bool)
func (_Erc20gateway *Erc20gatewayCallerSession) Pause() (bool, error) {
	return _Erc20gateway.Contract.Pause(&_Erc20gateway.CallOpts)
}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gateway *Erc20gatewayCaller) SupportedChains(opts *bind.CallOpts, arg0 string) (bool, error) {
	var out []interface{}
	err := _Erc20gateway.contract.Call(opts, &out, "supportedChains", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gateway *Erc20gatewaySession) SupportedChains(arg0 string) (bool, error) {
	return _Erc20gateway.Contract.SupportedChains(&_Erc20gateway.CallOpts, arg0)
}

// SupportedChains is a free data retrieval call binding the contract method 0x6c30aaa2.
//
// Solidity: function supportedChains(string ) view returns(bool)
func (_Erc20gateway *Erc20gatewayCallerSession) SupportedChains(arg0 string) (bool, error) {
	return _Erc20gateway.Contract.SupportedChains(&_Erc20gateway.CallOpts, arg0)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xe5eb8a89.
//
// Solidity: function addSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewayTransactor) AddSupportedChain(opts *bind.TransactOpts, chain string) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "addSupportedChain", chain)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xe5eb8a89.
//
// Solidity: function addSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewaySession) AddSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.AddSupportedChain(&_Erc20gateway.TransactOpts, chain)
}

// AddSupportedChain is a paid mutator transaction binding the contract method 0xe5eb8a89.
//
// Solidity: function addSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) AddSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.AddSupportedChain(&_Erc20gateway.TransactOpts, chain)
}

// PauseGateway is a paid mutator transaction binding the contract method 0xdf773a74.
//
// Solidity: function pauseGateway() returns()
func (_Erc20gateway *Erc20gatewayTransactor) PauseGateway(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "pauseGateway")
}

// PauseGateway is a paid mutator transaction binding the contract method 0xdf773a74.
//
// Solidity: function pauseGateway() returns()
func (_Erc20gateway *Erc20gatewaySession) PauseGateway() (*types.Transaction, error) {
	return _Erc20gateway.Contract.PauseGateway(&_Erc20gateway.TransactOpts)
}

// PauseGateway is a paid mutator transaction binding the contract method 0xdf773a74.
//
// Solidity: function pauseGateway() returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) PauseGateway() (*types.Transaction, error) {
	return _Erc20gateway.Contract.PauseGateway(&_Erc20gateway.TransactOpts)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x46560023.
//
// Solidity: function removeSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewayTransactor) RemoveSupportedChain(opts *bind.TransactOpts, chain string) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "removeSupportedChain", chain)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x46560023.
//
// Solidity: function removeSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewaySession) RemoveSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.RemoveSupportedChain(&_Erc20gateway.TransactOpts, chain)
}

// RemoveSupportedChain is a paid mutator transaction binding the contract method 0x46560023.
//
// Solidity: function removeSupportedChain(string chain) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) RemoveSupportedChain(chain string) (*types.Transaction, error) {
	return _Erc20gateway.Contract.RemoveSupportedChain(&_Erc20gateway.TransactOpts, chain)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gateway *Erc20gatewayTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gateway *Erc20gatewaySession) RenounceOwnership() (*types.Transaction, error) {
	return _Erc20gateway.Contract.RenounceOwnership(&_Erc20gateway.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Erc20gateway.Contract.RenounceOwnership(&_Erc20gateway.TransactOpts)
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xd0a26451.
//
// Solidity: function resumeGateway() returns()
func (_Erc20gateway *Erc20gatewayTransactor) ResumeGateway(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "resumeGateway")
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xd0a26451.
//
// Solidity: function resumeGateway() returns()
func (_Erc20gateway *Erc20gatewaySession) ResumeGateway() (*types.Transaction, error) {
	return _Erc20gateway.Contract.ResumeGateway(&_Erc20gateway.TransactOpts)
}

// ResumeGateway is a paid mutator transaction binding the contract method 0xd0a26451.
//
// Solidity: function resumeGateway() returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) ResumeGateway() (*types.Transaction, error) {
	return _Erc20gateway.Contract.ResumeGateway(&_Erc20gateway.TransactOpts)
}

// SetLiquidAddress is a paid mutator transaction binding the contract method 0x5b8948f3.
//
// Solidity: function setLiquidAddress(address _newLpPool) returns()
func (_Erc20gateway *Erc20gatewayTransactor) SetLiquidAddress(opts *bind.TransactOpts, _newLpPool common.Address) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "setLiquidAddress", _newLpPool)
}

// SetLiquidAddress is a paid mutator transaction binding the contract method 0x5b8948f3.
//
// Solidity: function setLiquidAddress(address _newLpPool) returns()
func (_Erc20gateway *Erc20gatewaySession) SetLiquidAddress(_newLpPool common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.SetLiquidAddress(&_Erc20gateway.TransactOpts, _newLpPool)
}

// SetLiquidAddress is a paid mutator transaction binding the contract method 0x5b8948f3.
//
// Solidity: function setLiquidAddress(address _newLpPool) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) SetLiquidAddress(_newLpPool common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.SetLiquidAddress(&_Erc20gateway.TransactOpts, _newLpPool)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferIn(opts *bind.TransactOpts, _token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferIn", _token, _recipient, _amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferIn(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferIn(&_Erc20gateway.TransactOpts, _token, _recipient, _amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address _token, address _recipient, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferIn(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferIn(&_Erc20gateway.TransactOpts, _token, _recipient, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x1dd29c89.
//
// Solidity: function transferOut(string _destChain, string _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferOut(opts *bind.TransactOpts, _destChain string, _recipient string, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferOut", _destChain, _recipient, _tokenOut, _tokenIn, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x1dd29c89.
//
// Solidity: function transferOut(string _destChain, string _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferOut(_destChain string, _recipient string, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOut(&_Erc20gateway.TransactOpts, _destChain, _recipient, _tokenOut, _tokenIn, _amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x1dd29c89.
//
// Solidity: function transferOut(string _destChain, string _recipient, address _tokenOut, address _tokenIn, uint256 _amount) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferOut(_destChain string, _recipient string, _tokenOut common.Address, _tokenIn common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOut(&_Erc20gateway.TransactOpts, _destChain, _recipient, _tokenOut, _tokenIn, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gateway *Erc20gatewayTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gateway *Erc20gatewaySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOwnership(&_Erc20gateway.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Erc20gateway *Erc20gatewayTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Erc20gateway.Contract.TransferOwnership(&_Erc20gateway.TransactOpts, newOwner)
}

// Erc20gatewayAddSupportedChainEventIterator is returned from FilterAddSupportedChainEvent and is used to iterate over the raw logs and unpacked data for AddSupportedChainEvent events raised by the Erc20gateway contract.
type Erc20gatewayAddSupportedChainEventIterator struct {
	Event *Erc20gatewayAddSupportedChainEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayAddSupportedChainEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayAddSupportedChainEvent)
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
		it.Event = new(Erc20gatewayAddSupportedChainEvent)
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
func (it *Erc20gatewayAddSupportedChainEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayAddSupportedChainEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayAddSupportedChainEvent represents a AddSupportedChainEvent event raised by the Erc20gateway contract.
type Erc20gatewayAddSupportedChainEvent struct {
	Chain common.Hash
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterAddSupportedChainEvent is a free log retrieval operation binding the contract event 0x7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a7.
//
// Solidity: event AddSupportedChainEvent(string indexed chain)
func (_Erc20gateway *Erc20gatewayFilterer) FilterAddSupportedChainEvent(opts *bind.FilterOpts, chain []string) (*Erc20gatewayAddSupportedChainEventIterator, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "AddSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayAddSupportedChainEventIterator{contract: _Erc20gateway.contract, event: "AddSupportedChainEvent", logs: logs, sub: sub}, nil
}

// WatchAddSupportedChainEvent is a free log subscription operation binding the contract event 0x7fa5b6d08b213cf08846553aed6553e01273440fcfb334111e8376b02ed434a7.
//
// Solidity: event AddSupportedChainEvent(string indexed chain)
func (_Erc20gateway *Erc20gatewayFilterer) WatchAddSupportedChainEvent(opts *bind.WatchOpts, sink chan<- *Erc20gatewayAddSupportedChainEvent, chain []string) (event.Subscription, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "AddSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayAddSupportedChainEvent)
				if err := _Erc20gateway.contract.UnpackLog(event, "AddSupportedChainEvent", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseAddSupportedChainEvent(log types.Log) (*Erc20gatewayAddSupportedChainEvent, error) {
	event := new(Erc20gatewayAddSupportedChainEvent)
	if err := _Erc20gateway.contract.UnpackLog(event, "AddSupportedChainEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Erc20gateway contract.
type Erc20gatewayOwnershipTransferredIterator struct {
	Event *Erc20gatewayOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayOwnershipTransferred)
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
		it.Event = new(Erc20gatewayOwnershipTransferred)
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
func (it *Erc20gatewayOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayOwnershipTransferred represents a OwnershipTransferred event raised by the Erc20gateway contract.
type Erc20gatewayOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Erc20gateway *Erc20gatewayFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Erc20gatewayOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayOwnershipTransferredIterator{contract: _Erc20gateway.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Erc20gateway *Erc20gatewayFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Erc20gatewayOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayOwnershipTransferred)
				if err := _Erc20gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseOwnershipTransferred(log types.Log) (*Erc20gatewayOwnershipTransferred, error) {
	event := new(Erc20gatewayOwnershipTransferred)
	if err := _Erc20gateway.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayRemoveSupportedChainEventIterator is returned from FilterRemoveSupportedChainEvent and is used to iterate over the raw logs and unpacked data for RemoveSupportedChainEvent events raised by the Erc20gateway contract.
type Erc20gatewayRemoveSupportedChainEventIterator struct {
	Event *Erc20gatewayRemoveSupportedChainEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayRemoveSupportedChainEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayRemoveSupportedChainEvent)
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
		it.Event = new(Erc20gatewayRemoveSupportedChainEvent)
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
func (it *Erc20gatewayRemoveSupportedChainEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayRemoveSupportedChainEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayRemoveSupportedChainEvent represents a RemoveSupportedChainEvent event raised by the Erc20gateway contract.
type Erc20gatewayRemoveSupportedChainEvent struct {
	Chain common.Hash
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterRemoveSupportedChainEvent is a free log retrieval operation binding the contract event 0xf300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca635.
//
// Solidity: event RemoveSupportedChainEvent(string indexed chain)
func (_Erc20gateway *Erc20gatewayFilterer) FilterRemoveSupportedChainEvent(opts *bind.FilterOpts, chain []string) (*Erc20gatewayRemoveSupportedChainEventIterator, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "RemoveSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayRemoveSupportedChainEventIterator{contract: _Erc20gateway.contract, event: "RemoveSupportedChainEvent", logs: logs, sub: sub}, nil
}

// WatchRemoveSupportedChainEvent is a free log subscription operation binding the contract event 0xf300fb61ffb72cae02d1183cefa3fd9604388876c9dae6eab266d6a2a69ca635.
//
// Solidity: event RemoveSupportedChainEvent(string indexed chain)
func (_Erc20gateway *Erc20gatewayFilterer) WatchRemoveSupportedChainEvent(opts *bind.WatchOpts, sink chan<- *Erc20gatewayRemoveSupportedChainEvent, chain []string) (event.Subscription, error) {

	var chainRule []interface{}
	for _, chainItem := range chain {
		chainRule = append(chainRule, chainItem)
	}

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "RemoveSupportedChainEvent", chainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayRemoveSupportedChainEvent)
				if err := _Erc20gateway.contract.UnpackLog(event, "RemoveSupportedChainEvent", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseRemoveSupportedChainEvent(log types.Log) (*Erc20gatewayRemoveSupportedChainEvent, error) {
	event := new(Erc20gatewayRemoveSupportedChainEvent)
	if err := _Erc20gateway.contract.UnpackLog(event, "RemoveSupportedChainEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewaySetLiquidPoolAddressIterator is returned from FilterSetLiquidPoolAddress and is used to iterate over the raw logs and unpacked data for SetLiquidPoolAddress events raised by the Erc20gateway contract.
type Erc20gatewaySetLiquidPoolAddressIterator struct {
	Event *Erc20gatewaySetLiquidPoolAddress // Event containing the contract specifics and raw log

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
func (it *Erc20gatewaySetLiquidPoolAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewaySetLiquidPoolAddress)
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
		it.Event = new(Erc20gatewaySetLiquidPoolAddress)
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
func (it *Erc20gatewaySetLiquidPoolAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewaySetLiquidPoolAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewaySetLiquidPoolAddress represents a SetLiquidPoolAddress event raised by the Erc20gateway contract.
type Erc20gatewaySetLiquidPoolAddress struct {
	NewLpPool common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSetLiquidPoolAddress is a free log retrieval operation binding the contract event 0x5f45b5214c6238b9219374183843360ca6fa9f8bce8f39c211452057a09ef03b.
//
// Solidity: event SetLiquidPoolAddress(address indexed newLpPool)
func (_Erc20gateway *Erc20gatewayFilterer) FilterSetLiquidPoolAddress(opts *bind.FilterOpts, newLpPool []common.Address) (*Erc20gatewaySetLiquidPoolAddressIterator, error) {

	var newLpPoolRule []interface{}
	for _, newLpPoolItem := range newLpPool {
		newLpPoolRule = append(newLpPoolRule, newLpPoolItem)
	}

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "SetLiquidPoolAddress", newLpPoolRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewaySetLiquidPoolAddressIterator{contract: _Erc20gateway.contract, event: "SetLiquidPoolAddress", logs: logs, sub: sub}, nil
}

// WatchSetLiquidPoolAddress is a free log subscription operation binding the contract event 0x5f45b5214c6238b9219374183843360ca6fa9f8bce8f39c211452057a09ef03b.
//
// Solidity: event SetLiquidPoolAddress(address indexed newLpPool)
func (_Erc20gateway *Erc20gatewayFilterer) WatchSetLiquidPoolAddress(opts *bind.WatchOpts, sink chan<- *Erc20gatewaySetLiquidPoolAddress, newLpPool []common.Address) (event.Subscription, error) {

	var newLpPoolRule []interface{}
	for _, newLpPoolItem := range newLpPool {
		newLpPoolRule = append(newLpPoolRule, newLpPoolItem)
	}

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "SetLiquidPoolAddress", newLpPoolRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewaySetLiquidPoolAddress)
				if err := _Erc20gateway.contract.UnpackLog(event, "SetLiquidPoolAddress", log); err != nil {
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

// ParseSetLiquidPoolAddress is a log parse operation binding the contract event 0x5f45b5214c6238b9219374183843360ca6fa9f8bce8f39c211452057a09ef03b.
//
// Solidity: event SetLiquidPoolAddress(address indexed newLpPool)
func (_Erc20gateway *Erc20gatewayFilterer) ParseSetLiquidPoolAddress(log types.Log) (*Erc20gatewaySetLiquidPoolAddress, error) {
	event := new(Erc20gatewaySetLiquidPoolAddress)
	if err := _Erc20gateway.contract.UnpackLog(event, "SetLiquidPoolAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferInEventIterator is returned from FilterTransferInEvent and is used to iterate over the raw logs and unpacked data for TransferInEvent events raised by the Erc20gateway contract.
type Erc20gatewayTransferInEventIterator struct {
	Event *Erc20gatewayTransferInEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferInEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferInEvent)
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
		it.Event = new(Erc20gatewayTransferInEvent)
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
func (it *Erc20gatewayTransferInEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferInEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferInEvent represents a TransferInEvent event raised by the Erc20gateway contract.
type Erc20gatewayTransferInEvent struct {
	Token     common.Address
	Recipient common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferInEvent is a free log retrieval operation binding the contract event 0xb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c.
//
// Solidity: event TransferInEvent(address indexed token, address indexed recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferInEvent(opts *bind.FilterOpts, token []common.Address, recipient []common.Address) (*Erc20gatewayTransferInEventIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferInEvent", tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferInEventIterator{contract: _Erc20gateway.contract, event: "TransferInEvent", logs: logs, sub: sub}, nil
}

// WatchTransferInEvent is a free log subscription operation binding the contract event 0xb7275fa1625b051238c95d6354c70b3ab71046400d703334de68a46923e6274c.
//
// Solidity: event TransferInEvent(address indexed token, address indexed recipient, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferInEvent(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferInEvent, token []common.Address, recipient []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var recipientRule []interface{}
	for _, recipientItem := range recipient {
		recipientRule = append(recipientRule, recipientItem)
	}

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferInEvent", tokenRule, recipientRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferInEvent)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferInEvent", log); err != nil {
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
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferInEvent(log types.Log) (*Erc20gatewayTransferInEvent, error) {
	event := new(Erc20gatewayTransferInEvent)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferInEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Erc20gatewayTransferOutEventIterator is returned from FilterTransferOutEvent and is used to iterate over the raw logs and unpacked data for TransferOutEvent events raised by the Erc20gateway contract.
type Erc20gatewayTransferOutEventIterator struct {
	Event *Erc20gatewayTransferOutEvent // Event containing the contract specifics and raw log

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
func (it *Erc20gatewayTransferOutEventIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Erc20gatewayTransferOutEvent)
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
		it.Event = new(Erc20gatewayTransferOutEvent)
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
func (it *Erc20gatewayTransferOutEventIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Erc20gatewayTransferOutEventIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Erc20gatewayTransferOutEvent represents a TransferOutEvent event raised by the Erc20gateway contract.
type Erc20gatewayTransferOutEvent struct {
	DestChain common.Hash
	Recipient common.Hash
	TokenOut  common.Address
	TokenIn   common.Address
	Sender    common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterTransferOutEvent is a free log retrieval operation binding the contract event 0xe1e85420cf59cbd542dddac06f75047b4763b5744628dec1e2e312807f6633b9.
//
// Solidity: event TransferOutEvent(string indexed destChain, string indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) FilterTransferOutEvent(opts *bind.FilterOpts, destChain []string, recipient []string, tokenOut []common.Address) (*Erc20gatewayTransferOutEventIterator, error) {

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

	logs, sub, err := _Erc20gateway.contract.FilterLogs(opts, "TransferOutEvent", destChainRule, recipientRule, tokenOutRule)
	if err != nil {
		return nil, err
	}
	return &Erc20gatewayTransferOutEventIterator{contract: _Erc20gateway.contract, event: "TransferOutEvent", logs: logs, sub: sub}, nil
}

// WatchTransferOutEvent is a free log subscription operation binding the contract event 0xe1e85420cf59cbd542dddac06f75047b4763b5744628dec1e2e312807f6633b9.
//
// Solidity: event TransferOutEvent(string indexed destChain, string indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) WatchTransferOutEvent(opts *bind.WatchOpts, sink chan<- *Erc20gatewayTransferOutEvent, destChain []string, recipient []string, tokenOut []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Erc20gateway.contract.WatchLogs(opts, "TransferOutEvent", destChainRule, recipientRule, tokenOutRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Erc20gatewayTransferOutEvent)
				if err := _Erc20gateway.contract.UnpackLog(event, "TransferOutEvent", log); err != nil {
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

// ParseTransferOutEvent is a log parse operation binding the contract event 0xe1e85420cf59cbd542dddac06f75047b4763b5744628dec1e2e312807f6633b9.
//
// Solidity: event TransferOutEvent(string indexed destChain, string indexed recipient, address indexed tokenOut, address tokenIn, address sender, uint256 amount)
func (_Erc20gateway *Erc20gatewayFilterer) ParseTransferOutEvent(log types.Log) (*Erc20gatewayTransferOutEvent, error) {
	event := new(Erc20gatewayTransferOutEvent)
	if err := _Erc20gateway.contract.UnpackLog(event, "TransferOutEvent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
