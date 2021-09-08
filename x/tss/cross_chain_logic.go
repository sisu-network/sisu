package tss

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/sisu-network/sisu/contracts/eth/dummy"
)

// This is the struct that prepares output transaction that will be signed by TSS before delivering
// to target chain.
type CrossChainLogic struct {
}

func NewCrossChainLogic() *CrossChainLogic {
	return &CrossChainLogic{}
}

func (logic *CrossChainLogic) PrepareEthContractDeployment(chain string, nonceIndex int64) *eTypes.Transaction {
	// Create Tx for dummy contract
	byteCode := common.FromHex(dummy.DummyBin)
	var nonce uint64
	nonce = 0

	rawTx := eTypes.NewContractCreation(nonce, big.NewInt(nonceIndex), logic.getGasLimit(chain), logic.getGasPrice(chain), byteCode)
	return rawTx
}

func (logic *CrossChainLogic) getGasLimit(chain string) uint64 {
	// TODO: Make this dependent on different chains.
	return uint64(5000000)
}

func (logic *CrossChainLogic) getGasPrice(chain string) *big.Int {
	// TODO: Make this dependent on different chains.
	return big.NewInt(50)
}
