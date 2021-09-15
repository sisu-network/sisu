package tss

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/contracts/eth/dummy"
)

type EthDeployment struct {
}

func NewEthDeployment() *EthDeployment {
	return &EthDeployment{}
}

func (ed *EthDeployment) PrepareEthContractDeployment(chain string, nonceIndex int64) *eTypes.Transaction {
	// Create Tx for dummy contract
	byteCode := common.FromHex(dummy.DummyBin)
	var nonce uint64
	nonce = 0

	rawTx := eTypes.NewContractCreation(nonce, big.NewInt(nonceIndex), ed.getGasLimit(chain), ed.getGasPrice(chain), byteCode)
	return rawTx
}

func (ed *EthDeployment) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(5000000)
}

func (ed *EthDeployment) getGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	return big.NewInt(50)
}
