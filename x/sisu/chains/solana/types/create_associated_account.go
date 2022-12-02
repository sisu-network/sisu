package types

import (
	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
)

type CreateAssociatedAccountIx struct {
	accounts []*solana.AccountMeta
}

func NewCreateAssociatedAccountIx(feePayer, owner, associatedToken, tokenMint solanago.PublicKey) *CreateAssociatedAccountIx {
	// This is the code for accounts meta in the official JS code.
	// const keys = [
	// 	{ pubkey: payer, isSigner: true, isWritable: true },
	// 	{ pubkey: associatedToken, isSigner: false, isWritable: true },
	// 	{ pubkey: owner, isSigner: false, isWritable: false },
	// 	{ pubkey: mint, isSigner: false, isWritable: false },
	// 	{ pubkey: SystemProgram.programId, isSigner: false, isWritable: false },
	// 	{ pubkey: programId, isSigner: false, isWritable: false },
	// ];

	// owner is also the payer.
	accounts := []*solana.AccountMeta{
		solanago.NewAccountMeta(feePayer, false, true),
		solanago.NewAccountMeta(associatedToken, true, false),
		solanago.NewAccountMeta(owner, false, false),
		solanago.NewAccountMeta(tokenMint, false, false),
		solanago.NewAccountMeta(solanago.SystemProgramID, false, false),
		solanago.NewAccountMeta(solanago.TokenProgramID, false, false), // token program id
	}

	return &CreateAssociatedAccountIx{
		accounts: accounts,
	}
}

// ProgramID is the programID the instruction acts on.
func (ix *CreateAssociatedAccountIx) ProgramID() solanago.PublicKey {
	// Associated program id. This is different from the token program id.
	return solanago.MustPublicKeyFromBase58("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
}

// Accounts returns the list of accounts the instructions requires
func (ix *CreateAssociatedAccountIx) Accounts() []*solana.AccountMeta {
	return ix.accounts
}

func (ix *CreateAssociatedAccountIx) Data() ([]byte, error) {
	return []byte{}, nil
}
