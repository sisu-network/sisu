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
	ABI: "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"tokenAddrs\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"lpTokens\",\"type\":\"address[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"lpTokenAmt\",\"type\":\"uint256\"}],\"name\":\"AddLiquidity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"burnLPTokenAmt\",\"type\":\"uint256\"}],\"name\":\"RemoveLiquidity\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"addLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"calculateLPTokenDepositOrWithdraw\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"_tokens\",\"type\":\"address[]\"},{\"internalType\":\"address\",\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"emergencyWithdrawFunds\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gateway\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"liquidityPool\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"rewardDebt\",\"type\":\"int256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"lpTokenMapping\",\"outputs\":[{\"internalType\":\"contractLPToken\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"removeLiquidity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_gateway\",\"type\":\"address\"}],\"name\":\"setGateway\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040523480156200001157600080fd5b506040516200276938038062002769833981810160405281019062000037919062000404565b620000576200004b6200014a60201b60201c565b6200015260201b60201c565b80518251146200006657600080fd5b60005b8251811015620001415781818151811062000089576200008862000489565b5b602002602001015160036000858481518110620000ab57620000aa62000489565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080806200013890620004f1565b91505062000069565b5050506200053f565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6200027a826200022f565b810181811067ffffffffffffffff821117156200029c576200029b62000240565b5b80604052505050565b6000620002b162000216565b9050620002bf82826200026f565b919050565b600067ffffffffffffffff821115620002e257620002e162000240565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006200032582620002f8565b9050919050565b620003378162000318565b81146200034357600080fd5b50565b60008151905062000357816200032c565b92915050565b6000620003746200036e84620002c4565b620002a5565b905080838252602082019050602084028301858111156200039a5762000399620002f3565b5b835b81811015620003c75780620003b2888262000346565b8452602084019350506020810190506200039c565b5050509392505050565b600082601f830112620003e957620003e86200022a565b5b8151620003fb8482602086016200035d565b91505092915050565b600080604083850312156200041e576200041d62000220565b5b600083015167ffffffffffffffff8111156200043f576200043e62000225565b5b6200044d85828601620003d1565b925050602083015167ffffffffffffffff81111562000471576200047062000225565b5b6200047f85828601620003d1565b9150509250929050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b6000620004fe82620004e7565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415620005345762000533620004b8565b5b600182019050919050565b61221a806200054f6000396000f3fe608060405234801561001057600080fd5b50600436106100b45760003560e01c8063a201ccf611610071578063a201ccf614610167578063aea3da6c14610183578063beabacc81461019f578063cc8edc09146101bb578063e9e43f2c146101ec578063f2fde38b1461021c576100b4565b8063116191b6146100b957806331b2554d146100d75780635668870014610107578063715018a6146101235780638da5cb5b1461012d57806390646b4a1461014b575b600080fd5b6100c1610238565b6040516100ce91906114a2565b60405180910390f35b6100f160048036038101906100ec91906114fd565b61025e565b6040516100fe9190611589565b60405180910390f35b610121600480360381019061011c91906115da565b610291565b005b61012b610597565b005b61013561061f565b60405161014291906114a2565b60405180910390f35b610165600480360381019061016091906114fd565b610648565b005b610181600480360381019061017c91906115da565b610708565b005b61019d60048036038101906101989190611773565b610a1e565b005b6101b960048036038101906101b491906117cf565b610ba6565b005b6101d560048036038101906101d09190611822565b610ca9565b6040516101e392919061188a565b60405180910390f35b610206600480360381019061020191906115da565b610cda565b60405161021391906118b3565b60405180910390f35b610236600480360381019061023191906114fd565b610ec9565b005b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60036020528060005260406000206000915054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600081116102d4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102cb9061192b565b60405180910390fd5b60008273ffffffffffffffffffffffffffffffffffffffff166370a08231336040518263ffffffff1660e01b815260040161030f91906114a2565b602060405180830381865afa15801561032c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906103509190611960565b905081811015610395576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161038c906119ff565b60405180910390fd5b6103c23330848673ffffffffffffffffffffffffffffffffffffffff16610fc1909392919063ffffffff16565b6000600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020905061045983826000015461104a90919063ffffffff16565b8160000181905550600061046d8585610cda565b9050600081111561054b576000600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff166340c10f1933846040518363ffffffff1660e01b8152600401610517929190611a1f565b600060405180830381600087803b15801561053157600080fd5b505af1158015610545573d6000803e3d6000fd5b50505050505b80848673ffffffffffffffffffffffffffffffffffffffff167f06239653922ac7bea6aa2b19dc486b9361821d37712eb796adfd38d81de278ca60405160405180910390a45050505050565b61059f611060565b73ffffffffffffffffffffffffffffffffffffffff166105bd61061f565b73ffffffffffffffffffffffffffffffffffffffff1614610613576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161060a90611a94565b60405180910390fd5b61061d6000611068565b565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b610650611060565b73ffffffffffffffffffffffffffffffffffffffff1661066e61061f565b73ffffffffffffffffffffffffffffffffffffffff16146106c4576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106bb90611a94565b60405180910390fd5b80600160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b6000811161074b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016107429061192b565b60405180910390fd5b60008273ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b815260040161078691906114a2565b602060405180830381865afa1580156107a3573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906107c79190611960565b90508181101561080c576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161080390611b26565b60405180910390fd5b6000600260008573ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020905082816000015410156108d3576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108ca90611bb8565b60405180910390fd5b60006108df8585610cda565b90506000600360008773ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690508073ffffffffffffffffffffffffffffffffffffffff16639dc29fac33846040518363ffffffff1660e01b8152600401610980929190611a1f565b600060405180830381600087803b15801561099a57600080fd5b505af11580156109ae573d6000803e3d6000fd5b505050506109c985846000015461112c90919063ffffffff16565b836000018190555081858773ffffffffffffffffffffffffffffffffffffffff167f0fbf06c058b90cb038a618f8c2acbf6145f8b3570fd1fa56abb8f0f3f05b36e860405160405180910390a4505050505050565b610a26611060565b73ffffffffffffffffffffffffffffffffffffffff16610a4461061f565b73ffffffffffffffffffffffffffffffffffffffff1614610a9a576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610a9190611a94565b60405180910390fd5b60005b8251811015610ba1576000838281518110610abb57610aba611bd8565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610afb91906114a2565b602060405180830381865afa158015610b18573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b3c9190611960565b90506000811115610b8d57610b8c8382868581518110610b5f57610b5e611bd8565b5b602002602001015173ffffffffffffffffffffffffffffffffffffffff166111429092919063ffffffff16565b5b508080610b9990611c36565b915050610a9d565b505050565b600160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610c36576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c2d90611ccb565b60405180910390fd5b60008111610c79576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610c709061192b565b60405180910390fd5b610ca482828573ffffffffffffffffffffffffffffffffffffffff166111429092919063ffffffff16565b505050565b6002602052816000526040600020602052806000526040600020600091509150508060000154908060010154905082565b6000808211610d1e576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610d159061192b565b60405180910390fd5b60008373ffffffffffffffffffffffffffffffffffffffff166370a08231306040518263ffffffff1660e01b8152600401610d5991906114a2565b602060405180830381865afa158015610d76573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d9a9190611960565b90506000600360008673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905060008173ffffffffffffffffffffffffffffffffffffffff166318160ddd6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610e4d573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610e719190611960565b90506000831480610e825750600081145b15610e9257849350505050610ec3565b6000610eb984610eab84896111c890919063ffffffff16565b6111de90919063ffffffff16565b9050809450505050505b92915050565b610ed1611060565b73ffffffffffffffffffffffffffffffffffffffff16610eef61061f565b73ffffffffffffffffffffffffffffffffffffffff1614610f45576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610f3c90611a94565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff161415610fb5576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610fac90611d5d565b60405180910390fd5b610fbe81611068565b50565b611044846323b872dd60e01b858585604051602401610fe293929190611d7d565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506111f4565b50505050565b600081836110589190611db4565b905092915050565b600033905090565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff169050816000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508173ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e060405160405180910390a35050565b6000818361113a9190611e0a565b905092915050565b6111c38363a9059cbb60e01b8484604051602401611161929190611a1f565b604051602081830303815290604052907bffffffffffffffffffffffffffffffffffffffffffffffffffffffff19166020820180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff83818316178352505050506111f4565b505050565b600081836111d69190611e3e565b905092915050565b600081836111ec9190611ec7565b905092915050565b6000611256826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166112bb9092919063ffffffff16565b90506000815111156112b657808060200190518101906112769190611f30565b6112b5576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016112ac90611fcf565b60405180910390fd5b5b505050565b60606112ca84846000856112d3565b90509392505050565b606082471015611318576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161130f90612061565b60405180910390fd5b611321856113e7565b611360576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401611357906120cd565b60405180910390fd5b6000808673ffffffffffffffffffffffffffffffffffffffff1685876040516113899190612167565b60006040518083038185875af1925050503d80600081146113c6576040519150601f19603f3d011682016040523d82523d6000602084013e6113cb565b606091505b50915091506113db8282866113fa565b92505050949350505050565b600080823b905060008111915050919050565b6060831561140a5782905061145a565b60008351111561141d5782518084602001fd5b816040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161145191906121c2565b60405180910390fd5b9392505050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b600061148c82611461565b9050919050565b61149c81611481565b82525050565b60006020820190506114b76000830184611493565b92915050565b6000604051905090565b600080fd5b600080fd5b6114da81611481565b81146114e557600080fd5b50565b6000813590506114f7816114d1565b92915050565b600060208284031215611513576115126114c7565b5b6000611521848285016114e8565b91505092915050565b6000819050919050565b600061154f61154a61154584611461565b61152a565b611461565b9050919050565b600061156182611534565b9050919050565b600061157382611556565b9050919050565b61158381611568565b82525050565b600060208201905061159e600083018461157a565b92915050565b6000819050919050565b6115b7816115a4565b81146115c257600080fd5b50565b6000813590506115d4816115ae565b92915050565b600080604083850312156115f1576115f06114c7565b5b60006115ff858286016114e8565b9250506020611610858286016115c5565b9150509250929050565b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6116688261161f565b810181811067ffffffffffffffff8211171561168757611686611630565b5b80604052505050565b600061169a6114bd565b90506116a6828261165f565b919050565b600067ffffffffffffffff8211156116c6576116c5611630565b5b602082029050602081019050919050565b600080fd5b60006116ef6116ea846116ab565b611690565b90508083825260208201905060208402830185811115611712576117116116d7565b5b835b8181101561173b578061172788826114e8565b845260208401935050602081019050611714565b5050509392505050565b600082601f83011261175a5761175961161a565b5b813561176a8482602086016116dc565b91505092915050565b6000806040838503121561178a576117896114c7565b5b600083013567ffffffffffffffff8111156117a8576117a76114cc565b5b6117b485828601611745565b92505060206117c5858286016114e8565b9150509250929050565b6000806000606084860312156117e8576117e76114c7565b5b60006117f6868287016114e8565b9350506020611807868287016114e8565b9250506040611818868287016115c5565b9150509250925092565b60008060408385031215611839576118386114c7565b5b6000611847858286016114e8565b9250506020611858858286016114e8565b9150509250929050565b61186b816115a4565b82525050565b6000819050919050565b61188481611871565b82525050565b600060408201905061189f6000830185611862565b6118ac602083018461187b565b9392505050565b60006020820190506118c86000830184611862565b92915050565b600082825260208201905092915050565b7f616d6f756e74206d7573742067726561746572207468616e2030000000000000600082015250565b6000611915601a836118ce565b9150611920826118df565b602082019050919050565b6000602082019050818103600083015261194481611908565b9050919050565b60008151905061195a816115ae565b92915050565b600060208284031215611976576119756114c7565b5b60006119848482850161194b565b91505092915050565b7f7573657227732062616c616e6365206973206c657373207468616e207265717560008201527f6972656420616d6f756e74000000000000000000000000000000000000000000602082015250565b60006119e9602b836118ce565b91506119f48261198d565b604082019050919050565b60006020820190508181036000830152611a18816119dc565b9050919050565b6000604082019050611a346000830185611493565b611a416020830184611862565b9392505050565b7f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572600082015250565b6000611a7e6020836118ce565b9150611a8982611a48565b602082019050919050565b60006020820190508181036000830152611aad81611a71565b9050919050565b7f6761746577617927732062616c616e6365206973206c657373207468616e207260008201527f6571756972656420616d6f756e74000000000000000000000000000000000000602082015250565b6000611b10602e836118ce565b9150611b1b82611ab4565b604082019050919050565b60006020820190508181036000830152611b3f81611b03565b9050919050565b7f6465706f736974656420746f6b656e20616d6f756e74206973206c657373207460008201527f68616e20776974686472617720746f6b656e20616d6f756e7400000000000000602082015250565b6000611ba26039836118ce565b9150611bad82611b46565b604082019050919050565b60006020820190508181036000830152611bd181611b95565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000611c41826115a4565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff821415611c7457611c73611c07565b5b600182019050919050565b7f4f6e6c7920676174657761790000000000000000000000000000000000000000600082015250565b6000611cb5600c836118ce565b9150611cc082611c7f565b602082019050919050565b60006020820190508181036000830152611ce481611ca8565b9050919050565b7f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160008201527f6464726573730000000000000000000000000000000000000000000000000000602082015250565b6000611d476026836118ce565b9150611d5282611ceb565b604082019050919050565b60006020820190508181036000830152611d7681611d3a565b9050919050565b6000606082019050611d926000830186611493565b611d9f6020830185611493565b611dac6040830184611862565b949350505050565b6000611dbf826115a4565b9150611dca836115a4565b9250827fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff03821115611dff57611dfe611c07565b5b828201905092915050565b6000611e15826115a4565b9150611e20836115a4565b925082821015611e3357611e32611c07565b5b828203905092915050565b6000611e49826115a4565b9150611e54836115a4565b9250817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0483118215151615611e8d57611e8c611c07565b5b828202905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000611ed2826115a4565b9150611edd836115a4565b925082611eed57611eec611e98565b5b828204905092915050565b60008115159050919050565b611f0d81611ef8565b8114611f1857600080fd5b50565b600081519050611f2a81611f04565b92915050565b600060208284031215611f4657611f456114c7565b5b6000611f5484828501611f1b565b91505092915050565b7f5361666545524332303a204552433230206f7065726174696f6e20646964206e60008201527f6f74207375636365656400000000000000000000000000000000000000000000602082015250565b6000611fb9602a836118ce565b9150611fc482611f5d565b604082019050919050565b60006020820190508181036000830152611fe881611fac565b9050919050565b7f416464726573733a20696e73756666696369656e742062616c616e636520666f60008201527f722063616c6c0000000000000000000000000000000000000000000000000000602082015250565b600061204b6026836118ce565b915061205682611fef565b604082019050919050565b6000602082019050818103600083015261207a8161203e565b9050919050565b7f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000600082015250565b60006120b7601d836118ce565b91506120c282612081565b602082019050919050565b600060208201905081810360008301526120e6816120aa565b9050919050565b600081519050919050565b600081905092915050565b60005b83811015612121578082015181840152602081019050612106565b83811115612130576000848401525b50505050565b6000612141826120ed565b61214b81856120f8565b935061215b818560208601612103565b80840191505092915050565b60006121738284612136565b915081905092915050565b600081519050919050565b60006121948261217e565b61219e81856118ce565b93506121ae818560208601612103565b6121b78161161f565b840191505092915050565b600060208201905081810360008301526121dc8184612189565b90509291505056fea26469706673582212201441408302a78ab9601454e8c5139576c52356835dfc42d11fa8c515d9c088a464736f6c634300080a0033",
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

