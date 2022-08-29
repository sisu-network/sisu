// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vault

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

// VaultABI is the input ABI used to generate the binding from.
const VaultABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Code501\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Code502\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"addSpender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"changeAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"depositFor\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"depositNativeFor\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"removeSpender\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"retryTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"retryTransferNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferIn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"tos\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"transferInMultiple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferInNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOut\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"dstChains\",\"type\":\"string[]\"},{\"internalType\":\"address[]\",\"name\":\"tos\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"transferOutMultiple\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokens\",\"type\":\"address[]\"},{\"internalType\":\"string[]\",\"name\":\"dstChains\",\"type\":\"string[]\"},{\"internalType\":\"string[]\",\"name\":\"tos\",\"type\":\"string[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"transferOutMultipleNonEvm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"}],\"name\":\"transferOutNative\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"dstChain\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"to\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferOutNonEvm\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawNative\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// VaultBin is the compiled bytecode used for deploying new contracts.
var VaultBin = "0x608060405234801561001057600080fd5b5033600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550612b27806100616000396000f3fe60806040526004361061011f5760003560e01c80638f283970116100a0578063db6b524611610064578063db6b524614610371578063e4652f491461037b578063e7e31e7a146103a4578063f7888aec146103cd578063fdd886b51461040a5761011f565b80638f283970146102b157806391a775b7146102da578063b3db428b146102f6578063c0d863711461031f578063d9caed12146103485761011f565b806347e7ef24116100e757806347e7ef24146101e457806351c7cf3f1461020d578063754ff725146102365780637eab231a1461025f5780638ce5877c146102885761011f565b80630603e1c51461012457806307b18bde146101405780630bae283d146101695780633683f9ab1461019257806346cf2e7c146101bb575b600080fd5b61013e60048036038101906101399190611ba3565b610433565b005b34801561014c57600080fd5b5061016760048036038101906101629190611caf565b610456565b005b34801561017557600080fd5b50610190600480360381019061018b9190611caf565b610483565b005b34801561019e57600080fd5b506101b960048036038101906101b49190611caf565b610687565b005b3480156101c757600080fd5b506101e260048036038101906101dd9190611cef565b610778565b005b3480156101f057600080fd5b5061020b60048036038101906102069190611caf565b61083f565b005b34801561021957600080fd5b50610234600480360381019061022f9190611fde565b61084e565b005b34801561024257600080fd5b5061025d600480360381019061025891906120b5565b610906565b005b34801561026b57600080fd5b506102866004803603810190610281919061215c565b610a27565b005b34801561029457600080fd5b506102af60048036038101906102aa91906121af565b610c7f565b005b3480156102bd57600080fd5b506102d860048036038101906102d391906121af565b610d69565b005b6102f460048036038101906102ef91906121af565b610e3d565b005b34801561030257600080fd5b5061031d6004803603810190610318919061215c565b610e5f565b005b34801561032b57600080fd5b50610346600480360381019061034191906121dc565b610e9c565b005b34801561035457600080fd5b5061036f600480360381019061036a919061215c565b610f54565b005b610379610f65565b005b34801561038757600080fd5b506103a2600480360381019061039d919061215c565b610f86565b005b3480156103b057600080fd5b506103cb60048036038101906103c691906121af565b6110fe565b005b3480156103d957600080fd5b506103f460048036038101906103ef91906122b3565b6111e8565b6040516104019190612302565b60405180910390f35b34801561041657600080fd5b50610431600480360381019061042c919061231d565b61126f565b005b6104527340000000000000000000000000000000000000003334611336565b5050565b610475734000000000000000000000000000000000000000338361143e565b61047f8282611605565b5050565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff1661050e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161050590612419565b60405180910390fd5b60008190506002600073400000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205481111561063f576002600073400000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490505b804710610655576106508382611605565b610682565b7fe1ef72f1796988294cc9a59a2e3073a1079dc556e353a1ebf43b456b191161a960405160405180910390a15b505050565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610712576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161070990612419565b60405180910390fd5b804710610728576107238282611605565b610774565b6107477340000000000000000000000000000000000000008383611336565b7f6a036bee1b01306b61370b348f57c9a7038a7bf8cb5d8a4cfbfb197b3f329e8360405160405180910390a15b5050565b80600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541061080b5761080684338361143e565b610839565b6108383330838773ffffffffffffffffffffffffffffffffffffffff16611705909392919063ffffffff16565b5b50505050565b61084a82338361178e565b5050565b60005b84518163ffffffff1610156108ff576108ec858263ffffffff168151811061087c5761087b612439565b5b6020026020010151858363ffffffff168151811061089d5761089c612439565b5b6020026020010151858463ffffffff16815181106108be576108bd612439565b5b6020026020010151858563ffffffff16815181106108df576108de612439565b5b6020026020010151610778565b80806108f7906124a7565b915050610851565b5050505050565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610991576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161098890612419565b60405180910390fd5b60005b83518163ffffffff161015610a2157610a0e848263ffffffff16815181106109bf576109be612439565b5b6020026020010151848363ffffffff16815181106109e0576109df612439565b5b6020026020010151848463ffffffff1681518110610a0157610a00612439565b5b6020026020010151610f86565b8080610a19906124a7565b915050610994565b50505050565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16610ab2576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610aa990612419565b60405180910390fd5b6000819050600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054811115610bbb57600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205490505b808473ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610bf591906124e3565b602060405180830381865afa158015610c12573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c369190612513565b10610c4c57610c47848485846117cb565b610c79565b7fe1ef72f1796988294cc9a59a2e3073a1079dc556e353a1ebf43b456b191161a960405160405180910390a15b50505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610d0f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d069061258c565b60405180910390fd5b60008060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610df9576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610df09061258c565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b610e5c7340000000000000000000000000000000000000008234611336565b50565b610e8c3330838673ffffffffffffffffffffffffffffffffffffffff16611705909392919063ffffffff16565b610e97838383611336565b505050565b60005b84518163ffffffff161015610f4d57610f3a858263ffffffff1681518110610eca57610ec9612439565b5b6020026020010151858363ffffffff1681518110610eeb57610eea612439565b5b6020026020010151858463ffffffff1681518110610f0c57610f0b612439565b5b6020026020010151858563ffffffff1681518110610f2d57610f2c612439565b5b602002602001015161126f565b8080610f45906124a7565b915050610e9f565b5050505050565b610f60833384846117cb565b505050565b610f847340000000000000000000000000000000000000003334611336565b565b6000803373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900460ff16611011576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161100890612419565b60405180910390fd5b808373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161104b91906124e3565b602060405180830381865afa158015611068573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061108c9190612513565b106110c1576110bc82828573ffffffffffffffffffffffffffffffffffffffff166118079092919063ffffffff16565b6110f9565b6110cc838383611336565b7f6a036bee1b01306b61370b348f57c9a7038a7bf8cb5d8a4cfbfb197b3f329e8360405160405180910390a15b505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff161461118e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111859061258c565b60405180910390fd5b60016000808373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548160ff02191690831515021790555050565b6000600260008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054905092915050565b80600260008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205410611302576112fd84338361143e565b611330565b61132f3330838773ffffffffffffffffffffffffffffffffffffffff16611705909392919063ffffffff16565b5b50505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156113a6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161139d906125f8565b60405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546114329190612618565b92505081905550505050565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff1614156114ae576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016114a5906126ba565b60405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054101561156d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161156490612726565b60405180910390fd5b80600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060008282546115f99190612746565b92505081905550505050565b60008273ffffffffffffffffffffffffffffffffffffffff1682600067ffffffffffffffff81111561163a57611639611a78565b5b6040519080825280601f01601f19166020018201604052801561166c5781602001600182028036833780820191505090505b5060405161167a91906127f4565b60006040518083038185875af1925050503d80600081146116b7576040519150601f19603f3d011682016040523d82523d6000602084013e6116bc565b606091505b5050905080611700576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016116f79061287d565b60405180910390fd5b505050565b611788846323b872dd60e01b8585856040516024016117269392919061289d565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061188d565b50505050565b6117bb8230838673ffffffffffffffffffffffffffffffffffffffff16611705909392919063ffffffff16565b6117c6838383611336565b505050565b6117d684848361143e565b61180182828673ffffffffffffffffffffffffffffffffffffffff166118079092919063ffffffff16565b50505050565b6118888363a9059cbb60e01b84846040516024016118269291906128d4565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff838183161783525050505061188d565b505050565b6118ac8273ffffffffffffffffffffffffffffffffffffffff166119fe565b6118eb576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016118e290612949565b60405180910390fd5b6000808373ffffffffffffffffffffffffffffffffffffffff168360405161191391906127f4565b6000604051808303816000865af19150503d8060008114611950576040519150601f19603f3d011682016040523d82523d6000602084013e611955565b606091505b50915091508161199a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611991906129b5565b60405180910390fd5b6000815111156119f857808060200190518101906119b89190612a0d565b6119f7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016119ee90612aac565b60405180910390fd5b5b50505050565b60008060007fc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a47060001b9050833f91506000801b8214158015611a405750808214155b92505050919050565b6000604051905090565b600080fd5b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b611ab082611a67565b810181811067ffffffffffffffff82111715611acf57611ace611a78565b5b80604052505050565b6000611ae2611a49565b9050611aee8282611aa7565b919050565b600067ffffffffffffffff821115611b0e57611b0d611a78565b5b611b1782611a67565b9050602081019050919050565b82818337600083830152505050565b6000611b46611b4184611af3565b611ad8565b905082815260208101848484011115611b6257611b61611a62565b5b611b6d848285611b24565b509392505050565b600082601f830112611b8a57611b89611a5d565b5b8135611b9a848260208601611b33565b91505092915050565b60008060408385031215611bba57611bb9611a53565b5b600083013567ffffffffffffffff811115611bd857611bd7611a58565b5b611be485828601611b75565b925050602083013567ffffffffffffffff811115611c0557611c04611a58565b5b611c1185828601611b75565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000611c4682611c1b565b9050919050565b611c5681611c3b565b8114611c6157600080fd5b50565b600081359050611c7381611c4d565b92915050565b6000819050919050565b611c8c81611c79565b8114611c9757600080fd5b50565b600081359050611ca981611c83565b92915050565b60008060408385031215611cc657611cc5611a53565b5b6000611cd485828601611c64565b9250506020611ce585828601611c9a565b9150509250929050565b60008060008060808587031215611d0957611d08611a53565b5b6000611d1787828801611c64565b945050602085013567ffffffffffffffff811115611d3857611d37611a58565b5b611d4487828801611b75565b9350506040611d5587828801611c64565b9250506060611d6687828801611c9a565b91505092959194509250565b600067ffffffffffffffff821115611d8d57611d8c611a78565b5b602082029050602081019050919050565b600080fd5b6000611db6611db184611d72565b611ad8565b90508083825260208201905060208402830185811115611dd957611dd8611d9e565b5b835b81811015611e025780611dee8882611c64565b845260208401935050602081019050611ddb565b5050509392505050565b600082601f830112611e2157611e20611a5d565b5b8135611e31848260208601611da3565b91505092915050565b600067ffffffffffffffff821115611e5557611e54611a78565b5b602082029050602081019050919050565b6000611e79611e7484611e3a565b611ad8565b90508083825260208201905060208402830185811115611e9c57611e9b611d9e565b5b835b81811015611ee357803567ffffffffffffffff811115611ec157611ec0611a5d565b5b808601611ece8982611b75565b85526020850194505050602081019050611e9e565b5050509392505050565b600082601f830112611f0257611f01611a5d565b5b8135611f12848260208601611e66565b91505092915050565b600067ffffffffffffffff821115611f3657611f35611a78565b5b602082029050602081019050919050565b6000611f5a611f5584611f1b565b611ad8565b90508083825260208201905060208402830185811115611f7d57611f7c611d9e565b5b835b81811015611fa65780611f928882611c9a565b845260208401935050602081019050611f7f565b5050509392505050565b600082601f830112611fc557611fc4611a5d565b5b8135611fd5848260208601611f47565b91505092915050565b60008060008060808587031215611ff857611ff7611a53565b5b600085013567ffffffffffffffff81111561201657612015611a58565b5b61202287828801611e0c565b945050602085013567ffffffffffffffff81111561204357612042611a58565b5b61204f87828801611eed565b935050604085013567ffffffffffffffff8111156120705761206f611a58565b5b61207c87828801611e0c565b925050606085013567ffffffffffffffff81111561209d5761209c611a58565b5b6120a987828801611fb0565b91505092959194509250565b6000806000606084860312156120ce576120cd611a53565b5b600084013567ffffffffffffffff8111156120ec576120eb611a58565b5b6120f886828701611e0c565b935050602084013567ffffffffffffffff81111561211957612118611a58565b5b61212586828701611e0c565b925050604084013567ffffffffffffffff81111561214657612145611a58565b5b61215286828701611fb0565b9150509250925092565b60008060006060848603121561217557612174611a53565b5b600061218386828701611c64565b935050602061219486828701611c64565b92505060406121a586828701611c9a565b9150509250925092565b6000602082840312156121c5576121c4611a53565b5b60006121d384828501611c64565b91505092915050565b600080600080608085870312156121f6576121f5611a53565b5b600085013567ffffffffffffffff81111561221457612213611a58565b5b61222087828801611e0c565b945050602085013567ffffffffffffffff81111561224157612240611a58565b5b61224d87828801611eed565b935050604085013567ffffffffffffffff81111561226e5761226d611a58565b5b61227a87828801611eed565b925050606085013567ffffffffffffffff81111561229b5761229a611a58565b5b6122a787828801611fb0565b91505092959194509250565b600080604083850312156122ca576122c9611a53565b5b60006122d885828601611c64565b92505060206122e985828601611c64565b9150509250929050565b6122fc81611c79565b82525050565b600060208201905061231760008301846122f3565b92915050565b6000806000806080858703121561233757612336611a53565b5b600061234587828801611c64565b945050602085013567ffffffffffffffff81111561236657612365611a58565b5b61237287828801611b75565b935050604085013567ffffffffffffffff81111561239357612392611a58565b5b61239f87828801611b75565b92505060606123b087828801611c9a565b91505092959194509250565b600082825260208201905092915050565b7f4e6f74207370656e6465723a20464f5242494444454e00000000000000000000600082015250565b60006124036016836123bc565b915061240e826123cd565b602082019050919050565b60006020820190508181036000830152612432816123f6565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b600063ffffffff82169050919050565b60006124b282612497565b915063ffffffff8214156124c9576124c8612468565b5b600182019050919050565b6124dd81611c3b565b82525050565b60006020820190506124f860008301846124d4565b92915050565b60008151905061250d81611c83565b92915050565b60006020828403121561252957612528611a53565b5b6000612537848285016124fe565b91505092915050565b7f4e6f742061646d696e3a20464f5242494444454e000000000000000000000000600082015250565b60006125766014836123bc565b915061258182612540565b602082019050919050565b600060208201905081810360008301526125a581612569565b9050919050565b7f696e633a20616464726573732069732030000000000000000000000000000000600082015250565b60006125e26011836123bc565b91506125ed826125ac565b602082019050919050565b60006020820190508181036000830152612611816125d5565b9050919050565b600061262382611c79565b915061262e83611c79565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0382111561266357612662612468565b5b828201905092915050565b7f6465633a20616464726573732069732030000000000000000000000000000000600082015250565b60006126a46011836123bc565b91506126af8261266e565b602082019050919050565b600060208201905081810360008301526126d381612697565b9050919050565b7f6465633a20616d6f756e7420657863656564732062616c616e63650000000000600082015250565b6000612710601b836123bc565b915061271b826126da565b602082019050919050565b6000602082019050818103600083015261273f81612703565b9050919050565b600061275182611c79565b915061275c83611c79565b92508282101561276f5761276e612468565b5b828203905092915050565b600081519050919050565b600081905092915050565b60005b838110156127ae578082015181840152602081019050612793565b838111156127bd576000848401525b50505050565b60006127ce8261277a565b6127d88185612785565b93506127e8818560208601612790565b80840191505092915050565b600061280082846127c3565b915081905092915050565b7f5472616e7366657248656c7065723a204e41544956455f5452414e534645525f60008201527f4641494c45440000000000000000000000000000000000000000000000000000602082015250565b60006128676026836123bc565b91506128728261280b565b604082019050919050565b600060208201905081810360008301526128968161285a565b9050919050565b60006060820190506128b260008301866124d4565b6128bf60208301856124d4565b6128cc60408301846122f3565b949350505050565b60006040820190506128e960008301856124d4565b6128f660208301846122f3565b9392505050565b7f5361666545524332303a2063616c6c20746f206e6f6e2d636f6e747261637400600082015250565b6000612933601f836123bc565b915061293e826128fd565b602082019050919050565b6000602082019050818103600083015261296281612926565b9050919050565b7f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564600082015250565b600061299f6020836123bc565b91506129aa82612969565b602082019050919050565b600060208201905081810360008301526129ce81612992565b9050919050565b60008115159050919050565b6129ea816129d5565b81146129f557600080fd5b50565b600081519050612a07816129e1565b92915050565b600060208284031215612a2357612a22611a53565b5b6000612a31848285016129f8565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000612a96602a836123bc565b9150612aa182612a3a565b604082019050919050565b60006020820190508181036000830152612ac581612a89565b905091905056fea2646970667358221220d65a06bd8b8384ae52ad57f7f5610a3775229990559ca695aaff009e49aa95a764736f6c637827302e382e31322d646576656c6f702e323032322e322e382b636f6d6d69742e35633362636236630058"

