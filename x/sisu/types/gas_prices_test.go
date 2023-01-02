package types

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestGasPrice_Validation(t *testing.T) {
	signerAcc, err := sdk.AccAddressFromBech32("cosmos1zf2ssujzp6y577gzwn457tnxy7yj44yq37t05z")
	signer := signerAcc.String()
	require.NoError(t, err)

	chain := "ganache1"

	// Valid message
	var gasPrice *GasPriceMsg
	gasPrice = NewGasPriceMsg(signer, []string{chain}, []int64{10}, nil, nil)
	require.Nil(t, gasPrice.ValidateBasic())

	gasPrice = NewGasPriceMsg(signer, []string{chain}, nil, []int64{20}, []int64{2})
	require.Nil(t, gasPrice.ValidateBasic())

	// Invalid message
	// nil chain array
	gasPrice = NewGasPriceMsg(signer, nil, []int64{10}, nil, nil)
	require.NotNil(t, gasPrice.ValidateBasic())

	// All gas prices, base fee & tips arrays are nil
	gasPrice = NewGasPriceMsg(signer, []string{chain}, nil, nil, nil)
	require.NotNil(t, gasPrice.ValidateBasic())

	// Chains & gas prices do not have the same length
	gasPrice = NewGasPriceMsg(signer, []string{chain}, []int64{10, 20}, nil, nil)
	require.NotNil(t, gasPrice.ValidateBasic())

	// Chains & base fee & tips do not have the same length
	gasPrice = NewGasPriceMsg(signer, []string{chain}, nil, []int64{10, 20}, []int64{1, 2})
	require.NotNil(t, gasPrice.ValidateBasic())
}
