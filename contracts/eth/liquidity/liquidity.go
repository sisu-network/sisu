// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package liquidity

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

// LiquidityMetaData contains all meta data concerning the Liquidity contract.
var LiquidityMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenAddrs\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"lpTokens\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"lpTokenAmt\",\"type\":\"uint256\"}],\"name\":\"AddLiquidity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"burnLPTokenAmt\",\"type\":\"uint256\"}],\"name\":\"RemoveLiquidity\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"calculateLPTokenDepositOrWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"liquidityPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"rewardDebt\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lpTokenMapping\",\"outputs\":[{\"internalType\":\"contractLPToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gateway\",\"type\":\"address\"}],\"name\":\"setGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200238f3803806200238f833981810160405281019062000037919062000404565b620000576200004b6200014a60201b60201c565b6200015260201b60201c565b80518251146200006657600080fd5b60005b8251811015620001415781818151811062000089576200008862000489565b5b602002602001015160036000858481518110620000ab57620000aa62000489565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080806200013890620004f1565b91505062000069565b5050506200053f565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200027a826200022f565b810181811067ffffffffffffffff821117156200029c576200029b62000240565b5b80604052505050565b6000620002b162000216565b9050620002bf82826200026f565b919050565b600067ffffffffffffffff821115620002e257620002e162000240565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200032582620002f8565b9050919050565b620003378162000318565b81146200034357600080fd5b50565b60008151905062000357816200032c565b92915050565b6000620003746200036e84620002c4565b620002a5565b905080838252602082019050602084028301858111156200039a5762000399620002f3565b5b835b81811015620003c75780620003b2888262000346565b8452602084019350506020810190506200039c565b5050509392505050565b600082601f830112620003e957620003e86200022a565b5b8151620003fb8482602086016200035d565b91505092915050565b600080604083850312156200041e576200041d62000220565b5b600083015167ffffffffffffffff8111156200043f576200043e62000225565b5b6200044d85828601620003d1565b925050602083015167ffffffffffffffff81111562000471576200047062000225565b5b6200047f85828601620003d1565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b6000620004fe82620004e7565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620005345762000533620004b8565b5b600182019050919050565b611e40806200054f6000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806390646b4a1161007157806390646b4a14610140578063a201ccf61461015c578063beabacc814610178578063cc8edc0914610194578063e9e43f2c146101c5578063f2fde38b146101f5576100a9565b8063116191b6146100ae57806331b2554d146100cc57806356688700146100fc578063715018a6146101185780638da5cb5b14610122575b600080fd5b6100b6610211565b6040516100c391906112f3565b60405180910390f35b6100e660048036038101906100e1919061133f565b610237565b6040516100f391906113cb565b60405180910390f35b6101166004803603810190610111919061141c565b61026a565b005b610120610570565b005b61012a6105f8565b60405161013791906112f3565b60405180910390f35b61015a6004803603810190610155919061133f565b610621565b005b6101766004803603810190610171919061141c565b6106e1565b005b610192600480360381019061018d919061145c565b6109f7565b005b6101ae60048036038101906101a991906114af565b610afa565b6040516101bc929190611517565b60405180910390f35b6101df60048036038101906101da919061141c565b610b2b565b6040516101ec9190611540565b60405180910390f35b61020f600480360381019061020a919061133f565b610d1a565b005b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60036020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600081116102ad576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102a4906115b8565b60405180910390fd5b60008273ffffffffffffffffffffffffffffffffffffffff166370a08231336040518263ffffffff1660e01b81526004016102e891906112f3565b602060405180830381865afa158015610305573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061032991906115ed565b90508181101561036e576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016103659061168c565b60405180910390fd5b61039b3330848673ffffffffffffffffffffffffffffffffffffffff16610e12909392919063ffffffff16565b6000600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050610432838260000154610e9b90919063ffffffff16565b816000018190555060006104468585610b2b565b90506000811115610524576000600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff166340c10f1933846040518363ffffffff1660e01b81526004016104f09291906116ac565b600060405180830381600087803b15801561050a57600080fd5b505af115801561051e573d6000803e3d6000fd5b50505050505b80848673ffffffffffffffffffffffffffffffffffffffff167f06239653922ac7bea6aa2b19dc486b9361821d37712eb796adfd38d81de278ca60405160405180910390a45050505050565b610578610eb1565b73ffffffffffffffffffffffffffffffffffffffff166105966105f8565b73ffffffffffffffffffffffffffffffffffffffff16146105ec576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016105e390611721565b60405180910390fd5b6105f66000610eb9565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610629610eb1565b73ffffffffffffffffffffffffffffffffffffffff166106476105f8565b73ffffffffffffffffffffffffffffffffffffffff161461069d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161069490611721565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b60008111610724576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161071b906115b8565b60405180910390fd5b60008273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161075f91906112f3565b602060405180830381865afa15801561077c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107a091906115ed565b9050818110156107e5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107dc906117b3565b60405180910390fd5b6000600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020905082816000015410156108ac576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108a390611845565b60405180910390fd5b60006108b88585610b2b565b90506000600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff16639dc29fac33846040518363ffffffff1660e01b81526004016109599291906116ac565b600060405180830381600087803b15801561097357600080fd5b505af1158015610987573d6000803e3d6000fd5b505050506109a2858460000154610f7d90919063ffffffff16565b836000018190555081858773ffffffffffffffffffffffffffffffffffffffff167f0fbf06c058b90cb038a618f8c2acbf6145f8b3570fd1fa56abb8f0f3f05b36e860405160405180910390a4505050505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610a87576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a7e906118b1565b60405180910390fd5b60008111610aca576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610ac1906115b8565b60405180910390fd5b610af582828573ffffffffffffffffffffffffffffffffffffffff16610f939092919063ffffffff16565b505050565b6002602052816000526040600020602052806000526040600020600091509150508060000154908060010154905082565b6000808211610b6f576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610b66906115b8565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610baa91906112f3565b602060405180830381865afa158015610bc7573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610beb91906115ed565b90506000600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610c9e573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610cc291906115ed565b90506000831480610cd35750600081145b15610ce357849350505050610d14565b6000610d0a84610cfc848961101990919063ffffffff16565b61102f90919063ffffffff16565b9050809450505050505b92915050565b610d22610eb1565b73ffffffffffffffffffffffffffffffffffffffff16610d406105f8565b73ffffffffffffffffffffffffffffffffffffffff1614610d96576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d8d90611721565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610e06576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dfd90611943565b60405180910390fd5b610e0f81610eb9565b50565b610e95846323b872dd60e01b858585604051602401610e3393929190611963565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611045565b50505050565b60008183610ea991906119c9565b905092915050565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b60008183610f8b9190611a1f565b905092915050565b6110148363a9059cbb60e01b8484604051602401610fb29291906116ac565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff8381831617835250505050611045565b505050565b600081836110279190611a53565b905092915050565b6000818361103d9190611adc565b905092915050565b60006110a7826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff1661110c9092919063ffffffff16565b905060008151111561110757808060200190518101906110c79190611b45565b611106576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016110fd90611be4565b60405180910390fd5b5b505050565b606061111b8484600085611124565b90509392505050565b606082471015611169576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161116090611c76565b60405180910390fd5b61117285611238565b6111b1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016111a890611ce2565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516111da9190611d7c565b60006040518083038185875af1925050503d8060008114611217576040519150601f19603f3d011682016040523d82523d6000602084013e61121c565b606091505b509150915061122c82828661124b565b92505050949350505050565b600080823b905060008111915050919050565b6060831561125b578290506112ab565b60008351111561126e5782518084602001fd5b816040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112a29190611de8565b60405180910390fd5b9392505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006112dd826112b2565b9050919050565b6112ed816112d2565b82525050565b600060208201905061130860008301846112e4565b92915050565b600080fd5b61131c816112d2565b811461132757600080fd5b50565b60008135905061133981611313565b92915050565b6000602082840312156113555761135461130e565b5b60006113638482850161132a565b91505092915050565b6000819050919050565b600061139161138c611387846112b2565b61136c565b6112b2565b9050919050565b60006113a382611376565b9050919050565b60006113b582611398565b9050919050565b6113c5816113aa565b82525050565b60006020820190506113e060008301846113bc565b92915050565b6000819050919050565b6113f9816113e6565b811461140457600080fd5b50565b600081359050611416816113f0565b92915050565b600080604083850312156114335761143261130e565b5b60006114418582860161132a565b925050602061145285828601611407565b9150509250929050565b6000806000606084860312156114755761147461130e565b5b60006114838682870161132a565b93505060206114948682870161132a565b92505060406114a586828701611407565b9150509250925092565b600080604083850312156114c6576114c561130e565b5b60006114d48582860161132a565b92505060206114e58582860161132a565b9150509250929050565b6114f8816113e6565b82525050565b6000819050919050565b611511816114fe565b82525050565b600060408201905061152c60008301856114ef565b6115396020830184611508565b9392505050565b600060208201905061155560008301846114ef565b92915050565b600082825260208201905092915050565b7f616d6f756e74206d7573742067726561746572207468616e2030000000000000600082015250565b60006115a2601a8361155b565b91506115ad8261156c565b602082019050919050565b600060208201905081810360008301526115d181611595565b9050919050565b6000815190506115e7816113f0565b92915050565b6000602082840312156116035761160261130e565b5b6000611611848285016115d8565b91505092915050565b7f7573657227732062616c616e6365206973206c657373207468616e207265717560008201527f6972656420616d6f756e74000000000000000000000000000000000000000000602082015250565b6000611676602b8361155b565b91506116818261161a565b604082019050919050565b600060208201905081810360008301526116a581611669565b9050919050565b60006040820190506116c160008301856112e4565b6116ce60208301846114ef565b9392505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b600061170b60208361155b565b9150611716826116d5565b602082019050919050565b6000602082019050818103600083015261173a816116fe565b9050919050565b7f6761746577617927732062616c616e6365206973206c657373207468616e207260008201527f6571756972656420616d6f756e74000000000000000000000000000000000000602082015250565b600061179d602e8361155b565b91506117a882611741565b604082019050919050565b600060208201905081810360008301526117cc81611790565b9050919050565b7f6465706f736974656420746f6b656e20616d6f756e74206973206c657373207460008201527f68616e20776974686472617720746f6b656e20616d6f756e7400000000000000602082015250565b600061182f60398361155b565b915061183a826117d3565b604082019050919050565b6000602082019050818103600083015261185e81611822565b9050919050565b7f4f6e6c7920676174657761790000000000000000000000000000000000000000600082015250565b600061189b600c8361155b565b91506118a682611865565b602082019050919050565b600060208201905081810360008301526118ca8161188e565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b600061192d60268361155b565b9150611938826118d1565b604082019050919050565b6000602082019050818103600083015261195c81611920565b9050919050565b600060608201905061197860008301866112e4565b61198560208301856112e4565b61199260408301846114ef565b949350505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006119d4826113e6565b91506119df836113e6565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611a1457611a1361199a565b5b828201905092915050565b6000611a2a826113e6565b9150611a35836113e6565b925082821015611a4857611a4761199a565b5b828203905092915050565b6000611a5e826113e6565b9150611a69836113e6565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611aa257611aa161199a565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611ae7826113e6565b9150611af2836113e6565b925082611b0257611b01611aad565b5b828204905092915050565b60008115159050919050565b611b2281611b0d565b8114611b2d57600080fd5b50565b600081519050611b3f81611b19565b92915050565b600060208284031215611b5b57611b5a61130e565b5b6000611b6984828501611b30565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000611bce602a8361155b565b9150611bd982611b72565b604082019050919050565b60006020820190508181036000830152611bfd81611bc1565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b6000611c6060268361155b565b9150611c6b82611c04565b604082019050919050565b60006020820190508181036000830152611c8f81611c53565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b6000611ccc601d8361155b565b9150611cd782611c96565b602082019050919050565b60006020820190508181036000830152611cfb81611cbf565b9050919050565b600081519050919050565b600081905092915050565b60005b83811015611d36578082015181840152602081019050611d1b565b83811115611d45576000848401525b50505050565b6000611d5682611d02565b611d608185611d0d565b9350611d70818560208601611d18565b80840191505092915050565b6000611d888284611d4b565b915081905092915050565b600081519050919050565b6000601f19601f8301169050919050565b6000611dba82611d93565b611dc4818561155b565b9350611dd4818560208601611d18565b611ddd81611d9e565b840191505092915050565b60006020820190508181036000830152611e028184611daf565b90509291505056fea264697066735822122033bf8d2b10a53ba1cc8f9c4c1f9ad41cb5c096427072de467db2aa688a249c9864736f6c634300080a0033",
}