// DeployVault deploys a new Ethereum contract, binding an instance of Vault to it.
func DeployVault(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Vault, error) {
	parsed, err := abi.JSON(strings.NewReader(VaultABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(VaultBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Vault{VaultCaller: VaultCaller{contract: contract}, VaultTransactor: VaultTransactor{contract: contract}, VaultFilterer: VaultFilterer{contract: contract}}, nil
}

// Vault is an auto generated Go binding around an Ethereum contract.
type Vault struct {
	VaultCaller     // Read-only binding to the contract
	VaultTransactor // Write-only binding to the contract
	VaultFilterer   // Log filterer for contract events
}

// VaultCaller is an auto generated read-only Go binding around an Ethereum contract.
type VaultCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VaultTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VaultFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VaultSession struct {
	Contract     *Vault            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VaultCallerSession struct {
	Contract *VaultCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// VaultTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VaultTransactorSession struct {
	Contract     *VaultTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultRaw is an auto generated low-level Go binding around an Ethereum contract.
type VaultRaw struct {
	Contract *Vault // Generic contract binding to access the raw methods on
}

// VaultCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VaultCallerRaw struct {
	Contract *VaultCaller // Generic read-only contract binding to access the raw methods on
}

// VaultTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VaultTransactorRaw struct {
	Contract *VaultTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVault creates a new instance of Vault, bound to a specific deployed contract.
func NewVault(address common.Address, backend bind.ContractBackend) (*Vault, error) {
	contract, err := bindVault(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vault{VaultCaller: VaultCaller{contract: contract}, VaultTransactor: VaultTransactor{contract: contract}, VaultFilterer: VaultFilterer{contract: contract}}, nil
}

// NewVaultCaller creates a new read-only instance of Vault, bound to a specific deployed contract.
func NewVaultCaller(address common.Address, caller bind.ContractCaller) (*VaultCaller, error) {
	contract, err := bindVault(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VaultCaller{contract: contract}, nil
}

// NewVaultTransactor creates a new write-only instance of Vault, bound to a specific deployed contract.
func NewVaultTransactor(address common.Address, transactor bind.ContractTransactor) (*VaultTransactor, error) {
	contract, err := bindVault(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VaultTransactor{contract: contract}, nil
}

// NewVaultFilterer creates a new log filterer instance of Vault, bound to a specific deployed contract.
func NewVaultFilterer(address common.Address, filterer bind.ContractFilterer) (*VaultFilterer, error) {
	contract, err := bindVault(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VaultFilterer{contract: contract}, nil
}

// bindVault binds a generic wrapper to an already deployed contract.
func bindVault(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(VaultABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.VaultCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.VaultTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vault *VaultCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vault.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vault *VaultTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vault *VaultTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vault.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address token, address account) view returns(uint256)
func (_Vault *VaultCaller) BalanceOf(opts *bind.CallOpts, token common.Address, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vault.contract.Call(opts, &out, "balanceOf", token, account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address token, address account) view returns(uint256)
func (_Vault *VaultSession) BalanceOf(token common.Address, account common.Address) (*big.Int, error) {
	return _Vault.Contract.BalanceOf(&_Vault.CallOpts, token, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0xf7888aec.
//
// Solidity: function balanceOf(address token, address account) view returns(uint256)
func (_Vault *VaultCallerSession) BalanceOf(token common.Address, account common.Address) (*big.Int, error) {
	return _Vault.Contract.BalanceOf(&_Vault.CallOpts, token, account)
}

// AddSpender is a paid mutator transaction binding the contract method 0xe7e31e7a.
//
// Solidity: function addSpender(address spender) returns()
func (_Vault *VaultTransactor) AddSpender(opts *bind.TransactOpts, spender common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "addSpender", spender)
}

// AddSpender is a paid mutator transaction binding the contract method 0xe7e31e7a.
//
// Solidity: function addSpender(address spender) returns()
func (_Vault *VaultSession) AddSpender(spender common.Address) (*types.Transaction, error) {
	return _Vault.Contract.AddSpender(&_Vault.TransactOpts, spender)
}

// AddSpender is a paid mutator transaction binding the contract method 0xe7e31e7a.
//
// Solidity: function addSpender(address spender) returns()
func (_Vault *VaultTransactorSession) AddSpender(spender common.Address) (*types.Transaction, error) {
	return _Vault.Contract.AddSpender(&_Vault.TransactOpts, spender)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_Vault *VaultTransactor) ChangeAdmin(opts *bind.TransactOpts, newAdmin common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "changeAdmin", newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_Vault *VaultSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Vault.Contract.ChangeAdmin(&_Vault.TransactOpts, newAdmin)
}

// ChangeAdmin is a paid mutator transaction binding the contract method 0x8f283970.
//
// Solidity: function changeAdmin(address newAdmin) returns()
func (_Vault *VaultTransactorSession) ChangeAdmin(newAdmin common.Address) (*types.Transaction, error) {
	return _Vault.Contract.ChangeAdmin(&_Vault.TransactOpts, newAdmin)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_Vault *VaultTransactor) Deposit(opts *bind.TransactOpts, token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "deposit", token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_Vault *VaultSession) Deposit(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Deposit(&_Vault.TransactOpts, token, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x47e7ef24.
//
// Solidity: function deposit(address token, uint256 amount) returns()
func (_Vault *VaultTransactorSession) Deposit(token common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Deposit(&_Vault.TransactOpts, token, amount)
}

// DepositFor is a paid mutator transaction binding the contract method 0xb3db428b.
//
// Solidity: function depositFor(address token, address receiver, uint256 amount) returns()
func (_Vault *VaultTransactor) DepositFor(opts *bind.TransactOpts, token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "depositFor", token, receiver, amount)
}

// DepositFor is a paid mutator transaction binding the contract method 0xb3db428b.
//
// Solidity: function depositFor(address token, address receiver, uint256 amount) returns()
func (_Vault *VaultSession) DepositFor(token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.DepositFor(&_Vault.TransactOpts, token, receiver, amount)
}

// DepositFor is a paid mutator transaction binding the contract method 0xb3db428b.
//
// Solidity: function depositFor(address token, address receiver, uint256 amount) returns()
func (_Vault *VaultTransactorSession) DepositFor(token common.Address, receiver common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.DepositFor(&_Vault.TransactOpts, token, receiver, amount)
}

// DepositNative is a paid mutator transaction binding the contract method 0xdb6b5246.
//
// Solidity: function depositNative() payable returns()
func (_Vault *VaultTransactor) DepositNative(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "depositNative")
}

// DepositNative is a paid mutator transaction binding the contract method 0xdb6b5246.
//
// Solidity: function depositNative() payable returns()
func (_Vault *VaultSession) DepositNative() (*types.Transaction, error) {
	return _Vault.Contract.DepositNative(&_Vault.TransactOpts)
}

// DepositNative is a paid mutator transaction binding the contract method 0xdb6b5246.
//
// Solidity: function depositNative() payable returns()
func (_Vault *VaultTransactorSession) DepositNative() (*types.Transaction, error) {
	return _Vault.Contract.DepositNative(&_Vault.TransactOpts)
}

// DepositNativeFor is a paid mutator transaction binding the contract method 0x91a775b7.
//
// Solidity: function depositNativeFor(address receiver) payable returns()
func (_Vault *VaultTransactor) DepositNativeFor(opts *bind.TransactOpts, receiver common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "depositNativeFor", receiver)
}

// DepositNativeFor is a paid mutator transaction binding the contract method 0x91a775b7.
//
// Solidity: function depositNativeFor(address receiver) payable returns()
func (_Vault *VaultSession) DepositNativeFor(receiver common.Address) (*types.Transaction, error) {
	return _Vault.Contract.DepositNativeFor(&_Vault.TransactOpts, receiver)
}

// DepositNativeFor is a paid mutator transaction binding the contract method 0x91a775b7.
//
// Solidity: function depositNativeFor(address receiver) payable returns()
func (_Vault *VaultTransactorSession) DepositNativeFor(receiver common.Address) (*types.Transaction, error) {
	return _Vault.Contract.DepositNativeFor(&_Vault.TransactOpts, receiver)
}

// RemoveSpender is a paid mutator transaction binding the contract method 0x8ce5877c.
//
// Solidity: function removeSpender(address spender) returns()
func (_Vault *VaultTransactor) RemoveSpender(opts *bind.TransactOpts, spender common.Address) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "removeSpender", spender)
}

// RemoveSpender is a paid mutator transaction binding the contract method 0x8ce5877c.
//
// Solidity: function removeSpender(address spender) returns()
func (_Vault *VaultSession) RemoveSpender(spender common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RemoveSpender(&_Vault.TransactOpts, spender)
}

// RemoveSpender is a paid mutator transaction binding the contract method 0x8ce5877c.
//
// Solidity: function removeSpender(address spender) returns()
func (_Vault *VaultTransactorSession) RemoveSpender(spender common.Address) (*types.Transaction, error) {
	return _Vault.Contract.RemoveSpender(&_Vault.TransactOpts, spender)
}

// RetryTransfer is a paid mutator transaction binding the contract method 0x7eab231a.
//
// Solidity: function retryTransfer(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactor) RetryTransfer(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "retryTransfer", token, to, amount)
}

// RetryTransfer is a paid mutator transaction binding the contract method 0x7eab231a.
//
// Solidity: function retryTransfer(address token, address to, uint256 amount) returns()
func (_Vault *VaultSession) RetryTransfer(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.RetryTransfer(&_Vault.TransactOpts, token, to, amount)
}

// RetryTransfer is a paid mutator transaction binding the contract method 0x7eab231a.
//
// Solidity: function retryTransfer(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) RetryTransfer(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.RetryTransfer(&_Vault.TransactOpts, token, to, amount)
}

// RetryTransferNative is a paid mutator transaction binding the contract method 0x0bae283d.
//
// Solidity: function retryTransferNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactor) RetryTransferNative(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "retryTransferNative", to, amount)
}

// RetryTransferNative is a paid mutator transaction binding the contract method 0x0bae283d.
//
// Solidity: function retryTransferNative(address to, uint256 amount) returns()
func (_Vault *VaultSession) RetryTransferNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.RetryTransferNative(&_Vault.TransactOpts, to, amount)
}

// RetryTransferNative is a paid mutator transaction binding the contract method 0x0bae283d.
//
// Solidity: function retryTransferNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) RetryTransferNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.RetryTransferNative(&_Vault.TransactOpts, to, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactor) TransferIn(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferIn", token, to, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address token, address to, uint256 amount) returns()
func (_Vault *VaultSession) TransferIn(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferIn(&_Vault.TransactOpts, token, to, amount)
}

// TransferIn is a paid mutator transaction binding the contract method 0xe4652f49.
//
// Solidity: function transferIn(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) TransferIn(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferIn(&_Vault.TransactOpts, token, to, amount)
}

// TransferInMultiple is a paid mutator transaction binding the contract method 0x754ff725.
//
// Solidity: function transferInMultiple(address[] tokens, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactor) TransferInMultiple(opts *bind.TransactOpts, tokens []common.Address, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferInMultiple", tokens, tos, amounts)
}

// TransferInMultiple is a paid mutator transaction binding the contract method 0x754ff725.
//
// Solidity: function transferInMultiple(address[] tokens, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultSession) TransferInMultiple(tokens []common.Address, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferInMultiple(&_Vault.TransactOpts, tokens, tos, amounts)
}

// TransferInMultiple is a paid mutator transaction binding the contract method 0x754ff725.
//
// Solidity: function transferInMultiple(address[] tokens, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactorSession) TransferInMultiple(tokens []common.Address, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferInMultiple(&_Vault.TransactOpts, tokens, tos, amounts)
}

