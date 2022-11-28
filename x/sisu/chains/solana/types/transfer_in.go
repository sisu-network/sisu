package types

import (
	"github.com/gagliardetto/solana-go"
	solanago "github.com/gagliardetto/solana-go"
	"github.com/near/borsh-go"
)

type TransferInDataInner struct {
	Nonce   uint64
	Amounts []uint64
}

type TransferInData struct {
	Instruction byte
	InnerData   TransferInDataInner
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
	amounts []uint64,
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

	// Add all bridge & receiver ata
	for i, token := range tokens {
		// BridgeATa
		bridgeAta, err := GetAtaPubkey(bridgePda, token)
		if err != nil {
			return nil, err
		}
		accounts = append(
			accounts,
			solanago.NewAccountMeta(bridgeAta, true, false),
		)

		// Receivert Ata
		receiverAta, err := solanago.PublicKeyFromBase58(receiverAtas[i])
		if err != nil {
			return nil, err
		}
		accounts = append(
			accounts,
			solanago.NewAccountMeta(receiverAta, true, false),
		)
	}

	return &TransferInIx{
		bridgeProgramdId: solanago.MustPublicKeyFromBase58(bridgeProgramdId),
		accounts:         accounts,
		data: TransferInData{
			Instruction: TranserIn,
			InnerData: TransferInDataInner{
				Nonce:   nonce,
				Amounts: amounts,
			},
		},
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
