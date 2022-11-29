package types

import (
	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type AddSpenderData struct {
	Instruction byte
	Spender     solanago.PublicKey
}

type AddSpenderIx struct {
	bridgeProgramdId solana.PublicKey
	accounts         []*solanago.AccountMeta
	data             AddSpenderData
}

func NewAddSpenderIx(
	bridgeProgramdId string,
	mpcAddress string,
	bridgePda string,
	spender string,
) (*AddSpenderIx, error) {
	if err := verifySolanaAddress([]string{bridgeProgramdId, mpcAddress, spender}); err != nil {
		return nil, err
	}

	return &AddSpenderIx{
		bridgeProgramdId: solanago.MustPublicKeyFromBase58(bridgeProgramdId),
		accounts: []*solanago.AccountMeta{
			solanago.NewAccountMeta(solanago.MustPublicKeyFromBase58(mpcAddress), false, true),
			solanago.NewAccountMeta(solanago.MustPublicKeyFromBase58(bridgePda), true, false),
		},
		data: AddSpenderData{
			Instruction: AddSpender,
			Spender:     solanago.MustPublicKeyFromBase58(spender),
		},
	}, nil
}

// ProgramID is the programID the instruction acts on.
func (ix *AddSpenderIx) ProgramID() solanago.PublicKey {
	return ix.bridgeProgramdId
}

// Accounts returns the list of accounts the instructions requires
func (ix *AddSpenderIx) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *AddSpenderIx) Data() ([]byte, error) {
	return borsh.Serialize(ix.data)
}