// TransferInNative is a paid mutator transaction binding the contract method 0x3683f9ab.
//
// Solidity: function transferInNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactor) TransferInNative(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferInNative", to, amount)
}

// TransferInNative is a paid mutator transaction binding the contract method 0x3683f9ab.
//
// Solidity: function transferInNative(address to, uint256 amount) returns()
func (_Vault *VaultSession) TransferInNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferInNative(&_Vault.TransactOpts, to, amount)
}

// TransferInNative is a paid mutator transaction binding the contract method 0x3683f9ab.
//
// Solidity: function transferInNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) TransferInNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferInNative(&_Vault.TransactOpts, to, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x46cf2e7c.
//
// Solidity: function transferOut(address token, string dstChain, address to, uint256 amount) returns()
func (_Vault *VaultTransactor) TransferOut(opts *bind.TransactOpts, token common.Address, dstChain string, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOut", token, dstChain, to, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x46cf2e7c.
//
// Solidity: function transferOut(address token, string dstChain, address to, uint256 amount) returns()
func (_Vault *VaultSession) TransferOut(token common.Address, dstChain string, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOut(&_Vault.TransactOpts, token, dstChain, to, amount)
}

// TransferOut is a paid mutator transaction binding the contract method 0x46cf2e7c.
//
// Solidity: function transferOut(address token, string dstChain, address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) TransferOut(token common.Address, dstChain string, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOut(&_Vault.TransactOpts, token, dstChain, to, amount)
}

