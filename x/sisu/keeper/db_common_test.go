package keeper

import (
	"fmt"
	"math/big"
	"sort"
	"testing"

	memstore "github.com/cosmos/cosmos-sdk/store/mem"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/sisu/x/sisu/types"
	"github.com/stretchr/testify/require"
)

///// Keygen

func Test_saveKeygen(t *testing.T) {
	store := memstore.NewStore()

	keyType := libchain.KEY_TYPE_ECDSA
	index := 0
	address := "1234"

	pubkey := []byte("Pubkey")
	keygen := &types.Keygen{
		KeyType:     keyType,
		Index:       int32(index),
		PubKeyBytes: pubkey,
		Address:     address,
	}

	saveKeygen(store, keygen)

	// Check keygen existed
	require.Equal(t, true, isKeygenExisted(store, keyType, index))
	require.Equal(t, false, isKeygenExisted(store, keyType, index+1))

	// Check address
	require.Equal(t, true, isKeygenAddress(store, keyType, address))
	require.Equal(t, false, isKeygenAddress(store, keyType, fmt.Sprintf("another %s", address)))
}

func Test_getAllKeygenPubkeys(t *testing.T) {
	store := memstore.NewStore()

	pubkey := []byte("Pubkey")
	keygen := &types.Keygen{
		KeyType:     libchain.KEY_TYPE_ECDSA,
		Index:       0,
		PubKeyBytes: pubkey,
	}

	saveKeygen(store, keygen)

	allPubkeys := getAllKeygenPubkeys(store)
	require.Equal(t, len(allPubkeys), 1, "allPubkeys length does not match")
	require.Equal(t, allPubkeys[libchain.KEY_TYPE_ECDSA], pubkey, "Pubkey does not match")
}

///// Keygen Result

func Test_saveKeygenResult(t *testing.T) {
	store := memstore.NewStore()

	node := "node0"
	node1 := "node1"
	keyType := libchain.KEY_TYPE_ECDSA
	index := int32(0)

	signer := &types.KeygenResultWithSigner{
		Signer: node,
		Keygen: &types.Keygen{
			KeyType: keyType,
			Index:   index,
		},
		Data: &types.KeygenResult{
			From:   node,
			Result: types.KeygenResult_SUCCESS,
		},
	}

	saveKeygenResult(store, signer)

	signer1 := &types.KeygenResultWithSigner{
		Signer: node,
		Keygen: &types.Keygen{
			KeyType: keyType,
			Index:   index,
		},
		Data: &types.KeygenResult{
			From:   node1,
			Result: types.KeygenResult_SUCCESS,
		},
	}
	saveKeygenResult(store, signer1)

	results := getAllKeygenResult(store, keyType, index)
	require.Equal(t, 2, len(results))
}

///// Token Prices
func TestSaveTokenPrices(t *testing.T) {
	store := memstore.NewStore()

	token := "ETH"

	signer1 := "signer1"
	msg1Signer1 := &types.UpdateTokenPrice{
		Signer: signer1,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: big.NewInt(5_000_000_000).String(),
			},
		},
	}
	msg2Signer1 := &types.UpdateTokenPrice{
		Signer: signer1,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: big.NewInt(6_000_000_000).String(),
			},
		},
	}

	signer2 := "signer2"
	msg1Signer2 := &types.UpdateTokenPrice{
		Signer: signer2,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: big.NewInt(10_000_000_000).String(),
			},
		},
	}
	msg2Signer2 := &types.UpdateTokenPrice{
		Signer: signer2,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: big.NewInt(11_000_000_000).String(),
			},
		}}

	setTokenPrices(store, 1, msg1Signer1)
	setTokenPrices(store, 1, msg2Signer1)
	setTokenPrices(store, 1, msg1Signer2)
	setTokenPrices(store, 1, msg2Signer2)

	allPrices := getAllTokenPrices(store)

	require.Equal(t, 2, len(allPrices))

	allSigners := make([]string, 0)
	for savedSigner, record := range allPrices {
		allSigners = append(allSigners, savedSigner)
		require.Equal(t, 1, len(record.Records))
	}

	sort.Strings(allSigners)
	require.Equal(t, []string{signer1, signer2}, allSigners)

	record := allPrices[signer1]
	require.Equal(t, "6000000000", record.Records[0].Price)
	record = allPrices[signer2]
	require.Equal(t, "11000000000", record.Records[0].Price)
}

///// Node
func Test_SaveNode(t *testing.T) {
	store := memstore.NewStore()

	node1 := &types.Node{
		ConsensusKey: &types.Pubkey{
			Type:  "ed",
			Bytes: []byte("pubkey1"),
		},
		AccAddress:  "addr1",
		IsValidator: true,
	}
	node2 := &types.Node{
		ConsensusKey: &types.Pubkey{
			Type:  "ed",
			Bytes: []byte("pubkey2"),
		},
		AccAddress:  "addr2",
		IsValidator: true,
	}

	saveNode(store, node1)
	saveNode(store, node2)

	vals := loadValidators(store)

	require.Equal(t, 2, len(vals), "there should be 2 validators")
	require.Equal(t, vals[0].AccAddress, "addr1")
	require.Equal(t, vals[1].AccAddress, "addr2")
}

func TestDb_KeygenPubkey(t *testing.T) {
	store := memstore.NewStore()

	saveKeygen(store, &types.Keygen{
		KeyType:     libchain.KEY_TYPE_ECDSA,
		Index:       0,
		PubKeyBytes: []byte("ec_key"),
	})
	saveKeygen(store, &types.Keygen{
		KeyType:     libchain.KEY_TYPE_EDDSA,
		Index:       0,
		PubKeyBytes: []byte("ed_key"),
	})

	// Get Ecdsa pubkey
	pubkeyBytes := getKeygenPubkey(store, libchain.KEY_TYPE_ECDSA)
	require.Equal(t, []byte("ec_key"), pubkeyBytes)

	// Get Eddsa pubkey
	pubkeyBytes = getKeygenPubkey(store, libchain.KEY_TYPE_EDDSA)
	require.Equal(t, []byte("ed_key"), pubkeyBytes)

	// Test that if there are multiple keygen, the getKeygenPubkey returns pubkey of the latest.
	saveKeygen(store, &types.Keygen{
		KeyType:     libchain.KEY_TYPE_ECDSA,
		Index:       1,
		PubKeyBytes: []byte("ec_key_1"),
	})
	pubkeyBytes = getKeygenPubkey(store, libchain.KEY_TYPE_ECDSA)
	require.Equal(t, []byte("ec_key_1"), pubkeyBytes)
}
