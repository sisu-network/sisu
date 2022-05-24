package keeper

import (
	"fmt"
	"sort"
	"strconv"
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
	msg1Signer1 := &types.UpdateTokenPrice{
		Signer: signer1,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 5_000_000_000,
			},
		},
	}
	msg2Signer1 := &types.UpdateTokenPrice{
		Signer: signer1,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 6_000_000_000,
			},
		},
	}

	signer2 := "signer2"
	msg1Signer2 := &types.UpdateTokenPrice{
		Signer: signer2,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 10_000_000_000,
			},
		},
	}
	msg2Signer2 := &types.UpdateTokenPrice{
		Signer: signer2,
		TokenPrices: []*types.TokenPrice{
			{
				Id:    token,
				Price: 11_000_000_000,
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
	require.Equal(t, int64(6_000_000_000), record.Records[0].Price)
	record = allPrices[signer2]
	require.Equal(t, int64(11_000_000_000), record.Records[0].Price)
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
		Status:      types.NodeStatus_Candidate,
	}
	node2 := &types.Node{
		ConsensusKey: &types.Pubkey{
			Type:  "ed",
			Bytes: []byte("pubkey2"),
		},
		AccAddress:  "addr2",
		IsValidator: true,
		Status:      types.NodeStatus_Candidate,
	}

	saveNode(store, node1)
	saveNode(store, node2)

	vals := loadNodesByStatus(store, types.NodeStatus_Candidate)

	require.Equal(t, 2, len(vals), "there should be 2 validators")
	require.Equal(t, vals[0].AccAddress, "addr1")
	require.Equal(t, vals[1].AccAddress, "addr2")
}

func TestDefaultKeeper_SetValidators(t *testing.T) {
	t.Parallel()

	store := memstore.NewStore()
	oldValidatorSet := make([]*types.Node, 0)
	for i := 0; i < 2; i++ {
		oldValidatorSet = append(oldValidatorSet, &types.Node{
			Id: "old_val" + strconv.Itoa(i),
			ConsensusKey: &types.Pubkey{
				Type:  "ed",
				Bytes: []byte("old_pubkey" + strconv.Itoa(i)),
			},
			AccAddress:  "old_addr" + strconv.Itoa(i),
			IsValidator: true,
		})
	}

	for _, val := range oldValidatorSet {
		saveNode(store, val)
	}

	newValidatorSet := make([]*types.Node, 0)
	for i := 0; i < 2; i++ {
		newValidatorSet = append(newValidatorSet, &types.Node{
			Id: "new_val" + strconv.Itoa(i),
			ConsensusKey: &types.Pubkey{
				Type:  "ed",
				Bytes: []byte("new_pubkey" + strconv.Itoa(i)),
			},
			AccAddress:  "new_addr" + strconv.Itoa(i),
			IsValidator: true,
		})
	}

	validVals, err := setValidators(store, newValidatorSet)
	require.NoError(t, err)
	require.Len(t, validVals, 2)

	for i, val := range validVals {
		require.Equal(t, "new_val"+strconv.Itoa(i), val.Id)
		require.Equal(t, "new_addr"+strconv.Itoa(i), val.AccAddress)
		require.True(t, val.IsValidator)
		require.Equal(t, []byte("new_pubkey"+strconv.Itoa(i)), val.ConsensusKey.Bytes)
	}
}