// EmergencyWithdrawFunds is a paid mutator transaction binding the contract method 0xaea3da6c.
//
// Solidity: function emergencyWithdrawFunds(address[] _tokens, address _newOwner) returns()
func (_Liquidity *LiquidityTransactor) EmergencyWithdrawFunds(opts *bind.TransactOpts, _tokens []common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.contract.Transact(opts, "emergencyWithdrawFunds", _tokens, _newOwner)
}

// EmergencyWithdrawFunds is a paid mutator transaction binding the contract method 0xaea3da6c.
//
// Solidity: function emergencyWithdrawFunds(address[] _tokens, address _newOwner) returns()
func (_Liquidity *LiquiditySession) EmergencyWithdrawFunds(_tokens []common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.EmergencyWithdrawFunds(&_Liquidity.TransactOpts, _tokens, _newOwner)
}

// EmergencyWithdrawFunds is a paid mutator transaction binding the contract method 0xaea3da6c.
//
// Solidity: function emergencyWithdrawFunds(address[] _tokens, address _newOwner) returns()
func (_Liquidity *LiquidityTransactorSession) EmergencyWithdrawFunds(_tokens []common.Address, _newOwner common.Address) (*types.Transaction, error) {
	return _Liquidity.Contract.EmergencyWithdrawFunds(&_Liquidity.TransactOpts, _tokens, _newOwner)
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