// TransferOutMultiple is a paid mutator transaction binding the contract method 0x51c7cf3f.
//
// Solidity: function transferOutMultiple(address[] tokens, string[] dstChains, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactor) TransferOutMultiple(opts *bind.TransactOpts, tokens []common.Address, dstChains []string, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOutMultiple", tokens, dstChains, tos, amounts)
}

// TransferOutMultiple is a paid mutator transaction binding the contract method 0x51c7cf3f.
//
// Solidity: function transferOutMultiple(address[] tokens, string[] dstChains, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultSession) TransferOutMultiple(tokens []common.Address, dstChains []string, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutMultiple(&_Vault.TransactOpts, tokens, dstChains, tos, amounts)
}

// TransferOutMultiple is a paid mutator transaction binding the contract method 0x51c7cf3f.
//
// Solidity: function transferOutMultiple(address[] tokens, string[] dstChains, address[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactorSession) TransferOutMultiple(tokens []common.Address, dstChains []string, tos []common.Address, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutMultiple(&_Vault.TransactOpts, tokens, dstChains, tos, amounts)
}

// TransferOutMultipleNonEvm is a paid mutator transaction binding the contract method 0xc0d86371.
//
// Solidity: function transferOutMultipleNonEvm(address[] tokens, string[] dstChains, string[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactor) TransferOutMultipleNonEvm(opts *bind.TransactOpts, tokens []common.Address, dstChains []string, tos []string, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOutMultipleNonEvm", tokens, dstChains, tos, amounts)
}

