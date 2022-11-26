package types

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type TransferInData struct {
	Nonce   uint64
	Amounts []uint64
}

type TransferInIx struct {
	bridgeProgramdId solana.PublicKey
	accounts         []*solanago.AccountMeta
	data             TransferInData
}

func NewTransferInIx(
	bridgeProgramdId string,
	mpcAddress string,
	nonce uint64,
	tokenProgramId string,
	bridgePda string,
	tokens []string,
	receiverAtas []string,
	amounts []*big.Int,
) (*TransferInIx, error) {
	// Verify that all strings are valid.
	accountStrs := []string{bridgeProgramdId, mpcAddress, tokenProgramId, bridgePda}
	accountStrs = append(accountStrs, tokens...)
	accountStrs = append(accountStrs, receiverAtas...)
	for _, accountStr := range accountStrs {
		_, err := solanago.PublicKeyFromBase58(accountStr)
		if err != nil {
			return nil, err
		}
	}

	accounts := []*solanago.AccountMeta{
		solanago.NewAccountMeta(solanago.MustPublicKeyFromBase58(mpcAddress), false, true),
		solanago.NewAccountMeta(solanago.MustPublicKeyFromBase58(tokenProgramId), false, false),
		solanago.NewAccountMeta(solanago.MustPublicKeyFromBase58(bridgePda), false, false),
	}

	// Add all bridge ata
	for _, token := range tokens {
		bridgeAta, _, err := solanago.FindAssociatedTokenAddress(
			solanago.MustPublicKeyFromBase58(bridgePda),
			solanago.MustPublicKeyFromBase58(token),
		)
		if err != nil {
			return nil, err
		}

		accounts = append(
			accounts,
			solanago.NewAccountMeta(bridgeAta, false, true),
		)
	}

	return &TransferInIx{
		bridgeProgramdId: solanago.MustPublicKeyFromBase58(bridgeProgramdId),
		accounts:         accounts,
	}, nil
}

func verifySolanaAddress(addr string) bool {
	_, err := solanago.PublicKeyFromBase58(addr)
	return err == nil
}

// ProgramID is the programID the instruction acts on.
func (ix *TransferInIx) ProgramID() solanago.PublicKey {
	return ix.bridgeProgramdId
}

// Accounts returns the list of accounts the instructions requires
func (ix *TransferInIx) Accounts() []*solanago.AccountMeta {
	return ix.accounts
}

func (ix *TransferInIx) Data() ([]byte, error) {
	return borsh.Serialize(ix.data)
}
