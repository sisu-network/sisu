package types

import (
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type MintTokenData struct {
	Instruction byte
	Amount      uint64
	Decimals    byte
}

type MintTokenIx struct {
	accounts []*solanago.AccountMeta
	data     MintTokenData
}

func NewMintTokenIx(tokenAddr, receiverAta, authorPubkey solanago.PublicKey, decimals byte, amount uint64) *MintTokenIx {
	return &MintTokenIx{
		accounts: []*solanago.AccountMeta{
			solanago.NewAccountMeta(tokenAddr, true, false),
			solanago.NewAccountMeta(receiverAta, true, false),
			solanago.NewAccountMeta(authorPubkey, false, true),
		},
		data: MintTokenData{
			Instruction: 14,
			Amount:      amount,
			Decimals:    decimals,
		},
	}
}

// ProgramID is the programID the instruction acts on
func (ix *MintTokenIx) ProgramID() solanago.PublicKey {
	return solanago.MustPublicKeyFromBase58("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
}

// Accounts returns the list of accounts the instructions requires
func (ix *MintTokenIx) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *MintTokenIx) Data() ([]byte, error) {
	bz, err := borsh.Serialize(ix.data)

	return bz, err
}
