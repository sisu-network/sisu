package utils

import (
	"encoding/hex"
	"testing"

	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

func TestGetSortedValidators(t *testing.T) {
	// To generate a key:
	hexes := []string{
		"02a4ffd4fe0a384e01568f5e8ad4564298971c72409b5d211809bebe4bc83c2cdc",
		"03d017b54de73b3662bf3e442c7552a4922a11ac96aabcadbf1fa238da66beaf51",
		"03edc9389e4daab2a8b22b7f481699cc24b6be9f0edde69663a335ec47e775aca1",
		"02e8ee66ea937576fbc58621ec1af494fa630c8fb05f4656e7696a465a03a06db6",
		"0268521cc6ee2d6f79656b88d9d97d2bb82e6342c5b922fb7694e54c6367d7a1e0",
	}

	nodes := []*types.Node{}
	for _, h := range hexes {
		bz, err := hex.DecodeString(h)
		require.Nil(t, err)

		nodes = append(nodes, &types.Node{
			ValPubkey: &types.ValPubkey{
				Type:  "secp256k1",
				Bytes: bz,
			},
		})
	}

	// Msg1
	sorted := GetSortedValidators("msg1", nodes)
	sortedPub := make([]string, 0, len(sorted))
	for _, s := range sorted {
		sortedPub = append(sortedPub, hex.EncodeToString(s.ValPubkey.Bytes))
	}
	require.Equal(t, []string{
		"0268521cc6ee2d6f79656b88d9d97d2bb82e6342c5b922fb7694e54c6367d7a1e0",
		"03edc9389e4daab2a8b22b7f481699cc24b6be9f0edde69663a335ec47e775aca1",
		"02e8ee66ea937576fbc58621ec1af494fa630c8fb05f4656e7696a465a03a06db6",
		"03d017b54de73b3662bf3e442c7552a4922a11ac96aabcadbf1fa238da66beaf51",
		"02a4ffd4fe0a384e01568f5e8ad4564298971c72409b5d211809bebe4bc83c2cdc",
	}, sortedPub)

	// Msg2
	sorted = GetSortedValidators("msg2", nodes)
	sortedPub = make([]string, 0, len(sorted))
	for _, s := range sorted {
		sortedPub = append(sortedPub, hex.EncodeToString(s.ValPubkey.Bytes))
	}
	require.Equal(t, []string{
		"02e8ee66ea937576fbc58621ec1af494fa630c8fb05f4656e7696a465a03a06db6",
		"02a4ffd4fe0a384e01568f5e8ad4564298971c72409b5d211809bebe4bc83c2cdc",
		"0268521cc6ee2d6f79656b88d9d97d2bb82e6342c5b922fb7694e54c6367d7a1e0",
		"03edc9389e4daab2a8b22b7f481699cc24b6be9f0edde69663a335ec47e775aca1",
		"03d017b54de73b3662bf3e442c7552a4922a11ac96aabcadbf1fa238da66beaf51",
	}, sortedPub)
}