// LiquidityABI is the input ABI used to generate the binding from.
// Deprecated: Use LiquidityMetaData.ABI instead.
var LiquidityABI = LiquidityMetaData.ABI

// LiquidityBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LiquidityMetaData.Bin instead.
var LiquidityBin = LiquidityMetaData.Bin

// DeployLiquidity deploys a new Ethereum contract, binding an instance of Liquidity to it.
func DeployLiquidity(auth *bind.TransactOpts, backend bind.ContractBackend, tokenAddrs []common.Address, lpTokens []common.Address) (common.Address, *types.Transaction, *Liquidity, error) {
	parsed, err := LiquidityMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LiquidityBin), backend, tokenAddrs, lpTokens)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Liquidity{LiquidityCaller: LiquidityCaller{contract: contract}, LiquidityTransactor: LiquidityTransactor{contract: contract}, LiquidityFilterer: LiquidityFilterer{contract: contract}}, nil
}

// Liquidity is an auto generated Go binding around an Ethereum contract.
type Liquidity struct {
	LiquidityCaller     // Read-only binding to the contract
	LiquidityTransactor // Write-only binding to the contract
	LiquidityFilterer   // Log filterer for contract events
}

// LiquidityCaller is an auto generated read-only Go binding around an Ethereum contract.
type LiquidityCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LiquidityTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquidityFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LiquidityFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LiquiditySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LiquiditySession struct {
	Contract     *Liquidity        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LiquidityCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LiquidityCallerSession struct {
	Contract *LiquidityCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// LiquidityTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LiquidityTransactorSession struct {
	Contract     *LiquidityTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// LiquidityRaw is an auto generated low-level Go binding around an Ethereum contract.
type LiquidityRaw struct {
	Contract *Liquidity // Generic contract binding to access the raw methods on
}

// LiquidityCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LiquidityCallerRaw struct {
	Contract *LiquidityCaller // Generic read-only contract binding to access the raw methods on
}

// LiquidityTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LiquidityTransactorRaw struct {
	Contract *LiquidityTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLiquidity creates a new instance of Liquidity, bound to a specific deployed contract.
func NewLiquidity(address common.Address, backend bind.ContractBackend) (*Liquidity, error) {
	contract, err := bindLiquidity(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Liquidity{LiquidityCaller: LiquidityCaller{contract: contract}, LiquidityTransactor: LiquidityTransactor{contract: contract}, LiquidityFilterer: LiquidityFilterer{contract: contract}}, nil
}

// NewLiquidityCaller creates a new read-only instance of Liquidity, bound to a specific deployed contract.
func NewLiquidityCaller(address common.Address, caller bind.ContractCaller) (*LiquidityCaller, error) {
	contract, err := bindLiquidity(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityCaller{contract: contract}, nil
}

// NewLiquidityTransactor creates a new write-only instance of Liquidity, bound to a specific deployed contract.
func NewLiquidityTransactor(address common.Address, transactor bind.ContractTransactor) (*LiquidityTransactor, error) {
	contract, err := bindLiquidity(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LiquidityTransactor{contract: contract}, nil
}

// NewLiquidityFilterer creates a new log filterer instance of Liquidity, bound to a specific deployed contract.
func NewLiquidityFilterer(address common.Address, filterer bind.ContractFilterer) (*LiquidityFilterer, error) {
	contract, err := bindLiquidity(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LiquidityFilterer{contract: contract}, nil
}

// bindLiquidity binds a generic wrapper to an already deployed contract.
func bindLiquidity(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(LiquidityABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquidity *LiquidityRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquidity.Contract.LiquidityCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquidity *LiquidityRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquidity.Contract.LiquidityTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquidity *LiquidityRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquidity.Contract.LiquidityTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Liquidity *LiquidityCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Liquidity.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Liquidity *LiquidityTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquidity.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Liquidity *LiquidityTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Liquidity.Contract.contract.Transact(opts, method, params...)
}

// CalculateLPTokenDepositOrWithdraw is a free data retrieval call binding the contract method 0xe9e43f2c.
//
// Solidity: function calculateLPTokenDepositOrWithdraw(address _token, uint256 _amount) view returns(uint256)
func (_Liquidity *LiquidityCaller) CalculateLPTokenDepositOrWithdraw(opts *bind.CallOpts, _token common.Address, _amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Liquidity.contract.Call(opts, &out, "calculateLPTokenDepositOrWithdraw", _token, _amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateLPTokenDepositOrWithdraw is a free data retrieval call binding the contract method 0xe9e43f2c.
//
// Solidity: function calculateLPTokenDepositOrWithdraw(address _token, uint256 _amount) view returns(uint256)
func (_Liquidity *LiquiditySession) CalculateLPTokenDepositOrWithdraw(_token common.Address, _amount *big.Int) (*big.Int, error) {
	return _Liquidity.Contract.CalculateLPTokenDepositOrWithdraw(&_Liquidity.CallOpts, _token, _amount)
}

// CalculateLPTokenDepositOrWithdraw is a free data retrieval call binding the contract method 0xe9e43f2c.
//
// Solidity: function calculateLPTokenDepositOrWithdraw(address _token, uint256 _amount) view returns(uint256)
func (_Liquidity *LiquidityCallerSession) CalculateLPTokenDepositOrWithdraw(_token common.Address, _amount *big.Int) (*big.Int, error) {
	return _Liquidity.Contract.CalculateLPTokenDepositOrWithdraw(&_Liquidity.CallOpts, _token, _amount)
}

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_Liquidity *LiquidityCaller) Gateway(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Liquidity.contract.Call(opts, &out, "gateway")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_Liquidity *LiquiditySession) Gateway() (common.Address, error) {
	return _Liquidity.Contract.Gateway(&_Liquidity.CallOpts)
}

// Gateway is a free data retrieval call binding the contract method 0x116191b6.
//
// Solidity: function gateway() view returns(address)
func (_Liquidity *LiquidityCallerSession) Gateway() (common.Address, error) {
	return _Liquidity.Contract.Gateway(&_Liquidity.CallOpts)
}

// LiquidityPool is a free data retrieval call binding the contract method 0xcc8edc09.
//
// Solidity: function liquidityPool(address , address ) view returns(uint256 amount, int256 rewardDebt)
func (_Liquidity *LiquidityCaller) LiquidityPool(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (struct {
	Amount     *big.Int
	RewardDebt *big.Int
}, error) {
	var out []interface{}
	err := _Liquidity.contract.Call(opts, &out, "liquidityPool", arg0, arg1)

	outstruct := new(struct {
		Amount     *big.Int
		RewardDebt *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Amount = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.RewardDebt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// LiquidityPool is a free data retrieval call binding the contract method 0xcc8edc09.
//
// Solidity: function liquidityPool(address , address ) view returns(uint256 amount, int256 rewardDebt)
func (_Liquidity *LiquiditySession) LiquidityPool(arg0 common.Address, arg1 common.Address) (struct {
	Amount     *big.Int
	RewardDebt *big.Int
}, error) {
	return _Liquidity.Contract.LiquidityPool(&_Liquidity.CallOpts, arg0, arg1)
}

// LiquidityPool is a free data retrieval call binding the contract method 0xcc8edc09.
//
// Solidity: function liquidityPool(address , address ) view returns(uint256 amount, int256 rewardDebt)
func (_Liquidity *LiquidityCallerSession) LiquidityPool(arg0 common.Address, arg1 common.Address) (struct {
	Amount     *big.Int
	RewardDebt *big.Int
}, error) {
	return _Liquidity.Contract.LiquidityPool(&_Liquidity.CallOpts, arg0, arg1)
}

// LpTokenMapping is a free data retrieval call binding the contract method 0x31b2554d.
//
// Solidity: function lpTokenMapping(address ) view returns(address)
func (_Liquidity *LiquidityCaller) LpTokenMapping(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _Liquidity.contract.Call(opts, &out, "lpTokenMapping", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpTokenMapping is a free data retrieval call binding the contract method 0x31b2554d.
//
// Solidity: function lpTokenMapping(address ) view returns(address)
func (_Liquidity *LiquiditySession) LpTokenMapping(arg0 common.Address) (common.Address, error) {
	return _Liquidity.Contract.LpTokenMapping(&_Liquidity.CallOpts, arg0)
}

// LpTokenMapping is a free data retrieval call binding the contract method 0x31b2554d.
//
// Solidity: function lpTokenMapping(address ) view returns(address)
func (_Liquidity *LiquidityCallerSession) LpTokenMapping(arg0 common.Address) (common.Address, error) {
	return _Liquidity.Contract.LpTokenMapping(&_Liquidity.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Liquidity *LiquidityCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Liquidity.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Liquidity *LiquiditySession) Owner() (common.Address, error) {
	return _Liquidity.Contract.Owner(&_Liquidity.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Liquidity *LiquidityCallerSession) Owner() (common.Address, error) {
	return _Liquidity.Contract.Owner(&_Liquidity.CallOpts)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x56688700.
//
// Solidity: function addLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactor) AddLiquidity(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "addLiquidity", _token, _amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x56688700.
//
// Solidity: function addLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquiditySession) AddLiquidity(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.AddLiquidity(&_Liquidity.TransactOpts, _token, _amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x56688700.
//
// Solidity: function addLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactorSession) AddLiquidity(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.AddLiquidity(&_Liquidity.TransactOpts, _token, _amount)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xa201ccf6.
//
// Solidity: function removeLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactor) RemoveLiquidity(opts *bind.TransactOpts, _token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "removeLiquidity", _token, _amount)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xa201ccf6.
//
// Solidity: function removeLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquiditySession) RemoveLiquidity(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.RemoveLiquidity(&_Liquidity.TransactOpts, _token, _amount)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xa201ccf6.
//
// Solidity: function removeLiquidity(address _token, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactorSession) RemoveLiquidity(_token common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.RemoveLiquidity(&_Liquidity.TransactOpts, _token, _amount)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Liquidity *LiquidityTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Liquidity *LiquiditySession) RenounceOwnership() (*types.Transaction, error) {
	return _Liquidity.Contract.RenounceOwnership(&_Liquidity.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Liquidity *LiquidityTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Liquidity.Contract.RenounceOwnership(&_Liquidity.TransactOpts)
}

// SetGateway is a paid mutator transaction binding the contract method 0x90646b4a.
//
// Solidity: function setGateway(address _gateway) returns()
func (_Liquidity *LiquidityTransactor) SetGateway(opts *bind.TransactOpts, _gateway common.Address) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "setGateway", _gateway)
}

// SetGateway is a paid mutator transaction binding the contract method 0x90646b4a.
//
// Solidity: function setGateway(address _gateway) returns()
func (_Liquidity *LiquiditySession) SetGateway(_gateway common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.SetGateway(&_Liquidity.TransactOpts, _gateway)
}

// SetGateway is a paid mutator transaction binding the contract method 0x90646b4a.
//
// Solidity: function setGateway(address _gateway) returns()
func (_Liquidity *LiquidityTransactorSession) SetGateway(_gateway common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.SetGateway(&_Liquidity.TransactOpts, _gateway)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address _token, address _recipient, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactor) Transfer(opts *bind.TransactOpts, _token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "transfer", _token, _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address _token, address _recipient, uint256 _amount) returns()
func (_Liquidity *LiquiditySession) Transfer(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.Transfer(&_Liquidity.TransactOpts, _token, _recipient, _amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xbeabacc8.
//
// Solidity: function transfer(address _token, address _recipient, uint256 _amount) returns()
func (_Liquidity *LiquidityTransactorSession) Transfer(_token common.Address, _recipient common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Liquidity.Contract.Transfer(&_Liquidity.TransactOpts, _token, _recipient, _amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Liquidity *LiquidityTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Liquidity *LiquiditySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.TransferOwnership(&_Liquidity.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Liquidity *LiquidityTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.TransferOwnership(&_Liquidity.TransactOpts, newOwner)
}

// LiquidityAddLiquidityIterator is returned from FilterAddLiquidity and is used to iterate over the raw logs and unpacked data for AddLiquidity events raised by the Liquidity contract.
type LiquidityAddLiquidityIterator struct {
	Event *LiquidityAddLiquidity // Event containing the contract specifics and raw log

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
func (it *LiquidityAddLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityAddLiquidity)
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
		it.Event = new(LiquidityAddLiquidity)
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
func (it *LiquidityAddLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityAddLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityAddLiquidity represents a AddLiquidity event raised by the Liquidity contract.
type LiquidityAddLiquidity struct {
	Token      common.Address
	Amount     *big.Int
	LpTokenAmt *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterAddLiquidity is a free log retrieval operation binding the contract event 0x06239653922ac7bea6aa2b19dc486b9361821d37712eb796adfd38d81de278ca.
//
// Solidity: event AddLiquidity(address indexed token, uint256 indexed amount, uint256 indexed lpTokenAmt)
func (_Liquidity *LiquidityFilterer) FilterAddLiquidity(opts *bind.FilterOpts, token []common.Address, amount []*big.Int, lpTokenAmt []*big.Int) (*LiquidityAddLiquidityIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var lpTokenAmtRule []interface{}
	for _, lpTokenAmtItem := range lpTokenAmt {
		lpTokenAmtRule = append(lpTokenAmtRule, lpTokenAmtItem)
	}

	logs, sub, err := _Liquidity.contract.FilterLogs(opts, "AddLiquidity", tokenRule, amountRule, lpTokenAmtRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityAddLiquidityIterator{contract: _Liquidity.contract, event: "AddLiquidity", logs: logs, sub: sub}, nil
}

// WatchAddLiquidity is a free log subscription operation binding the contract event 0x06239653922ac7bea6aa2b19dc486b9361821d37712eb796adfd38d81de278ca.
//
// Solidity: event AddLiquidity(address indexed token, uint256 indexed amount, uint256 indexed lpTokenAmt)
func (_Liquidity *LiquidityFilterer) WatchAddLiquidity(opts *bind.WatchOpts, sink chan<- *LiquidityAddLiquidity, token []common.Address, amount []*big.Int, lpTokenAmt []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var lpTokenAmtRule []interface{}
	for _, lpTokenAmtItem := range lpTokenAmt {
		lpTokenAmtRule = append(lpTokenAmtRule, lpTokenAmtItem)
	}

	logs, sub, err := _Liquidity.contract.WatchLogs(opts, "AddLiquidity", tokenRule, amountRule, lpTokenAmtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityAddLiquidity)
				if err := _Liquidity.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
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

// ParseAddLiquidity is a log parse operation binding the contract event 0x06239653922ac7bea6aa2b19dc486b9361821d37712eb796adfd38d81de278ca.
//
// Solidity: event AddLiquidity(address indexed token, uint256 indexed amount, uint256 indexed lpTokenAmt)
func (_Liquidity *LiquidityFilterer) ParseAddLiquidity(log types.Log) (*LiquidityAddLiquidity, error) {
	event := new(LiquidityAddLiquidity)
	if err := _Liquidity.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Liquidity contract.
type LiquidityOwnershipTransferredIterator struct {
	Event *LiquidityOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *LiquidityOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityOwnershipTransferred)
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
		it.Event = new(LiquidityOwnershipTransferred)
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
func (it *LiquidityOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityOwnershipTransferred represents a OwnershipTransferred event raised by the Liquidity contract.
type LiquidityOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Liquidity *LiquidityFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LiquidityOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Liquidity.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityOwnershipTransferredIterator{contract: _Liquidity.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Liquidity *LiquidityFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LiquidityOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Liquidity.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityOwnershipTransferred)
				if err := _Liquidity.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Liquidity *LiquidityFilterer) ParseOwnershipTransferred(log types.Log) (*LiquidityOwnershipTransferred, error) {
	event := new(LiquidityOwnershipTransferred)
	if err := _Liquidity.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LiquidityRemoveLiquidityIterator is returned from FilterRemoveLiquidity and is used to iterate over the raw logs and unpacked data for RemoveLiquidity events raised by the Liquidity contract.
type LiquidityRemoveLiquidityIterator struct {
	Event *LiquidityRemoveLiquidity // Event containing the contract specifics and raw log

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
func (it *LiquidityRemoveLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LiquidityRemoveLiquidity)
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
		it.Event = new(LiquidityRemoveLiquidity)
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
func (it *LiquidityRemoveLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LiquidityRemoveLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LiquidityRemoveLiquidity represents a RemoveLiquidity event raised by the Liquidity contract.
type LiquidityRemoveLiquidity struct {
	Token          common.Address
	Amount         *big.Int
	BurnLPTokenAmt *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidity is a free log retrieval operation binding the contract event 0x0fbf06c058b90cb038a618f8c2acbf6145f8b3570fd1fa56abb8f0f3f05b36e8.
//
// Solidity: event RemoveLiquidity(address indexed token, uint256 indexed amount, uint256 indexed burnLPTokenAmt)
func (_Liquidity *LiquidityFilterer) FilterRemoveLiquidity(opts *bind.FilterOpts, token []common.Address, amount []*big.Int, burnLPTokenAmt []*big.Int) (*LiquidityRemoveLiquidityIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var burnLPTokenAmtRule []interface{}
	for _, burnLPTokenAmtItem := range burnLPTokenAmt {
		burnLPTokenAmtRule = append(burnLPTokenAmtRule, burnLPTokenAmtItem)
	}

	logs, sub, err := _Liquidity.contract.FilterLogs(opts, "RemoveLiquidity", tokenRule, amountRule, burnLPTokenAmtRule)
	if err != nil {
		return nil, err
	}
	return &LiquidityRemoveLiquidityIterator{contract: _Liquidity.contract, event: "RemoveLiquidity", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidity is a free log subscription operation binding the contract event 0x0fbf06c058b90cb038a618f8c2acbf6145f8b3570fd1fa56abb8f0f3f05b36e8.
//
// Solidity: event RemoveLiquidity(address indexed token, uint256 indexed amount, uint256 indexed burnLPTokenAmt)
func (_Liquidity *LiquidityFilterer) WatchRemoveLiquidity(opts *bind.WatchOpts, sink chan<- *LiquidityRemoveLiquidity, token []common.Address, amount []*big.Int, burnLPTokenAmt []*big.Int) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var burnLPTokenAmtRule []interface{}
	for _, burnLPTokenAmtItem := range burnLPTokenAmt {
		burnLPTokenAmtRule = append(burnLPTokenAmtRule, burnLPTokenAmtItem)
	}

	logs, sub, err := _Liquidity.contract.WatchLogs(opts, "RemoveLiquidity", tokenRule, amountRule, burnLPTokenAmtRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LiquidityRemoveLiquidity)
				if err := _Liquidity.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
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

// ParseRemoveLiquidity is a log parse operation binding the contract event 0x0fbf06c058b90cb038a618f8c2acbf6145f8b3570fd1fa56abb8f0f3f05b36e8.
//
// Solidity: event RemoveLiquidity(address indexed token, uint256 indexed amount, uint256 indexed burnLPTokenAmt)
func (_Liquidity *LiquidityFilterer) ParseRemoveLiquidity(log types.Log) (*LiquidityRemoveLiquidity, error) {
	event := new(LiquidityRemoveLiquidity)
	if err := _Liquidity.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
