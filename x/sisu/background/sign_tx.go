package background

import (
	"fmt"
	"strconv"

	cosmoscrypto "github.com/cosmos/cosmos-sdk/crypto/types"
	liskcrypto "github.com/sisu-network/deyes/chains/lisk/crypto"

	lisktypes "github.com/sisu-network/deyes/chains/lisk/types"

	ethtypes "github.com/ethereum/go-ethereum/core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/echovl/cardano-go"
	hTypes "github.com/sisu-network/dheart/types"
	libchain "github.com/sisu-network/lib/chain"
	"github.com/sisu-network/lib/log"
	"github.com/sisu-network/sisu/x/sisu/components"
	"github.com/sisu-network/sisu/x/sisu/external"
	"github.com/sisu-network/sisu/x/sisu/types"
)

func SignTxOut(ctx sdk.Context, dheartClient external.DheartClient,
	partyManager components.PartyManager, txOut *types.TxOut) {
	pubKeys := partyManager.GetActivePartyPubkeys()

	if libchain.IsETHBasedChain(txOut.Content.OutChain) {
		signEthTx(ctx, dheartClient, pubKeys, txOut)
	}

	if libchain.IsCardanoChain(txOut.Content.OutChain) {
		signCardanoTx(ctx, dheartClient, pubKeys, txOut)
	}

	if libchain.IsSolanaChain(txOut.Content.OutChain) {
		signSolana(ctx, dheartClient, pubKeys, txOut)
	}

	if libchain.IsLiskChain(txOut.Content.OutChain) {
		signLisk(ctx, dheartClient, pubKeys, txOut)
	}
}

// signEthTx sends a TxOut to dheart for TSS signing.
func signEthTx(ctx sdk.Context, dheartClient external.DheartClient, pubKeys []cosmoscrypto.PubKey,
	tx *types.TxOut) error {
	log.Info("Delivering TXOUT for chain ", tx.Content.OutChain, " tx hash = ", tx.Content.OutHash)
	ethTx := &ethtypes.Transaction{}
	err := ethTx.UnmarshalBinary(tx.Content.OutBytes)
	if err != nil {
		log.Error("signEthTx: cannot unmarshal tx, err =", err)
		return err
	}

	signer := libchain.GetEthChainSigner(tx.Content.OutChain)
	if signer == nil {
		return fmt.Errorf("cannot find signer for chain %s", tx.Content.OutChain)
	}

	hash := signer.Hash(ethTx)
	// 4. Send it to Dheart for signing.
	keysignReq := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_ECDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          getKeysignRequestId(tx.Content.OutChain, ctx.BlockHeight(), tx.Content.OutHash),
				OutChain:    tx.Content.OutChain,
				OutHash:     tx.Content.OutHash,
				Bytes:       tx.Content.OutBytes,
				BytesToSign: hash[:],
			},
		},
	}

	err = dheartClient.KeySign(keysignReq, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}

	return nil
}

func signCardanoTx(ctx sdk.Context, dheartClient external.DheartClient,
	pubKeys []cosmoscrypto.PubKey, txOut *types.TxOut) {
	tx := &cardano.Tx{}
	if err := tx.UnmarshalCBOR(txOut.Content.OutBytes); err != nil {
		log.Error("error when unmarshalling cardano tx out: ", err)
		return
	}

	txHash, err := tx.Hash()
	if err != nil {
		log.Error("error when getting cardano tx hash: ", err)
		return
	}

	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          getKeysignRequestId(txOut.Content.OutChain, ctx.BlockHeight(), txOut.Content.OutHash),
				OutChain:    txOut.Content.OutChain,
				OutHash:     txOut.Content.OutHash,
				BytesToSign: txHash[:],
			},
		},
	}

	err = dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("Keysign: err =", err)
	}
}

func signSolana(ctx sdk.Context, dheartClient external.DheartClient, pubKeys []cosmoscrypto.PubKey,
	txOut *types.TxOut) {
	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          getKeysignRequestId(txOut.Content.OutChain, ctx.BlockHeight(), txOut.Content.OutHash),
				OutChain:    txOut.Content.OutChain,
				OutHash:     txOut.Content.OutHash,
				BytesToSign: txOut.Content.OutBytes,
			},
		},
	}

	err := dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("signSolana: err =", err)
	}
}

func signLisk(ctx sdk.Context, dheartClient external.DheartClient, pubKeys []cosmoscrypto.PubKey,
	txOut *types.TxOut) {
	networkId := lisktypes.NetworkId[txOut.Content.OutChain]
	if len(networkId) == 0 {
		log.Errorf(fmt.Sprintf("cannot find lisk network id for chain %s", txOut.Content.OutChain))
		return
	}

	bytesToSign, err := liskcrypto.GetSigningBytes(networkId, txOut.Content.OutBytes)
	if err != nil {
		log.Errorf("Failed to get lisk bytes to sign, err = %s", err)
		return
	}

	signRequest := &hTypes.KeysignRequest{
		KeyType: libchain.KEY_TYPE_EDDSA,
		KeysignMessages: []*hTypes.KeysignMessage{
			{
				Id:          getKeysignRequestId(txOut.Content.OutChain, ctx.BlockHeight(), txOut.Content.OutHash),
				OutChain:    txOut.Content.OutChain,
				OutHash:     txOut.Content.OutHash,
				Bytes:       txOut.Content.OutBytes,
				BytesToSign: bytesToSign,
			},
		},
	}

	err = dheartClient.KeySign(signRequest, pubKeys)
	if err != nil {
		log.Error("signLisk: err =", err)
	}
}

func getKeysignRequestId(chain string, blockHeight int64, txHash string) string {
	return chain + "_" + strconv.Itoa(int(blockHeight)) + "_" + txHash
}
