package types

import (
	"fmt"

	solanago "github.com/gagliardetto/solana-go"
)

// GetAtaPubkey calculates and returns the associated token address for a wallet with a token.
func GetAtaPubkey(wallet, token string) (solanago.PublicKey, error) {
	walletPubkey, err := solanago.PublicKeyFromBase58(wallet)
	if err != nil {
		return solanago.PublicKey{}, fmt.Errorf("Address `%s` is not a valid solana address", wallet)
	}
	tokenPubkey, err := solanago.PublicKeyFromBase58(token)
	if err != nil {
		return solanago.PublicKey{}, fmt.Errorf("Address `%s` is not a valid solana address", token)
	}

	ata, _, err := solanago.FindAssociatedTokenAddress(walletPubkey, tokenPubkey)
	if err != nil {
		return solanago.PublicKey{}, err
	}

	return ata, nil
}