// TransferOutMultipleNonEvm is a paid mutator transaction binding the contract method 0xc0d86371.
//
// Solidity: function transferOutMultipleNonEvm(address[] tokens, string[] dstChains, string[] tos, uint256[] amounts) returns()
func (_Vault *VaultSession) TransferOutMultipleNonEvm(tokens []common.Address, dstChains []string, tos []string, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutMultipleNonEvm(&_Vault.TransactOpts, tokens, dstChains, tos, amounts)
}

// TransferOutMultipleNonEvm is a paid mutator transaction binding the contract method 0xc0d86371.
//
// Solidity: function transferOutMultipleNonEvm(address[] tokens, string[] dstChains, string[] tos, uint256[] amounts) returns()
func (_Vault *VaultTransactorSession) TransferOutMultipleNonEvm(tokens []common.Address, dstChains []string, tos []string, amounts []*big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutMultipleNonEvm(&_Vault.TransactOpts, tokens, dstChains, tos, amounts)
}

// TransferOutNative is a paid mutator transaction binding the contract method 0x0603e1c5.
//
// Solidity: function transferOutNative(string to, string dstChain) payable returns()
func (_Vault *VaultTransactor) TransferOutNative(opts *bind.TransactOpts, to string, dstChain string) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOutNative", to, dstChain)
}

