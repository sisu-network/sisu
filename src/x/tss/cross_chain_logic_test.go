package tss

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

func TestSignTx(t *testing.T) {
	o := NewCrossChainLogic()

	rawTxs := o.PrepareEthContractDeployment("eth", 0)
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")

	if err != nil {
		t.Fail()
		return
	}

	// We can at least sign rawTx
	_, err = types.SignTx(rawTxs[0], types.NewEIP155Signer(big.NewInt(1)), privateKey)
	if err != nil {
		t.Fail()
		return
	}
}
