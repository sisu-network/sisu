package keeper

import (
	"fmt"
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

	results := getAllKeygenResult(store, keyType, index)
	require.Equal(t, 1, len(results))
}

///// Token Prices
func TestSaveTokenPrices(t *testing.T) {
	store := memstore.NewStore()

	token := "ETH"

	signer1 := "signer1"
	msg1 := &types.UpdateTokenPrice{
		Signer: signer1,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 5_000_000_000,
			},
		},
	}

	signer2 := "signer2"
	msg2 := &types.UpdateTokenPrice{
		Signer: signer2,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 10_000_000_000,
			},
		},
	}

	setTokenPrices(store, 1, msg1)
	setTokenPrices(store, 1, msg2)

	allPrices := getAllTokenPrices(store)

	require.Equal(t, 2, len(allPrices))

	allSigners := make([]string, 0)
	for savedSigner, record := range allPrices {
		allSigners = append(allSigners, savedSigner)
		require.Equal(t, 1, len(record.Prices))
	}

	sort.Strings(allSigners)
	require.Equal(t, []string{signer1, signer2}, allSigners)

	record := allPrices[signer1]
	require.Equal(t, int64(5_000_000_000), record.Prices[token].Price)
	record = allPrices[signer2]
	require.Equal(t, int64(10_000_000_000), record.Prices[token].Price)
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