// TransferOutNative is a paid mutator transaction binding the contract method 0x0603e1c5.
//
// Solidity: function transferOutNative(string to, string dstChain) payable returns()
func (_Vault *VaultSession) TransferOutNative(to string, dstChain string) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutNative(&_Vault.TransactOpts, to, dstChain)
}

// TransferOutNative is a paid mutator transaction binding the contract method 0x0603e1c5.
//
// Solidity: function transferOutNative(string to, string dstChain) payable returns()
func (_Vault *VaultTransactorSession) TransferOutNative(to string, dstChain string) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutNative(&_Vault.TransactOpts, to, dstChain)
}

// TransferOutNonEvm is a paid mutator transaction binding the contract method 0xfdd886b5.
//
// Solidity: function transferOutNonEvm(address token, string dstChain, string to, uint256 amount) returns()
func (_Vault *VaultTransactor) TransferOutNonEvm(opts *bind.TransactOpts, token common.Address, dstChain string, to string, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "transferOutNonEvm", token, dstChain, to, amount)
}

// TransferOutNonEvm is a paid mutator transaction binding the contract method 0xfdd886b5.
//
// Solidity: function transferOutNonEvm(address token, string dstChain, string to, uint256 amount) returns()
func (_Vault *VaultSession) TransferOutNonEvm(token common.Address, dstChain string, to string, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutNonEvm(&_Vault.TransactOpts, token, dstChain, to, amount)
}

