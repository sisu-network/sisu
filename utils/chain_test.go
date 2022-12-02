package utils

import (
	"encoding/hex"
	"testing"

	"github.com/mr-tron/base58"
	"github.com/stretchr/testify/require"
)

func TestGetSolanaAddressFromPubkey(t *testing.T) {
	hexString := "ce8a3099749c74c3c086290362e0d513137cf36433afa43176d2900591ab82b4"
	bz, err := hex.DecodeString(hexString)
	require.Nil(t, err)

	base58String := base58.Encode(bz)
	require.Equal(t, "EuFCgxwUQMFoC8N1iFazZsn1C66wyDx4dTP4RSpMJscw", base58String)
}
