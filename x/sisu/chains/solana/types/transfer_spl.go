package types

import (
	"math/big"

	solanago "github.com/gagliardetto/solana-go"

	"github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

// TODO: We cannot use amount as big.Int as it will serialize to more than 8 bytes. The SPL program
// expects to have unsigned 8 bytes for Amount field and uint64 should be enough.
type TransferSplTokenData struct {
	Instruction byte
	Amount      uint64
	Decimals    byte
}

type TransferSplTokenIx struct {
	Instruction byte
	data        TransferSplTokenData

	accounts []*solana.AccountMeta
}

// This is the JS code for constructing instruction data.
// TokenInstruction.TransferChecked = 12
//
// const data = Buffer.alloc(transferCheckedInstructionData.span);
// transferCheckedInstructionData.encode(
// 		{
// 				instruction: TokenInstruction.TransferChecked,
// 				amount: BigInt(amount),
// 				decimals,
// 		},
// 		data
// );
func NewTransferTokenIx(sourceAta, token, dstAta, feePayerPubkey solanago.PublicKey, amount *big.Int, decimals byte) *TransferSplTokenIx {
	// This is the key source code in JS.
	// 	const keys = addSigners(
	// 		[
	// 				{ pubkey: source, isSigner: false, isWritable: true },
	// 				{ pubkey: mint, isSigner: false, isWritable: false },
	// 				{ pubkey: destination, isSigner: false, isWritable: true },
	// 		],
	// 		owner,
	// 		multiSigners
	// );

	accounts := []*solana.AccountMeta{
		solana.NewAccountMeta(sourceAta, true, false),
		solana.NewAccountMeta(token, false, false),
		solana.NewAccountMeta(dstAta, true, false),
		solana.NewAccountMeta(feePayerPubkey, false, true),
	}

	data := &TransferSplTokenData{
		Instruction: 12,
		Amount:      amount.Uint64(),
		Decimals:    decimals,
	}

	return &TransferSplTokenIx{
		accounts: accounts,
		data:     *data,
	}
}

// ProgramID is the programID the instruction acts on
func (ix *TransferSplTokenIx) ProgramID() solanago.PublicKey {
	return solanago.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
}

// Accounts returns the list of accounts the instructions requires
func (ix *TransferSplTokenIx) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *TransferSplTokenIx) Data() ([]byte, error) {
	bz, err := borsh.Serialize(ix.data)

	return bz, err
}
