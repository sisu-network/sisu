package types

import (
	"fmt"
	"math/big"

	"github.com/near/borsh-go"
)

type InstructionType byte

const (
	Initialize    InstructionType = 0
	TranserOut                    = 1
	TranserIn                     = 2
	AddSpender                    = 3
	RemoveSpender                 = 4
	ChangeAdmin                   = 5
)

type TransferOutData struct {
	Amount       big.Int
	TokenAddress string
	ChainId      uint64
	Recipient    string
}

type TransferOutInstruction struct {
	Instruction byte
	Data        TransferOutData
}

func NewTransferOutData(amount *big.Int) *TransferOutData {
	return &TransferOutData{}
}

func (ix *TransferOutInstruction) Serialize() ([]byte, error) {
	bz, err := borsh.Serialize(ix.Data)
	if err != nil {
		return nil, err
	}

	ret := []byte{ix.Instruction}
	ret = append(ret, bz...)

	return ret, nil
}

func (ix *TransferOutInstruction) Deserialize(bytesArr []byte) error {
	if len(bytesArr) == 0 {
		return fmt.Errorf("Byte array is nil")
	}

	borshBz := bytesArr[1:]
	Data := TransferOutData{}
	err := borsh.Deserialize(&Data, borshBz)
	if err != nil {
		return err
	}

	ix.Instruction = bytesArr[0]
	ix.Data = Data

	return nil
}