// TransferOutNonEvm is a paid mutator transaction binding the contract method 0xfdd886b5.
//
// Solidity: function transferOutNonEvm(address token, string dstChain, string to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) TransferOutNonEvm(token common.Address, dstChain string, to string, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.TransferOutNonEvm(&_Vault.TransactOpts, token, dstChain, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactor) Withdraw(opts *bind.TransactOpts, token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "withdraw", token, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address to, uint256 amount) returns()
func (_Vault *VaultSession) Withdraw(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Withdraw(&_Vault.TransactOpts, token, to, amount)
}

// Withdraw is a paid mutator transaction binding the contract method 0xd9caed12.
//
// Solidity: function withdraw(address token, address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) Withdraw(token common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.Withdraw(&_Vault.TransactOpts, token, to, amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x07b18bde.
//
// Solidity: function withdrawNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactor) WithdrawNative(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.contract.Transact(opts, "withdrawNative", to, amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x07b18bde.
//
// Solidity: function withdrawNative(address to, uint256 amount) returns()
func (_Vault *VaultSession) WithdrawNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.WithdrawNative(&_Vault.TransactOpts, to, amount)
}

// WithdrawNative is a paid mutator transaction binding the contract method 0x07b18bde.
//
// Solidity: function withdrawNative(address to, uint256 amount) returns()
func (_Vault *VaultTransactorSession) WithdrawNative(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Vault.Contract.WithdrawNative(&_Vault.TransactOpts, to, amount)
}

