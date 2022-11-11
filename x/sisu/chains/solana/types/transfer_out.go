package types

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type InstructionType byte

const (
	Initialize    InstructionType = 0
	TransferOut                   = 1
	TranserIn                     = 2
	AddSpender                    = 3
	RemoveSpender                 = 4
	ChangeAdmin                   = 5
)

type TransferOutData struct {
	Instruction  byte
	Amount       uint64
	TokenAddress string
	ChainId      uint64
	Recipient    string
}

func NewTransferOutData(
	amount uint64,
	tokenAddress string,
	chainId uint64,
	recipient string,
) TransferOutData {
	return TransferOutData{
		Instruction:  TransferOut,
		Amount:       amount,
		TokenAddress: tokenAddress,
		ChainId:      chainId,
		Recipient:    recipient,
	}
}

func (ix *TransferOutData) Serialize() ([]byte, error) {
	return borsh.Serialize(*ix)
}

func (ix *TransferOutData) Deserialize(bytesArr []byte) error {
	if len(bytesArr) == 0 {
		return fmt.Errorf("Byte array is nil")
	}

	return borsh.Deserialize(ix, bytesArr)
}

type TransferOutInstruction struct {
	bridgeProgramdId solana.PublicKey
	accounts         []*solanago.AccountMeta
	data             TransferOutData
}

func NewTransferOutInstruction(
	programId solana.PublicKey,
	owner solanago.PublicKey,
	ownerAta solanago.PublicKey,
	bridgeAta solanago.PublicKey,
	bridgePda solanago.PublicKey,
	data TransferOutData) *TransferOutInstruction {
	tokenProgramId := solana.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")

	accounts := []*solanago.AccountMeta{
		solanago.NewAccountMeta(owner, false, true),
		solanago.NewAccountMeta(tokenProgramId, false, false),
		solanago.NewAccountMeta(ownerAta, true, false),
		solanago.NewAccountMeta(bridgeAta, true, false),
		solanago.NewAccountMeta(bridgePda, true, false),
	}

	return &TransferOutInstruction{
		bridgeProgramdId: programId,
		accounts:         accounts,
		data:             data,
	}
}

// ProgramID is the programID the instruction acts on.
func (ix *TransferOutInstruction) ProgramID() solanago.PublicKey {
	// Associated program id. This is different from the token program id.
	return ix.bridgeProgramdId
}

// Accounts returns the list of accounts the instructions requires
func (ix *TransferOutInstruction) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *TransferOutInstruction) Data() ([]byte, error) {
	return borsh.Serialize(ix.data)
}
