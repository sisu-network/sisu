package keeper

import (
	"fmt"
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

///// Node
func Test_SaveNode(t *testing.T) {
	store := memstore.NewStore()

	node1 := &types.Node{
		ValPubkey: &types.ValPubkey{
			Type:  "ed25519",
			Bytes: []byte("pubkey1"),
		},
		AccAddress:  "addr1",
		IsValidator: true,
	}
	node2 := &types.Node{
		ValPubkey: &types.ValPubkey{
			Type:  "ed25519",
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

func TestDb_KeygenResult(t *testing.T) {
	store := memstore.NewStore()

	saveKeygenResult(store, &types.KeygenResultWithSigner{
		Signer: "signer1",
		Keygen: &types.Keygen{
			KeyType: "ecdsa",
			Index:   0,
		},
		Data: &types.KeygenResult{
			From: "signer1",
		},
	})
	saveKeygenResult(store, &types.KeygenResultWithSigner{
		Signer: "signer2",
		Keygen: &types.Keygen{
			KeyType: "ecdsa",
			Index:   0,
		},
		Data: &types.KeygenResult{
			From: "signer2",
		},
	})
	saveKeygenResult(store, &types.KeygenResultWithSigner{
		Signer: "signer3",
		Keygen: &types.Keygen{
			KeyType: "ecdsa",
			Index:   0,
		},
		Data: &types.KeygenResult{
			From: "signer3",
		},
	})

	results := getAllKeygenResult(store, "ecdsa", 0)
	require.Equal(t, 3, len(results))
}

func TestDb_getGateway(t *testing.T) {
	store := memstore.NewStore()

	setSisuAccount(store, "ganache1", "addr1")
	sisu := getSisuAccount(store, "ganache1")
	require.Equal(t, "addr1", sisu)

	setSisuAccount(store, "ganache2", "addr2")
	sisu = getSisuAccount(store, "ganache2")
	require.Equal(t, "addr2", sisu)
}