// VaultCode501Iterator is returned from FilterCode501 and is used to iterate over the raw logs and unpacked data for Code501 events raised by the Vault contract.
type VaultCode501Iterator struct {
	Event *VaultCode501 // Event containing the contract specifics and raw log

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
func (it *VaultCode501Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultCode501)
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
		it.Event = new(VaultCode501)
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
func (it *VaultCode501Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultCode501Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultCode501 represents a Code501 event raised by the Vault contract.
type VaultCode501 struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCode501 is a free log retrieval operation binding the contract event 0x6a036bee1b01306b61370b348f57c9a7038a7bf8cb5d8a4cfbfb197b3f329e83.
//
// Solidity: event Code501()
func (_Vault *VaultFilterer) FilterCode501(opts *bind.FilterOpts) (*VaultCode501Iterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Code501")
	if err != nil {
		return nil, err
	}
	return &VaultCode501Iterator{contract: _Vault.contract, event: "Code501", logs: logs, sub: sub}, nil
}

// WatchCode501 is a free log subscription operation binding the contract event 0x6a036bee1b01306b61370b348f57c9a7038a7bf8cb5d8a4cfbfb197b3f329e83.
//
// Solidity: event Code501()
func (_Vault *VaultFilterer) WatchCode501(opts *bind.WatchOpts, sink chan<- *VaultCode501) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Code501")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultCode501)
				if err := _Vault.contract.UnpackLog(event, "Code501", log); err != nil {
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

// ParseCode501 is a log parse operation binding the contract event 0x6a036bee1b01306b61370b348f57c9a7038a7bf8cb5d8a4cfbfb197b3f329e83.
//
// Solidity: event Code501()
func (_Vault *VaultFilterer) ParseCode501(log types.Log) (*VaultCode501, error) {
	event := new(VaultCode501)
	if err := _Vault.contract.UnpackLog(event, "Code501", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultCode502Iterator is returned from FilterCode502 and is used to iterate over the raw logs and unpacked data for Code502 events raised by the Vault contract.
type VaultCode502Iterator struct {
	Event *VaultCode502 // Event containing the contract specifics and raw log

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
func (it *VaultCode502Iterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultCode502)
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
		it.Event = new(VaultCode502)
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
func (it *VaultCode502Iterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultCode502Iterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultCode502 represents a Code502 event raised by the Vault contract.
type VaultCode502 struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterCode502 is a free log retrieval operation binding the contract event 0xe1ef72f1796988294cc9a59a2e3073a1079dc556e353a1ebf43b456b191161a9.
//
// Solidity: event Code502()
func (_Vault *VaultFilterer) FilterCode502(opts *bind.FilterOpts) (*VaultCode502Iterator, error) {

	logs, sub, err := _Vault.contract.FilterLogs(opts, "Code502")
	if err != nil {
		return nil, err
	}
	return &VaultCode502Iterator{contract: _Vault.contract, event: "Code502", logs: logs, sub: sub}, nil
}

// WatchCode502 is a free log subscription operation binding the contract event 0xe1ef72f1796988294cc9a59a2e3073a1079dc556e353a1ebf43b456b191161a9.
//
// Solidity: event Code502()
func (_Vault *VaultFilterer) WatchCode502(opts *bind.WatchOpts, sink chan<- *VaultCode502) (event.Subscription, error) {

	logs, sub, err := _Vault.contract.WatchLogs(opts, "Code502")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultCode502)
				if err := _Vault.contract.UnpackLog(event, "Code502", log); err != nil {
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

// ParseCode502 is a log parse operation binding the contract event 0xe1ef72f1796988294cc9a59a2e3073a1079dc556e353a1ebf43b456b191161a9.
//
// Solidity: event Code502()
func (_Vault *VaultFilterer) ParseCode502(log types.Log) (*VaultCode502, error) {
	event := new(VaultCode502)
	if err := _Vault.contract.UnpackLog(event, "Code502", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
