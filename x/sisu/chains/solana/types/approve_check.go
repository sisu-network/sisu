package types

import (
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type ApproveCheckedData struct {
	Instruction byte
	Amount      uint64
	Decimals    byte
}

type ApproveCheckedIx struct {
	accounts []*solanago.AccountMeta
	data     ApproveCheckedData
}

func NewApproveCheckedIx(
	owner solanago.PublicKey,
	ownerAta solanago.PublicKey,
	token solanago.PublicKey,
	delegate solanago.PublicKey,
	amount uint64,
	decimals byte,
) *ApproveCheckedIx {
	accounts := []*solanago.AccountMeta{
		solanago.NewAccountMeta(ownerAta, true, false),
		solanago.NewAccountMeta(token, false, false),
		solanago.NewAccountMeta(delegate, false, false),
		solanago.NewAccountMeta(owner, false, false),
	}

	return &ApproveCheckedIx{
		accounts: accounts,
		data: ApproveCheckedData{
			Instruction: 13, // https://github.com/solana-labs/solana-program-library/blob/374f5283ec90457113eb7a6eb868130cc544f990/token/js/src/instructions/types.ts#L16
			Amount:      amount,
			Decimals:    decimals,
		},
	}
}

// ProgramID is the programID the instruction acts on.
func (ix *ApproveCheckedIx) ProgramID() solanago.PublicKey {
	return solanago.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
}

// Accounts returns the list of accounts the instructions requires
func (ix *ApproveCheckedIx) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *ApproveCheckedIx) Data() ([]byte, error) {
	return borsh.Serialize(ix.data)
}
