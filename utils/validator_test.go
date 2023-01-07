package utils

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	tcrypto "github.com/tendermint/tendermint/crypto"
	tsecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestGetSortedValidators(t *testing.T) {
	// To generate a key:
	//  tsecp256k1.GenPrivKey().PubKey()
	hexes := []string{
		"02a4ffd4fe0a384e01568f5e8ad4564298971c72409b5d211809bebe4bc83c2cdc",
		"03d017b54de73b3662bf3e442c7552a4922a11ac96aabcadbf1fa238da66beaf51",
		"03edc9389e4daab2a8b22b7f481699cc24b6be9f0edde69663a335ec47e775aca1",
		"02e8ee66ea937576fbc58621ec1af494fa630c8fb05f4656e7696a465a03a06db6",
		"0268521cc6ee2d6f79656b88d9d97d2bb82e6342c5b922fb7694e54c6367d7a1e0",
	}

	keys := []tcrypto.PubKey{}
	for _, key := range keys {
		fmt.Printf("\"%s\",\n", hex.EncodeToString(key.Bytes()))
	}

	for _, h := range hexes {
		bz, err := hex.DecodeString(h)
		require.Nil(t, err)

		key := tsecp256k1.PrivKey(bz)
		keys = append(keys, key.PubKey())
	}

	// Msg1
	sorted := GetSortedValidators("msg1", keys)
	sortedPub := make([]string, 0, len(sorted))
	for _, s := range sorted {
		sortedPub = append(sortedPub, s.Address().String())
	}
	require.Equal(t, []string{
		"6F6BC678DEB7805FA6BD8B75BD1F82D0C7616E48",
		"F5BC6966FB35956014969369109F89B13403DDA8",
		"C6E8E6CD8211F2EF30E2E880FE6AE5AF0D036E65",
		"A63B5B638F81159C6EFD1484641DCAD3F119E59F",
		"F9BA47894233B03FFF3C72A6F835FA98412D007E",
	}, sortedPub)

	// Msg2
	sorted = GetSortedValidators("msg2", keys)
	sortedPub = make([]string, 0, len(sorted))
	for _, s := range sorted {
		sortedPub = append(sortedPub, s.Address().String())
	}
	require.Equal(t, []string{
		"6F6BC678DEB7805FA6BD8B75BD1F82D0C7616E48",
		"A63B5B638F81159C6EFD1484641DCAD3F119E59F",
		"F9BA47894233B03FFF3C72A6F835FA98412D007E",
		"C6E8E6CD8211F2EF30E2E880FE6AE5AF0D036E65",
		"F5BC6966FB35956014969369109F89B13403DDA8",
	}, sortedPub)
}
